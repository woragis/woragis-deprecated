package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/woragis/backend/email-worker/internal/config"
	"github.com/woragis/backend/email-worker/internal/queue"
	"github.com/woragis/backend/email-worker/internal/sender"
	"github.com/woragis/backend/email-worker/pkg/health"
	"github.com/woragis/backend/email-worker/pkg/logger"
	appmetrics "github.com/woragis/backend/email-worker/pkg/metrics"
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

	// Connect to RabbitMQ with retry logic
	logger.Info("Connecting to RabbitMQ", slog.String("url", rabbitmqCfg.URL))
	var conn *queue.Connection
	const maxRabbitMQAttempts = 5
	for attempt := 1; attempt <= maxRabbitMQAttempts; attempt++ {
		var err error
		conn, err = queue.NewConnection(rabbitmqCfg.URL)
		if err != nil {
			logger.Warn("RabbitMQ connection failed, retrying...",
				slog.Int("attempt", attempt),
				slog.Int("max_attempts", maxRabbitMQAttempts),
				slog.Any("error", err),
			)
			if attempt < maxRabbitMQAttempts {
				time.Sleep(time.Duration(attempt) * time.Second)
				continue
			}
			logger.Error("RabbitMQ connection failed after multiple attempts", slog.Any("error", err))
			os.Exit(1)
		}
		logger.Info("Connected to RabbitMQ")
		break
	}
	defer conn.Close()

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

	// Setup health check HTTP server
	healthChecker := health.NewHealthChecker(conn, logger)
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/healthz", healthChecker.Handler())
	healthMux.Handle("/metrics", promhttp.Handler()) // Prometheus metrics endpoint
	healthServer := &http.Server{
		Addr:    ":8080",
		Handler: healthMux,
	}

	go func() {
		logger.Info("Health check server starting", slog.String("addr", ":8080"))
		if err := healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Health check server failed", slog.Any("error", err))
		}
	}()

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start consuming in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- emailQueue.Consume(ctx, func(envelope queue.EmailEnvelope) error {
			start := time.Now()
			workerName := "email-worker"

			// Convert envelope to email message
			msg := sender.Message{
				To:       envelope.Destination,
				Subject:  envelope.Subject,
				TextBody: envelope.TextMessage,
				HTMLBody: envelope.HTMLMessage,
			}

			// Send email
			if err := emailSender.Send(ctx, msg); err != nil {
				duration := time.Since(start).Seconds()
				appmetrics.RecordJobProcessed(workerName, "failed", duration)
				appmetrics.RecordJobFailed(workerName, "send_error")
				logger.Error("Failed to send email",
					slog.String("user_id", envelope.UserID),
					slog.String("destination", envelope.Destination),
					slog.String("subject", envelope.Subject),
					slog.Any("error", err),
				)
				return err
			}

			duration := time.Since(start).Seconds()
			appmetrics.RecordJobProcessed(workerName, "success", duration)
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
		// Shutdown health check server
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := healthServer.Shutdown(shutdownCtx); err != nil {
			logger.Warn("Health check server shutdown error", slog.Any("error", err))
		}
	case err := <-errChan:
		if err != nil {
			logger.Error("Email queue consumer error", slog.Any("error", err))
			os.Exit(1)
		}
	}

	logger.Info("Email worker stopped")
}
