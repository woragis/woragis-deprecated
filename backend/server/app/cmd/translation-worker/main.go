package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	langchainservice "github.com/woragis/backend/server/app/internal/services/langchain"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	translationworker "github.com/woragis/backend/server/app/internal/workers/translations"
	appconfig "github.com/woragis/backend/server/app/pkg/config"
	applogger "github.com/woragis/backend/server/app/pkg/logger"
	"github.com/woragis/backend/server/app/pkg/rabbitmq"
)

func main() {
	// Initialize logger
	logger := applogger.New(os.Getenv("ENV"))
	logger.Info("starting translation worker")

	// Load configuration
	rabbitmqCfg := appconfig.LoadRabbitMQConfig()

	// Connect to database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Error("DATABASE_URL not set")
		os.Exit(1)
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	// Connect to RabbitMQ
	rabbitmqConn, err := rabbitmq.NewConnection(rabbitmqCfg.URL, logger)
	if err != nil {
		logger.Error("failed to connect to rabbitmq", slog.Any("error", err))
		os.Exit(1)
	}
	defer rabbitmqConn.Close()
	logger.Info("connected to rabbitmq", slog.String("url", rabbitmqCfg.URL))

	// Initialize services
	translationRepo := translationsdomain.NewGormRepository(db)
	translationQueue, err := translationsdomain.NewRabbitMQQueue(rabbitmqConn)
	if err != nil {
		logger.Error("failed to create translation queue", slog.Any("error", err))
		os.Exit(1)
	}
	aiClient := langchainservice.NewClient(logger)
	translationService := translationsdomain.NewService(translationRepo, translationQueue, aiClient, db, logger)

	// Create worker
	worker := translationworker.NewWorker(translationQueue, translationService, logger)

	// Start worker in background
	workerCtx, workerCancel := context.WithCancel(context.Background())
	go worker.Start(workerCtx)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("shutting down translation worker")

	// Stop worker
	worker.Stop()
	workerCancel()

	// Give worker time to finish current job
	time.Sleep(2 * time.Second)

	logger.Info("translation worker stopped")
}

