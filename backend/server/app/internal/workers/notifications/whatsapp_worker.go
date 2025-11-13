package notifications

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/redis/go-redis/v9"

	whatsappservice "github.com/woragis/backend/server/app/internal/services/whatsapp"
)

// StartWhatsAppWorker subscribes to WhatsApp report notifications.
func StartWhatsAppWorker(ctx context.Context, client *redis.Client, notifier whatsappservice.Notifier, logger *slog.Logger) error {
	if client == nil {
		return nil
	}

	sub := client.Subscribe(ctx, whatsappChannel)
	ch := sub.Channel()

	go func() {
		for {
			select {
			case msg := <-ch:
				if msg == nil {
					continue
				}
				var envelope ReportEnvelope
				if err := json.Unmarshal([]byte(msg.Payload), &envelope); err != nil {
					if logger != nil {
						logger.Error("whatsapp worker: invalid payload", slog.Any("error", err))
					}
					continue
				}
				if err := notifier.Send(ctx, envelope.Destination, envelope.TextMessage); err != nil && logger != nil {
					logger.Error("whatsapp worker: send failed", slog.String("user_id", envelope.UserID), slog.Any("error", err))
				}
			case <-ctx.Done():
				_ = sub.Close()
				return
			}
		}
	}()

	return nil
}
