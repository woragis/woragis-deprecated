package systemdesigns

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for system designs.
type Repository interface {
	CreateSystemDesign(ctx context.Context, systemDesign *SystemDesign) error
	UpdateSystemDesign(ctx context.Context, systemDesign *SystemDesign) error
	GetSystemDesign(ctx context.Context, systemDesignID uuid.UUID, userID uuid.UUID) (*SystemDesign, error)
	GetSystemDesignPublic(ctx context.Context, systemDesignID uuid.UUID) (*SystemDesign, error)
	ListSystemDesigns(ctx context.Context, userID uuid.UUID) ([]SystemDesign, error)
	ListFeaturedSystemDesigns(ctx context.Context) ([]SystemDesign, error)
	DeleteSystemDesign(ctx context.Context, systemDesignID uuid.UUID, userID uuid.UUID) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateSystemDesign(ctx context.Context, systemDesign *SystemDesign) error {
	if err := systemDesign.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(systemDesign).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateSystemDesign(ctx context.Context, systemDesign *SystemDesign) error {
	if err := systemDesign.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Save(systemDesign).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetSystemDesign(ctx context.Context, systemDesignID uuid.UUID, userID uuid.UUID) (*SystemDesign, error) {
	var systemDesign SystemDesign
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", systemDesignID, userID).
		First(&systemDesign).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrSystemDesignNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &systemDesign, nil
}

func (r *gormRepository) GetSystemDesignPublic(ctx context.Context, systemDesignID uuid.UUID) (*SystemDesign, error) {
	var systemDesign SystemDesign
	err := r.db.WithContext(ctx).
		Where("id = ?", systemDesignID).
		First(&systemDesign).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrSystemDesignNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &systemDesign, nil
}

func (r *gormRepository) ListSystemDesigns(ctx context.Context, userID uuid.UUID) ([]SystemDesign, error) {
	var systemDesigns []SystemDesign
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&systemDesigns).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return systemDesigns, nil
}

func (r *gormRepository) ListFeaturedSystemDesigns(ctx context.Context) ([]SystemDesign, error) {
	var systemDesigns []SystemDesign
	err := r.db.WithContext(ctx).
		Where("featured = ?", true).
		Order("created_at DESC").
		Find(&systemDesigns).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return systemDesigns, nil
}

func (r *gormRepository) DeleteSystemDesign(ctx context.Context, systemDesignID uuid.UUID, userID uuid.UUID) error {
	var systemDesign SystemDesign
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", systemDesignID, userID).
		First(&systemDesign).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return NewDomainError(ErrCodeNotFound, ErrSystemDesignNotFound)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	if err := r.db.WithContext(ctx).Delete(&systemDesign).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

