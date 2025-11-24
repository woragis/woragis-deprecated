package posts

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for posts.
type Repository interface {
	// Post operations
	CreatePost(ctx context.Context, post *Post) error
	UpdatePost(ctx context.Context, post *Post) error
	GetPost(ctx context.Context, postID uuid.UUID) (*Post, error)
	GetPostBySlug(ctx context.Context, slug string) (*Post, error)
	DeletePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error
	ListPosts(ctx context.Context, filters PostFilters) ([]Post, error)
	IsPostSlugTaken(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error)
	IncrementPostViews(ctx context.Context, postID uuid.UUID) error

	// Category operations
	CreateCategory(ctx context.Context, category *Category) error
	UpdateCategory(ctx context.Context, category *Category) error
	GetCategory(ctx context.Context, categoryID uuid.UUID) (*Category, error)
	GetCategoryBySlug(ctx context.Context, slug string) (*Category, error)
	ListCategories(ctx context.Context) ([]Category, error)
	IsCategorySlugTaken(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error)

	// Tag operations
	CreateTag(ctx context.Context, tag *Tag) error
	GetTag(ctx context.Context, tagID uuid.UUID) (*Tag, error)
	GetTagBySlug(ctx context.Context, slug string) (*Tag, error)
	GetOrCreateTag(ctx context.Context, name string) (*Tag, error)
	ListTags(ctx context.Context) ([]Tag, error)
	IsTagSlugTaken(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error)

	// Post-Skill relationship operations
	AttachSkillToPost(ctx context.Context, postID, skillID uuid.UUID) error
	DetachSkillFromPost(ctx context.Context, postID, skillID uuid.UUID) error
	GetPostSkills(ctx context.Context, postID uuid.UUID) ([]uuid.UUID, error) // Returns skill IDs
	PostHasSkill(ctx context.Context, postID, skillID uuid.UUID) (bool, error)

	// Post-Category relationship operations
	AttachCategoryToPost(ctx context.Context, postID, categoryID uuid.UUID) error
	DetachCategoryFromPost(ctx context.Context, postID, categoryID uuid.UUID) error
	GetPostCategories(ctx context.Context, postID uuid.UUID) ([]Category, error)
	PostHasCategory(ctx context.Context, postID, categoryID uuid.UUID) (bool, error)

	// Post-Tag relationship operations
	AttachTagToPost(ctx context.Context, postID, tagID uuid.UUID) error
	DetachTagFromPost(ctx context.Context, postID, tagID uuid.UUID) error
	GetPostTags(ctx context.Context, postID uuid.UUID) ([]Tag, error)
	PostHasTag(ctx context.Context, postID, tagID uuid.UUID) (bool, error)
}

// PostFilters represents filtering options for listing posts.
type PostFilters struct {
	UserID     *uuid.UUID
	Status     *PostStatus
	Featured   *bool
	CategoryID *uuid.UUID
	TagID      *uuid.UUID
	SkillID    *uuid.UUID
	Search     string
	Limit      int
	Offset     int
	OrderBy    string // "created_at", "updated_at", "published_at", "views_count"
	Order      string // "asc", "desc"
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// isUniqueConstraintError checks if the error is a unique constraint violation.
func isUniqueConstraintError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}

// Post operations

func (r *gormRepository) CreatePost(ctx context.Context, post *Post) error {
	if err := post.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeDuplicateSlug, ErrPostSlugTaken)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdatePost(ctx context.Context, post *Post) error {
	if err := post.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(post).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeDuplicateSlug, ErrPostSlugTaken)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetPost(ctx context.Context, postID uuid.UUID) (*Post, error) {
	var post Post
	err := r.db.WithContext(ctx).Where("id = ?", postID).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodePostNotFound, ErrPostNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &post, nil
}

func (r *gormRepository) GetPostBySlug(ctx context.Context, slug string) (*Post, error) {
	var post Post
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodePostNotFound, ErrPostNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &post, nil
}

func (r *gormRepository) DeletePost(ctx context.Context, postID uuid.UUID, userID uuid.UUID) error {
	// First verify ownership
	var post Post
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", postID, userID).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return NewDomainError(ErrCodePostNotFound, ErrPostNotFound)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	// Delete relationships first
	r.db.WithContext(ctx).Where("post_id = ?", postID).Delete(&PostSkill{})
	r.db.WithContext(ctx).Where("post_id = ?", postID).Delete(&PostCategory{})
	r.db.WithContext(ctx).Where("post_id = ?", postID).Delete(&PostTag{})

	// Delete the post
	if err := r.db.WithContext(ctx).Delete(&post).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) ListPosts(ctx context.Context, filters PostFilters) ([]Post, error) {
	var posts []Post
	query := r.db.WithContext(ctx).Model(&Post{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.Featured != nil {
		query = query.Where("featured = ?", *filters.Featured)
	}

	if filters.Search != "" {
		searchPattern := "%" + filters.Search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ? OR excerpt ILIKE ?", searchPattern, searchPattern, searchPattern)
	}

	if filters.CategoryID != nil {
		query = query.Joins("JOIN post_categories ON posts.id = post_categories.post_id").
			Where("post_categories.category_id = ?", *filters.CategoryID)
	}

	if filters.TagID != nil {
		query = query.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
			Where("post_tags.tag_id = ?", *filters.TagID)
	}

	if filters.SkillID != nil {
		query = query.Joins("JOIN post_skills ON posts.id = post_skills.post_id").
			Where("post_skills.skill_id = ?", *filters.SkillID)
	}

	// Ordering
	orderBy := normalizeOrderBy(filters.OrderBy)
	if orderBy == "" {
		orderBy = "created_at"
	}
	order := filters.Order
	if order == "" {
		order = "desc"
	}
	// Validate order direction
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	query = query.Order(orderBy + " " + order)

	// Pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	if err := query.Find(&posts).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return posts, nil
}

func (r *gormRepository) IsPostSlugTaken(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&Post{}).Where("slug = ?", slug)
	if excludeID != uuid.Nil {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return count > 0, nil
}

func (r *gormRepository) IncrementPostViews(ctx context.Context, postID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&Post{}).
		Where("id = ?", postID).
		UpdateColumn("views_count", gorm.Expr("views_count + 1")).Error
}

// Category operations

func (r *gormRepository) CreateCategory(ctx context.Context, category *Category) error {
	if err := category.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(category).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeDuplicateSlug, ErrCategorySlugTaken)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateCategory(ctx context.Context, category *Category) error {
	if err := category.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(category).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeDuplicateSlug, ErrCategorySlugTaken)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetCategory(ctx context.Context, categoryID uuid.UUID) (*Category, error) {
	var category Category
	err := r.db.WithContext(ctx).Where("id = ?", categoryID).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeCategoryNotFound, ErrCategoryNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &category, nil
}

func (r *gormRepository) GetCategoryBySlug(ctx context.Context, slug string) (*Category, error) {
	var category Category
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeCategoryNotFound, ErrCategoryNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &category, nil
}

func (r *gormRepository) ListCategories(ctx context.Context) ([]Category, error) {
	var categories []Category
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&categories).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return categories, nil
}

func (r *gormRepository) IsCategorySlugTaken(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&Category{}).Where("slug = ?", slug)
	if excludeID != uuid.Nil {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return count > 0, nil
}

// Tag operations

func (r *gormRepository) CreateTag(ctx context.Context, tag *Tag) error {
	if err := tag.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(tag).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeDuplicateSlug, ErrTagSlugTaken)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) GetTag(ctx context.Context, tagID uuid.UUID) (*Tag, error) {
	var tag Tag
	err := r.db.WithContext(ctx).Where("id = ?", tagID).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeTagNotFound, ErrTagNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &tag, nil
}

func (r *gormRepository) GetTagBySlug(ctx context.Context, slug string) (*Tag, error) {
	var tag Tag
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeTagNotFound, ErrTagNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &tag, nil
}

func (r *gormRepository) GetOrCreateTag(ctx context.Context, name string) (*Tag, error) {
	// Try to find existing tag by name
	var tag Tag
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error
	if err == nil {
		return &tag, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	// Create new tag
	newTag, err := NewTag(name)
	if err != nil {
		return nil, err
	}

	if err := r.CreateTag(ctx, newTag); err != nil {
		return nil, err
	}

	return newTag, nil
}

func (r *gormRepository) ListTags(ctx context.Context) ([]Tag, error) {
	var tags []Tag
	if err := r.db.WithContext(ctx).Order("name ASC").Find(&tags).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return tags, nil
}

func (r *gormRepository) IsTagSlugTaken(ctx context.Context, slug string, excludeID uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&Tag{}).Where("slug = ?", slug)
	if excludeID != uuid.Nil {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return count > 0, nil
}

// Post-Skill relationship operations

func (r *gormRepository) AttachSkillToPost(ctx context.Context, postID, skillID uuid.UUID) error {
	postSkill := PostSkill{
		PostID:    postID,
		SkillID:   skillID,
		CreatedAt: time.Now().UTC(),
	}

	if err := r.db.WithContext(ctx).Create(&postSkill).Error; err != nil {
		if isUniqueConstraintError(err) {
			return nil // Already attached, no error
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) DetachSkillFromPost(ctx context.Context, postID, skillID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("post_id = ? AND skill_id = ?", postID, skillID).
		Delete(&PostSkill{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetPostSkills(ctx context.Context, postID uuid.UUID) ([]uuid.UUID, error) {
	var skillIDs []uuid.UUID
	if err := r.db.WithContext(ctx).Model(&PostSkill{}).
		Where("post_id = ?", postID).
		Pluck("skill_id", &skillIDs).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return skillIDs, nil
}

func (r *gormRepository) PostHasSkill(ctx context.Context, postID, skillID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&PostSkill{}).
		Where("post_id = ? AND skill_id = ?", postID, skillID).
		Count(&count).Error; err != nil {
		return false, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return count > 0, nil
}

// Post-Category relationship operations

func (r *gormRepository) AttachCategoryToPost(ctx context.Context, postID, categoryID uuid.UUID) error {
	postCategory := PostCategory{
		PostID:     postID,
		CategoryID: categoryID,
		CreatedAt:  time.Now().UTC(),
	}

	if err := r.db.WithContext(ctx).Create(&postCategory).Error; err != nil {
		if isUniqueConstraintError(err) {
			return nil // Already attached, no error
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) DetachCategoryFromPost(ctx context.Context, postID, categoryID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("post_id = ? AND category_id = ?", postID, categoryID).
		Delete(&PostCategory{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetPostCategories(ctx context.Context, postID uuid.UUID) ([]Category, error) {
	var categories []Category
	if err := r.db.WithContext(ctx).
		Joins("JOIN post_categories ON categories.id = post_categories.category_id").
		Where("post_categories.post_id = ?", postID).
		Find(&categories).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return categories, nil
}

func (r *gormRepository) PostHasCategory(ctx context.Context, postID, categoryID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&PostCategory{}).
		Where("post_id = ? AND category_id = ?", postID, categoryID).
		Count(&count).Error; err != nil {
		return false, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return count > 0, nil
}

// Post-Tag relationship operations

func (r *gormRepository) AttachTagToPost(ctx context.Context, postID, tagID uuid.UUID) error {
	postTag := PostTag{
		PostID:    postID,
		TagID:     tagID,
		CreatedAt: time.Now().UTC(),
	}

	if err := r.db.WithContext(ctx).Create(&postTag).Error; err != nil {
		if isUniqueConstraintError(err) {
			return nil // Already attached, no error
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) DetachTagFromPost(ctx context.Context, postID, tagID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("post_id = ? AND tag_id = ?", postID, tagID).
		Delete(&PostTag{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetPostTags(ctx context.Context, postID uuid.UUID) ([]Tag, error) {
	var tags []Tag
	if err := r.db.WithContext(ctx).
		Joins("JOIN post_tags ON tags.id = post_tags.tag_id").
		Where("post_tags.post_id = ?", postID).
		Find(&tags).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return tags, nil
}

func (r *gormRepository) PostHasTag(ctx context.Context, postID, tagID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&PostTag{}).
		Where("post_id = ? AND tag_id = ?", postID, tagID).
		Count(&count).Error; err != nil {
		return false, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return count > 0, nil
}

// normalizeOrderBy converts camelCase orderBy values to snake_case database column names
// and validates that the column is allowed for ordering
func normalizeOrderBy(orderBy string) string {
	if orderBy == "" {
		return ""
	}

	// Map of allowed camelCase to snake_case conversions
	allowedColumns := map[string]string{
		"createdAt":   "created_at",
		"updatedAt":   "updated_at",
		"publishedAt": "published_at",
		"viewsCount":  "views_count",
		"created_at":  "created_at",
		"updated_at":  "updated_at",
		"published_at": "published_at",
		"views_count":  "views_count",
	}

	// Check if it's already in the map
	if normalized, ok := allowedColumns[orderBy]; ok {
		return normalized
	}

	// Convert camelCase to snake_case
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(orderBy, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ToLower(snake)

	// Validate the converted value is allowed
	if _, ok := allowedColumns[snake]; ok {
		return snake
	}

	// If not in allowed list, return empty string to use default
	return ""
}

