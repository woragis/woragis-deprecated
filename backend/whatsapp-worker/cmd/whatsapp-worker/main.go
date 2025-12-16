package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/woragis/backend/whatsapp-worker/internal/config"
	"github.com/woragis/backend/whatsapp-worker/internal/queue"
	"github.com/woragis/backend/whatsapp-worker/internal/notifier"
	"github.com/woragis/backend/whatsapp-worker/pkg/logger"
)

func main() {
	// Setup structured logger
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	logger := logger.New(env)

	logger.Info("Starting WhatsApp worker", "env", env)

	// Load configuration
	whatsappCfg := config.LoadWhatsAppConfig()
	rabbitmqCfg := config.LoadRabbitMQConfig()
	workerCfg := config.LoadWorkerConfig()

	// Initialize WhatsApp notifier
	whatsappNotifier, err := notifier.NewWhatsmeowNotifier(whatsappCfg.SessionPath, logger)
	if err != nil {
		logger.Error("Failed to create WhatsApp notifier", slog.Any("error", err))
		os.Exit(1)
	}

	// Connect to WhatsApp in background (non-blocking)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start WhatsApp connection in background
	// Note: WhatsApp connection may require QR code scan, so it runs asynchronously
	go func() {
		logger.Info("Initiating WhatsApp connection...")
		if err := whatsappNotifier.Connect(ctx); err != nil {
			logger.Warn("Failed to connect to WhatsApp (will retry on message send)", slog.Any("error", err))
		} else {
			logger.Info("WhatsApp connection established")
		}
	}()

	// Connect to RabbitMQ
	logger.Info("Connecting to RabbitMQ", slog.String("url", rabbitmqCfg.URL))
	conn, err := queue.NewConnection(rabbitmqCfg.URL)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", slog.Any("error", err))
		os.Exit(1)
	}
	defer conn.Close()

	logger.Info("Connected to RabbitMQ")

	// Create queue
	whatsappQueue, err := queue.NewQueue(
		conn,
		workerCfg.QueueName,
		workerCfg.Exchange,
		workerCfg.RoutingKey,
		logger,
	)
	if err != nil {
		logger.Error("Failed to create WhatsApp queue", slog.Any("error", err))
		os.Exit(1)
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start consuming in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- whatsappQueue.Consume(ctx, func(envelope queue.WhatsAppEnvelope) error {
			// Validate destination
			if envelope.Destination == "" {
				logger.Error("Missing destination in WhatsApp message",
					slog.String("user_id", envelope.UserID),
				)
				return fmt.Errorf("missing destination")
			}

			// Send WhatsApp message
			if err := whatsappNotifier.Send(ctx, envelope.Destination, envelope.TextMessage); err != nil {
				logger.Error("Failed to send WhatsApp message",
					slog.String("user_id", envelope.UserID),
					slog.String("destination", envelope.Destination),
					slog.Any("error", err),
				)
				return err
			}

			logger.Info("WhatsApp message sent successfully",
				slog.String("user_id", envelope.UserID),
				slog.String("destination", envelope.Destination),
			)

			return nil
		})
	}()

	// Wait for signal or error
	select {
	case sig := <-sigChan:
		logger.Info("Received shutdown signal", slog.String("signal", sig.String()))
		cancel()
		whatsappNotifier.Disconnect()
	case err := <-errChan:
		if err != nil {
			logger.Error("WhatsApp queue consumer error", slog.Any("error", err))
			os.Exit(1)
		}
	}

	logger.Info("WhatsApp worker stopped")
}
