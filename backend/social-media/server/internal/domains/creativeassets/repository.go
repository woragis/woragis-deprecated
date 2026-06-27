package creativeassets

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for creative assets
type Repository interface {
	CreateAsset(ctx context.Context, asset *CreativeAsset) error
	UpdateAsset(ctx context.Context, asset *CreativeAsset) error
	GetAsset(ctx context.Context, assetID uuid.UUID) (*CreativeAsset, error)
	GetAssetsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]CreativeAsset, error)
	GetAssetByEntityAndPurpose(ctx context.Context, entityType EntityType, entityID uuid.UUID, purpose AssetPurpose) (*CreativeAsset, error)
	DeleteAsset(ctx context.Context, assetID uuid.UUID, userID uuid.UUID) error
	DeleteAssetsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID, userID uuid.UUID) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewRepository creates a new GORM-based repository
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateAsset(ctx context.Context, asset *CreativeAsset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *gormRepository) UpdateAsset(ctx context.Context, asset *CreativeAsset) error {
	return r.db.WithContext(ctx).Save(asset).Error
}

func (r *gormRepository) GetAsset(ctx context.Context, assetID uuid.UUID) (*CreativeAsset, error) {
	var asset CreativeAsset
	if err := r.db.WithContext(ctx).First(&asset, "id = ?", assetID).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *gormRepository) GetAssetsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]CreativeAsset, error) {
	var assets []CreativeAsset
	err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Find(&assets).Error
	return assets, err
}

func (r *gormRepository) GetAssetByEntityAndPurpose(ctx context.Context, entityType EntityType, entityID uuid.UUID, purpose AssetPurpose) (*CreativeAsset, error) {
	var asset CreativeAsset
	err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ? AND purpose = ?", entityType, entityID, purpose).
		First(&asset).Error
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *gormRepository) DeleteAsset(ctx context.Context, assetID uuid.UUID, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", assetID, userID).
		Delete(&CreativeAsset{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrAssetNotFound)
	}
	return nil
}

func (r *gormRepository) DeleteAssetsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ? AND user_id = ?", entityType, entityID, userID).
		Delete(&CreativeAsset{}).Error
}

