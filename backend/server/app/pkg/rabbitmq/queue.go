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

func (q *TaskQueue) Conn() *Connection {
	return q.conn
}

func (q *TaskQueue) QueueName() string {
	return q.queueName
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
			Timestamp:    time.Now(),
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

// ConsumeOne consumes a single message with timeout, returns the message and delivery tag
func (q *TaskQueue) ConsumeOne(ctx context.Context, timeout time.Duration) ([]byte, *amqp.Delivery, error) {
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
		return nil, nil, err
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	select {
	case msg := <-msgs:
		return msg.Body, &msg, nil
	case <-ctxWithTimeout.Done():
		return nil, nil, nil // Timeout - no message available
	}
}
