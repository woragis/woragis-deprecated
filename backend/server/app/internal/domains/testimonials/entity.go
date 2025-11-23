package testimonials

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// TestimonialStatus represents the moderation status of a testimonial.
type TestimonialStatus string

const (
	TestimonialStatusPending  TestimonialStatus = "pending"
	TestimonialStatusApproved TestimonialStatus = "approved"
	TestimonialStatusRejected TestimonialStatus = "rejected"
	TestimonialStatusHidden   TestimonialStatus = "hidden"
)

// Testimonial represents a recommendation or testimonial from a colleague, client, or mentor.
type Testimonial struct {
	ID           uuid.UUID        `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID        `gorm:"column:user_id;type:uuid;index;not null" json:"userId"` // The user who received the testimonial
	AuthorName   string           `gorm:"column:author_name;size:120;not null" json:"authorName"`
	AuthorRole   string           `gorm:"column:author_role;size:120" json:"authorRole,omitempty"`
	AuthorCompany string          `gorm:"column:author_company;size:120" json:"authorCompany,omitempty"`
	AuthorPhoto  string           `gorm:"column:author_photo;size:512" json:"authorPhoto,omitempty"` // URL to photo
	Content      string           `gorm:"column:content;type:text;not null" json:"content"`
	Rating       *int             `gorm:"column:rating;check:rating >= 1 AND rating <= 5" json:"rating,omitempty"` // Optional 1-5 star rating
	LinkedInURL  string           `gorm:"column:linkedin_url;size:512" json:"linkedinUrl,omitempty"`
	Status       TestimonialStatus `gorm:"column:status;type:varchar(32);not null;default:'pending';index" json:"status"`
	DisplayOrder int              `gorm:"column:display_order;not null;default:0;index" json:"displayOrder"` // For ordering in carousel/list
	CreatedAt    time.Time        `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time        `gorm:"column:updated_at" json:"updatedAt"`
}

// NewTestimonial creates a new testimonial entity.
func NewTestimonial(userID uuid.UUID, authorName, content string) *Testimonial {
	return &Testimonial{
		ID:         uuid.New(),
		UserID:     userID,
		AuthorName: strings.TrimSpace(authorName),
		Content:    strings.TrimSpace(content),
		Status:     TestimonialStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// Validate validates the testimonial entity.
func (t *Testimonial) Validate() error {
	if t.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if strings.TrimSpace(t.AuthorName) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyAuthorName)
	}
	if strings.TrimSpace(t.Content) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyContent)
	}
	if t.Rating != nil && (*t.Rating < 1 || *t.Rating > 5) {
		return NewDomainError(ErrCodeInvalidPayload, ErrInvalidRating)
	}
	if t.Status != TestimonialStatusPending && t.Status != TestimonialStatusApproved && t.Status != TestimonialStatusRejected && t.Status != TestimonialStatusHidden {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}
	return nil
}

// IsApproved returns true if the testimonial is approved and visible.
func (t *Testimonial) IsApproved() bool {
	return t.Status == TestimonialStatusApproved
}

// IsVisible returns true if the testimonial should be visible to the public.
func (t *Testimonial) IsVisible() bool {
	return t.Status == TestimonialStatusApproved
}

