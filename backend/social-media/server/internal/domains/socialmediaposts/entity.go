package socialmediaposts

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Platform represents the social media platform.
type Platform string

const (
	PlatformLinkedIn  Platform = "linkedin"
	PlatformTwitter   Platform = "twitter"
	PlatformInstagram Platform = "instagram"
	PlatformMedium    Platform = "medium"
	PlatformSubstack  Platform = "substack"
	PlatformValete    Platform = "valete"
	PlatformWebsite   Platform = "website"
)

// ContentFormat represents the format of the social media post content.
type ContentFormat string

const (
	FormatLongForm   ContentFormat = "long-form"
	FormatThread     ContentFormat = "thread"
	FormatCarousel   ContentFormat = "carousel"
	FormatArticle    ContentFormat = "article"
	FormatNewsletter ContentFormat = "newsletter"
	FormatPost       ContentFormat = "post"
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

// PostStatus represents the status of a social media post in the workflow.
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusReady     PostStatus = "ready"
	PostStatusScheduled PostStatus = "scheduled"
	PostStatusPosted    PostStatus = "posted"
	PostStatusAnalyzed  PostStatus = "analyzed"
	PostStatusArchived  PostStatus = "archived"
)

// SocialMediaPost represents a social media post with its metadata.
type SocialMediaPost struct {
	ID            uuid.UUID      `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ContentPostID *uuid.UUID     `gorm:"column:content_post_id;type:uuid;index" json:"contentPostId,omitempty"` // Link to content post (for repurposing)
	Platform      Platform       `gorm:"column:platform;type:varchar(20);not null;index" json:"platform"`
	Format        ContentFormat  `gorm:"column:format;type:varchar(20);not null;index" json:"format"`
	Status        PostStatus     `gorm:"column:status;type:varchar(20);not null;default:'draft';index" json:"status"`
	
	// Content fields
	Title      string `gorm:"column:title;size:255;not null" json:"title"`
	Content    string `gorm:"column:content;type:text;not null" json:"content"` // Platform-adapted content
	WordCount  int    `gorm:"column:word_count;not null;default:0" json:"wordCount"`
	ImageCount int    `gorm:"column:image_count;not null;default:0" json:"imageCount"`
	
	// Scheduling fields
	ScheduledAt *time.Time `gorm:"column:scheduled_at;type:timestamp;index" json:"scheduledAt,omitempty"`
	PostedAt    *time.Time `gorm:"column:posted_at;type:timestamp;index" json:"postedAt,omitempty"`
	AnalyzedAt  *time.Time `gorm:"column:analyzed_at;type:timestamp" json:"analyzedAt,omitempty"`
	
	// Posting metadata (optional, set after posting)
	URL            *string `gorm:"column:url;size:512;uniqueIndex:idx_social_post_url" json:"url,omitempty"` // Optional, generated after posting
	PlatformPostID string  `gorm:"column:platform_post_id;size:255" json:"platformPostId,omitempty"`        // External platform post ID
	
	// Engagement metrics
	Likes    *int64 `gorm:"column:likes;type:bigint" json:"likes,omitempty"`
	Shares   *int64 `gorm:"column:shares;type:bigint" json:"shares,omitempty"`
	Comments *int64 `gorm:"column:comments;type:bigint" json:"comments,omitempty"`
	Views    *int64 `gorm:"column:views;type:bigint" json:"views,omitempty"`
	
	// Additional metadata
	Metadata datatypes.JSON `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	
	// Timestamps
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
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
func NewSocialMediaPost(platform Platform, format ContentFormat, title, content string) (*SocialMediaPost, error) {
	post := &SocialMediaPost{
		ID:        uuid.New(),
		Platform:  platform,
		Format:    format,
		Title:     strings.TrimSpace(title),
		Content:   content,
		Status:    PostStatusDraft,
		WordCount: 0,
		ImageCount: 0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	// Calculate word count
	post.WordCount = countWords(post.Content)

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

	if !isValidPlatform(p.Platform) {
		return NewDomainError(ErrCodeInvalidPlatform, ErrUnsupportedPlatform)
	}

	if !isValidContentFormat(p.Format) {
		return NewDomainError(ErrCodeInvalidFormat, ErrUnsupportedFormat)
	}

	if !isValidPostStatus(p.Status) {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}

	if strings.TrimSpace(p.Title) == "" {
		return NewDomainError(ErrCodeInvalidTitle, ErrEmptyPostTitle)
	}

	if strings.TrimSpace(p.Content) == "" {
		return NewDomainError(ErrCodeInvalidContent, ErrEmptyPostContent)
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

// UpdateStatus updates the post status with workflow validation.
func (p *SocialMediaPost) UpdateStatus(status PostStatus) error {
	if !isValidPostStatus(status) {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}

	// Validate status transitions
	if err := p.validateStatusTransition(status); err != nil {
		return err
	}

	p.Status = status
	now := time.Now().UTC()
	p.UpdatedAt = now

	// Set timestamps based on status
	switch status {
	case PostStatusScheduled:
		if p.ScheduledAt == nil {
			p.ScheduledAt = &now
		}
	case PostStatusPosted:
		if p.PostedAt == nil {
			p.PostedAt = &now
		}
	case PostStatusAnalyzed:
		if p.AnalyzedAt == nil {
			p.AnalyzedAt = &now
		}
	}

	return nil
}

// validateStatusTransition ensures status changes follow the workflow.
func (p *SocialMediaPost) validateStatusTransition(newStatus PostStatus) error {
	// Define valid transitions
	validTransitions := map[PostStatus][]PostStatus{
		PostStatusDraft:     {PostStatusReady, PostStatusArchived},
		PostStatusReady:     {PostStatusScheduled, PostStatusDraft, PostStatusArchived},
		PostStatusScheduled: {PostStatusPosted, PostStatusDraft, PostStatusArchived},
		PostStatusPosted:    {PostStatusAnalyzed, PostStatusArchived},
		PostStatusAnalyzed:  {PostStatusArchived},
		PostStatusArchived:  {}, // Terminal state, no transitions allowed
	}

	allowed, ok := validTransitions[p.Status]
	if !ok {
		return NewDomainError(ErrCodeInvalidStatus, "invalid current status")
	}

	for _, allowedStatus := range allowed {
		if newStatus == allowedStatus {
			return nil
		}
	}

	return NewDomainError(ErrCodeInvalidStatusTransition, ErrInvalidStatusTransition)
}

// SetContentPostID links this social post to a content post.
func (p *SocialMediaPost) SetContentPostID(contentPostID uuid.UUID) {
	p.ContentPostID = &contentPostID
	p.UpdatedAt = time.Now().UTC()
}

// SetScheduledTime sets the scheduled time for the post.
func (p *SocialMediaPost) SetScheduledTime(scheduledAt time.Time) error {
	if p.Status != PostStatusReady && p.Status != PostStatusScheduled {
		return NewDomainError(ErrCodeInvalidStatusTransition, "can only schedule posts with status 'ready' or 'scheduled'")
	}
	p.ScheduledAt = &scheduledAt
	if p.Status != PostStatusScheduled {
		return p.UpdateStatus(PostStatusScheduled)
	}
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// MarkAsPosted marks the post as posted and sets the URL and platform post ID.
func (p *SocialMediaPost) MarkAsPosted(url, platformPostID string) error {
	if p.Status != PostStatusScheduled && p.Status != PostStatusReady {
		return NewDomainError(ErrCodeInvalidStatusTransition, "can only mark scheduled or ready posts as posted")
	}

	if url != "" {
		p.URL = &url
	}
	if platformPostID != "" {
		p.PlatformPostID = platformPostID
	}

	now := time.Now().UTC()
	p.PostedAt = &now

	return p.UpdateStatus(PostStatusPosted)
}

// SetURL sets the URL for the post (typically after posting).
func (p *SocialMediaPost) SetURL(url string) {
	if url != "" {
		p.URL = &url
		p.UpdatedAt = time.Now().UTC()
	}
}

// UpdateContent updates the post content and recalculates word count.
func (p *SocialMediaPost) UpdateContent(content string) error {
	if strings.TrimSpace(content) == "" {
		return NewDomainError(ErrCodeInvalidContent, ErrEmptyPostContent)
	}

	if p.Status == PostStatusPosted || p.Status == PostStatusAnalyzed || p.Status == PostStatusArchived {
		return NewDomainError(ErrCodeInvalidStatusTransition, "cannot update content of posted/analyzed/archived posts")
	}

	p.Content = content
	p.WordCount = countWords(content)
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateTitle updates the post title.
func (p *SocialMediaPost) UpdateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return NewDomainError(ErrCodeInvalidTitle, ErrEmptyPostTitle)
	}

	if p.Status == PostStatusPosted || p.Status == PostStatusAnalyzed || p.Status == PostStatusArchived {
		return NewDomainError(ErrCodeInvalidStatusTransition, "cannot update title of posted/analyzed/archived posts")
	}

	p.Title = strings.TrimSpace(title)
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
	case PlatformLinkedIn, PlatformTwitter, PlatformInstagram,
		PlatformMedium, PlatformSubstack, PlatformValete, PlatformWebsite:
		return true
	}
	return false
}

func isValidContentFormat(f ContentFormat) bool {
	switch f {
	case FormatLongForm, FormatThread, FormatCarousel,
		FormatArticle, FormatNewsletter, FormatPost:
		return true
	}
	return false
}

func isValidEntityType(et EntityType) bool {
	switch et {
	case EntityTypePost, EntityTypeCaseStudy, EntityTypeProjectCaseStudy,
		EntityTypeSkill, EntityTypeProject, EntityTypeProblemSolution,
		EntityTypeSystemDesign, EntityTypeInterest, EntityTypeTestimonial,
		EntityTypeCertification:
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
	case PostStatusDraft, PostStatusReady, PostStatusScheduled,
		PostStatusPosted, PostStatusAnalyzed, PostStatusArchived:
		return true
	}
	return false
}

// countWords calculates the approximate word count of a text string.
func countWords(text string) int {
	if text == "" {
		return 0
	}
	words := strings.Fields(text)
	return len(words)
}
