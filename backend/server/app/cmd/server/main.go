package main

import (
	"context"
	"fmt"
	stdlog "log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	emailservice "github.com/woragis/backend/server/app/internal/services/email"
	appconfig "github.com/woragis/backend/server/app/pkg/config"
	applogger "github.com/woragis/backend/server/app/pkg/logger"
)

func main() {
	cfg, err := appconfig.Load()
	if err != nil {
		stdlog.Fatalf("config: %v", err)
	}

	slogLogger := applogger.New(cfg.Env)
	slogLogger.Info("starting woragis backend",
		slog.String("env", cfg.Env),
		slog.Int("port", cfg.Port),
	)

	db, err := connectDatabase(slogLogger)
	if err != nil {
		slogLogger.Error("database connection failed", slog.Any("error", err))
		os.Exit(1)
	}

	if err := migrate(db); err != nil {
		slogLogger.Error("database migration failed", slog.Any("error", err))
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	app.Use(recover.New())
	app.Use(fiberlogger.New())

	api := app.Group("/api")

	emailSender := emailservice.NewNoopSender(slogLogger)
	authRepo := authdomain.NewGormRepository(db)
	authService := authdomain.NewService(authRepo, emailSender, slogLogger)
	authHandler := authdomain.NewHandler(authService, slogLogger)
	authdomain.SetupRoutes(api, authHandler)

	go func() {
		addr := fmt.Sprintf(":%d", cfg.Port)
		slogLogger.Info("http server listening", slog.String("addr", addr))
		if err := app.Listen(addr); err != nil {
			slogLogger.Error("fiber shutdown", slog.Any("error", err))
		}
	}()

	waitForShutdown(slogLogger, app)
}

func connectDatabase(log *slog.Logger) (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	gormLogger := gormlogger.New(
		stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
		gormlogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormlogger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	config := &gorm.Config{
		Logger: gormLogger,
	}

	const maxAttempts = 5
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err := openDatabase(dsn, config)
		if err == nil {
			if err := configurePool(dsn, db); err != nil {
				return nil, err
			}
			return db, nil
		}

		lastErr = err
		if log != nil {
			log.Warn("database connection attempt failed",
				slog.Int("attempt", attempt),
				slog.Int("max_attempts", maxAttempts),
				slog.Any("error", err),
			)
		}

		time.Sleep(time.Duration(attempt) * time.Second)
	}

	return nil, fmt.Errorf("connect database: %w", lastErr)
}

func openDatabase(dsn string, config *gorm.Config) (*gorm.DB, error) {
	switch {
	case strings.HasPrefix(dsn, "sqlite://"):
		sqliteDSN := strings.TrimPrefix(dsn, "sqlite://")
		return gorm.Open(sqlite.Open(sqliteDSN), config)
	default:
		return gorm.Open(postgres.Open(dsn), config)
	}
}

func configurePool(dsn string, db *gorm.DB) error {
	if strings.HasPrefix(dsn, "sqlite://") {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	sqlDB.SetConnMaxLifetime(60 * time.Minute)

	return nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&authdomain.User{},
	)
}

func waitForShutdown(log *slog.Logger, app *fiber.App) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Error("error during fiber shutdown", slog.Any("error", err))
	}
}
