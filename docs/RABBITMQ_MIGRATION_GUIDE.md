# RabbitMQ Migration Guide

## Overview

This document provides a comprehensive guide for migrating from Redis to RabbitMQ for queue management and messaging between workers in the Woragis backend system.

## Current Architecture

### Redis Usage

The current system uses Redis for three main purposes:

1. **Task Queues** (Redis Lists)
   - Job Application Queue (`job-applications:queue`)
   - Translation Queue (`translations:queue`)
   - Resume Generation Queue (`resumes:queue`)

2. **Pub/Sub Messaging** (Redis Channels)
   - Email notifications (`reports.email`)
   - WhatsApp notifications (`reports.whatsapp`)
   - Resume events (`resumes:events`)

3. **Token Storage** (Redis Key-Value)
   - Password reset tokens (`auth:reset:token:*`)
   - User token mapping (`auth:reset:user:*`)

### Current Workers

1. **Resume Worker** (Python) - `backend/server/resume-worker/`
2. **Job Application Worker** (Node.js) - `backend/server/job-application-worker/`
3. **Translation Worker** (Go) - `backend/server/app/cmd/translation-worker/`
4. **Email Worker** (Go) - `backend/server/app/internal/workers/notifications/`
5. **WhatsApp Worker** (Go) - `backend/server/app/internal/workers/notifications/`

---

## RabbitMQ Architecture Design

### Exchange and Queue Strategy

#### 1. Task Queues (Work Queue Pattern)

Use **Direct Exchanges** with durable queues for task processing:

```
Exchange: woragis.tasks (direct, durable)
├── Queue: job-applications.queue (durable, auto-delete=false)
├── Queue: translations.queue (durable, auto-delete=false)
└── Queue: resumes.queue (durable, auto-delete=false)

Routing Keys:
- job-applications.queue -> job-applications.queue
- translations.queue -> translations.queue
- resumes.queue -> resumes.queue
```

**Benefits:**
- Message persistence (survives broker restarts)
- Automatic acknowledgment with retry support
- Dead letter queues for failed messages
- Priority queues for urgent tasks

#### 2. Pub/Sub Notifications (Fanout Exchange)

Use **Fanout Exchanges** for broadcasting notifications:

```
Exchange: woragis.notifications (fanout, durable)
├── Queue: email.notifications (durable, auto-delete=false)
└── Queue: whatsapp.notifications (durable, auto-delete=false)

Exchange: woragis.resume.events (fanout, durable)
└── Queue: resume.events (durable, auto-delete=false)
```

**Benefits:**
- Multiple consumers can receive the same message
- Decoupled publishers and subscribers
- Message persistence

#### 3. Token Storage

**Keep Redis for token storage** - RabbitMQ is not ideal for key-value storage. Continue using Redis for:
- Password reset tokens
- Session tokens
- Other temporary key-value data

**Alternative:** If you want to remove Redis entirely, consider:
- PostgreSQL with TTL cleanup jobs
- In-memory cache with periodic cleanup
- External token service

---

## Implementation Plan

### Phase 1: Infrastructure Setup

#### 1.1 Docker Compose Changes

**File:** `backend/server/docker-compose.yml`

```yaml
# Add RabbitMQ service
rabbitmq:
  image: rabbitmq:3.13-management-alpine
  container_name: woragis-rabbitmq
  environment:
    RABBITMQ_DEFAULT_USER: ${RABBITMQ_USER:-woragis}
    RABBITMQ_DEFAULT_PASS: ${RABBITMQ_PASSWORD:-woragis}
    RABBITMQ_DEFAULT_VHOST: ${RABBITMQ_VHOST:-/woragis}
  healthcheck:
    test: ["CMD", "rabbitmq-diagnostics", "ping"]
    interval: 10s
    timeout: 5s
    retries: 5
  ports:
    - "5672:5672"   # AMQP port
    - "15672:15672" # Management UI
  volumes:
    - rabbitmq-data:/var/lib/rabbitmq
  restart: unless-stopped

# Update service dependencies
services:
  app:
    depends_on:
      rabbitmq:
        condition: service_healthy
    environment:
      RABBITMQ_URL: amqp://woragis:woragis@rabbitmq:5672/woragis
      # Keep REDIS_URL for token storage (optional)
      REDIS_URL: redis://redis:6379/0

  translation-worker:
    depends_on:
      rabbitmq:
        condition: service_healthy
    environment:
      RABBITMQ_URL: amqp://woragis:woragis@rabbitmq:5672/woragis

  job-application-worker:
    depends_on:
      rabbitmq:
        condition: service_healthy
    environment:
      RABBITMQ_URL: amqp://woragis:woragis@rabbitmq:5672/woragis

  resume-worker:
    depends_on:
      rabbitmq:
        condition: service_healthy
    environment:
      RABBITMQ_URL: amqp://woragis:woragis@rabbitmq:5672/woragis

volumes:
  rabbitmq-data:
```

#### 1.2 Environment Variables

**File:** `backend/server/env.sample`

```bash
# RabbitMQ Configuration
RABBITMQ_URL=amqp://woragis:woragis@rabbitmq:5672/woragis
RABBITMQ_USER=woragis
RABBITMQ_PASSWORD=woragis
RABBITMQ_VHOST=/woragis

# Redis (keep for token storage)
REDIS_URL=redis://redis:6379/0
```

---

### Phase 2: Go Backend Implementation

#### 2.1 RabbitMQ Configuration

**File:** `backend/server/app/pkg/config/rabbitmq.go` (NEW)

```go
package config

import "os"

// RabbitMQConfig holds connection details for RabbitMQ.
type RabbitMQConfig struct {
	URL      string
	User     string
	Password string
	VHost    string
}

// LoadRabbitMQConfig reads RabbitMQ configuration from environment variables.
func LoadRabbitMQConfig() RabbitMQConfig {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	user := os.Getenv("RABBITMQ_USER")
	if user == "" {
		user = "guest"
	}

	password := os.Getenv("RABBITMQ_PASSWORD")
	if password == "" {
		password = "guest"
	}

	vhost := os.Getenv("RABBITMQ_VHOST")
	if vhost == "" {
		vhost = "/"
	}

	return RabbitMQConfig{
		URL:      url,
		User:     user,
		Password: password,
		VHost:    vhost,
	}
}
```

#### 2.2 RabbitMQ Connection Manager

**File:** `backend/server/app/pkg/rabbitmq/connection.go` (NEW)

```go
package rabbitmq

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *slog.Logger
}

func NewConnection(url string, logger *slog.Logger) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Connection{
		conn:    conn,
		channel: ch,
		logger:  logger,
	}, nil
}

func (c *Connection) Channel() *amqp.Channel {
	return c.channel
}

func (c *Connection) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Connection) DeclareExchange(name, kind string) error {
	return c.channel.ExchangeDeclare(
		name,  // name
		kind,  // kind (direct, fanout, topic, headers)
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

func (c *Connection) DeclareQueue(name string) (amqp.Queue, error) {
	return c.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    "woragis.dlx",
			"x-dead-letter-routing-key": name + ".failed",
		}, // arguments
	)
}

func (c *Connection) BindQueue(queue, exchange, routingKey string) error {
	return c.channel.QueueBind(
		queue,      // queue name
		routingKey, // routing key
		exchange,  // exchange
		false,      // no-wait
		nil,        // arguments
	)
}
```

#### 2.3 Task Queue Interface (Generic)

**File:** `backend/server/app/pkg/rabbitmq/queue.go` (NEW)

```go
package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type TaskQueue struct {
	conn      *Connection
	queueName string
	exchange  string
}

func NewTaskQueue(conn *Connection, queueName, exchange string) (*TaskQueue, error) {
	// Declare exchange
	if err := conn.DeclareExchange(exchange, "direct"); err != nil {
		return nil, err
	}

	// Declare queue
	_, err := conn.DeclareQueue(queueName)
	if err != nil {
		return nil, err
	}

	// Bind queue to exchange
	if err := conn.BindQueue(queueName, exchange, queueName); err != nil {
		return nil, err
	}

	return &TaskQueue{
		conn:      conn,
		queueName: queueName,
		exchange:  exchange,
	}, nil
}

func (q *TaskQueue) Publish(ctx context.Context, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return q.conn.Channel().PublishWithContext(
		ctx,
		q.exchange,   // exchange
		q.queueName,  // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
			Timestamp:   time.Now(),
		},
	)
}

func (q *TaskQueue) Consume(ctx context.Context, handler func([]byte) error) error {
	msgs, err := q.conn.Channel().Consume(
		q.queueName, // queue
		"",          // consumer
		false,       // auto-ack (manual ack for reliability)
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-msgs:
			if err := handler(msg.Body); err != nil {
				// Reject and requeue on error
				msg.Nack(false, true)
				continue
			}
			// Acknowledge successful processing
			msg.Ack(false)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
```

#### 2.4 Update Job Applications Queue

**File:** `backend/server/app/internal/domains/jobapplications/queue.go`

Replace Redis implementation with RabbitMQ:

```go
package jobapplications

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/backend/server/app/pkg/rabbitmq"
)

const (
	jobApplicationQueueName = "job-applications.queue"
	jobApplicationExchange  = "woragis.tasks"
)

// Queue manages job application jobs in RabbitMQ.
type Queue interface {
	EnqueueJob(ctx context.Context, job *JobApplicationJob) error
	DequeueJob(ctx context.Context, timeout time.Duration) (*JobApplicationJob, error)
	GetJob(ctx context.Context, jobID string) (*JobApplicationJob, error)
	MarkJobComplete(ctx context.Context, jobID string) error
	MarkJobFailed(ctx context.Context, jobID string, errorMsg string) error
}

type rabbitmqQueue struct {
	queue *rabbitmq.TaskQueue
	// Optional: Keep Redis for job storage lookup
	// Or use RabbitMQ message properties for metadata
}

// NewRabbitMQQueue creates a new RabbitMQ-backed job application queue.
func NewRabbitMQQueue(conn *rabbitmq.Connection) (Queue, error) {
	queue, err := rabbitmq.NewTaskQueue(conn, jobApplicationQueueName, jobApplicationExchange)
	if err != nil {
		return nil, err
	}

	return &rabbitmqQueue{
		queue: queue,
	}, nil
}

func (q *rabbitmqQueue) EnqueueJob(ctx context.Context, job *JobApplicationJob) error {
	if job.ID == "" {
		job.ID = uuid.New().String()
	}

	return q.queue.Publish(ctx, job)
}

func (q *rabbitmqQueue) DequeueJob(ctx context.Context, timeout time.Duration) (*JobApplicationJob, error) {
	var job *JobApplicationJob
	var err error

	done := make(chan bool)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	go func() {
		err = q.queue.Consume(ctxWithTimeout, func(body []byte) error {
			var j JobApplicationJob
			if err := json.Unmarshal(body, &j); err != nil {
				return err
			}
			job = &j
			done <- true
			return nil
		})
	}()

	select {
	case <-done:
		return job, err
	case <-ctxWithTimeout.Done():
		return nil, nil // Timeout - no job available
	}
}

func (q *rabbitmqQueue) GetJob(ctx context.Context, jobID string) (*JobApplicationJob, error) {
	// Note: RabbitMQ doesn't support direct job lookup by ID
	// Options:
	// 1. Store job metadata in database
	// 2. Use RabbitMQ management API
	// 3. Keep Redis for job lookup (hybrid approach)
	return nil, fmt.Errorf("direct job lookup not supported in RabbitMQ - use database or hybrid approach")
}

func (q *rabbitmqQueue) MarkJobComplete(ctx context.Context, jobID string) error {
	// Job is automatically removed after acknowledgment
	// If you need job history, store in database
	return nil
}

func (q *rabbitMQQueue) MarkJobFailed(ctx context.Context, jobID string, errorMsg string) error {
	// Failed jobs go to dead letter queue automatically
	// You can also publish to a failed jobs exchange
	return nil
}
```

#### 2.5 Update Translation Queue

**File:** `backend/server/app/internal/domains/translations/queue.go`

Similar pattern to job applications queue:

```go
package translations

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/backend/server/app/pkg/rabbitmq"
)

const (
	translationQueueName = "translations.queue"
	translationExchange  = "woragis.tasks"
)

type rabbitmqQueue struct {
	queue *rabbitmq.TaskQueue
}

func NewRabbitMQQueue(conn *rabbitmq.Connection) (Queue, error) {
	queue, err := rabbitmq.NewTaskQueue(conn, translationQueueName, translationExchange)
	if err != nil {
		return nil, err
	}
	return &rabbitmqQueue{queue: queue}, nil
}

// Implementation similar to job applications queue
```

#### 2.6 Update Resume Queue

**File:** `backend/server/app/internal/domains/resumes/queue.go`

Similar pattern, but also handle event publishing:

```go
package resumes

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/woragis/backend/server/app/pkg/rabbitmq"
)

const (
	resumeQueueName   = "resumes.queue"
	resumeExchange     = "woragis.tasks"
	resumeEventsExchange = "woragis.resume.events"
)

type rabbitmqQueue struct {
	queue      *rabbitmq.TaskQueue
	eventQueue *rabbitmq.TaskQueue
}

func NewRabbitMQQueue(conn *rabbitmq.Connection) (Queue, error) {
	queue, err := rabbitmq.NewTaskQueue(conn, resumeQueueName, resumeExchange)
	if err != nil {
		return nil, err
	}

	// Declare fanout exchange for events
	if err := conn.DeclareExchange(resumeEventsExchange, "fanout"); err != nil {
		return nil, err
	}

	return &rabbitmqQueue{
		queue: queue,
	}, nil
}

func (q *rabbitmqQueue) UpdateJobStatus(ctx context.Context, jobID string, status string, errorMsg *string, errorType *string, retryCount *int, result *ResumeJobResult) error {
	// Update job in database or Redis cache
	
	// Publish event to fanout exchange
	event := map[string]interface{}{
		"job_id":   jobID,
		"status":   status,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	if errorMsg != nil {
		event["error"] = *errorMsg
	}
	if result != nil {
		event["result"] = result
	}

	eventBody, _ := json.Marshal(event)
	return q.conn.Channel().PublishWithContext(
		ctx,
		resumeEventsExchange,
		"", // fanout doesn't use routing key
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         eventBody,
			DeliveryMode: amqp.Persistent,
		},
	)
}
```

#### 2.7 Update Notification Publisher

**File:** `backend/server/app/internal/workers/notifications/publisher.go`

```go
package notifications

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	notificationsExchange = "woragis.notifications"
)

type Publisher struct {
	channel *amqp.Channel
}

func NewPublisher(channel *amqp.Channel) *Publisher {
	return &Publisher{channel: channel}
}

func (p *Publisher) PublishEmailReport(ctx context.Context, env ReportEnvelope) error {
	return p.publish(ctx, env, "email")
}

func (p *Publisher) PublishWhatsAppReport(ctx context.Context, env ReportEnvelope) error {
	return p.publish(ctx, env, "whatsapp")
}

func (p *Publisher) publish(ctx context.Context, env ReportEnvelope, routingKey string) error {
	payload, err := json.Marshal(env)
	if err != nil {
		return err
	}

	return p.channel.PublishWithContext(
		ctx,
		notificationsExchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent,
		},
	)
}
```

#### 2.8 Update Notification Workers

**File:** `backend/server/app/internal/workers/notifications/email_worker.go`

```go
package notifications

import (
	"context"
	"encoding/json"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	emailservice "github.com/woragis/backend/server/app/internal/services/email"
)

const (
	emailQueueName = "email.notifications"
)

func StartEmailWorker(ctx context.Context, channel *amqp.Channel, sender emailservice.Sender, logger *slog.Logger) error {
	// Declare queue
	_, err := channel.QueueDeclare(
		emailQueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return err
	}

	// Bind to notifications exchange
	err = channel.QueueBind(
		emailQueueName,
		"email",
		"woragis.notifications",
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := channel.Consume(
		emailQueueName,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case msg := <-msgs:
				var envelope ReportEnvelope
				if err := json.Unmarshal(msg.Body, &envelope); err != nil {
					if logger != nil {
						logger.Error("email worker: invalid payload", slog.Any("error", err))
					}
					msg.Nack(false, false) // Don't requeue invalid messages
					continue
				}

				message := emailservice.Message{
					To:       envelope.Destination,
					Subject:  envelope.Subject,
					TextBody: envelope.TextMessage,
					HTMLBody: envelope.HTMLMessage,
				}

				if err := sender.Send(ctx, message); err != nil {
					if logger != nil {
						logger.Error("email worker: send failed", slog.String("user_id", envelope.UserID), slog.Any("error", err))
					}
					msg.Nack(false, true) // Requeue on failure
					continue
				}

				msg.Ack(false) // Acknowledge successful processing
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}
```

#### 2.9 Update Main Server

**File:** `backend/server/app/cmd/server/main.go`

```go
// Add RabbitMQ connection
rabbitmqCfg := appconfig.LoadRabbitMQConfig()
rabbitmqConn, err := rabbitmq.NewConnection(rabbitmqCfg.URL, slogLogger)
if err != nil {
	slogLogger.Error("failed to connect to RabbitMQ", slog.Any("error", err))
	return
}
defer rabbitmqConn.Close()

// Initialize exchanges and queues
if err := rabbitmqConn.DeclareExchange("woragis.tasks", "direct"); err != nil {
	slogLogger.Error("failed to declare tasks exchange", slog.Any("error", err))
}

if err := rabbitmqConn.DeclareExchange("woragis.notifications", "fanout"); err != nil {
	slogLogger.Error("failed to declare notifications exchange", slog.Any("error", err))
}

// Update queue initializations
applicationQueue, err := jobapplicationsdomain.NewRabbitMQQueue(rabbitmqConn)
if err != nil {
	slogLogger.Error("failed to create application queue", slog.Any("error", err))
}

translationQueue, err := translationsdomain.NewRabbitMQQueue(rabbitmqConn)
if err != nil {
	slogLogger.Error("failed to create translation queue", slog.Any("error", err))
}

resumeQueue, err := resumesdomain.NewRabbitMQQueue(rabbitmqConn)
if err != nil {
	slogLogger.Error("failed to create resume queue", slog.Any("error", err))
}

// Update notification publisher
publisher := notifications.NewPublisher(rabbitmqConn.Channel())

// Update workers
if err := notifications.StartEmailWorker(workerCtx, rabbitmqConn.Channel(), emailSender, slogLogger); err != nil {
	slogLogger.Error("failed to start email worker", slog.Any("error", err))
}
```

---

### Phase 3: Node.js Worker Implementation

#### 3.1 Update Job Application Worker

**File:** `backend/server/job-application-worker/src/queue.js`

Replace Redis with RabbitMQ using `amqplib`:

```javascript
import amqp from 'amqplib';
import { logger } from './utils/logger.js';

const QUEUE_NAME = 'job-applications.queue';
const EXCHANGE = 'woragis.tasks';

export class Queue {
  constructor() {
    this.connection = null;
    this.channel = null;
  }

  async connect() {
    const rabbitmqUrl = process.env.RABBITMQ_URL || 'amqp://guest:guest@localhost:5672/';
    
    try {
      this.connection = await amqp.connect(rabbitmqUrl);
      this.channel = await this.connection.createChannel();
      
      // Declare exchange
      await this.channel.assertExchange(EXCHANGE, 'direct', { durable: true });
      
      // Declare queue
      await this.channel.assertQueue(QUEUE_NAME, {
        durable: true,
        arguments: {
          'x-dead-letter-exchange': 'woragis.dlx',
          'x-dead-letter-routing-key': QUEUE_NAME + '.failed'
        }
      });
      
      // Bind queue to exchange
      await this.channel.bindQueue(QUEUE_NAME, EXCHANGE, QUEUE_NAME);
      
      logger.info('Connected to RabbitMQ', { url: rabbitmqUrl });
    } catch (err) {
      logger.error('RabbitMQ connection error', { error: err.message });
      throw err;
    }
  }

  async enqueueJob(job) {
    if (!job.id) {
      job.id = `${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    }

    const message = Buffer.from(JSON.stringify(job));
    
    return this.channel.publish(EXCHANGE, QUEUE_NAME, message, {
      persistent: true,
      messageId: job.id,
      timestamp: Date.now()
    });
  }

  async dequeueJob(timeout = 5000) {
    return new Promise((resolve, reject) => {
      const timer = setTimeout(() => {
        this.channel.cancel(consumerTag);
        resolve(null);
      }, timeout);

      const consumerTag = this.channel.consume(QUEUE_NAME, (msg) => {
        clearTimeout(timer);
        this.channel.cancel(consumerTag);
        
        if (!msg) {
          resolve(null);
          return;
        }

        try {
          const job = JSON.parse(msg.content.toString());
          // Store message for acknowledgment
          job._message = msg;
          resolve(job);
        } catch (err) {
          logger.error('Failed to parse job', { error: err.message });
          this.channel.nack(msg, false, false); // Don't requeue invalid messages
          resolve(null);
        }
      }, { noAck: false });
    });
  }

  async markJobComplete(job) {
    if (job._message) {
      this.channel.ack(job._message);
      delete job._message;
    }
  }

  async markJobFailed(job, errorMsg) {
    logger.error('Job failed', { jobId: job.id, errorMsg });
    if (job._message) {
      this.channel.nack(job._message, false, true); // Requeue for retry
      delete job._message;
    }
  }

  async disconnect() {
    if (this.channel) {
      await this.channel.close();
    }
    if (this.connection) {
      await this.connection.close();
    }
  }
}
```

**File:** `backend/server/job-application-worker/src/worker.js`

Update to use new queue methods:

```javascript
async processApplication(job) {
  // ... existing code ...
  
  try {
    await this.scraper.applyToJob(job, coverLetter);
    
    await this.db.updateApplication(application.id, {
      status: 'applied',
      coverLetter: coverLetter,
      appliedAt: new Date(),
    });
    
    // Acknowledge job
    await this.queue.markJobComplete(job);
  } catch (error) {
    await this.db.updateApplication(application.id, {
      status: 'failed',
      errorMessage: error.message,
    });
    
    // Reject and requeue
    await this.queue.markJobFailed(job, error.message);
    throw error;
  }
}
```

**File:** `backend/server/job-application-worker/package.json`

Add dependency:

```json
{
  "dependencies": {
    "amqplib": "^0.10.3"
  }
}
```

---

### Phase 4: Python Worker Implementation

#### 4.1 Update Resume Worker

**File:** `backend/server/resume-worker/src/queue.py` (NEW)

```python
import os
import json
import pika
import logging
from typing import Optional, Dict, Any

logger = logging.getLogger(__name__)

QUEUE_NAME = 'resumes.queue'
EXCHANGE = 'woragis.tasks'

class Queue:
    def __init__(self):
        self.connection = None
        self.channel = None
    
    def connect(self):
        rabbitmq_url = os.getenv('RABBITMQ_URL', 'amqp://guest:guest@localhost:5672/')
        
        try:
            params = pika.URLParameters(rabbitmq_url)
            self.connection = pika.BlockingConnection(params)
            self.channel = self.connection.channel()
            
            # Declare exchange
            self.channel.exchange_declare(
                exchange=EXCHANGE,
                exchange_type='direct',
                durable=True
            )
            
            # Declare queue
            self.channel.queue_declare(
                queue=QUEUE_NAME,
                durable=True,
                arguments={
                    'x-dead-letter-exchange': 'woragis.dlx',
                    'x-dead-letter-routing-key': QUEUE_NAME + '.failed'
                }
            )
            
            # Bind queue
            self.channel.queue_bind(
                exchange=EXCHANGE,
                queue=QUEUE_NAME,
                routing_key=QUEUE_NAME
            )
            
            logger.info(f'Connected to RabbitMQ: {rabbitmq_url}')
        except Exception as e:
            logger.error(f'RabbitMQ connection error: {e}')
            raise
    
    def enqueue_job(self, job: Dict[str, Any]) -> str:
        if 'id' not in job:
            import uuid
            job['id'] = str(uuid.uuid4())
        
        message = json.dumps(job)
        
        self.channel.basic_publish(
            exchange=EXCHANGE,
            routing_key=QUEUE_NAME,
            body=message,
            properties=pika.BasicProperties(
                delivery_mode=2,  # Make message persistent
                message_id=job['id'],
                timestamp=int(time.time())
            )
        )
        
        return job['id']
    
    def dequeue_job(self, timeout: int = 5) -> Optional[Dict[str, Any]]:
        method_frame, header_frame, body = self.channel.basic_get(
            queue=QUEUE_NAME,
            auto_ack=False
        )
        
        if method_frame:
            try:
                job = json.loads(body)
                job['_delivery_tag'] = method_frame.delivery_tag
                return job
            except json.JSONDecodeError as e:
                logger.error(f'Failed to parse job: {e}')
                self.channel.basic_nack(
                    delivery_tag=method_frame.delivery_tag,
                    requeue=False
                )
                return None
        
        return None
    
    def mark_job_complete(self, job: Dict[str, Any]):
        if '_delivery_tag' in job:
            self.channel.basic_ack(delivery_tag=job['_delivery_tag'])
            del job['_delivery_tag']
    
    def mark_job_failed(self, job: Dict[str, Any], error_msg: str):
        logger.error(f'Job failed: {job.get("id")}, error: {error_msg}')
        if '_delivery_tag' in job:
            self.channel.basic_nack(
                delivery_tag=job['_delivery_tag'],
                requeue=True  # Requeue for retry
            )
            del job['_delivery_tag']
    
    def disconnect(self):
        if self.connection and not self.connection.is_closed:
            self.connection.close()
```

**File:** `backend/server/resume-worker/requirements.txt`

Add dependency:

```
pika==1.3.2
```

**Note:** The resume worker currently doesn't use a queue - it's called directly. You may want to add queue support or keep it as-is.

---

## Migration Strategy

### Option 1: Big Bang Migration (Not Recommended)

- Stop all services
- Deploy RabbitMQ
- Update all code
- Restart services
- **Risk:** High downtime, potential data loss

### Option 2: Gradual Migration (Recommended)

#### Step 1: Deploy RabbitMQ alongside Redis
- Add RabbitMQ to docker-compose
- Keep Redis running
- No code changes yet

#### Step 2: Dual Write (Hybrid Approach)
- Write to both Redis and RabbitMQ
- Read from Redis (existing workers)
- New workers read from RabbitMQ
- Gradually migrate workers one by one

#### Step 3: Migrate Workers One by One
1. **Translation Worker** (simplest, good test)
2. **Job Application Worker**
3. **Resume Worker** (if adding queue support)
4. **Notification Workers** (Email, WhatsApp)

#### Step 4: Remove Redis Dependencies
- Once all workers migrated
- Keep Redis only for token storage (or migrate to DB)
- Remove Redis queue code

---

## Testing Strategy

### 1. Unit Tests

Test queue implementations independently:

```go
func TestRabbitMQQueue_EnqueueDequeue(t *testing.T) {
	// Test enqueue and dequeue operations
}
```

### 2. Integration Tests

Test with real RabbitMQ instance:

```go
func TestJobApplicationQueue_Integration(t *testing.T) {
	// Test full flow: enqueue -> worker processes -> complete
}
```

### 3. Load Testing

- Test message throughput
- Test with multiple workers
- Test failure scenarios (worker crashes, network issues)

### 4. Migration Testing

- Test dual-write mode
- Test failover scenarios
- Test message loss prevention

---

## Monitoring and Observability

### RabbitMQ Management UI

Access at `http://localhost:15672` (default credentials: guest/guest)

**Monitor:**
- Queue lengths
- Message rates
- Consumer connections
- Failed messages

### Metrics to Track

1. **Queue Depth**: Messages waiting to be processed
2. **Message Rate**: Messages/second published/consumed
3. **Acknowledgment Rate**: Success vs failure rate
4. **Consumer Lag**: Time between publish and consume
5. **Dead Letter Queue**: Failed messages

### Logging

Add structured logging for:
- Message published (with job ID)
- Message consumed (with job ID)
- Processing success/failure
- Queue connection issues

---

## Dead Letter Queue (DLQ) Setup

Configure DLQ for failed messages:

```go
// Declare dead letter exchange
conn.DeclareExchange("woragis.dlx", "direct")

// Declare failed queues
conn.DeclareQueue("job-applications.queue.failed")
conn.DeclareQueue("translations.queue.failed")
conn.DeclareQueue("resumes.queue.failed")

// Bind failed queues
conn.BindQueue("job-applications.queue.failed", "woragis.dlx", "job-applications.queue.failed")
```

**Benefits:**
- Failed messages don't get lost
- Can inspect and retry manually
- Alert on DLQ growth

---

## Performance Considerations

### Prefetch Count

Limit unacknowledged messages per worker:

```go
// In worker initialization
channel.Qos(
	10,     // prefetch count
	0,      // prefetch size (0 = unlimited)
	false,  // global (false = per consumer)
)
```

**Benefits:**
- Better load distribution
- Prevents one slow worker from hogging messages

### Message Persistence

All messages should be persistent:

```go
amqp.Publishing{
	DeliveryMode: amqp.Persistent, // Survives broker restart
}
```

### Connection Pooling

Reuse connections and channels:
- One connection per application
- One channel per worker/queue
- Don't create new connections for each message

---

## Troubleshooting

### Common Issues

1. **Messages not being consumed**
   - Check queue bindings
   - Verify exchange type matches
   - Check consumer is running

2. **Messages lost**
   - Verify message persistence enabled
   - Check durable queues
   - Verify acknowledgments

3. **High memory usage**
   - Check queue lengths
   - Set message TTL
   - Monitor prefetch counts

4. **Connection failures**
   - Implement reconnection logic
   - Use connection recovery callbacks
   - Monitor connection health

### Health Checks

Add health check endpoints:

```go
func (c *Connection) HealthCheck() error {
	if c.conn.IsClosed() {
		return fmt.Errorf("connection closed")
	}
	return nil
}
```

---

## Rollback Plan

If migration fails:

1. **Immediate Rollback**
   - Revert code changes
   - Restart services with Redis
   - Messages in RabbitMQ will be lost (if not dual-writing)

2. **Gradual Rollback**
   - Stop writing to RabbitMQ
   - Drain remaining messages
   - Switch workers back to Redis

3. **Data Recovery**
   - If dual-writing: no data loss
   - If not: messages in RabbitMQ queues may need manual processing

---

## Benefits of Migration

1. **Reliability**
   - Message persistence
   - Automatic retries
   - Dead letter queues

2. **Scalability**
   - Better load distribution
   - Multiple consumers per queue
   - Priority queues

3. **Observability**
   - Management UI
   - Better monitoring
   - Message tracing

4. **Features**
   - Message TTL
   - Delayed messages
   - Routing flexibility

---

## Next Steps

1. **Review this document** with the team
2. **Set up RabbitMQ** in development environment
3. **Implement one queue** (start with translations - simplest)
4. **Test thoroughly** before migrating others
5. **Monitor** queue performance and errors
6. **Gradually migrate** remaining queues
7. **Remove Redis** queue dependencies (keep for tokens if needed)

---

## Additional Resources

- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html)
- [AMQP 0-9-1 Model](https://www.rabbitmq.com/tutorials/amqp-concepts.html)
- [RabbitMQ Best Practices](https://www.rabbitmq.com/best-practices.html)
- [Go RabbitMQ Client](https://github.com/rabbitmq/amqp091-go)
- [Node.js RabbitMQ Client](https://github.com/amqp-node/amqplib)
- [Python RabbitMQ Client](https://pika.readthedocs.io/)

---

## Questions to Consider

1. **Do you need to keep Redis for token storage?**
   - If yes: Keep Redis service, remove queue/pub-sub code
   - If no: Migrate tokens to PostgreSQL or external service

2. **Do you need job lookup by ID?**
   - RabbitMQ doesn't support this natively
   - Options: Database, Redis cache, or RabbitMQ management API

3. **What's your message retention policy?**
   - Set TTL on messages
   - Configure queue max length
   - Archive old messages

4. **Do you need message priorities?**
   - RabbitMQ supports priority queues
   - Configure when declaring queues

---

**Document Version:** 1.0  
**Last Updated:** 2025-01-XX  
**Author:** AI Assistant  
**Status:** Draft - Ready for Review

