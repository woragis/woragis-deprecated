package notifications

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	emailChannel    = "reports.email"
	whatsappChannel = "reports.whatsapp"
)

// ReportEnvelope represents the payload sent to workers.
type ReportEnvelope struct {
	UserID      string `json:"user_id"`
	Subject     string `json:"subject"`
	TextMessage string `json:"text_message"`
	HTMLMessage string `json:"html_message,omitempty"`
	Destination string `json:"destination,omitempty"`
}

// Publisher publishes notifications to Redis channels.
type Publisher struct {
	client *redis.Client
}

// NewPublisher constructs a Publisher.
func NewPublisher(client *redis.Client) *Publisher {
	return &Publisher{client: client}
}

// PublishEmailReport publishes an email summary.
func (p *Publisher) PublishEmailReport(ctx context.Context, env ReportEnvelope) error {
	return p.publish(ctx, emailChannel, env)
}

// PublishWhatsAppReport publishes a WhatsApp summary.
func (p *Publisher) PublishWhatsAppReport(ctx context.Context, env ReportEnvelope) error {
	return p.publish(ctx, whatsappChannel, env)
}

func (p *Publisher) publish(ctx context.Context, channel string, env ReportEnvelope) error {
	if p.client == nil {
		return fmt.Errorf("redis client not configured")
	}

	payload, err := json.Marshal(env)
	if err != nil {
		return err
	}

	return p.client.Publish(ctx, channel, payload).Err()
}
