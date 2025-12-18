package analytics

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for post analytics.
type Repository interface {
	CreateAnalytics(ctx context.Context, analytics *PostAnalytics) error
	UpdateAnalytics(ctx context.Context, analytics *PostAnalytics) error
	GetAnalytics(ctx context.Context, analyticsID uuid.UUID) (*PostAnalytics, error)
	GetAnalyticsForPost(ctx context.Context, socialPostID uuid.UUID, startDate, endDate *time.Time) ([]PostAnalytics, error)
	GetLatestAnalyticsForPost(ctx context.Context, socialPostID uuid.UUID) (*PostAnalytics, error)
	GetAnalyticsSummary(ctx context.Context, filters AnalyticsFilters) (*AnalyticsSummary, error)
	GetPlatformAnalytics(ctx context.Context, platformID uuid.UUID, startDate, endDate time.Time) ([]PostAnalytics, error)
	GetTopPosts(ctx context.Context, limit int, metric string) ([]TopPost, error)
}

// AnalyticsFilters for querying analytics.
type AnalyticsFilters struct {
	SocialPostID *uuid.UUID
	StartDate    *time.Time
	EndDate      *time.Time
}

// AnalyticsSummary represents aggregated analytics data.
type AnalyticsSummary struct {
	TotalLikes          int64   `json:"totalLikes"`
	TotalComments       int64   `json:"totalComments"`
	TotalShares         int64   `json:"totalShares"`
	TotalViews          int64   `json:"totalViews"`
	TotalClicks         int64   `json:"totalClicks"`
	TotalReach          int64   `json:"totalReach"`
	TotalImpressions    int64   `json:"totalImpressions"`
	AverageEngagementRate float64 `json:"averageEngagementRate"`
	PostCount           int     `json:"postCount"`
}

// TopPost represents a top performing post.
type TopPost struct {
	SocialPostID uuid.UUID `json:"socialPostId"`
	MetricValue  int64     `json:"metricValue"`
	MetricName   string    `json:"metricName"`
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateAnalytics(ctx context.Context, analytics *PostAnalytics) error {
	if err := analytics.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(analytics).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrPostAnalyticsExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateAnalytics(ctx context.Context, analytics *PostAnalytics) error {
	if err := analytics.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(analytics).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetAnalytics(ctx context.Context, analyticsID uuid.UUID) (*PostAnalytics, error) {
	var analytics PostAnalytics
	err := r.db.WithContext(ctx).Where("id = ?", analyticsID).First(&analytics).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrPostAnalyticsNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &analytics, nil
}

func (r *gormRepository) GetAnalyticsForPost(ctx context.Context, socialPostID uuid.UUID, startDate, endDate *time.Time) ([]PostAnalytics, error) {
	var analytics []PostAnalytics
	query := r.db.WithContext(ctx).Where("social_post_id = ?", socialPostID)

	if startDate != nil {
		query = query.Where("metric_date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("metric_date <= ?", endDate)
	}

	if err := query.Order("metric_date desc").Find(&analytics).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return analytics, nil
}

func (r *gormRepository) GetLatestAnalyticsForPost(ctx context.Context, socialPostID uuid.UUID) (*PostAnalytics, error) {
	var analytics PostAnalytics
	err := r.db.WithContext(ctx).
		Where("social_post_id = ?", socialPostID).
		Order("metric_date desc, created_at desc").
		First(&analytics).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrPostAnalyticsNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &analytics, nil
}

func (r *gormRepository) GetAnalyticsSummary(ctx context.Context, filters AnalyticsFilters) (*AnalyticsSummary, error) {
	var summary AnalyticsSummary
	query := r.db.WithContext(ctx).Model(&PostAnalytics{})

	if filters.SocialPostID != nil {
		query = query.Where("social_post_id = ?", *filters.SocialPostID)
	}
	if filters.StartDate != nil {
		query = query.Where("metric_date >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		query = query.Where("metric_date <= ?", filters.EndDate)
	}

	// Get latest analytics for each post (to avoid double counting)
	subQuery := r.db.WithContext(ctx).
		Table("post_analytics").
		Select("social_post_id, MAX(metric_date) as max_date").
		Group("social_post_id")

	if filters.StartDate != nil {
		subQuery = subQuery.Where("metric_date >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		subQuery = subQuery.Where("metric_date <= ?", filters.EndDate)
	}

	err := query.
		Joins("INNER JOIN (?) as latest ON post_analytics.social_post_id = latest.social_post_id AND post_analytics.metric_date = latest.max_date", subQuery).
		Select(`
			COALESCE(SUM(likes), 0) as total_likes,
			COALESCE(SUM(comments), 0) as total_comments,
			COALESCE(SUM(shares), 0) as total_shares,
			COALESCE(SUM(views), 0) as total_views,
			COALESCE(SUM(clicks), 0) as total_clicks,
			COALESCE(SUM(reach), 0) as total_reach,
			COALESCE(SUM(impressions), 0) as total_impressions,
			COALESCE(AVG(engagement_rate), 0) as average_engagement_rate,
			COUNT(DISTINCT social_post_id) as post_count
		`).
		Scan(&summary).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &summary, nil
}

func (r *gormRepository) GetPlatformAnalytics(ctx context.Context, platformID uuid.UUID, startDate, endDate time.Time) ([]PostAnalytics, error) {
	// This would require joining with social_media_posts table
	// For now, return empty - would need to implement proper join
	var analytics []PostAnalytics
	return analytics, nil
}

func (r *gormRepository) GetTopPosts(ctx context.Context, limit int, metric string) ([]TopPost, error) {
	if limit <= 0 {
		limit = 10
	}

	var topPosts []TopPost
	query := r.db.WithContext(ctx).
		Table("post_analytics").
		Select("social_post_id, MAX(?) as metric_value", metric).
		Group("social_post_id").
		Order("metric_value desc").
		Limit(limit)

	err := query.Scan(&topPosts).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	// Set metric name
	for i := range topPosts {
		topPosts[i].MetricName = metric
	}

	return topPosts, nil
}

// isUniqueConstraintError checks if the error is a unique constraint violation.
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// Check for PostgreSQL unique constraint violation (code 23505)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}

	// Check for SQLite unique constraint violation
	if err.Error() == "UNIQUE constraint failed" {
		return true
	}

	return false
}
