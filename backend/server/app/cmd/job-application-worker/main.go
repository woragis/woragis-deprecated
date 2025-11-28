package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	jobapplicationsdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications"
	jobwebsitesdomain "github.com/woragis/backend/server/app/internal/domains/jobwebsites"
	aiservice "github.com/woragis/backend/server/app/internal/services/ai"
	playwrightservice "github.com/woragis/backend/server/app/internal/services/playwright"
	langchainservice "github.com/woragis/backend/server/app/internal/services/langchain"
	jobapplicationworker "github.com/woragis/backend/server/app/internal/workers/jobapplications"
	appconfig "github.com/woragis/backend/server/app/pkg/config"
	applogger "github.com/woragis/backend/server/app/pkg/logger"
)

func main() {
	// Initialize logger
	logger := applogger.New(os.Getenv("APP_ENV"))
	logger.Info("starting job application worker")

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

	// Initialize repositories and services
	applicationRepo := jobapplicationsdomain.NewGormRepository(db)
	applicationQueue := jobapplicationsdomain.NewRedisQueue(redisClient)
	websiteRepo := jobwebsitesdomain.NewGormRepository(db)
	websiteService := jobwebsitesdomain.NewService(websiteRepo, logger)

	// Initialize AI service
	aiClient := langchainservice.NewClient(logger)
	coverLetterService := aiservice.NewCoverLetterService(aiClient, logger)

	// Initialize Playwright service
	playwrightOpts := playwrightservice.BrowserOptions{
		Headless:    getEnvBool("PLAYWRIGHT_HEADLESS", true),
		SlowMo:      getEnvInt("PLAYWRIGHT_SLOW_MO", 100),
		Timeout:     getEnvInt("PLAYWRIGHT_TIMEOUT", 30000),
		BrowserPath: os.Getenv("PLAYWRIGHT_BROWSER_PATH"),
	}
	browserManager := playwrightservice.NewBrowserManager(playwrightOpts, logger)
	scraper := playwrightservice.NewScraper(browserManager, logger)

	// Initialize orchestrator
	orchestrator := jobapplicationworker.NewOrchestrator(websiteService, logger)

	// Create worker
	worker := jobapplicationworker.NewWorker(
		applicationQueue,
		applicationRepo,
		websiteService,
		orchestrator,
		scraper,
		coverLetterService,
		db,
		logger,
	)

	// Start worker in background
	workerCtx, workerCancel := context.WithCancel(context.Background())
	go worker.Start(workerCtx)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("shutting down job application worker")

	// Stop worker
	worker.Stop()
	workerCancel()

	// Give worker time to finish current job
	time.Sleep(2 * time.Second)

	logger.Info("job application worker stopped")
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true" || value == "1"
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	var result int
	if _, err := fmt.Sscanf(value, "%d", &result); err != nil {
		return defaultValue
	}
	return result
}

