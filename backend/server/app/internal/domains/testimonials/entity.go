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

// TestimonialType represents the type/category of a testimonial.
type TestimonialType string

const (
	TestimonialTypeGeneral        TestimonialType = "general"
	TestimonialTypeProjectSpecific TestimonialType = "project_specific"
	TestimonialTypeSkillSpecific   TestimonialType = "skill_specific"
)

// EntityType represents the type of entity being linked to a testimonial.
type EntityType string

const (
	EntityTypeProject EntityType = "project"
	EntityTypeSkill   EntityType = "skill"
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
	Context      string           `gorm:"column:context;type:text" json:"context,omitempty"` // When/where/why the testimonial was given
	VideoURL     string           `gorm:"column:video_url;size:512" json:"videoUrl,omitempty"` // Optional video testimonial URL
	Type         TestimonialType  `gorm:"column:type;type:varchar(32);not null;default:'general';index" json:"type"`
	Rating       *int             `gorm:"column:rating;check:rating >= 1 AND rating <= 5" json:"rating,omitempty"` // Optional 1-5 star rating
	LinkedInURL  string           `gorm:"column:linkedin_url;size:512" json:"linkedinUrl,omitempty"`
	Status       TestimonialStatus `gorm:"column:status;type:varchar(32);not null;default:'pending';index" json:"status"`
	DisplayOrder int              `gorm:"column:display_order;not null;default:0;index" json:"displayOrder"` // For ordering in carousel/list
	CreatedAt    time.Time        `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time        `gorm:"column:updated_at" json:"updatedAt"`
}

// TestimonialEntityLink represents the relationship between a testimonial and an entity (project or skill).
type TestimonialEntityLink struct {
	ID           uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	TestimonialID uuid.UUID `gorm:"column:testimonial_id;type:uuid;not null;index:idx_testimonial_entity_link" json:"testimonialId"`
	EntityType   EntityType `gorm:"column:entity_type;type:varchar(50);not null;index:idx_testimonial_entity_link" json:"entityType"`
	EntityID     uuid.UUID  `gorm:"column:entity_id;type:uuid;not null;index:idx_testimonial_entity_link" json:"entityId"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for TestimonialEntityLink.
func (TestimonialEntityLink) TableName() string {
	return "testimonial_entity_links"
}

// NewTestimonial creates a new testimonial entity.
func NewTestimonial(userID uuid.UUID, authorName, content string) *Testimonial {
	return &Testimonial{
		ID:         uuid.New(),
		UserID:     userID,
		AuthorName: strings.TrimSpace(authorName),
		Content:    strings.TrimSpace(content),
		Type:       TestimonialTypeGeneral,
		Status:     TestimonialStatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// SetContext updates the context field.
func (t *Testimonial) SetContext(context string) {
	t.Context = strings.TrimSpace(context)
	t.UpdatedAt = time.Now()
}

// SetVideoURL updates the video URL field.
func (t *Testimonial) SetVideoURL(videoURL string) {
	t.VideoURL = strings.TrimSpace(videoURL)
	t.UpdatedAt = time.Now()
}

// SetType updates the testimonial type.
func (t *Testimonial) SetType(testimonialType TestimonialType) error {
	if !isValidTestimonialType(testimonialType) {
		return NewDomainError(ErrCodeInvalidType, ErrUnsupportedTestimonialType)
	}
	t.Type = testimonialType
	t.UpdatedAt = time.Now()
	return nil
}

// NewTestimonialEntityLink creates a new link between a testimonial and an entity.
func NewTestimonialEntityLink(testimonialID uuid.UUID, entityType EntityType, entityID uuid.UUID) (*TestimonialEntityLink, error) {
	link := &TestimonialEntityLink{
		ID:           uuid.New(),
		TestimonialID: testimonialID,
		EntityType:   entityType,
		EntityID:     entityID,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	return link, link.Validate()
}

// Validate ensures testimonial entity link invariants hold.
func (l *TestimonialEntityLink) Validate() error {
	if l == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilLink)
	}
	if l.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyLinkID)
	}
	if l.TestimonialID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}
	if l.EntityID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyEntityID)
	}
	if !isValidEntityType(l.EntityType) {
		return NewDomainError(ErrCodeInvalidEntityType, ErrUnsupportedEntityType)
	}
	return nil
}

// Validation helpers

func isValidTestimonialType(tt TestimonialType) bool {
	switch tt {
	case TestimonialTypeGeneral, TestimonialTypeProjectSpecific, TestimonialTypeSkillSpecific:
		return true
	}
	return false
}

func isValidEntityType(et EntityType) bool {
	switch et {
	case EntityTypeProject, EntityTypeSkill:
		return true
	}
	return false
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

