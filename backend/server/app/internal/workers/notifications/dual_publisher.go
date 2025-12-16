package notifications

import (
	"context"
	"log/slog"
)

// DualPublisher publishes to both Redis (for backward compatibility) and RabbitMQ (new).
// During migration, this ensures messages reach both systems.
type DualPublisher struct {
	redisPublisher    *Publisher
	rabbitmqPublisher *RabbitMQPublisher
	logger            *slog.Logger
}

// NewDualPublisher creates a publisher that publishes to both Redis and RabbitMQ.
func NewDualPublisher(redisPublisher *Publisher, rabbitmqPublisher *RabbitMQPublisher, logger *slog.Logger) *DualPublisher {
	return &DualPublisher{
		redisPublisher:    redisPublisher,
		rabbitmqPublisher: rabbitmqPublisher,
		logger:            logger,
	}
}

// PublishEmailReport publishes to both Redis and RabbitMQ.
func (p *DualPublisher) PublishEmailReport(ctx context.Context, env ReportEnvelope) error {
	// Publish to RabbitMQ (primary)
	if p.rabbitmqPublisher != nil {
		if err := p.rabbitmqPublisher.PublishEmailReport(ctx, env); err != nil {
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
			if p.rabbitmqPublisher == nil {
				return err
			}
		}
	}

	return nil
}

// PublishWhatsAppReport publishes to both Redis and RabbitMQ.
func (p *DualPublisher) PublishWhatsAppReport(ctx context.Context, env ReportEnvelope) error {
	// For now, only publish to Redis (WhatsApp worker not migrated yet)
	if p.redisPublisher != nil {
		return p.redisPublisher.PublishWhatsAppReport(ctx, env)
	}
	return nil
}
