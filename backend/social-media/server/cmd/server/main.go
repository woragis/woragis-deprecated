package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"woragis-social-media-service/internal/config"
	"woragis-social-media-service/internal/database"
	"woragis-social-media-service/pkg/health"
	// authPkg "woragis-social-media-service/pkg/auth" // Not used
	applogger "woragis-social-media-service/pkg/logger"
	appmetrics "woragis-social-media-service/pkg/metrics"
	apptracing "woragis-social-media-service/pkg/tracing"
	appsecurity "woragis-social-media-service/pkg/security"
	apptimeout "woragis-social-media-service/pkg/timeout"
	
	socialmediadomain "woragis-social-media-service/internal/domains"
)

func main() {
	// Load configuration first to get environment
	cfg := config.Load()
	env := cfg.Env
	if env == "" {
		env = os.Getenv("ENV")
		if env == "" {
			env = "development"
		}
	}

	// Setup structured logger with trace ID support
	slogLogger := applogger.New(env)

	// Initialize OpenTelemetry tracing
	tracingShutdown, err := apptracing.Init(apptracing.Config{
		ServiceName:    cfg.AppName,
		ServiceVersion: "1.0.0", // TODO: Get from build info
		Environment:    env,
		JaegerEndpoint: os.Getenv("JAEGER_ENDPOINT"), // Defaults to http://jaeger:4318
	})
	if err != nil {
		slogLogger.Warn("failed to initialize tracing", "error", err)
	} else {
		slogLogger.Info("tracing initialized", "service", cfg.AppName)
		defer func() {
			if tracingShutdown != nil {
				tracingShutdown()
			}
		}()
	}

	// Load database and Redis configs
	dbCfg := config.LoadDatabaseConfig()
	redisCfg := config.LoadRedisConfig()
	
	// Initialize database manager
	dbManager, err := database.NewFromConfig(dbCfg, redisCfg)
	if err != nil {
		slogLogger.Error("failed to initialize database manager", "error", err)
		os.Exit(1)
	}
	defer dbManager.Close()
	
	// Perform initial health check
	if err := dbManager.HealthCheck(); err != nil {
		slogLogger.Warn("Database health check failed", "error", err)
	} else {
		slogLogger.Info("All database connections are healthy")
	}

	// Run migrations
	if err := socialmediadomain.MigrateSocialMediaTables(dbManager.GetPostgres()); err != nil {
		slogLogger.Error("failed to run social media migrations", "error", err)
		os.Exit(1)
	}

	// Create Fiber app
	app := config.CreateFiberApp(cfg)

	// Recovery middleware (early in chain)
	app.Use(recover.New())

	// Security headers middleware (must be early, before other middlewares)
	app.Use(appsecurity.SecurityHeadersMiddleware())

	// Request timeout middleware (30 seconds default)
	app.Use(apptimeout.Middleware(apptimeout.DefaultConfig()))

	// Add OpenTelemetry tracing middleware (must be first to extract trace context)
	app.Use(apptracing.Middleware(cfg.AppName))
	// Add request ID middleware for distributed tracing (works with tracing, preserves trace_id)
	app.Use(applogger.RequestIDMiddleware(slogLogger))
	// Add structured request logging middleware
	app.Use(applogger.RequestLoggerMiddleware(slogLogger))
	// Add Prometheus metrics middleware
	app.Use(appmetrics.Middleware())

	// Request size limit (10MB)
	app.Use(appsecurity.RequestSizeLimitMiddleware(10 * 1024 * 1024))

	// Input sanitization
	app.Use(appsecurity.InputSanitizationMiddleware())

	// CSRF protection (for state-changing requests)
	csrfCfg := appsecurity.DefaultCSRFConfig(dbManager.GetRedis())
	app.Use(appsecurity.CSRFMiddleware(csrfCfg))

	// Rate limiting (100 requests per minute per IP/user)
	app.Use(appsecurity.RateLimitMiddleware(100, time.Minute))

	// CORS middleware (if enabled)
	corsCfg := config.LoadCORSConfig()
	if corsCfg.Enabled {
		config.SetupCORS(app, corsCfg)
	}

	// Initialize health checker
	healthChecker := health.NewHealthChecker(dbManager.GetPostgres(), dbManager.GetRedis(), slogLogger)

	// Health check endpoints (before API routes, no auth required)
	app.Get("/healthz", healthChecker.Handler())           // Combined health check
	app.Get("/healthz/live", healthChecker.LivenessHandler())   // Liveness probe (Kubernetes)
	app.Get("/healthz/ready", healthChecker.ReadinessHandler()) // Readiness probe (Kubernetes)

	// Prometheus metrics endpoint (before API routes, no auth required)
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// API routes group
	api := app.Group("/api/v1")

	// Load auth service URL for JWT validation
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://auth-service:3000"
	}

	// Setup social media domain routes
	socialmediadomain.SetupRoutes(api, dbManager.GetPostgres(), authServiceURL, slogLogger)

	// Setup graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%s", cfg.Port)
		slogLogger.Info("starting social media service", "addr", addr, "env", env)
		if err := app.Listen(addr); err != nil {
			slogLogger.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	slogLogger.Info("shutting down social media service gracefully")

	// Give ongoing requests time to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		slogLogger.Error("error during shutdown", "error", err)
	}

	slogLogger.Info("social media service stopped")
}

