package userpreferences

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for user preferences.
type Repository interface {
	GetPreferences(ctx context.Context, userID uuid.UUID) (*UserPreferences, error)
	CreatePreferences(ctx context.Context, preferences *UserPreferences) error
	UpdatePreferences(ctx context.Context, preferences *UserPreferences) error
	UpsertPreferences(ctx context.Context, preferences *UserPreferences) error
}

type gormRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewGormRepository instantiates the repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{
		db:     db,
		logger: slog.Default(),
	}
}

func (r *gormRepository) GetPreferences(ctx context.Context, userID uuid.UUID) (*UserPreferences, error) {
	var preferences UserPreferences
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&preferences).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrPreferencesNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &preferences, nil
}

func (r *gormRepository) CreatePreferences(ctx context.Context, preferences *UserPreferences) error {
	if err := preferences.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(preferences).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) UpdatePreferences(ctx context.Context, preferences *UserPreferences) error {
	if err := preferences.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(preferences).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	return nil
}

func (r *gormRepository) UpsertPreferences(ctx context.Context, preferences *UserPreferences) error {
	if err := preferences.Validate(); err != nil {
		return err
	}

	// Try to get existing preferences
	existing, err := r.GetPreferences(ctx, preferences.UserID)
	if err != nil {
		// If not found, create new
		if domainErr, ok := AsDomainError(err); ok && domainErr.Code == ErrCodeNotFound {
			return r.CreatePreferences(ctx, preferences)
		}
		return err
	}

	// Update existing
	preferences.ID = existing.ID
	preferences.CreatedAt = existing.CreatedAt
	return r.UpdatePreferences(ctx, preferences)
}

