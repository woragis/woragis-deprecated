package jobapplications

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitmqQueue struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
	exchange  string
	routingKey string
}

// NewRabbitMQQueue creates a new RabbitMQ-backed job application queue.
func NewRabbitMQQueue(url, queueName, exchange, routingKey string) (Queue, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange
	if err := ch.ExchangeDeclare(
		exchange, // name
		"direct", // kind
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue with dead letter exchange
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			"x-dead-letter-exchange":    "woragis.dlx",
			"x-dead-letter-routing-key": queueName + ".failed",
		}, // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	if err := ch.QueueBind(
		queueName,  // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // arguments
	); err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	return &rabbitmqQueue{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
		exchange:  exchange,
		routingKey: routingKey,
	}, nil
}

func (q *rabbitmqQueue) EnqueueJob(ctx context.Context, job *JobApplicationJob) error {
	jobData, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	err = q.channel.PublishWithContext(
		ctx,
		q.exchange,   // exchange
		q.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         jobData,
			DeliveryMode: amqp.Persistent, // Make message persistent
			Timestamp:    time.Now(),
			MessageId:    job.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish job: %w", err)
	}

	return nil
}

func (q *rabbitmqQueue) DequeueJob(ctx context.Context, timeout time.Duration) (*JobApplicationJob, error) {
	// RabbitMQ doesn't support blocking dequeue in the same way as Redis
	// This method is kept for interface compatibility but won't be used by the worker
	// The worker will consume messages directly
	return nil, fmt.Errorf("DequeueJob not supported for RabbitMQ - use Consume instead")
}

func (q *rabbitmqQueue) GetJob(ctx context.Context, jobID string) (*JobApplicationJob, error) {
	// For RabbitMQ, the job data is in the message body, not stored separately
	// This method is kept for interface compatibility
	return nil, fmt.Errorf("GetJob not supported for RabbitMQ - job data is in message body")
}

func (q *rabbitmqQueue) MarkJobComplete(ctx context.Context, jobID string) error {
	// For RabbitMQ, acknowledgment is handled by the consumer
	// This method is kept for interface compatibility
	return nil
}

func (q *rabbitmqQueue) MarkJobFailed(ctx context.Context, jobID string, errorMsg string) error {
	// For RabbitMQ, failed jobs go to DLQ automatically
	// This method is kept for interface compatibility
	return nil
}

// Close closes the RabbitMQ connection and channel.
func (q *rabbitmqQueue) Close() error {
	if q.channel != nil {
		q.channel.Close()
	}
	if q.conn != nil {
		return q.conn.Close()
	}
	return nil
}
