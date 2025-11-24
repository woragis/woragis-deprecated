package technicalwritings

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

// WritingType represents the type of technical writing.
type WritingType string

const (
	WritingTypeArticle       WritingType = "article"
	WritingTypeDocumentation WritingType = "documentation"
	WritingTypeTutorial      WritingType = "tutorial"
	WritingTypeGuide         WritingType = "guide"
	WritingTypeBlogPost      WritingType = "blog_post"
	WritingTypeCaseStudy     WritingType = "case_study"
	WritingTypeOther         WritingType = "other"
)

// PublicationPlatform represents where the writing was published.
type PublicationPlatform string

const (
	PlatformMedium      PublicationPlatform = "medium"
	PlatformDevTo       PublicationPlatform = "dev_to"
	PlatformHashnode    PublicationPlatform = "hashnode"
	PlatformPersonalBlog PublicationPlatform = "personal_blog"
	PlatformGitHub       PublicationPlatform = "github"
	PlatformCompanyBlog  PublicationPlatform = "company_blog"
	PlatformSubstack     PublicationPlatform = "substack"
	PlatformLinkedIn     PublicationPlatform = "linkedin"
	PlatformOther        PublicationPlatform = "other"
)

// TechnicalWriting represents a technical writing portfolio entry.
type TechnicalWriting struct {
	ID                uuid.UUID          `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID            uuid.UUID          `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Title             string             `gorm:"column:title;size:255;not null" json:"title"`
	Description       string             `gorm:"column:description;type:text;not null" json:"description"`
	Type              WritingType        `gorm:"column:type;type:varchar(50);not null;index" json:"type"`
	Platform          PublicationPlatform `gorm:"column:platform;type:varchar(50);not null;index" json:"platform"`
	// Content and URLs
	Content           string             `gorm:"column:content;type:text" json:"content,omitempty"` // Full content or excerpt
	URL               string             `gorm:"column:url;size:500;not null" json:"url"`
	CanonicalURL      string             `gorm:"column:canonical_url;size:500" json:"canonicalUrl,omitempty"`
	// Publication details
	PublishedAt       *time.Time         `gorm:"column:published_at;type:timestamp" json:"publishedAt,omitempty"`
	ReadingTime       int                `gorm:"column:reading_time;default:0" json:"readingTime,omitempty"` // in minutes
	// Topics and technologies
	Topics             JSONArray          `gorm:"column:topics;type:jsonb" json:"topics,omitempty"` // Array of topic tags
	Technologies       JSONArray          `gorm:"column:technologies;type:jsonb" json:"technologies,omitempty"` // Array of technology names
	// Metrics (optional, can be updated over time)
	Views              *int               `gorm:"column:views" json:"views,omitempty"`
	Likes              *int               `gorm:"column:likes" json:"likes,omitempty"`
	Shares             *int               `gorm:"column:shares" json:"shares,omitempty"`
	Comments           *int               `gorm:"column:comments" json:"comments,omitempty"`
	// Links to other entities
	ProjectID          *uuid.UUID         `gorm:"column:project_id;type:uuid;index" json:"projectId,omitempty"`
	CaseStudyID        *uuid.UUID         `gorm:"column:case_study_id;type:uuid;index" json:"caseStudyId,omitempty"`
	// Display and organization
	Featured           bool               `gorm:"column:featured;not null;default:false;index" json:"featured"`
	DisplayOrder       int                `gorm:"column:display_order;not null;default:0;index" json:"displayOrder"`
	// Metadata
	Excerpt            string             `gorm:"column:excerpt;type:text" json:"excerpt,omitempty"` // Short excerpt for previews
	CoverImageURL      string             `gorm:"column:cover_image_url;size:500" json:"coverImageUrl,omitempty"`
	// Timestamps
	CreatedAt          time.Time          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          time.Time          `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for TechnicalWriting.
func (TechnicalWriting) TableName() string {
	return "technical_writings"
}

// NewTechnicalWriting creates a new technical writing entity.
func NewTechnicalWriting(userID uuid.UUID, title, description string, writingType WritingType, platform PublicationPlatform, url string) (*TechnicalWriting, error) {
	writing := &TechnicalWriting{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Type:        writingType,
		Platform:    platform,
		URL:         strings.TrimSpace(url),
		Topics:      JSONArray{},
		Technologies: JSONArray{},
		Featured:    false,
		DisplayOrder: 0,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	return writing, writing.Validate()
}

// Validate ensures technical writing invariants hold.
func (t *TechnicalWriting) Validate() error {
	if t == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilWriting)
	}
	if t.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyWritingID)
	}
	if t.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if strings.TrimSpace(t.Title) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTitle)
	}
	if strings.TrimSpace(t.Description) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDescription)
	}
	if strings.TrimSpace(t.URL) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyURL)
	}
	if !isValidWritingType(t.Type) {
		return NewDomainError(ErrCodeInvalidType, ErrUnsupportedWritingType)
	}
	if !isValidPlatform(t.Platform) {
		return NewDomainError(ErrCodeInvalidPlatform, ErrUnsupportedPlatform)
	}
	return nil
}

// UpdateDetails updates writing details.
func (t *TechnicalWriting) UpdateDetails(title, description, content, excerpt string) {
	if title != "" {
		t.Title = strings.TrimSpace(title)
	}
	if description != "" {
		t.Description = strings.TrimSpace(description)
	}
	if content != "" {
		t.Content = strings.TrimSpace(content)
	}
	if excerpt != "" {
		t.Excerpt = strings.TrimSpace(excerpt)
	}
	t.UpdatedAt = time.Now().UTC()
}

// SetPublicationInfo updates publication information.
func (t *TechnicalWriting) SetPublicationInfo(publishedAt *time.Time, readingTime int) {
	if publishedAt != nil {
		t.PublishedAt = publishedAt
	}
	if readingTime > 0 {
		t.ReadingTime = readingTime
	}
	t.UpdatedAt = time.Now().UTC()
}

// SetMetrics updates engagement metrics.
func (t *TechnicalWriting) SetMetrics(views, likes, shares, comments *int) {
	if views != nil {
		t.Views = views
	}
	if likes != nil {
		t.Likes = likes
	}
	if shares != nil {
		t.Shares = shares
	}
	if comments != nil {
		t.Comments = comments
	}
	t.UpdatedAt = time.Now().UTC()
}

// SetTopics updates the topics array.
func (t *TechnicalWriting) SetTopics(topics []string) {
	t.Topics = JSONArray(topics)
	t.UpdatedAt = time.Now().UTC()
}

// SetTechnologies updates the technologies array.
func (t *TechnicalWriting) SetTechnologies(technologies []string) {
	t.Technologies = JSONArray(technologies)
	t.UpdatedAt = time.Now().UTC()
}

// SetProjectLink updates the project link.
func (t *TechnicalWriting) SetProjectLink(projectID uuid.UUID) {
	t.ProjectID = &projectID
	t.UpdatedAt = time.Now().UTC()
}

// ClearProjectLink removes the project link.
func (t *TechnicalWriting) ClearProjectLink() {
	t.ProjectID = nil
	t.UpdatedAt = time.Now().UTC()
}

// SetCaseStudyLink updates the case study link.
func (t *TechnicalWriting) SetCaseStudyLink(caseStudyID uuid.UUID) {
	t.CaseStudyID = &caseStudyID
	t.UpdatedAt = time.Now().UTC()
}

// ClearCaseStudyLink removes the case study link.
func (t *TechnicalWriting) ClearCaseStudyLink() {
	t.CaseStudyID = nil
	t.UpdatedAt = time.Now().UTC()
}

// SetFeatured updates the featured flag.
func (t *TechnicalWriting) SetFeatured(featured bool) {
	t.Featured = featured
	t.UpdatedAt = time.Now().UTC()
}

// SetDisplayOrder updates the display order.
func (t *TechnicalWriting) SetDisplayOrder(order int) {
	t.DisplayOrder = order
	t.UpdatedAt = time.Now().UTC()
}

// SetURLs updates URL fields.
func (t *TechnicalWriting) SetURLs(url, canonicalURL, coverImageURL string) {
	if url != "" {
		t.URL = strings.TrimSpace(url)
	}
	if canonicalURL != "" {
		t.CanonicalURL = strings.TrimSpace(canonicalURL)
	}
	if coverImageURL != "" {
		t.CoverImageURL = strings.TrimSpace(coverImageURL)
	}
	t.UpdatedAt = time.Now().UTC()
}

// Validation helpers

func isValidWritingType(wt WritingType) bool {
	switch wt {
	case WritingTypeArticle, WritingTypeDocumentation, WritingTypeTutorial,
		WritingTypeGuide, WritingTypeBlogPost, WritingTypeCaseStudy, WritingTypeOther:
		return true
	}
	return false
}

func isValidPlatform(p PublicationPlatform) bool {
	switch p {
	case PlatformMedium, PlatformDevTo, PlatformHashnode, PlatformPersonalBlog,
		PlatformGitHub, PlatformCompanyBlog, PlatformSubstack, PlatformLinkedIn, PlatformOther:
		return true
	}
	return false
}

// JSONArray is a custom type for storing JSON arrays in PostgreSQL.
type JSONArray []string

// Value implements the driver.Valuer interface.
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface.
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), j)
	}
	return json.Unmarshal(bytes, j)
}

