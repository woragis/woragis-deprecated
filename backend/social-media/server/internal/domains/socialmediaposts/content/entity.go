package content

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ContentPostStatus represents the status of a content post in the repurposing workflow.
type ContentPostStatus string

const (
	ContentPostStatusPending   ContentPostStatus = "pending"
	ContentPostStatusInProgress ContentPostStatus = "in_progress"
	ContentPostStatusCompleted  ContentPostStatus = "completed"
	ContentPostStatusArchived   ContentPostStatus = "archived"
)

// ContentPostPriority represents the priority level for repurposing.
type ContentPostPriority string

const (
	ContentPostPriorityLow    ContentPostPriority = "low"
	ContentPostPriorityMedium ContentPostPriority = "medium"
	ContentPostPriorityHigh   ContentPostPriority = "high"
)

// ContentPost represents a post from the backend that can be repurposed to social media.
type ContentPost struct {
	ID          uuid.UUID         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	PostID      uuid.UUID          `gorm:"column:post_id;type:uuid;not null;uniqueIndex" json:"postId"` // Links to backend Post
	ContentType string             `gorm:"column:content_type;type:varchar(50)" json:"contentType,omitempty"`
	Project     *string             `gorm:"column:project;type:varchar(255)" json:"project,omitempty"`
	Priority    ContentPostPriority `gorm:"column:priority;type:varchar(20);not null;default:'medium';index" json:"priority"`
	Status      ContentPostStatus   `gorm:"column:status;type:varchar(20);not null;default:'pending';index" json:"status"`
	Metadata    datatypes.JSON     `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	CreatedAt   time.Time          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time          `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for ContentPost.
func (ContentPost) TableName() string {
	return "content_posts"
}

// ContentRepurposing tracks the relationship between a content post and social media posts.
type ContentRepurposing struct {
	ID              uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ContentPostID   uuid.UUID `gorm:"column:content_post_id;type:uuid;not null;index:idx_content_repurposing" json:"contentPostId"`
	SocialPostID    uuid.UUID `gorm:"column:social_post_id;type:uuid;not null;index:idx_content_repurposing" json:"socialPostId"`
	RepurposedAt     time.Time `gorm:"column:repurposed_at;not null" json:"repurposedAt"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"createdAt"`
}

// TableName specifies the table name for ContentRepurposing.
func (ContentRepurposing) TableName() string {
	return "content_repurposings"
}

// NewContentPost creates a new content post entity.
func NewContentPost(postID uuid.UUID, contentType string) (*ContentPost, error) {
	post := &ContentPost{
		ID:          uuid.New(),
		PostID:      postID,
		ContentType: strings.TrimSpace(contentType),
		Priority:    ContentPostPriorityMedium,
		Status:      ContentPostStatusPending,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	return post, post.Validate()
}

// Validate ensures content post invariants hold.
func (p *ContentPost) Validate() error {
	if p == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilContentPost)
	}

	if p.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyContentPostID)
	}

	if p.PostID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyPostID)
	}

	if !isValidContentPostStatus(p.Status) {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}

	if !isValidContentPostPriority(p.Priority) {
		return NewDomainError(ErrCodeInvalidPriority, ErrUnsupportedPriority)
	}

	return nil
}

// UpdatePriority updates the priority of the content post.
func (p *ContentPost) UpdatePriority(priority ContentPostPriority) error {
	if !isValidContentPostPriority(priority) {
		return NewDomainError(ErrCodeInvalidPriority, ErrUnsupportedPriority)
	}
	p.Priority = priority
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateStatus updates the status of the content post.
func (p *ContentPost) UpdateStatus(status ContentPostStatus) error {
	if !isValidContentPostStatus(status) {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}
	p.Status = status
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// SetProject sets the project name for the content post.
func (p *ContentPost) SetProject(project string) {
	if project != "" {
		p.Project = &project
		p.UpdatedAt = time.Now().UTC()
	}
}

// NewContentRepurposing creates a new repurposing relationship.
func NewContentRepurposing(contentPostID, socialPostID uuid.UUID) (*ContentRepurposing, error) {
	repurposing := &ContentRepurposing{
		ID:            uuid.New(),
		ContentPostID: contentPostID,
		SocialPostID:  socialPostID,
		RepurposedAt:  time.Now().UTC(),
		CreatedAt:     time.Now().UTC(),
	}

	return repurposing, repurposing.Validate()
}

// Validate ensures content repurposing invariants hold.
func (r *ContentRepurposing) Validate() error {
	if r == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilRepurposing)
	}

	if r.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyRepurposingID)
	}

	if r.ContentPostID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyContentPostID)
	}

	if r.SocialPostID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySocialPostID)
	}

	return nil
}

// Validation helpers

func isValidContentPostStatus(s ContentPostStatus) bool {
	switch s {
	case ContentPostStatusPending, ContentPostStatusInProgress,
		ContentPostStatusCompleted, ContentPostStatusArchived:
		return true
	}
	return false
}

func isValidContentPostPriority(p ContentPostPriority) bool {
	switch p {
	case ContentPostPriorityLow, ContentPostPriorityMedium, ContentPostPriorityHigh:
		return true
	}
	return false
}
