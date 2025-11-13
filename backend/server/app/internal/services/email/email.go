package email

import (
	"context"
	"log/slog"
)

// Sender defines the contract for sending e-mails.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}

// Message represents an outbound email payload.
type Message struct {
	To       string
	Subject  string
	TextBody string
	HTMLBody string
}

// NoopSender logs e-mail dispatch attempts without actually sending them.
type NoopSender struct {
	logger *slog.Logger
}

// NewNoopSender returns a new NoopSender instance.
func NewNoopSender(logger *slog.Logger) *NoopSender {
	return &NoopSender{logger: logger}
}

// Send logs the e-mail dispatch details for observability.
func (s *NoopSender) Send(ctx context.Context, msg Message) error {
	if s.logger != nil {
		s.logger.Info("email: noop send",
			slog.String("to", msg.To),
			slog.String("subject", msg.Subject),
		)
	}

	return nil
}
