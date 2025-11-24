package socialmediaposts

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Platform represents the social media platform.
type Platform string

const (
	PlatformLinkedIn Platform = "linkedin"
	PlatformTwitter   Platform = "twitter"
	PlatformInstagram Platform = "instagram"
)

// EntityType represents the type of entity being linked.
type EntityType string

const (
	EntityTypePost             EntityType = "post"
	EntityTypeCaseStudy        EntityType = "case_study"
	EntityTypeProjectCaseStudy EntityType = "project_case_study"
	EntityTypeSkill            EntityType = "skill"
	EntityTypeProject          EntityType = "project"
	EntityTypeProblemSolution  EntityType = "problem_solution"
	EntityTypeSystemDesign     EntityType = "system_design"
	EntityTypeInterest         EntityType = "interest"
	EntityTypeTestimonial      EntityType = "testimonial"
	EntityTypeCertification    EntityType = "certification"
)

// RelationshipType represents how prominently an entity is mentioned in a social media post.
type RelationshipType string

const (
	RelationshipTypeMainTopic        RelationshipType = "main_topic"
	RelationshipTypeSecondaryTopic   RelationshipType = "secondary_topic"
	RelationshipTypeMentionedBriefly RelationshipType = "mentioned_briefly"
	RelationshipTypeExample          RelationshipType = "example"
	RelationshipTypeComparison       RelationshipType = "comparison"
)

// PostStatus represents the status of a social media post.
type PostStatus string

const (
	PostStatusActive      PostStatus = "active"
	PostStatusDeleted     PostStatus = "deleted"
	PostStatusUnavailable PostStatus = "unavailable"
)

// SocialMediaPost represents a social media post with its metadata.
type SocialMediaPost struct {
	ID              uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	URL             string    `gorm:"column:url;size:512;not null;uniqueIndex:idx_social_post_url" json:"url"`
	Platform        Platform  `gorm:"column:platform;type:varchar(20);not null;index" json:"platform"`
	Title           string    `gorm:"column:title;size:255" json:"title,omitempty"` // Optional title/preview
	ContentPreview  string    `gorm:"column:content_preview;type:text" json:"contentPreview,omitempty"` // Optional content preview
	PublishedDate   *time.Time `gorm:"column:published_date;type:timestamp" json:"publishedDate,omitempty"`
	Likes           *int64    `gorm:"column:likes;type:bigint" json:"likes,omitempty"`
	Shares          *int64    `gorm:"column:shares;type:bigint" json:"shares,omitempty"`
	Comments        *int64    `gorm:"column:comments;type:bigint" json:"comments,omitempty"`
	Views           *int64    `gorm:"column:views;type:bigint" json:"views,omitempty"`
	Status          PostStatus `gorm:"column:status;type:varchar(20);not null;default:'active';index" json:"status"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for SocialMediaPost.
func (SocialMediaPost) TableName() string {
	return "social_media_posts"
}

// SocialMediaEntityLink represents the relationship between a social media post and an entity.
type SocialMediaEntityLink struct {
	ID               uuid.UUID        `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	SocialMediaPostID uuid.UUID       `gorm:"column:social_media_post_id;type:uuid;not null;index:idx_social_entity_link" json:"socialMediaPostId"`
	EntityType       EntityType       `gorm:"column:entity_type;type:varchar(50);not null;index:idx_social_entity_link" json:"entityType"`
	EntityID         uuid.UUID        `gorm:"column:entity_id;type:uuid;not null;index:idx_social_entity_link" json:"entityId"`
	RelationshipType RelationshipType `gorm:"column:relationship_type;type:varchar(30);not null;index" json:"relationshipType"`
	CreatedAt        time.Time        `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time        `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for SocialMediaEntityLink.
func (SocialMediaEntityLink) TableName() string {
	return "social_media_entity_links"
}

// NewSocialMediaPost creates a new social media post entity.
func NewSocialMediaPost(url string, platform Platform, title, contentPreview string, publishedDate *time.Time) (*SocialMediaPost, error) {
	post := &SocialMediaPost{
		ID:             uuid.New(),
		URL:            strings.TrimSpace(url),
		Platform:       platform,
		Title:          strings.TrimSpace(title),
		ContentPreview: strings.TrimSpace(contentPreview),
		PublishedDate:  publishedDate,
		Status:         PostStatusActive,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	return post, post.Validate()
}

// Validate ensures social media post invariants hold.
func (p *SocialMediaPost) Validate() error {
	if p == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilPost)
	}

	if p.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyPostID)
	}

	if strings.TrimSpace(p.URL) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyURL)
	}

	if !isValidPlatform(p.Platform) {
		return NewDomainError(ErrCodeInvalidPlatform, ErrUnsupportedPlatform)
	}

	if !isValidPostStatus(p.Status) {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}

	return nil
}

// UpdateEngagement updates engagement metrics.
func (p *SocialMediaPost) UpdateEngagement(likes, shares, comments, views *int64) {
	if likes != nil {
		p.Likes = likes
	}
	if shares != nil {
		p.Shares = shares
	}
	if comments != nil {
		p.Comments = comments
	}
	if views != nil {
		p.Views = views
	}
	p.UpdatedAt = time.Now().UTC()
}

// UpdateStatus updates the post status.
func (p *SocialMediaPost) UpdateStatus(status PostStatus) error {
	if !isValidPostStatus(status) {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}
	p.Status = status
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// NewSocialMediaEntityLink creates a new link between a social media post and an entity.
func NewSocialMediaEntityLink(postID uuid.UUID, entityType EntityType, entityID uuid.UUID, relationshipType RelationshipType) (*SocialMediaEntityLink, error) {
	link := &SocialMediaEntityLink{
		ID:               uuid.New(),
		SocialMediaPostID: postID,
		EntityType:       entityType,
		EntityID:         entityID,
		RelationshipType: relationshipType,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	return link, link.Validate()
}

// Validate ensures social media entity link invariants hold.
func (l *SocialMediaEntityLink) Validate() error {
	if l == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilLink)
	}

	if l.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyLinkID)
	}

	if l.SocialMediaPostID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyPostID)
	}

	if l.EntityID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyEntityID)
	}

	if !isValidEntityType(l.EntityType) {
		return NewDomainError(ErrCodeInvalidEntityType, ErrUnsupportedEntityType)
	}

	if !isValidRelationshipType(l.RelationshipType) {
		return NewDomainError(ErrCodeInvalidRelationshipType, ErrUnsupportedRelationshipType)
	}

	return nil
}

// Validation helpers

func isValidPlatform(p Platform) bool {
	switch p {
	case PlatformLinkedIn, PlatformTwitter, PlatformInstagram:
		return true
	}
	return false
}

func isValidEntityType(et EntityType) bool {
	switch et {
	case EntityTypePost, EntityTypeCaseStudy, EntityTypeProjectCaseStudy,
		EntityTypeSkill, EntityTypeProject, EntityTypeProblemSolution,
		EntityTypeSystemDesign, EntityTypeInterest, EntityTypeTestimonial:
		return true
	}
	return false
}

func isValidRelationshipType(rt RelationshipType) bool {
	switch rt {
	case RelationshipTypeMainTopic, RelationshipTypeSecondaryTopic,
		RelationshipTypeMentionedBriefly, RelationshipTypeExample,
		RelationshipTypeComparison:
		return true
	}
	return false
}

func isValidPostStatus(s PostStatus) bool {
	switch s {
	case PostStatusActive, PostStatusDeleted, PostStatusUnavailable:
		return true
	}
	return false
}

