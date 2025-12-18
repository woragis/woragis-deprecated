package scheduling

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ScheduledPostStatus represents the status of a scheduled post.
type ScheduledPostStatus string

const (
	ScheduledPostStatusPending   ScheduledPostStatus = "pending"
	ScheduledPostStatusScheduled ScheduledPostStatus = "scheduled"
	ScheduledPostStatusPosted    ScheduledPostStatus = "posted"
	ScheduledPostStatusCancelled ScheduledPostStatus = "cancelled"
)

// ScheduledPost represents a scheduled social media post.
type ScheduledPost struct {
	ID             uuid.UUID          `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	SocialPostID   uuid.UUID          `gorm:"column:social_post_id;type:uuid;not null;index" json:"socialPostId"`
	PlatformID     *uuid.UUID         `gorm:"column:platform_id;type:uuid;index" json:"platformId,omitempty"` // Optional link to platform config
	ScheduledDate  time.Time          `gorm:"column:scheduled_date;type:date;not null;index" json:"scheduledDate"`
	ScheduledTime  time.Time          `gorm:"column:scheduled_time;type:time;not null" json:"scheduledTime"`
	ScheduledAt    time.Time          `gorm:"column:scheduled_at;type:timestamp;not null;index" json:"scheduledAt"` // Combined date + time
	Status         ScheduledPostStatus `gorm:"column:status;type:varchar(20);not null;default:'pending';index" json:"status"`
	Metadata       datatypes.JSON      `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	CreatedAt      time.Time           `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time           `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for ScheduledPost.
func (ScheduledPost) TableName() string {
	return "scheduled_posts"
}

// NewScheduledPost creates a new scheduled post entity.
func NewScheduledPost(socialPostID uuid.UUID, scheduledAt time.Time) (*ScheduledPost, error) {
	scheduled := &ScheduledPost{
		ID:            uuid.New(),
		SocialPostID:  socialPostID,
		ScheduledDate: scheduledAt.Truncate(24 * time.Hour),
		ScheduledTime: time.Date(0, 1, 1, scheduledAt.Hour(), scheduledAt.Minute(), scheduledAt.Second(), 0, time.UTC),
		ScheduledAt:   scheduledAt.UTC(),
		Status:        ScheduledPostStatusPending,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}

	return scheduled, scheduled.Validate()
}

// Validate ensures scheduled post invariants hold.
func (s *ScheduledPost) Validate() error {
	if s == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilScheduledPost)
	}

	if s.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyScheduledPostID)
	}

	if s.SocialPostID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySocialPostID)
	}

	if !isValidScheduledPostStatus(s.Status) {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}

	return nil
}

// UpdateStatus updates the status of the scheduled post.
func (s *ScheduledPost) UpdateStatus(status ScheduledPostStatus) error {
	if !isValidScheduledPostStatus(status) {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}
	s.Status = status
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// SetPlatformID sets the platform configuration ID.
func (s *ScheduledPost) SetPlatformID(platformID uuid.UUID) {
	s.PlatformID = &platformID
	s.UpdatedAt = time.Now().UTC()
}

// MarkAsPosted marks the scheduled post as posted.
func (s *ScheduledPost) MarkAsPosted() error {
	if s.Status != ScheduledPostStatusScheduled && s.Status != ScheduledPostStatusPending {
		return NewDomainError(ErrCodeInvalidStatusTransition, "can only mark scheduled or pending posts as posted")
	}
	return s.UpdateStatus(ScheduledPostStatusPosted)
}

// Cancel cancels the scheduled post.
func (s *ScheduledPost) Cancel() error {
	if s.Status == ScheduledPostStatusPosted {
		return NewDomainError(ErrCodeInvalidStatusTransition, "cannot cancel already posted posts")
	}
	return s.UpdateStatus(ScheduledPostStatusCancelled)
}

// Validation helpers

func isValidScheduledPostStatus(s ScheduledPostStatus) bool {
	switch s {
	case ScheduledPostStatusPending, ScheduledPostStatusScheduled,
		ScheduledPostStatusPosted, ScheduledPostStatusCancelled:
		return true
	}
	return false
}
