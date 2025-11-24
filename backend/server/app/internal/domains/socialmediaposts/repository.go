package socialmediaposts

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for social media posts and entity links.
type Repository interface {
	// Social Media Post operations
	CreatePost(ctx context.Context, post *SocialMediaPost) error
	UpdatePost(ctx context.Context, post *SocialMediaPost) error
	GetPost(ctx context.Context, postID uuid.UUID) (*SocialMediaPost, error)
	GetPostByURL(ctx context.Context, url string) (*SocialMediaPost, error)
	ListPosts(ctx context.Context, filters PostFilters) ([]SocialMediaPost, error)
	DeletePost(ctx context.Context, postID uuid.UUID) error

	// Entity Link operations
	CreateLink(ctx context.Context, link *SocialMediaEntityLink) error
	UpdateLink(ctx context.Context, link *SocialMediaEntityLink) error
	GetLink(ctx context.Context, linkID uuid.UUID) (*SocialMediaEntityLink, error)
	DeleteLink(ctx context.Context, linkID uuid.UUID) error
	GetLinksByPost(ctx context.Context, postID uuid.UUID) ([]SocialMediaEntityLink, error)
	GetLinksByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]SocialMediaEntityLink, error)
	GetPostsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID, relationshipType *RelationshipType) ([]SocialMediaPostWithLink, error)
	GetEntitiesByPost(ctx context.Context, postID uuid.UUID) ([]EntityLinkInfo, error)
}

// PostFilters for querying posts.
type PostFilters struct {
	Platform *Platform
	Status   *PostStatus
}

// SocialMediaPostWithLink represents a post with its relationship to an entity.
type SocialMediaPostWithLink struct {
	SocialMediaPost
	RelationshipType RelationshipType `json:"relationshipType"`
	LinkID           uuid.UUID        `json:"linkId"`
}

// EntityLinkInfo represents entity information from a link.
type EntityLinkInfo struct {
	EntityType       EntityType       `json:"entityType"`
	EntityID         uuid.UUID        `json:"entityId"`
	RelationshipType RelationshipType `json:"relationshipType"`
	LinkID           uuid.UUID        `json:"linkId"`
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// Social Media Post operations

func (r *gormRepository) CreatePost(ctx context.Context, post *SocialMediaPost) error {
	if err := post.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrPostAlreadyExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdatePost(ctx context.Context, post *SocialMediaPost) error {
	if err := post.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(post).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrPostAlreadyExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetPost(ctx context.Context, postID uuid.UUID) (*SocialMediaPost, error) {
	var post SocialMediaPost
	err := r.db.WithContext(ctx).Where("id = ?", postID).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrPostNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &post, nil
}

func (r *gormRepository) GetPostByURL(ctx context.Context, url string) (*SocialMediaPost, error) {
	var post SocialMediaPost
	err := r.db.WithContext(ctx).Where("url = ?", url).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrPostNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &post, nil
}

func (r *gormRepository) ListPosts(ctx context.Context, filters PostFilters) ([]SocialMediaPost, error) {
	var posts []SocialMediaPost
	query := r.db.WithContext(ctx)

	if filters.Platform != nil {
		query = query.Where("platform = ?", *filters.Platform)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if err := query.Order("published_date desc, created_at desc").Find(&posts).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return posts, nil
}

func (r *gormRepository) DeletePost(ctx context.Context, postID uuid.UUID) error {
	// Soft delete by setting status to deleted
	post, err := r.GetPost(ctx, postID)
	if err != nil {
		return err
	}
	post.UpdateStatus(PostStatusDeleted)
	return r.UpdatePost(ctx, post)
}

// Entity Link operations

func (r *gormRepository) CreateLink(ctx context.Context, link *SocialMediaEntityLink) error {
	if err := link.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(link).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrLinkAlreadyExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateLink(ctx context.Context, link *SocialMediaEntityLink) error {
	if err := link.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(link).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetLink(ctx context.Context, linkID uuid.UUID) (*SocialMediaEntityLink, error) {
	var link SocialMediaEntityLink
	err := r.db.WithContext(ctx).Where("id = ?", linkID).First(&link).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrLinkNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &link, nil
}

func (r *gormRepository) DeleteLink(ctx context.Context, linkID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&SocialMediaEntityLink{}, "id = ?", linkID).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetLinksByPost(ctx context.Context, postID uuid.UUID) ([]SocialMediaEntityLink, error) {
	var links []SocialMediaEntityLink
	if err := r.db.WithContext(ctx).
		Where("social_media_post_id = ?", postID).
		Order("created_at desc").
		Find(&links).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return links, nil
}

func (r *gormRepository) GetLinksByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]SocialMediaEntityLink, error) {
	var links []SocialMediaEntityLink
	if err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at desc").
		Find(&links).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return links, nil
}

func (r *gormRepository) GetPostsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID, relationshipType *RelationshipType) ([]SocialMediaPostWithLink, error) {
	var results []SocialMediaPostWithLink
	query := r.db.WithContext(ctx).
		Table("social_media_posts").
		Select("social_media_posts.*, social_media_entity_links.relationship_type, social_media_entity_links.id as link_id").
		Joins("INNER JOIN social_media_entity_links ON social_media_entity_links.social_media_post_id = social_media_posts.id").
		Where("social_media_entity_links.entity_type = ? AND social_media_entity_links.entity_id = ?", entityType, entityID)

	if relationshipType != nil {
		query = query.Where("social_media_entity_links.relationship_type = ?", *relationshipType)
	}

	if err := query.Order("social_media_posts.published_date desc, social_media_posts.created_at desc").Scan(&results).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return results, nil
}

func (r *gormRepository) GetEntitiesByPost(ctx context.Context, postID uuid.UUID) ([]EntityLinkInfo, error) {
	var entities []EntityLinkInfo
	if err := r.db.WithContext(ctx).
		Table("social_media_entity_links").
		Select("entity_type, entity_id, relationship_type, id as link_id").
		Where("social_media_post_id = ?", postID).
		Order("created_at desc").
		Scan(&entities).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return entities, nil
}

// isUniqueConstraintError checks if the error is a unique constraint violation.
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// Check for PostgreSQL unique constraint violation (code 23505)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}

	// Check for SQLite unique constraint violation
	if err.Error() == "UNIQUE constraint failed" {
		return true
	}

	return false
}

