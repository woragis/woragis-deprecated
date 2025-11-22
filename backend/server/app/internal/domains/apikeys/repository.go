package apikeys

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for API keys.
type Repository interface {
	CreateAPIKey(ctx context.Context, apiKey *APIKey) error
	UpdateAPIKey(ctx context.Context, apiKey *APIKey) error
	DeleteAPIKey(ctx context.Context, apiKeyID uuid.UUID, userID uuid.UUID) error
	GetAPIKey(ctx context.Context, apiKeyID uuid.UUID, userID uuid.UUID) (*APIKey, error)
	GetAPIKeyByHash(ctx context.Context, keyHash string) (*APIKey, error)
	ListAPIKeys(ctx context.Context, userID uuid.UUID) ([]APIKey, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateAPIKey(ctx context.Context, apiKey *APIKey) error {
	if err := apiKey.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(apiKey).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateAPIKey(ctx context.Context, apiKey *APIKey) error {
	if err := apiKey.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(apiKey).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) DeleteAPIKey(ctx context.Context, apiKeyID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var apiKey APIKey
		if err := tx.Where("id = ? AND user_id = ?", apiKeyID, userID).First(&apiKey).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return NewDomainError(ErrCodeNotFound, ErrAPIKeyNotFound)
			}
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
		}

		if err := tx.Delete(&apiKey).Error; err != nil {
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}
		return nil
	})
}

func (r *gormRepository) GetAPIKey(ctx context.Context, apiKeyID uuid.UUID, userID uuid.UUID) (*APIKey, error) {
	var apiKey APIKey
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", apiKeyID, userID).First(&apiKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrAPIKeyNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &apiKey, nil
}

func (r *gormRepository) GetAPIKeyByHash(ctx context.Context, keyHash string) (*APIKey, error) {
	var apiKey APIKey
	err := r.db.WithContext(ctx).Where("key_hash = ?", keyHash).First(&apiKey).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrAPIKeyNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &apiKey, nil
}

func (r *gormRepository) ListAPIKeys(ctx context.Context, userID uuid.UUID) ([]APIKey, error) {
	var apiKeys []APIKey
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&apiKeys).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return apiKeys, nil
}

