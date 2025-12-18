package assets

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for content assets.
type Repository interface {
	CreateAsset(ctx context.Context, asset *ContentAsset) error
	UpdateAsset(ctx context.Context, asset *ContentAsset) error
	GetAsset(ctx context.Context, assetID uuid.UUID) (*ContentAsset, error)
	GetAssetsByContentPost(ctx context.Context, contentPostID uuid.UUID) ([]ContentAsset, error)
	GetAssetsBySocialPost(ctx context.Context, socialPostID uuid.UUID) ([]ContentAsset, error)
	DeleteAsset(ctx context.Context, assetID uuid.UUID) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateAsset(ctx context.Context, asset *ContentAsset) error {
	if err := asset.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(asset).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrContentAssetExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateAsset(ctx context.Context, asset *ContentAsset) error {
	if err := asset.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(asset).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetAsset(ctx context.Context, assetID uuid.UUID) (*ContentAsset, error) {
	var asset ContentAsset
	err := r.db.WithContext(ctx).Where("id = ?", assetID).First(&asset).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrContentAssetNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &asset, nil
}

func (r *gormRepository) GetAssetsByContentPost(ctx context.Context, contentPostID uuid.UUID) ([]ContentAsset, error) {
	var assets []ContentAsset
	err := r.db.WithContext(ctx).
		Where("content_post_id = ?", contentPostID).
		Order("created_at desc").
		Find(&assets).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return assets, nil
}

func (r *gormRepository) GetAssetsBySocialPost(ctx context.Context, socialPostID uuid.UUID) ([]ContentAsset, error) {
	var assets []ContentAsset
	err := r.db.WithContext(ctx).
		Where("social_post_id = ?", socialPostID).
		Order("created_at desc").
		Find(&assets).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return assets, nil
}

func (r *gormRepository) DeleteAsset(ctx context.Context, assetID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&ContentAsset{}, "id = ?", assetID).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
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
