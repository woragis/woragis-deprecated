package posts

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

// PostStatus represents the publication status of a post.
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

// Post represents a blog post.
type Post struct {
	ID              uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID  `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Title           string     `gorm:"column:title;size:255;not null" json:"title"`
	Slug            string     `gorm:"column:slug;size:255;not null;uniqueIndex:idx_post_slug" json:"slug"`
	Content         string     `gorm:"column:content;type:text;not null" json:"content"` // Markdown content
	Excerpt         string     `gorm:"column:excerpt;type:text" json:"excerpt,omitempty"`
	Status          PostStatus `gorm:"column:status;type:varchar(32);not null;default:'draft';index" json:"status"`
	PublishedAt     *time.Time `gorm:"column:published_at;index" json:"publishedAt,omitempty"`
	FeaturedImage   string     `gorm:"column:featured_image;size:512" json:"featuredImage,omitempty"`
	MetaTitle       string     `gorm:"column:meta_title;size:255" json:"metaTitle,omitempty"`
	MetaDescription string     `gorm:"column:meta_description;size:512" json:"metaDescription,omitempty"`
	MetaKeywords    string     `gorm:"column:meta_keywords;size:255" json:"metaKeywords,omitempty"`
	OGTitle         string     `gorm:"column:og_title;size:255" json:"ogTitle,omitempty"`
	OGDescription   string     `gorm:"column:og_description;size:512" json:"ogDescription,omitempty"`
	OGImage         string     `gorm:"column:og_image;size:512" json:"ogImage,omitempty"`
	Featured        bool       `gorm:"column:featured;not null;default:false;index" json:"featured"`
	ViewsCount      int64      `gorm:"column:views_count;not null;default:0" json:"viewsCount"`
	CreatedAt       time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"column:updated_at" json:"updatedAt"`

	// Relationships (not stored in DB, loaded via joins in repository)
	// Skills, Categories, Tags are loaded separately via repository methods
}

// PostSkill represents the many-to-many relationship between posts and skills.
type PostSkill struct {
	PostID    uuid.UUID `gorm:"column:post_id;type:uuid;primaryKey;index" json:"postId"`
	SkillID   uuid.UUID `gorm:"column:skill_id;type:uuid;primaryKey;index" json:"skillId"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

// TableName specifies the table name for PostSkill.
func (PostSkill) TableName() string {
	return "post_skills"
}

// Category represents a post category.
type Category struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"column:name;size:120;not null;uniqueIndex:idx_category_name" json:"name"`
	Slug        string    `gorm:"column:slug;size:160;not null;uniqueIndex:idx_category_slug" json:"slug"`
	Description string    `gorm:"column:description;size:512" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// PostCategory represents the many-to-many relationship between posts and categories.
type PostCategory struct {
	PostID     uuid.UUID `gorm:"column:post_id;type:uuid;primaryKey;index" json:"postId"`
	CategoryID uuid.UUID `gorm:"column:category_id;type:uuid;primaryKey;index" json:"categoryId"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt"`
}

// TableName specifies the table name for PostCategory.
func (PostCategory) TableName() string {
	return "post_categories"
}

// Tag represents a post tag.
type Tag struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"column:name;size:80;not null;uniqueIndex:idx_tag_name" json:"name"`
	Slug      string    `gorm:"column:slug;size:120;not null;uniqueIndex:idx_tag_slug" json:"slug"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// PostTag represents the many-to-many relationship between posts and tags.
type PostTag struct {
	PostID    uuid.UUID `gorm:"column:post_id;type:uuid;primaryKey;index" json:"postId"`
	TagID     uuid.UUID `gorm:"column:tag_id;type:uuid;primaryKey;index" json:"tagId"`
	CreatedAt time.Time `gorm:"column:created_at" json:"createdAt"`
}

// TableName specifies the table name for PostTag.
func (PostTag) TableName() string {
	return "post_tags"
}

// NewPost creates a new post entity.
func NewPost(userID uuid.UUID, title, content, excerpt string, status PostStatus) (*Post, error) {
	post := &Post{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     strings.TrimSpace(title),
		Content:   content,
		Excerpt:   strings.TrimSpace(excerpt),
		Status:    status,
		Featured:  false,
		ViewsCount: 0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	post.Slug = generatePostSlug(post.Title)

	if status == PostStatusPublished {
		now := time.Now().UTC()
		post.PublishedAt = &now
	}

	return post, post.Validate()
}

// Validate ensures post invariants hold.
func (p *Post) Validate() error {
	if p == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilPost)
	}

	if p.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyPostID)
	}

	if p.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if p.Title == "" {
		return NewDomainError(ErrCodeInvalidTitle, ErrEmptyPostTitle)
	}

	if strings.TrimSpace(p.Slug) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyPostSlug)
	}

	if p.Content == "" {
		return NewDomainError(ErrCodeInvalidContent, ErrEmptyPostContent)
	}

	switch p.Status {
	case PostStatusDraft, PostStatusPublished, PostStatusArchived:
	default:
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedPostStatus)
	}

	return nil
}

// Publish marks the post as published and sets the published_at timestamp.
func (p *Post) Publish() error {
	if p.Status == PostStatusPublished {
		return nil // Already published
	}

	p.Status = PostStatusPublished
	now := time.Now().UTC()
	p.PublishedAt = &now
	p.UpdatedAt = time.Now().UTC()

	return nil
}

// Unpublish marks the post as draft and clears the published_at timestamp.
func (p *Post) Unpublish() {
	p.Status = PostStatusDraft
	p.PublishedAt = nil
	p.UpdatedAt = time.Now().UTC()
}

// Archive marks the post as archived.
func (p *Post) Archive() {
	p.Status = PostStatusArchived
	p.UpdatedAt = time.Now().UTC()
}

// IncrementViews increments the view count.
func (p *Post) IncrementViews() {
	p.ViewsCount++
	p.UpdatedAt = time.Now().UTC()
}

// UpdateContent updates the post content.
func (p *Post) UpdateContent(content string) error {
	if content == "" {
		return NewDomainError(ErrCodeInvalidContent, ErrEmptyPostContent)
	}
	p.Content = content
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateTitle updates the post title and regenerates the slug.
func (p *Post) UpdateTitle(title string) error {
	if title == "" {
		return NewDomainError(ErrCodeInvalidTitle, ErrEmptyPostTitle)
	}
	p.Title = strings.TrimSpace(title)
	p.Slug = generatePostSlug(p.Title)
	p.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateSEO updates SEO-related fields.
func (p *Post) UpdateSEO(metaTitle, metaDescription, metaKeywords, ogTitle, ogDescription, ogImage string) {
	if metaTitle != "" {
		p.MetaTitle = strings.TrimSpace(metaTitle)
	}
	if metaDescription != "" {
		p.MetaDescription = strings.TrimSpace(metaDescription)
	}
	if metaKeywords != "" {
		p.MetaKeywords = strings.TrimSpace(metaKeywords)
	}
	if ogTitle != "" {
		p.OGTitle = strings.TrimSpace(ogTitle)
	}
	if ogDescription != "" {
		p.OGDescription = strings.TrimSpace(ogDescription)
	}
	if ogImage != "" {
		p.OGImage = strings.TrimSpace(ogImage)
	}
	p.UpdatedAt = time.Now().UTC()
}

// SetFeatured sets the featured flag.
func (p *Post) SetFeatured(featured bool) {
	p.Featured = featured
	p.UpdatedAt = time.Now().UTC()
}

var slugSanitizer = regexp.MustCompile(`[^a-z0-9]+`)

func generatePostSlug(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = slugSanitizer.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "post"
	}
	return slug
}

// NewCategory creates a new category entity.
func NewCategory(name, description string) (*Category, error) {
	category := &Category{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	category.Slug = generateCategorySlug(category.Name)

	return category, category.Validate()
}

// Validate ensures category invariants hold.
func (c *Category) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilCategory)
	}

	if c.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCategoryID)
	}

	if c.Name == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyCategoryName)
	}

	if strings.TrimSpace(c.Slug) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCategorySlug)
	}

	return nil
}

func generateCategorySlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = slugSanitizer.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "category"
	}
	return slug
}

// NewTag creates a new tag entity.
func NewTag(name string) (*Tag, error) {
	tag := &Tag{
		ID:        uuid.New(),
		Name:      strings.TrimSpace(name),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	tag.Slug = generateTagSlug(tag.Name)

	return tag, tag.Validate()
}

// Validate ensures tag invariants hold.
func (t *Tag) Validate() error {
	if t == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilTag)
	}

	if t.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTagID)
	}

	if t.Name == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyTagName)
	}

	if strings.TrimSpace(t.Slug) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTagSlug)
	}

	return nil
}

func generateTagSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = slugSanitizer.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		slug = "tag"
	}
	return slug
}

