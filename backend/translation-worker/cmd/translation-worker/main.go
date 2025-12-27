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

	"github.com/google/uuid"
	"github.com/woragis/backend/translation-worker/internal/config"
	"github.com/woragis/backend/translation-worker/internal/database"
	"github.com/woragis/backend/translation-worker/internal/queue"
	"github.com/woragis/backend/translation-worker/internal/translator"
	"github.com/woragis/backend/translation-worker/pkg/health"
	"github.com/woragis/backend/translation-worker/pkg/logger"
	appmetrics "github.com/woragis/backend/translation-worker/pkg/metrics"
)

func main() {
	// Setup structured logger
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	logger := logger.New(env)

	logger.Info("Starting translation worker", "env", env)

	// Load configuration
	rabbitmqCfg := config.LoadRabbitMQConfig()
	dbCfg := config.LoadDatabaseConfig()
	translationCfg := config.LoadTranslationConfig()
	workerCfg := config.LoadWorkerConfig()

	// Initialize database repository
	dbRepo, err := database.NewRepository(dbCfg.URL, logger)
	if err != nil {
		logger.Error("Failed to create database repository", slog.Any("error", err))
		os.Exit(1)
	}

	// Initialize translator
	translatorClient, err := translator.NewTranslator(translationCfg, logger)
	if err != nil {
		logger.Error("Failed to create translator", slog.Any("error", err))
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
	translationQueue, err := queue.NewQueue(
		conn,
		workerCfg.QueueName,
		workerCfg.Exchange,
		workerCfg.RoutingKey,
		workerCfg.PrefetchCount,
		logger,
	)
	if err != nil {
		logger.Error("Failed to create translation queue", slog.Any("error", err))
		os.Exit(1)
	}

	// Setup health check HTTP server
	healthChecker := health.NewHealthChecker(conn, logger)
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/healthz", healthChecker.Handler())
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
		errChan <- translationQueue.Consume(ctx, func(job queue.TranslationJob) error {
			start := time.Now()
			workerName := "translation-worker"

			// Validate job
			if err := queue.ValidateTranslationJob(job); err != nil {
				logger.Error("Invalid translation job", slog.Any("error", err), slog.String("job_id", job.ID))
				appmetrics.RecordJobProcessed(workerName, "failed", time.Since(start).Seconds())
				appmetrics.RecordJobFailed(workerName, "validation_error")
				return fmt.Errorf("validation failed: %w", err)
			}

			err := processTranslationJob(ctx, job, dbRepo, translatorClient, logger)
			duration := time.Since(start).Seconds()

			if err != nil {
				appmetrics.RecordJobProcessed(workerName, "failed", duration)
				appmetrics.RecordJobFailed(workerName, "processing_error")
			} else {
				appmetrics.RecordJobProcessed(workerName, "success", duration)
			}

			return err
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
			logger.Error("Translation queue consumer error", slog.Any("error", err))
			os.Exit(1)
		}
	}

	logger.Info("Translation worker stopped")
}

// processTranslationJob processes a single translation job.
func processTranslationJob(
	ctx context.Context,
	job queue.TranslationJob,
	dbRepo database.Repository,
	translatorClient translator.Translator,
	logger *slog.Logger,
) error {
	logger.Info("Processing translation job",
		slog.String("job_id", job.ID),
		slog.String("entity_type", job.EntityType),
		slog.String("entity_id", job.EntityID),
		slog.String("language", job.Language),
	)

	// Parse entity ID
	entityID, err := uuid.Parse(job.EntityID)
	if err != nil {
		return fmt.Errorf("invalid entity ID: %w", err)
	}

	// Get or create translation record
	translation, err := dbRepo.GetTranslationByEntity(ctx, job.EntityType, entityID, job.Language)
	if err != nil {
		return fmt.Errorf("failed to get translation: %w", err)
	}

	// Create new translation if it doesn't exist
	if translation == nil {
		translation = &database.Translation{
			ID:         uuid.New(),
			EntityType: job.EntityType,
			EntityID:   entityID,
			Language:   job.Language,
			Status:     "processing",
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  time.Now().UTC(),
		}
		// Initialize empty fields
		emptyFields := make(map[string]string)
		if err := translation.SetFields(emptyFields); err != nil {
			return fmt.Errorf("failed to initialize fields: %w", err)
		}
		if err := dbRepo.CreateTranslation(ctx, translation); err != nil {
			return fmt.Errorf("failed to create translation: %w", err)
		}
	} else {
		// Update status to processing
		translation.Status = "processing"
		translation.ErrorMessage = ""
		if err := dbRepo.UpdateTranslation(ctx, translation); err != nil {
			return fmt.Errorf("failed to update translation status: %w", err)
		}
	}

	// Fetch source text if not provided
	sourceText := job.SourceText
	if len(sourceText) == 0 {
		sourceText, err = dbRepo.FetchSourceTextFromEntity(ctx, job.EntityType, entityID, job.Fields)
		if err != nil {
			translation.Status = "failed"
			translation.ErrorMessage = fmt.Sprintf("Failed to fetch source text: %v", err)
			dbRepo.UpdateTranslation(ctx, translation)
			return fmt.Errorf("failed to fetch source text: %w", err)
		}
	}

	// Translate each field
	translatedFields := make(map[string]string)
	for _, field := range job.Fields {
		text, ok := sourceText[field]
		if !ok || text == "" {
			logger.Warn("Skipping field with empty source text",
				slog.String("field", field),
				slog.String("entity_id", job.EntityID),
			)
			continue
		}

		translated, err := translatorClient.Translate(ctx, text, job.Language)
		if err != nil {
			translation.Status = "failed"
			translation.ErrorMessage = fmt.Sprintf("Failed to translate field %s: %v", field, err)
			dbRepo.UpdateTranslation(ctx, translation)
			return fmt.Errorf("failed to translate field %s: %w", field, err)
		}

		translatedFields[field] = translated
	}

	// Update translation with results
	if err := translation.SetFields(translatedFields); err != nil {
		return fmt.Errorf("failed to set translated fields: %w", err)
	}
	translation.Status = "completed"
	translation.ErrorMessage = ""

	if err := dbRepo.UpdateTranslation(ctx, translation); err != nil {
		return fmt.Errorf("failed to update translation: %w", err)
	}

	logger.Info("Translation job completed successfully",
		slog.String("job_id", job.ID),
		slog.String("entity_type", job.EntityType),
		slog.String("entity_id", job.EntityID),
		slog.String("language", job.Language),
	)

	return nil
}
