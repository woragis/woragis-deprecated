package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/woragis/backend/whatsapp-worker/internal/config"
	"github.com/woragis/backend/whatsapp-worker/internal/queue"
	"github.com/woragis/backend/whatsapp-worker/internal/notifier"
	"github.com/woragis/backend/whatsapp-worker/pkg/health"
	"github.com/woragis/backend/whatsapp-worker/pkg/logger"
	appmetrics "github.com/woragis/backend/whatsapp-worker/pkg/metrics"
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
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start consuming in goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- whatsappQueue.Consume(ctx, func(envelope queue.WhatsAppEnvelope) error {
			start := time.Now()
			workerName := "whatsapp-worker"

			// Validate envelope
			if err := queue.ValidateWhatsAppEnvelope(envelope); err != nil {
				duration := time.Since(start).Seconds()
				appmetrics.RecordJobProcessed(workerName, "failed", duration)
				appmetrics.RecordJobFailed(workerName, "validation_error")
				logger.Error("Invalid WhatsApp message",
					slog.String("user_id", envelope.UserID),
					slog.Any("error", err),
				)
				return fmt.Errorf("validation failed: %w", err)
			}

			// Send WhatsApp message
			if err := whatsappNotifier.Send(ctx, envelope.Destination, envelope.TextMessage); err != nil {
				duration := time.Since(start).Seconds()
				appmetrics.RecordJobProcessed(workerName, "failed", duration)
				appmetrics.RecordJobFailed(workerName, "send_error")
				logger.Error("Failed to send WhatsApp message",
					slog.String("user_id", envelope.UserID),
					slog.String("destination", envelope.Destination),
					slog.Any("error", err),
				)
				return err
			}

			duration := time.Since(start).Seconds()
			appmetrics.RecordJobProcessed(workerName, "success", duration)
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
		// Shutdown health check server
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := healthServer.Shutdown(shutdownCtx); err != nil {
			logger.Warn("Health check server shutdown error", slog.Any("error", err))
		}
	case err := <-errChan:
		if err != nil {
			logger.Error("WhatsApp queue consumer error", slog.Any("error", err))
			os.Exit(1)
		}
	}

	logger.Info("WhatsApp worker stopped")
}
