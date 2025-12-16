package health

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
)

// RabbitMQConnection is an interface for checking RabbitMQ connection status
type RabbitMQConnection interface {
	IsClosed() bool
}

// HealthChecker manages health checks for the WhatsApp worker
type HealthChecker struct {
	rabbitmqConn RabbitMQConnection
	logger       *slog.Logger
	mu           sync.RWMutex
	cache        *HealthResponse
	lastCheck    time.Time
	cacheTTL     time.Duration
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status string        `json:"status"`
	Checks []CheckResult `json:"checks"`
}

// CheckResult represents the result of a health check
type CheckResult struct {
	Name    string `json:"name"`
	Status  string `json:"status"` // "ok" or "error"
	Message string `json:"message,omitempty"`
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(rabbitmqConn RabbitMQConnection, logger *slog.Logger) *HealthChecker {
	return &HealthChecker{
		rabbitmqConn: rabbitmqConn,
		logger:       logger,
		cacheTTL:     5 * time.Second,
	}
}

// Check performs all health checks
func (h *HealthChecker) Check(ctx context.Context) HealthResponse {
	h.mu.RLock()
	if h.cache != nil && time.Since(h.lastCheck) < h.cacheTTL {
		cached := *h.cache
		h.mu.RUnlock()
		return cached
	}
	h.mu.RUnlock()

	checks := []CheckResult{
		h.checkRabbitMQ(ctx),
	}

	status := StatusHealthy
	for _, check := range checks {
		if check.Status == "error" {
			status = StatusUnhealthy
			break
		}
	}

	response := HealthResponse{
		Status: status,
		Checks: checks,
	}

	h.mu.Lock()
	h.cache = &response
	h.lastCheck = time.Now()
	h.mu.Unlock()

	return response
}

func (h *HealthChecker) checkRabbitMQ(ctx context.Context) CheckResult {
	if h.rabbitmqConn == nil {
		return CheckResult{
			Name:   "rabbitmq",
			Status: "error",
			Message: "not configured",
		}
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if h.rabbitmqConn.IsClosed() {
		return CheckResult{
			Name:   "rabbitmq",
			Status: "error",
			Message: "connection closed",
		}
	}

	select {
	case <-pingCtx.Done():
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

// Handler returns an HTTP handler for health checks
func (h *HealthChecker) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		result := h.Check(ctx)

		statusCode := http.StatusOK
		if result.Status == StatusUnhealthy {
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(result)
	}
}
