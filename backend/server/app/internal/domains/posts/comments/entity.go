package comments

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// CommentStatus represents the status of a comment.
type CommentStatus string

const (
	CommentStatusPending  CommentStatus = "pending"  // Awaiting moderation
	CommentStatusApproved CommentStatus = "approved"
	CommentStatusRejected CommentStatus = "rejected"
	CommentStatusSpam     CommentStatus = "spam"
)

// Comment represents a comment on a blog post.
type Comment struct {
	ID        uuid.UUID    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	PostID    uuid.UUID    `gorm:"column:post_id;type:uuid;index;not null" json:"postId"`
	UserID    *uuid.UUID   `gorm:"column:user_id;type:uuid;index" json:"userId,omitempty"` // Optional for anonymous comments
	ParentID  *uuid.UUID   `gorm:"column:parent_id;type:uuid;index" json:"parentId,omitempty"` // For nested/reply comments
	AuthorName string      `gorm:"column:author_name;size:120" json:"authorName"` // Required for anonymous comments
	AuthorEmail string     `gorm:"column:author_email;size:255;index" json:"authorEmail,omitempty"`
	AuthorURL   string     `gorm:"column:author_url;size:512" json:"authorUrl,omitempty"`
	Content     string     `gorm:"column:content;type:text;not null" json:"content"`
	Status      CommentStatus `gorm:"column:status;type:varchar(32);not null;default:'pending';index" json:"status"`
	IPAddress   string     `gorm:"column:ip_address;size:45" json:"-"` // IPv6 compatible, not exposed in API
	UserAgent   string     `gorm:"column:user_agent;size:512" json:"-"` // Not exposed in API
	CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

// NewComment creates a new comment entity.
func NewComment(postID uuid.UUID, content, authorName string, userID *uuid.UUID) (*Comment, error) {
	comment := &Comment{
		ID:         uuid.New(),
		PostID:     postID,
		UserID:     userID,
		AuthorName: strings.TrimSpace(authorName),
		Content:    strings.TrimSpace(content),
		Status:     CommentStatusPending,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	return comment, comment.Validate()
}

// Validate ensures comment invariants hold.
func (c *Comment) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilComment)
	}

	if c.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCommentID)
	}

	if c.PostID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyPostID)
	}

	if c.AuthorName == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyAuthorName)
	}

	if c.Content == "" {
		return NewDomainError(ErrCodeInvalidContent, ErrEmptyCommentContent)
	}

	switch c.Status {
	case CommentStatusPending, CommentStatusApproved, CommentStatusRejected, CommentStatusSpam:
	default:
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedCommentStatus)
	}

	return nil
}

// Approve marks the comment as approved.
func (c *Comment) Approve() {
	c.Status = CommentStatusApproved
	c.UpdatedAt = time.Now().UTC()
}

// Reject marks the comment as rejected.
func (c *Comment) Reject() {
	c.Status = CommentStatusRejected
	c.UpdatedAt = time.Now().UTC()
}

// MarkAsSpam marks the comment as spam.
func (c *Comment) MarkAsSpam() {
	c.Status = CommentStatusSpam
	c.UpdatedAt = time.Now().UTC()
}

// UpdateContent updates the comment content.
func (c *Comment) UpdateContent(content string) error {
	if content == "" {
		return NewDomainError(ErrCodeInvalidContent, ErrEmptyCommentContent)
	}
	c.Content = strings.TrimSpace(content)
	c.UpdatedAt = time.Now().UTC()
	return nil
}

// SetReplyTo sets the parent comment ID for nested comments.
func (c *Comment) SetReplyTo(parentID uuid.UUID) {
	c.ParentID = &parentID
	c.UpdatedAt = time.Now().UTC()
}

