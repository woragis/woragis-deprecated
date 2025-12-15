package translations

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/google/uuid"
	"github.com/woragis/backend/server/app/pkg/rabbitmq"
)

const (
	translationQueueName = "translations.queue"
	translationExchange  = "woragis.tasks"
)

type rabbitmqQueue struct {
	queue      *rabbitmq.TaskQueue
	pendingAck map[string]*amqp.Delivery // Store pending deliveries by job ID
	mu         sync.Mutex                 // Protect pendingAck map
}

// NewRabbitMQQueue creates a new RabbitMQ-backed translation queue.
func NewRabbitMQQueue(conn *rabbitmq.Connection) (Queue, error) {
	queue, err := rabbitmq.NewTaskQueue(conn, translationQueueName, translationExchange)
	if err != nil {
		return nil, err
	}

	return &rabbitmqQueue{
		queue:      queue,
		pendingAck: make(map[string]*amqp.Delivery),
	}, nil
}

func (q *rabbitmqQueue) EnqueueJob(ctx context.Context, job *TranslationJob) error {
	if job.ID == "" {
		job.ID = uuid.New().String()
	}

	return q.queue.Publish(ctx, job)
}

func (q *rabbitmqQueue) DequeueJob(ctx context.Context, timeout time.Duration) (*TranslationJob, error) {
	channel := q.queue.Conn().Channel()
	queueName := q.queue.QueueName()
	
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Use basic_get with polling - this doesn't create persistent consumers
	// Poll every 500ms until timeout
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctxWithTimeout.Done():
			return nil, nil // Timeout - no job available
		case <-ticker.C:
			// Try to get a message (non-blocking)
			msg, ok, err := channel.Get(queueName, false) // false = don't auto-ack
			if err != nil {
				return nil, NewDomainError(ErrCodeJobQueueFailure, ErrJobQueueUnavailable)
			}
			
			if !ok {
				// No message available, continue polling
				continue
			}

			// Message received - parse it
			var j TranslationJob
			if err := json.Unmarshal(msg.Body, &j); err != nil {
				msg.Nack(false, false) // Don't requeue invalid messages
				return nil, fmt.Errorf("failed to unmarshal job: %w", err)
			}

			// Store delivery for later acknowledgment after processing
			// Create a copy of the delivery to store (delivery contains channel reference)
			deliveryCopy := msg
			q.mu.Lock()
			q.pendingAck[j.ID] = &deliveryCopy
			q.mu.Unlock()

			return &j, nil
		}
	}
}

func (q *rabbitmqQueue) GetJob(ctx context.Context, jobID string) (*TranslationJob, error) {
	// Note: RabbitMQ doesn't support direct job lookup by ID
	// In RabbitMQ, jobs are consumed from the queue and processed immediately
	// If you need job lookup, consider storing job metadata in the database
	return nil, NewDomainError(ErrCodeNotFound, "translations: job lookup not supported in RabbitMQ - jobs are processed immediately")
}

func (q *rabbitmqQueue) MarkJobComplete(ctx context.Context, jobID string) error {
	// Acknowledge the message
	q.mu.Lock()
	delivery, ok := q.pendingAck[jobID]
	if ok {
		delete(q.pendingAck, jobID)
	}
	q.mu.Unlock()
	
	if ok && delivery != nil {
		delivery.Ack(false)
	}
	return nil
}

func (q *rabbitmqQueue) MarkJobFailed(ctx context.Context, jobID string, errorMsg string) error {
	// Nack the message and requeue for retry
	q.mu.Lock()
	delivery, ok := q.pendingAck[jobID]
	if ok {
		delete(q.pendingAck, jobID)
	}
	q.mu.Unlock()
	
	if ok && delivery != nil {
		delivery.Nack(false, true) // Requeue for retry
	}
	return nil
}
