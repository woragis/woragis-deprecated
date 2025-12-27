package notifications

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/redis/go-redis/v9"

	emailservice "github.com/woragis/backend/server/app/internal/services/email"
)

// StartEmailWorker subscribes to email report notifications.
func StartEmailWorker(ctx context.Context, client *redis.Client, sender emailservice.Sender, logger *slog.Logger) error {
	if client == nil {
		return nil
	}

	sub := client.Subscribe(ctx, emailChannel)
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
						logger.Error("email worker: invalid payload", slog.Any("error", err))
					}
					continue
				}

				// Validate envelope
				if err := ValidateReportEnvelope(&envelope); err != nil {
					if logger != nil {
						logger.Error("email worker: validation failed", slog.Any("error", err))
					}
					continue
				}

				message := emailservice.Message{
					To:       envelope.Destination,
					Subject:  envelope.Subject,
					TextBody: envelope.TextMessage,
					HTMLBody: envelope.HTMLMessage,
				}

				// Validate message
				if err := emailservice.ValidateMessage(message); err != nil {
					if logger != nil {
						logger.Error("email worker: message validation failed", slog.Any("error", err))
					}
					continue
				}

				if err := sender.Send(ctx, message); err != nil && logger != nil {
					logger.Error("email worker: send failed", slog.String("user_id", envelope.UserID), slog.Any("error", err))
				}
			case <-ctx.Done():
				_ = sub.Close()
				return
			}
		}
	}()

	return nil
}
