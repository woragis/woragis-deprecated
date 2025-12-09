package userprofiles

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines data access operations for user profiles.
type Repository interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
	CreateProfile(ctx context.Context, profile *UserProfile) error
	UpdateProfile(ctx context.Context, profile *UserProfile) error
	UpsertProfile(ctx context.Context, profile *UserProfile) error
}

// repository implements Repository using GORM.
type repository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM-based repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// GetProfile retrieves a user profile by user ID.
func (r *repository) GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	var profile UserProfile
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrProfileNotFound)
		}
		return nil, err
	}
	return &profile, nil
}

// CreateProfile creates a new user profile.
func (r *repository) CreateProfile(ctx context.Context, profile *UserProfile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

// UpdateProfile updates an existing user profile.
func (r *repository) UpdateProfile(ctx context.Context, profile *UserProfile) error {
	return r.db.WithContext(ctx).Save(profile).Error
}

// UpsertProfile creates or updates a user profile.
func (r *repository) UpsertProfile(ctx context.Context, profile *UserProfile) error {
	var existing UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", profile.UserID).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new profile
		return r.db.WithContext(ctx).Create(profile).Error
	} else if err != nil {
		return err
	}
	
	// Update existing profile
	profile.ID = existing.ID
	profile.CreatedAt = existing.CreatedAt
	return r.db.WithContext(ctx).Save(profile).Error
}

