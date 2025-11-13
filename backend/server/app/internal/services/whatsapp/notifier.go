package whatsapp

import (
	"context"
	"log/slog"
)

// Notifier defines WhatsApp notification contract.
type Notifier interface {
	Send(ctx context.Context, to string, message string) error
}

// NoopNotifier logs WhatsApp notifications without sending them.
type NoopNotifier struct {
	logger *slog.Logger
}

// NewNoopNotifier builds a new no-op notifier.
func NewNoopNotifier(logger *slog.Logger) *NoopNotifier {
	return &NoopNotifier{logger: logger}
}

// Send logs the outgoing message.
func (n *NoopNotifier) Send(ctx context.Context, to string, message string) error {
	if n.logger != nil {
		n.logger.Info("whatsapp: noop send", slog.String("to", to), slog.String("message", message))
	}
	return nil
}

