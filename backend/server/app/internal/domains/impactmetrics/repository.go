package impactmetrics

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for impact metrics.
type Repository interface {
	CreateImpactMetric(ctx context.Context, metric *ImpactMetric) error
	UpdateImpactMetric(ctx context.Context, metric *ImpactMetric) error
	GetImpactMetric(ctx context.Context, metricID uuid.UUID, userID uuid.UUID) (*ImpactMetric, error)
	ListImpactMetrics(ctx context.Context, filters ImpactMetricFilters) ([]ImpactMetric, error)
	ListFeaturedImpactMetrics(ctx context.Context) ([]ImpactMetric, error)
	GetMetricsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]ImpactMetric, error)
	DeleteImpactMetric(ctx context.Context, metricID uuid.UUID, userID uuid.UUID) error
	// Dashboard aggregation methods
	GetDashboardMetrics(ctx context.Context, userID uuid.UUID) (*DashboardMetrics, error)
	GetMetricsByType(ctx context.Context, userID uuid.UUID, metricType MetricType) ([]ImpactMetric, error)
	GetTotalValueByType(ctx context.Context, userID uuid.UUID, metricType MetricType) (float64, error)
}

// ImpactMetricFilters represents filtering options for listing metrics.
type ImpactMetricFilters struct {
	UserID     *uuid.UUID
	Type       *MetricType
	EntityType *EntityType
	EntityID   *uuid.UUID
	Featured   *bool
	PeriodStart *time.Time
	PeriodEnd   *time.Time
	Limit      int
	Offset     int
	OrderBy    string // "created_at", "value", "display_order"
	Order      string // "asc", "desc"
}

// DashboardMetrics represents aggregated dashboard data.
type DashboardMetrics struct {
	ProjectsDelivered    *MetricSummary `json:"projectsDelivered,omitempty"`
	UsersImpacted       *MetricSummary `json:"usersImpacted,omitempty"`
	PerformanceImprovement *MetricSummary `json:"performanceImprovement,omitempty"`
	CostSavings         *MetricSummary `json:"costSavings,omitempty"`
	TimeSaved           *MetricSummary `json:"timeSaved,omitempty"`
	TotalMetrics        int            `json:"totalMetrics"`
	LastUpdated         time.Time      `json:"lastUpdated"`
}

// MetricSummary represents aggregated data for a metric type.
type MetricSummary struct {
	Type        MetricType `json:"type"`
	TotalValue  float64    `json:"totalValue"`
	Unit        MetricUnit `json:"unit"`
	Count       int        `json:"count"`
	Average     float64    `json:"average"`
	Min         float64    `json:"min"`
	Max         float64    `json:"max"`
	LatestValue float64    `json:"latestValue"`
	LatestDate  *time.Time `json:"latestDate,omitempty"`
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateImpactMetric(ctx context.Context, metric *ImpactMetric) error {
	if metric == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilMetric)
	}

	if err := metric.Validate(); err != nil {
		return err
	}

	now := time.Now().UTC()
	metric.CreatedAt = now
	metric.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(metric).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return NewDomainError(ErrCodeConflict, ErrMetricAlreadyExists)
			}
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) UpdateImpactMetric(ctx context.Context, metric *ImpactMetric) error {
	if metric == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilMetric)
	}

	if metric.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyMetricID)
	}

	if err := metric.Validate(); err != nil {
		return err
	}

	metric.UpdatedAt = time.Now().UTC()

	result := r.db.WithContext(ctx).Model(&ImpactMetric{}).
		Where("id = ?", metric.ID).
		Updates(metric)

	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrMetricNotFound)
	}

	return nil
}

func (r *gormRepository) GetImpactMetric(ctx context.Context, metricID uuid.UUID, userID uuid.UUID) (*ImpactMetric, error) {
	if metricID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyMetricID)
	}

	var metric ImpactMetric
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", metricID, userID).
		First(&metric).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrMetricNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &metric, nil
}

func (r *gormRepository) ListImpactMetrics(ctx context.Context, filters ImpactMetricFilters) ([]ImpactMetric, error) {
	query := r.db.WithContext(ctx).Model(&ImpactMetric{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.EntityType != nil {
		query = query.Where("entity_type = ?", *filters.EntityType)
	}

	if filters.EntityID != nil {
		query = query.Where("entity_id = ?", *filters.EntityID)
	}

	if filters.Featured != nil {
		query = query.Where("featured = ?", *filters.Featured)
	}

	if filters.PeriodStart != nil {
		query = query.Where("period_start >= ?", *filters.PeriodStart)
	}

	if filters.PeriodEnd != nil {
		query = query.Where("period_end <= ?", *filters.PeriodEnd)
	}

	// Default ordering
	orderBy := filters.OrderBy
	if orderBy == "" {
		orderBy = "created_at"
	}
	order := filters.Order
	if order == "" {
		order = "desc"
	}
	query = query.Order(orderBy + " " + order)

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var metrics []ImpactMetric
	if err := query.Find(&metrics).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return metrics, nil
}

func (r *gormRepository) ListFeaturedImpactMetrics(ctx context.Context) ([]ImpactMetric, error) {
	var metrics []ImpactMetric
	err := r.db.WithContext(ctx).
		Where("featured = ?", true).
		Order("display_order ASC, created_at DESC").
		Find(&metrics).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return metrics, nil
}

func (r *gormRepository) GetMetricsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]ImpactMetric, error) {
	var metrics []ImpactMetric
	err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Find(&metrics).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return metrics, nil
}

func (r *gormRepository) DeleteImpactMetric(ctx context.Context, metricID uuid.UUID, userID uuid.UUID) error {
	if metricID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyMetricID)
	}

	// Verify ownership
	var metric ImpactMetric
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", metricID, userID).
		First(&metric).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewDomainError(ErrCodeNotFound, ErrMetricNotFound)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	if err := r.db.WithContext(ctx).Delete(&metric).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

// Dashboard aggregation methods

func (r *gormRepository) GetDashboardMetrics(ctx context.Context, userID uuid.UUID) (*DashboardMetrics, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	dashboard := &DashboardMetrics{
		LastUpdated: time.Now().UTC(),
	}

	// Get summary for each metric type
	projectsDelivered, err := r.getMetricSummary(ctx, userID, MetricTypeProjectsDelivered)
	if err == nil {
		dashboard.ProjectsDelivered = projectsDelivered
	}

	usersImpacted, err := r.getMetricSummary(ctx, userID, MetricTypeUsersImpacted)
	if err == nil {
		dashboard.UsersImpacted = usersImpacted
	}

	performanceImprovement, err := r.getMetricSummary(ctx, userID, MetricTypePerformanceImprovement)
	if err == nil {
		dashboard.PerformanceImprovement = performanceImprovement
	}

	costSavings, err := r.getMetricSummary(ctx, userID, MetricTypeCostSavings)
	if err == nil {
		dashboard.CostSavings = costSavings
	}

	timeSaved, err := r.getMetricSummary(ctx, userID, MetricTypeTimeSaved)
	if err == nil {
		dashboard.TimeSaved = timeSaved
	}

	// Get total count
	var totalCount int64
	r.db.WithContext(ctx).Model(&ImpactMetric{}).
		Where("user_id = ?", userID).
		Count(&totalCount)
	dashboard.TotalMetrics = int(totalCount)

	return dashboard, nil
}

func (r *gormRepository) getMetricSummary(ctx context.Context, userID uuid.UUID, metricType MetricType) (*MetricSummary, error) {
	var metrics []ImpactMetric
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, metricType).
		Order("created_at DESC").
		Find(&metrics).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	if len(metrics) == 0 {
		return nil, nil
	}

	summary := &MetricSummary{
		Type:   metricType,
		Count:  len(metrics),
		Unit:   metrics[0].Unit, // Assume all metrics of same type have same unit
		Min:    metrics[0].Value,
		Max:    metrics[0].Value,
		LatestValue: metrics[0].Value,
		LatestDate:  &metrics[0].CreatedAt,
	}

	var total float64
	for _, m := range metrics {
		total += m.Value
		if m.Value < summary.Min {
			summary.Min = m.Value
		}
		if m.Value > summary.Max {
			summary.Max = m.Value
		}
	}

	summary.TotalValue = total
	if len(metrics) > 0 {
		summary.Average = total / float64(len(metrics))
	}

	return summary, nil
}

func (r *gormRepository) GetMetricsByType(ctx context.Context, userID uuid.UUID, metricType MetricType) ([]ImpactMetric, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	var metrics []ImpactMetric
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, metricType).
		Order("created_at DESC").
		Find(&metrics).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return metrics, nil
}

func (r *gormRepository) GetTotalValueByType(ctx context.Context, userID uuid.UUID, metricType MetricType) (float64, error) {
	if userID == uuid.Nil {
		return 0, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	var total float64
	err := r.db.WithContext(ctx).Model(&ImpactMetric{}).
		Where("user_id = ? AND type = ?", userID, metricType).
		Select("COALESCE(SUM(value), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return total, nil
}

