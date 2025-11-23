package translations

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for translations.
type Repository interface {
	CreateTranslation(ctx context.Context, translation *Translation) error
	UpdateTranslation(ctx context.Context, translation *Translation) error
	GetTranslation(ctx context.Context, translationID uuid.UUID) (*Translation, error)
	GetTranslationByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID, language Language) (*Translation, error)
	ListTranslations(ctx context.Context, filters TranslationFilters) ([]Translation, error)
	DeleteTranslation(ctx context.Context, translationID uuid.UUID) error
}

// TranslationFilters represents filtering options for listing translations.
type TranslationFilters struct {
	EntityType *EntityType
	EntityID   *uuid.UUID
	Language   *Language
	Status     *TranslationStatus
	Limit      int
	Offset     int
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateTranslation(ctx context.Context, translation *Translation) error {
	if err := translation.Validate(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(translation).Error
}

func (r *gormRepository) UpdateTranslation(ctx context.Context, translation *Translation) error {
	if err := translation.Validate(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(translation).Error
}

func (r *gormRepository) GetTranslation(ctx context.Context, translationID uuid.UUID) (*Translation, error) {
	var translation Translation
	if err := r.db.WithContext(ctx).Where("id = ?", translationID).First(&translation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrTranslationNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &translation, nil
}

func (r *gormRepository) GetTranslationByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID, language Language) (*Translation, error) {
	var translation Translation
	if err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ? AND language = ?", entityType, entityID, language).
		First(&translation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrTranslationNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &translation, nil
}

func (r *gormRepository) ListTranslations(ctx context.Context, filters TranslationFilters) ([]Translation, error) {
	var translations []Translation
	query := r.db.WithContext(ctx).Model(&Translation{})

	if filters.EntityType != nil {
		query = query.Where("entity_type = ?", *filters.EntityType)
	}
	if filters.EntityID != nil {
		query = query.Where("entity_id = ?", *filters.EntityID)
	}
	if filters.Language != nil {
		query = query.Where("language = ?", *filters.Language)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	if err := query.Find(&translations).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return translations, nil
}

func (r *gormRepository) DeleteTranslation(ctx context.Context, translationID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&Translation{}, translationID)
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrTranslationNotFound)
	}
	return nil
}

