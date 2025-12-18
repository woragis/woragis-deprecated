package analytics

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// PostAnalytics represents analytics metrics for a social media post at a specific point in time.
type PostAnalytics struct {
	ID             uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	SocialPostID   uuid.UUID  `gorm:"column:social_post_id;type:uuid;not null;index:idx_post_analytics" json:"socialPostId"`
	MetricDate     time.Time  `gorm:"column:metric_date;type:date;not null;index:idx_post_analytics" json:"metricDate"`
	Likes          *int64     `gorm:"column:likes;type:bigint" json:"likes,omitempty"`
	Comments       *int64     `gorm:"column:comments;type:bigint" json:"comments,omitempty"`
	Shares         *int64     `gorm:"column:shares;type:bigint" json:"shares,omitempty"`
	Views          *int64     `gorm:"column:views;type:bigint" json:"views,omitempty"`
	Clicks         *int64     `gorm:"column:clicks;type:bigint" json:"clicks,omitempty"`
	EngagementRate *float64   `gorm:"column:engagement_rate;type:decimal(10,4)" json:"engagementRate,omitempty"`
	Reach          *int64     `gorm:"column:reach;type:bigint" json:"reach,omitempty"`
	Impressions    *int64     `gorm:"column:impressions;type:bigint" json:"impressions,omitempty"`
	Metadata       datatypes.JSON `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for PostAnalytics.
func (PostAnalytics) TableName() string {
	return "post_analytics"
}

// NewPostAnalytics creates a new post analytics entry.
func NewPostAnalytics(socialPostID uuid.UUID, metricDate time.Time) (*PostAnalytics, error) {
	analytics := &PostAnalytics{
		ID:           uuid.New(),
		SocialPostID: socialPostID,
		MetricDate:   metricDate.Truncate(24 * time.Hour),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	return analytics, analytics.Validate()
}

// Validate ensures post analytics invariants hold.
func (a *PostAnalytics) Validate() error {
	if a == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilPostAnalytics)
	}

	if a.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyPostAnalyticsID)
	}

	if a.SocialPostID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySocialPostID)
	}

	return nil
}

// UpdateMetrics updates the analytics metrics.
func (a *PostAnalytics) UpdateMetrics(likes, comments, shares, views, clicks, reach, impressions *int64) {
	if likes != nil {
		a.Likes = likes
	}
	if comments != nil {
		a.Comments = comments
	}
	if shares != nil {
		a.Shares = shares
	}
	if views != nil {
		a.Views = views
	}
	if clicks != nil {
		a.Clicks = clicks
	}
	if reach != nil {
		a.Reach = reach
	}
	if impressions != nil {
		a.Impressions = impressions
	}
	a.UpdatedAt = time.Now().UTC()
	
	// Calculate engagement rate
	a.calculateEngagementRate()
}

// calculateEngagementRate calculates engagement rate based on available metrics.
func (a *PostAnalytics) calculateEngagementRate() {
	if a.Impressions == nil || *a.Impressions == 0 {
		return
	}

	var engagement int64
	if a.Likes != nil {
		engagement += *a.Likes
	}
	if a.Comments != nil {
		engagement += *a.Comments
	}
	if a.Shares != nil {
		engagement += *a.Shares
	}
	if a.Clicks != nil {
		engagement += *a.Clicks
	}

	if engagement > 0 {
		rate := float64(engagement) / float64(*a.Impressions) * 100.0
		a.EngagementRate = &rate
	}
}
