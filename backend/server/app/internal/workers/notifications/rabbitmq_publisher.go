package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQPublisher publishes notifications to RabbitMQ queues.
type RabbitMQPublisher struct {
	channel    *amqp.Channel
	exchange   string
	routingKey string
}

// NewRabbitMQPublisher creates a new RabbitMQ publisher.
func NewRabbitMQPublisher(channel *amqp.Channel, exchange, routingKey string) (*RabbitMQPublisher, error) {
	// Declare exchange if it doesn't exist
	if err := channel.ExchangeDeclare(
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

	return &RabbitMQPublisher{
		channel:    channel,
		exchange:   exchange,
		routingKey: routingKey,
	}, nil
}

// PublishEmailReport publishes an email message to RabbitMQ.
func (p *RabbitMQPublisher) PublishEmailReport(ctx context.Context, env ReportEnvelope) error {
	payload, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("failed to marshal email envelope: %w", err)
	}

	return p.channel.PublishWithContext(
		ctx,
		p.exchange,   // exchange
		p.routingKey, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent, // Make message persistent
			Timestamp:    time.Now(),
		},
	)
}
