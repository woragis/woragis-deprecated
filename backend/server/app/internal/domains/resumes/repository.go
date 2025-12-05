package resumes

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for resumes.
type Repository interface {
	CreateResume(ctx context.Context, resume *Resume) error
	UpdateResume(ctx context.Context, resume *Resume) error
	DeleteResume(ctx context.Context, resumeID uuid.UUID, userID uuid.UUID) error
	GetResume(ctx context.Context, resumeID uuid.UUID, userID uuid.UUID) (*Resume, error)
	ListResumes(ctx context.Context, userID uuid.UUID) ([]Resume, error)
	GetMainResume(ctx context.Context, userID uuid.UUID) (*Resume, error)
	GetFeaturedResume(ctx context.Context, userID uuid.UUID) (*Resume, error)
	UnmarkAllAsMain(ctx context.Context, userID uuid.UUID) error
}

// gormRepository implements Repository using GORM.
type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM-based repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// CreateResume creates a new resume.
func (r *gormRepository) CreateResume(ctx context.Context, resume *Resume) error {
	return r.db.WithContext(ctx).Create(resume).Error
}

// UpdateResume updates an existing resume.
func (r *gormRepository) UpdateResume(ctx context.Context, resume *Resume) error {
	return r.db.WithContext(ctx).Save(resume).Error
}

// DeleteResume deletes a resume.
func (r *gormRepository) DeleteResume(ctx context.Context, resumeID uuid.UUID, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", resumeID, userID).
		Delete(&Resume{})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrResumeNotFound)
	}
	
	return nil
}

// GetResume retrieves a resume by ID.
func (r *gormRepository) GetResume(ctx context.Context, resumeID uuid.UUID, userID uuid.UUID) (*Resume, error) {
	var resume Resume
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", resumeID, userID).
		First(&resume).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrResumeNotFound)
		}
		return nil, err
	}
	
	return &resume, nil
}

// ListResumes lists all resumes for a user.
func (r *gormRepository) ListResumes(ctx context.Context, userID uuid.UUID) ([]Resume, error) {
	var resumes []Resume
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_main DESC, is_featured DESC, created_at DESC").
		Find(&resumes).Error
	
	return resumes, err
}

// GetMainResume retrieves the main resume for a user.
func (r *gormRepository) GetMainResume(ctx context.Context, userID uuid.UUID) (*Resume, error) {
	var resume Resume
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_main = ?", userID, true).
		First(&resume).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrNoMainResume)
		}
		return nil, err
	}
	
	return &resume, nil
}

// GetFeaturedResume retrieves a featured resume for a user (fallback if no main).
func (r *gormRepository) GetFeaturedResume(ctx context.Context, userID uuid.UUID) (*Resume, error) {
	var resume Resume
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_featured = ?", userID, true).
		Order("created_at DESC").
		First(&resume).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrResumeNotFound)
		}
		return nil, err
	}
	
	return &resume, nil
}

// UnmarkAllAsMain unmarks all resumes as main for a user.
func (r *gormRepository) UnmarkAllAsMain(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&Resume{}).
		Where("user_id = ?", userID).
		Update("is_main", false).Error
}

