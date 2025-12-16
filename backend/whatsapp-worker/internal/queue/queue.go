package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

// WhatsAppEnvelope represents the WhatsApp message payload from RabbitMQ.
type WhatsAppEnvelope struct {
	UserID      string `json:"user_id"`
	TextMessage string `json:"text_message"`
	Destination string `json:"destination,omitempty"`
}

// Queue handles RabbitMQ operations for WhatsApp messages.
type Queue struct {
	conn      *Connection
	queueName string
	exchange  string
	routingKey string
	logger    *slog.Logger
}

// Connection wraps RabbitMQ connection and channel.
type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewConnection creates a new RabbitMQ connection.
func NewConnection(url string) (*Connection, error) {
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
	}, nil
}

// Channel returns the AMQP channel.
func (c *Connection) Channel() *amqp.Channel {
	return c.channel
}

// Close closes the connection and channel.
func (c *Connection) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsClosed checks if the RabbitMQ connection is closed
func (c *Connection) IsClosed() bool {
	if c.conn == nil {
		return true
	}
	return c.conn.IsClosed()
}

// NewQueue creates a new WhatsApp queue.
func NewQueue(conn *Connection, queueName, exchange, routingKey string, logger *slog.Logger) (*Queue, error) {
	// Declare exchange
	if err := conn.channel.ExchangeDeclare(
		exchange, // name
		"direct", // kind
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	); err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue with dead letter exchange
	_, err := conn.channel.QueueDeclare(
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
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	if err := conn.channel.QueueBind(
		queueName,  // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	// Set QoS to process one message at a time (configurable)
	prefetchCount := 1
	if err := conn.channel.Qos(prefetchCount, 0, false); err != nil {
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return &Queue{
		conn:      conn,
		queueName: queueName,
		exchange:  exchange,
		routingKey: routingKey,
		logger:    logger,
	}, nil
}

// Consume starts consuming messages from the queue.
func (q *Queue) Consume(ctx context.Context, handler func(WhatsAppEnvelope) error) error {
	msgs, err := q.conn.channel.Consume(
		q.queueName, // queue
		"",          // consumer tag
		false,       // auto-ack (manual ack for reliability)
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	q.logger.Info("Started consuming WhatsApp messages",
		slog.String("queue", q.queueName),
		slog.String("exchange", q.exchange),
	)

	for {
		select {
		case <-ctx.Done():
			q.logger.Info("Stopping WhatsApp queue consumer")
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				q.logger.Warn("Message channel closed")
				return fmt.Errorf("message channel closed")
			}

			// Parse message
			var envelope WhatsAppEnvelope
			if err := json.Unmarshal(msg.Body, &envelope); err != nil {
				q.logger.Error("Failed to parse WhatsApp message",
					slog.Any("error", err),
				)
				// Reject message without requeue (invalid format)
				msg.Nack(false, false)
				continue
			}

			// Process message
			if err := handler(envelope); err != nil {
				q.logger.Error("Failed to process WhatsApp message",
					slog.String("user_id", envelope.UserID),
					slog.String("destination", envelope.Destination),
					slog.Any("error", err),
				)
				// Reject and requeue for retry
				msg.Nack(false, true)
				continue
			}

			// Acknowledge successful processing
			if err := msg.Ack(false); err != nil {
				q.logger.Error("Failed to acknowledge message",
					slog.Any("error", err),
				)
			} else {
				q.logger.Debug("WhatsApp message processed successfully",
					slog.String("user_id", envelope.UserID),
					slog.String("destination", envelope.Destination),
				)
			}
		}
	}
}
