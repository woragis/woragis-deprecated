package notifications

import (
	"context"
	"log/slog"
)

// DualPublisher publishes to both Redis (for backward compatibility) and RabbitMQ (new).
// During migration, this ensures messages reach both systems.
type DualPublisher struct {
	redisPublisher         *Publisher
	rabbitmqEmailPublisher *RabbitMQPublisher
	rabbitmqWhatsAppPublisher *RabbitMQPublisher
	logger                 *slog.Logger
}

// NewDualPublisher creates a publisher that publishes to both Redis and RabbitMQ.
func NewDualPublisher(redisPublisher *Publisher, rabbitmqEmailPublisher *RabbitMQPublisher, rabbitmqWhatsAppPublisher *RabbitMQPublisher, logger *slog.Logger) *DualPublisher {
	return &DualPublisher{
		redisPublisher:          redisPublisher,
		rabbitmqEmailPublisher:  rabbitmqEmailPublisher,
		rabbitmqWhatsAppPublisher: rabbitmqWhatsAppPublisher,
		logger:                  logger,
	}
}

// PublishEmailReport publishes to both Redis and RabbitMQ.
func (p *DualPublisher) PublishEmailReport(ctx context.Context, env ReportEnvelope) error {
	// Publish to RabbitMQ (primary)
	if p.rabbitmqEmailPublisher != nil {
		if err := p.rabbitmqEmailPublisher.PublishEmailReport(ctx, env); err != nil {
			if p.logger != nil {
				p.logger.Error("failed to publish email to RabbitMQ", slog.Any("error", err))
			}
			// Continue to try Redis as fallback
		}
	}

	// Publish to Redis (for backward compatibility during migration)
	// TODO: Remove this after migration is complete
	if p.redisPublisher != nil {
		if err := p.redisPublisher.PublishEmailReport(ctx, env); err != nil {
			if p.logger != nil {
				p.logger.Error("failed to publish email to Redis", slog.Any("error", err))
			}
			// If both fail, return the last error
			if p.rabbitmqEmailPublisher == nil {
				return err
			}
		}
	}

	return nil
}

// PublishWhatsAppReport publishes to both Redis and RabbitMQ.
func (p *DualPublisher) PublishWhatsAppReport(ctx context.Context, env ReportEnvelope) error {
	// Publish to RabbitMQ (primary)
	if p.rabbitmqWhatsAppPublisher != nil {
		if err := p.rabbitmqWhatsAppPublisher.PublishWhatsAppReport(ctx, env); err != nil {
			if p.logger != nil {
				p.logger.Error("failed to publish whatsapp to RabbitMQ", slog.Any("error", err))
			}
			// Continue to try Redis as fallback
		}
	}

	// Publish to Redis (for backward compatibility during migration)
	// TODO: Remove this after migration is complete
	if p.redisPublisher != nil {
		if err := p.redisPublisher.PublishWhatsAppReport(ctx, env); err != nil {
			if p.logger != nil {
				p.logger.Error("failed to publish whatsapp to Redis", slog.Any("error", err))
			}
			// If both fail, return the last error
			if p.rabbitmqWhatsAppPublisher == nil {
				return err
			}
		}
	}

	return nil
}
