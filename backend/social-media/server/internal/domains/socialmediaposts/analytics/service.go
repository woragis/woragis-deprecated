package analytics

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates analytics workflows.
type Service interface {
	RecordAnalytics(ctx context.Context, req RecordAnalyticsRequest) (*PostAnalytics, error)
	GetPostAnalytics(ctx context.Context, socialPostID uuid.UUID, startDate, endDate *time.Time) ([]PostAnalytics, error)
	GetAnalyticsSummary(ctx context.Context, filters AnalyticsFilters) (*AnalyticsSummary, error)
	GetTopPosts(ctx context.Context, limit int, metric string) ([]TopPost, error)
	CalculateEngagementRate(ctx context.Context, socialPostID uuid.UUID) (*float64, error)
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

type RecordAnalyticsRequest struct {
	SocialPostID uuid.UUID `json:"socialPostId"`
	MetricDate   time.Time `json:"metricDate"`
	Likes         *int64    `json:"likes,omitempty"`
	Comments      *int64    `json:"comments,omitempty"`
	Shares        *int64    `json:"shares,omitempty"`
	Views         *int64    `json:"views,omitempty"`
	Clicks        *int64    `json:"clicks,omitempty"`
	Reach         *int64    `json:"reach,omitempty"`
	Impressions   *int64    `json:"impressions,omitempty"`
}

// RecordAnalytics records analytics metrics for a post.
func (s *service) RecordAnalytics(ctx context.Context, req RecordAnalyticsRequest) (*PostAnalytics, error) {
	// Check if analytics already exists for this post and date
	existing, _ := s.repo.GetLatestAnalyticsForPost(ctx, req.SocialPostID)
	if existing != nil && existing.MetricDate.Equal(req.MetricDate.Truncate(24 * time.Hour)) {
		// Update existing
		existing.UpdateMetrics(req.Likes, req.Comments, req.Shares, req.Views, req.Clicks, req.Reach, req.Impressions)
		if err := s.repo.UpdateAnalytics(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// Create new
	analytics, err := NewPostAnalytics(req.SocialPostID, req.MetricDate)
	if err != nil {
		return nil, err
	}

	analytics.UpdateMetrics(req.Likes, req.Comments, req.Shares, req.Views, req.Clicks, req.Reach, req.Impressions)

	if err := s.repo.CreateAnalytics(ctx, analytics); err != nil {
		return nil, err
	}

	return analytics, nil
}

// GetPostAnalytics retrieves analytics for a specific post.
func (s *service) GetPostAnalytics(ctx context.Context, socialPostID uuid.UUID, startDate, endDate *time.Time) ([]PostAnalytics, error) {
	return s.repo.GetAnalyticsForPost(ctx, socialPostID, startDate, endDate)
}

// GetAnalyticsSummary retrieves aggregated analytics.
func (s *service) GetAnalyticsSummary(ctx context.Context, filters AnalyticsFilters) (*AnalyticsSummary, error) {
	return s.repo.GetAnalyticsSummary(ctx, filters)
}

// GetTopPosts retrieves top performing posts by metric.
func (s *service) GetTopPosts(ctx context.Context, limit int, metric string) ([]TopPost, error) {
	if metric == "" {
		metric = "likes"
	}
	return s.repo.GetTopPosts(ctx, limit, metric)
}

// CalculateEngagementRate calculates engagement rate for a post.
func (s *service) CalculateEngagementRate(ctx context.Context, socialPostID uuid.UUID) (*float64, error) {
	analytics, err := s.repo.GetLatestAnalyticsForPost(ctx, socialPostID)
	if err != nil {
		return nil, err
	}

	return analytics.EngagementRate, nil
}
