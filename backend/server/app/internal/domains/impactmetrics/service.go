package impactmetrics

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates impact metric workflows.
type Service interface {
	CreateImpactMetric(ctx context.Context, userID uuid.UUID, req CreateImpactMetricRequest) (*ImpactMetric, error)
	UpdateImpactMetric(ctx context.Context, userID, metricID uuid.UUID, req UpdateImpactMetricRequest) (*ImpactMetric, error)
	GetImpactMetric(ctx context.Context, metricID uuid.UUID, userID uuid.UUID) (*ImpactMetric, error)
	ListImpactMetrics(ctx context.Context, filters ListImpactMetricsFilters) ([]ImpactMetric, error)
	ListFeaturedImpactMetrics(ctx context.Context) ([]ImpactMetric, error)
	GetMetricsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]ImpactMetric, error)
	DeleteImpactMetric(ctx context.Context, userID, metricID uuid.UUID) error
	// Dashboard methods
	GetDashboardMetrics(ctx context.Context, userID uuid.UUID) (*DashboardMetrics, error)
	GetMetricsByType(ctx context.Context, userID uuid.UUID, metricType MetricType) ([]ImpactMetric, error)
	GetTotalValueByType(ctx context.Context, userID uuid.UUID, metricType MetricType) (float64, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

var _ Service = (*service)(nil)

// NewService constructs a Service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// Request payloads

type CreateImpactMetricRequest struct {
	Type        MetricType  `json:"type"`
	Value       float64     `json:"value"`
	Unit        MetricUnit  `json:"unit"`
	Description string      `json:"description,omitempty"`
	EntityType  *EntityType `json:"entityType,omitempty"`
	EntityID    *string     `json:"entityId,omitempty"`
	PeriodStart *time.Time  `json:"periodStart,omitempty"`
	PeriodEnd   *time.Time  `json:"periodEnd,omitempty"`
	Featured    bool        `json:"featured,omitempty"`
	DisplayOrder int        `json:"displayOrder,omitempty"`
}

type UpdateImpactMetricRequest struct {
	Type        *MetricType  `json:"type,omitempty"`
	Value       *float64     `json:"value,omitempty"`
	Unit        *MetricUnit  `json:"unit,omitempty"`
	Description *string      `json:"description,omitempty"`
	EntityType  *EntityType  `json:"entityType,omitempty"`
	EntityID    *string      `json:"entityId,omitempty"`
	PeriodStart *time.Time   `json:"periodStart,omitempty"`
	PeriodEnd   *time.Time   `json:"periodEnd,omitempty"`
	Featured    *bool        `json:"featured,omitempty"`
	DisplayOrder *int        `json:"displayOrder,omitempty"`
}

type ListImpactMetricsFilters struct {
	UserID      *uuid.UUID
	Type        *MetricType
	EntityType  *EntityType
	EntityID    *uuid.UUID
	Featured    *bool
	PeriodStart *time.Time
	PeriodEnd   *time.Time
	Limit       int
	Offset      int
	OrderBy     string
	Order       string
}

// Service methods

func (s *service) CreateImpactMetric(ctx context.Context, userID uuid.UUID, req CreateImpactMetricRequest) (*ImpactMetric, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	metric, err := NewImpactMetric(userID, req.Type, req.Value, req.Unit)
	if err != nil {
		return nil, err
	}

	if req.Description != "" {
		metric.Description = req.Description
	}
	if req.EntityType != nil && req.EntityID != nil {
		entityID, err := uuid.Parse(*req.EntityID)
		if err != nil {
			return nil, NewDomainError(ErrCodeInvalidPayload, "invalid entity id format")
		}
		if err := metric.SetEntityLink(*req.EntityType, entityID); err != nil {
			return nil, err
		}
	}
	if req.PeriodStart != nil || req.PeriodEnd != nil {
		metric.UpdateDetails(metric.Description, req.PeriodStart, req.PeriodEnd)
	}
	metric.Featured = req.Featured
	metric.DisplayOrder = req.DisplayOrder

	if err := metric.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateImpactMetric(ctx, metric); err != nil {
		return nil, err
	}

	return metric, nil
}

func (s *service) UpdateImpactMetric(ctx context.Context, userID, metricID uuid.UUID, req UpdateImpactMetricRequest) (*ImpactMetric, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if metricID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyMetricID)
	}

	// Get existing metric
	metric, err := s.repo.GetImpactMetric(ctx, metricID, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Type != nil {
		metric.Type = *req.Type
	}
	if req.Value != nil {
		if err := metric.SetValue(*req.Value); err != nil {
			return nil, err
		}
	}
	if req.Unit != nil {
		metric.Unit = *req.Unit
	}
	if req.Description != nil {
		metric.Description = *req.Description
	}
	if req.EntityType != nil && req.EntityID != nil {
		entityID, err := uuid.Parse(*req.EntityID)
		if err != nil {
			return nil, NewDomainError(ErrCodeInvalidPayload, "invalid entity id format")
		}
		if err := metric.SetEntityLink(*req.EntityType, entityID); err != nil {
			return nil, err
		}
	} else if req.EntityType == nil && req.EntityID == nil {
		// Clear entity link if both are nil (explicit clearing)
		metric.ClearEntityLink()
	}
	if req.PeriodStart != nil || req.PeriodEnd != nil {
		periodStart := req.PeriodStart
		periodEnd := req.PeriodEnd
		if periodStart == nil {
			periodStart = metric.PeriodStart
		}
		if periodEnd == nil {
			periodEnd = metric.PeriodEnd
		}
		metric.UpdateDetails(metric.Description, periodStart, periodEnd)
	}
	if req.Featured != nil {
		metric.SetFeatured(*req.Featured)
	}
	if req.DisplayOrder != nil {
		metric.SetDisplayOrder(*req.DisplayOrder)
	}

	if err := metric.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateImpactMetric(ctx, metric); err != nil {
		return nil, err
	}

	return metric, nil
}

func (s *service) GetImpactMetric(ctx context.Context, metricID uuid.UUID, userID uuid.UUID) (*ImpactMetric, error) {
	if metricID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyMetricID)
	}

	return s.repo.GetImpactMetric(ctx, metricID, userID)
}

func (s *service) ListImpactMetrics(ctx context.Context, filters ListImpactMetricsFilters) ([]ImpactMetric, error) {
	repoFilters := ImpactMetricFilters{
		UserID:      filters.UserID,
		Type:        filters.Type,
		EntityType:  filters.EntityType,
		EntityID:    filters.EntityID,
		Featured:    filters.Featured,
		PeriodStart: filters.PeriodStart,
		PeriodEnd:   filters.PeriodEnd,
		Limit:       filters.Limit,
		Offset:      filters.Offset,
		OrderBy:     filters.OrderBy,
		Order:       filters.Order,
	}

	return s.repo.ListImpactMetrics(ctx, repoFilters)
}

func (s *service) ListFeaturedImpactMetrics(ctx context.Context) ([]ImpactMetric, error) {
	return s.repo.ListFeaturedImpactMetrics(ctx)
}

func (s *service) GetMetricsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]ImpactMetric, error) {
	if entityID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, "entity id cannot be empty")
	}

	return s.repo.GetMetricsByEntity(ctx, entityType, entityID)
}

func (s *service) DeleteImpactMetric(ctx context.Context, userID, metricID uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if metricID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyMetricID)
	}

	return s.repo.DeleteImpactMetric(ctx, metricID, userID)
}

// Dashboard methods

func (s *service) GetDashboardMetrics(ctx context.Context, userID uuid.UUID) (*DashboardMetrics, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	return s.repo.GetDashboardMetrics(ctx, userID)
}

func (s *service) GetMetricsByType(ctx context.Context, userID uuid.UUID, metricType MetricType) ([]ImpactMetric, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	return s.repo.GetMetricsByType(ctx, userID, metricType)
}

func (s *service) GetTotalValueByType(ctx context.Context, userID uuid.UUID, metricType MetricType) (float64, error) {
	if userID == uuid.Nil {
		return 0, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	return s.repo.GetTotalValueByType(ctx, userID, metricType)
}

