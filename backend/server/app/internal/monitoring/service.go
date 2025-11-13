package monitoring

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	appconfig "github.com/woragis/backend/server/app/pkg/config"
)

// Tracker describes observability hooks used by other domains.
type Tracker interface {
	RecordUserRegistration(ctx context.Context, userID uuid.UUID)
}

// Service wires Prometheus metrics with optional persistence for Grafana dashboards.
type Service struct {
	metrics *metricsRegistry
	repo    Repository
	logger  *slog.Logger
	enabled bool
}

// NewService creates a monitoring service and auto-migrates the repository.
func NewService(cfg appconfig.MonitoringConfig, repo Repository, logger *slog.Logger) *Service {
	svc := &Service{
		metrics: newMetricsRegistry(cfg.MetricsNamespace),
		repo:    repo,
		logger:  logger,
		enabled: cfg.Enabled,
	}

	return svc
}

// RecordUserRegistration increments counters and persists an event when possible.
func (s *Service) RecordUserRegistration(ctx context.Context, userID uuid.UUID) {
	s.metrics.incrementRegistrations()

	if s.repo == nil {
		return
	}

	event := Event{
		Type:      "user_registration",
		Reference: userID.String(),
	}

	if err := s.repo.StoreEvent(ctx, event); err != nil && s.logger != nil {
		s.logger.Warn("monitoring: failed to persist registration event", slog.Any("error", err))
	}
}

// ListRecentEvents returns persisted events when repository is available.
func (s *Service) ListRecentEvents(ctx context.Context, limit int) ([]Event, error) {
	if s.repo == nil {
		return []Event{}, nil
	}

	return s.repo.ListEvents(ctx, limit)
}

// MetricsMiddleware exposes a Fiber middleware to capture HTTP metrics.
func (s *Service) MetricsMiddleware() fiber.Handler {
	return s.metrics.metricsMiddleware()
}

// MetricsHandler returns the HTTP handler for Prometheus scraping.
func (s *Service) MetricsHandler() http.Handler {
	return s.metrics.metricsHandler()
}
