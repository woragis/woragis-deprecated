package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	langchainservice "github.com/woragis/backend/server/app/internal/services/langchain"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	translationworker "github.com/woragis/backend/server/app/internal/workers/translations"
	appconfig "github.com/woragis/backend/server/app/pkg/config"
	applogger "github.com/woragis/backend/server/app/pkg/logger"
)

func main() {
	// Initialize logger
	logger := applogger.New(os.Getenv("ENV"))
	logger.Info("starting translation worker")

	// Load configuration
	redisCfg := appconfig.LoadRedisConfig()

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

	// Connect to Redis
	redisOpts, err := redis.ParseURL(redisCfg.URL)
	if err != nil {
		logger.Error("failed to parse redis url", slog.Any("error", err))
		os.Exit(1)
	}

	redisClient := redis.NewClient(redisOpts)
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		logger.Error("failed to connect to redis", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("connected to redis", slog.String("url", redisCfg.URL))

	// Initialize services
	translationRepo := translationsdomain.NewGormRepository(db)
	translationQueue := translationsdomain.NewRedisQueue(redisClient)
	aiClient := langchainservice.NewClient(logger)
	translationService := translationsdomain.NewService(translationRepo, translationQueue, aiClient, logger)

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

