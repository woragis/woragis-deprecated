package content

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for content posts and repurposing.
type Repository interface {
	// Content Post operations
	CreateContentPost(ctx context.Context, post *ContentPost) error
	UpdateContentPost(ctx context.Context, post *ContentPost) error
	GetContentPost(ctx context.Context, contentPostID uuid.UUID) (*ContentPost, error)
	GetContentPostByPostID(ctx context.Context, postID uuid.UUID) (*ContentPost, error)
	ListContentPosts(ctx context.Context, filters ContentPostFilters) ([]ContentPost, error)
	GetContentBacklog(ctx context.Context) ([]ContentPost, error)
	DeleteContentPost(ctx context.Context, contentPostID uuid.UUID) error

	// Repurposing operations
	CreateRepurposing(ctx context.Context, repurposing *ContentRepurposing) error
	GetRepurposingsByContentPost(ctx context.Context, contentPostID uuid.UUID) ([]ContentRepurposing, error)
	GetRepurposingsBySocialPost(ctx context.Context, socialPostID uuid.UUID) ([]ContentRepurposing, error)
}

// ContentPostFilters for querying content posts.
type ContentPostFilters struct {
	Status   *ContentPostStatus
	Priority *ContentPostPriority
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// Content Post operations

func (r *gormRepository) CreateContentPost(ctx context.Context, post *ContentPost) error {
	if err := post.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrContentPostExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateContentPost(ctx context.Context, post *ContentPost) error {
	if err := post.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(post).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetContentPost(ctx context.Context, contentPostID uuid.UUID) (*ContentPost, error) {
	var post ContentPost
	err := r.db.WithContext(ctx).Where("id = ?", contentPostID).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrContentPostNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &post, nil
}

func (r *gormRepository) GetContentPostByPostID(ctx context.Context, postID uuid.UUID) (*ContentPost, error) {
	var post ContentPost
	err := r.db.WithContext(ctx).Where("post_id = ?", postID).First(&post).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrContentPostNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &post, nil
}

func (r *gormRepository) ListContentPosts(ctx context.Context, filters ContentPostFilters) ([]ContentPost, error) {
	var posts []ContentPost
	query := r.db.WithContext(ctx)

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.Priority != nil {
		query = query.Where("priority = ?", *filters.Priority)
	}

	if err := query.Order("priority desc, created_at desc").Find(&posts).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return posts, nil
}

func (r *gormRepository) GetContentBacklog(ctx context.Context) ([]ContentPost, error) {
	var posts []ContentPost
	err := r.db.WithContext(ctx).
		Where("status = ?", ContentPostStatusPending).
		Order("priority desc, created_at asc").
		Find(&posts).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return posts, nil
}

func (r *gormRepository) DeleteContentPost(ctx context.Context, contentPostID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&ContentPost{}, "id = ?", contentPostID).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

// Repurposing operations

func (r *gormRepository) CreateRepurposing(ctx context.Context, repurposing *ContentRepurposing) error {
	if err := repurposing.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(repurposing).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) GetRepurposingsByContentPost(ctx context.Context, contentPostID uuid.UUID) ([]ContentRepurposing, error) {
	var repurposings []ContentRepurposing
	err := r.db.WithContext(ctx).
		Where("content_post_id = ?", contentPostID).
		Order("repurposed_at desc").
		Find(&repurposings).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return repurposings, nil
}

func (r *gormRepository) GetRepurposingsBySocialPost(ctx context.Context, socialPostID uuid.UUID) ([]ContentRepurposing, error) {
	var repurposings []ContentRepurposing
	err := r.db.WithContext(ctx).
		Where("social_post_id = ?", socialPostID).
		Order("repurposed_at desc").
		Find(&repurposings).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return repurposings, nil
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
