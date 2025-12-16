package health

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	// StatusHealthy indicates all checks passed
	StatusHealthy = "healthy"
	// StatusDegraded indicates some non-critical checks failed
	StatusDegraded = "degraded"
	// StatusUnhealthy indicates critical checks failed
	StatusUnhealthy = "unhealthy"
)

// CheckResult represents the result of a health check
type CheckResult struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // "ok" or "error"
	Message string `json:"message,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string        `json:"status"` // "healthy", "degraded", or "unhealthy"
	Checks []CheckResult `json:"checks"`
}

// RabbitMQChecker is an interface for checking RabbitMQ connection
type RabbitMQChecker interface {
	IsConnected() bool
}

// HealthChecker manages health checks for dependencies
type HealthChecker struct {
	db            *gorm.DB
	redisClient   *redis.Client
	rabbitmqCheck RabbitMQChecker
	logger        *slog.Logger
	mu            sync.RWMutex
	cache         *HealthResponse
	lastCheck     time.Time
	cacheTTL      time.Duration
}

// NewHealthChecker creates a new health checker instance
func NewHealthChecker(db *gorm.DB, redisClient *redis.Client, logger *slog.Logger) *HealthChecker {
	return &HealthChecker{
		db:          db,
		redisClient: redisClient,
		logger:      logger,
		cacheTTL:    5 * time.Second, // Cache results for 5 seconds
	}
}

// SetRabbitMQChecker sets the RabbitMQ checker (optional)
func (h *HealthChecker) SetRabbitMQChecker(checker RabbitMQChecker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.rabbitmqCheck = checker
	// Invalidate cache when RabbitMQ checker is set
	h.cache = nil
}

// Check performs all health checks
func (h *HealthChecker) Check(ctx context.Context) HealthResponse {
	h.mu.RLock()
	// Return cached result if still valid
	if h.cache != nil && time.Since(h.lastCheck) < h.cacheTTL {
		cached := *h.cache
		h.mu.RUnlock()
		return cached
	}
	h.mu.RUnlock()

	// Perform checks
	checks := []CheckResult{
		h.checkDatabase(ctx),
		h.checkRedis(ctx),
	}

	// Add RabbitMQ check if checker is available
	if h.rabbitmqCheck != nil {
		checks = append(checks, h.checkRabbitMQ(ctx))
	}

	// Determine overall status
	status := StatusHealthy
	hasErrors := false
	hasWarnings := false

	for _, check := range checks {
		if check.Status == "error" {
			hasErrors = true
		}
	}

	if hasErrors {
		status = StatusUnhealthy
	} else if hasWarnings {
		status = StatusDegraded
	}

	response := HealthResponse{
		Status: status,
		Checks: checks,
	}

	// Cache the result
	h.mu.Lock()
	h.cache = &response
	h.lastCheck = time.Now()
	h.mu.Unlock()

	return response
}

// checkDatabase checks database connectivity
func (h *HealthChecker) checkDatabase(ctx context.Context) CheckResult {
	if h.db == nil {
		return CheckResult{
			Name:   "database",
			Status: "error",
			Message: "database not configured",
		}
	}

	sqlDB, err := h.db.DB()
	if err != nil {
		return CheckResult{
			Name:   "database",
			Status: "error",
			Message: err.Error(),
		}
	}

	// Use a timeout context for the ping
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		return CheckResult{
			Name:   "database",
			Status: "error",
			Message: err.Error(),
		}
	}

	return CheckResult{
		Name:   "database",
		Status: "ok",
	}
}

// checkRedis checks Redis connectivity
func (h *HealthChecker) checkRedis(ctx context.Context) CheckResult {
	if h.redisClient == nil {
		return CheckResult{
			Name:   "redis",
			Status: "error",
			Message: "redis not configured",
		}
	}

	// Use a timeout context for the ping
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := h.redisClient.Ping(pingCtx).Err(); err != nil {
		return CheckResult{
			Name:   "redis",
			Status: "error",
			Message: err.Error(),
		}
	}

	return CheckResult{
		Name:   "redis",
		Status: "ok",
	}
}

// checkRabbitMQ checks RabbitMQ connectivity
func (h *HealthChecker) checkRabbitMQ(ctx context.Context) CheckResult {
	if h.rabbitmqCheck == nil {
		return CheckResult{
			Name:   "rabbitmq",
			Status: "ok", // Not configured, but not an error
			Message: "not configured",
		}
	}

	// Use a timeout context
	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Check if connected (non-blocking)
	connected := h.rabbitmqCheck.IsConnected()
	if !connected {
		return CheckResult{
			Name:   "rabbitmq",
			Status: "error",
			Message: "not connected",
		}
	}

	// If we have a context that's done, return error
	select {
	case <-checkCtx.Done():
		return CheckResult{
			Name:   "rabbitmq",
			Status: "error",
			Message: "check timeout",
		}
	default:
	}

	return CheckResult{
		Name:   "rabbitmq",
		Status: "ok",
	}
}

// Handler returns a Fiber handler for the health check endpoint
func (h *HealthChecker) Handler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		result := h.Check(ctx)

		// Determine HTTP status code
		statusCode := fiber.StatusOK
		if result.Status == StatusUnhealthy {
			statusCode = fiber.StatusServiceUnavailable
		} else if result.Status == StatusDegraded {
			statusCode = fiber.StatusOK // Still 200, but status indicates degradation
		}

		return c.Status(statusCode).JSON(result)
	}
}
