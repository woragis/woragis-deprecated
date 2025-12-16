package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/woragis/backend/email-worker/internal/config"
	"github.com/woragis/backend/email-worker/internal/queue"
	"github.com/woragis/backend/email-worker/internal/sender"
	"github.com/woragis/backend/email-worker/pkg/logger"
)

func main() {
	// Setup structured logger
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	logger := logger.New(env)

	logger.Info("Starting email worker", "env", env)

	// Load configuration
	emailCfg, err := config.LoadEmailConfig()
	if err != nil {
		logger.Error("Failed to load email config", slog.Any("error", err))
		os.Exit(1)
	}

	if !emailCfg.Enabled() {
		logger.Error("Email configuration not enabled (SMTP_HOST and SMTP_FROM required)")
		os.Exit(1)
	}

	rabbitmqCfg := config.LoadRabbitMQConfig()
	workerCfg := config.LoadWorkerConfig()

	// Initialize email sender
	emailSender, err := sender.NewSMTPSender(emailCfg, logger)
	if err != nil {
		logger.Error("Failed to create email sender", slog.Any("error", err))
		os.Exit(1)
	}

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
	emailQueue, err := queue.NewQueue(
		conn,
		workerCfg.QueueName,
		workerCfg.Exchange,
		workerCfg.RoutingKey,
		logger,
	)
	if err != nil {
		logger.Error("Failed to create email queue", slog.Any("error", err))
		os.Exit(1)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start consuming in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- emailQueue.Consume(ctx, func(envelope queue.EmailEnvelope) error {
			// Convert envelope to email message
			msg := sender.Message{
				To:       envelope.Destination,
				Subject:  envelope.Subject,
				TextBody: envelope.TextMessage,
				HTMLBody: envelope.HTMLMessage,
			}

			// Send email
			if err := emailSender.Send(ctx, msg); err != nil {
				logger.Error("Failed to send email",
					slog.String("user_id", envelope.UserID),
					slog.String("destination", envelope.Destination),
					slog.String("subject", envelope.Subject),
					slog.Any("error", err),
				)
				return err
			}

			logger.Info("Email sent successfully",
				slog.String("user_id", envelope.UserID),
				slog.String("destination", envelope.Destination),
				slog.String("subject", envelope.Subject),
			)

			return nil
		})
	}()

	// Wait for signal or error
	select {
	case sig := <-sigChan:
		logger.Info("Received shutdown signal", slog.String("signal", sig.String()))
		cancel()
	case err := <-errChan:
		if err != nil {
			logger.Error("Email queue consumer error", slog.Any("error", err))
			os.Exit(1)
		}
	}

	logger.Info("Email worker stopped")
}
