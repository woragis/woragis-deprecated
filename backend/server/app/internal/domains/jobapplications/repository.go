package jobapplications

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for job applications.
type Repository interface {
	CreateJobApplication(ctx context.Context, application *JobApplication) error
	UpdateJobApplication(ctx context.Context, application *JobApplication) error
	GetJobApplication(ctx context.Context, applicationID uuid.UUID) (*JobApplication, error)
	ListJobApplications(ctx context.Context, filters JobApplicationFilters) ([]JobApplication, error)
	DeleteJobApplication(ctx context.Context, applicationID uuid.UUID) error
}

// JobApplicationFilters represents filtering options for listing job applications.
type JobApplicationFilters struct {
	UserID           *uuid.UUID
	Website          *string
	Status           *ApplicationStatus
	ResumeID         *uuid.UUID
	InterestLevel    *string
	Source           *string
	ApplicationMethod *string
	Language         *string
	Limit            int
	Offset           int
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateJobApplication(ctx context.Context, application *JobApplication) error {
	if err := application.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(application).Error; err != nil {
		return handleDatabaseError(err)
	}
	return nil
}

func (r *gormRepository) UpdateJobApplication(ctx context.Context, application *JobApplication) error {
	if err := application.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Save(application).Error; err != nil {
		return handleDatabaseError(err)
	}
	return nil
}

func (r *gormRepository) GetJobApplication(ctx context.Context, applicationID uuid.UUID) (*JobApplication, error) {
	var application JobApplication
	if err := r.db.WithContext(ctx).Where("id = ?", applicationID).First(&application).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrApplicationNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &application, nil
}

func (r *gormRepository) ListJobApplications(ctx context.Context, filters JobApplicationFilters) ([]JobApplication, error) {
	var applications []JobApplication
	query := r.db.WithContext(ctx).Model(&JobApplication{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}
	if filters.Website != nil {
		query = query.Where("website = ?", *filters.Website)
	}
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	if filters.ResumeID != nil {
		query = query.Where("resume_id = ?", *filters.ResumeID)
	}
	if filters.InterestLevel != nil {
		query = query.Where("interest_level = ?", *filters.InterestLevel)
	}
	if filters.Source != nil {
		query = query.Where("source = ?", *filters.Source)
	}
	if filters.ApplicationMethod != nil {
		query = query.Where("application_method = ?", *filters.ApplicationMethod)
	}
	if filters.Language != nil {
		query = query.Where("language = ?", *filters.Language)
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	query = query.Order("created_at DESC")

	if err := query.Find(&applications).Error; err != nil {
		return nil, handleDatabaseError(err)
	}

	return applications, nil
}

func (r *gormRepository) DeleteJobApplication(ctx context.Context, applicationID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&JobApplication{}, applicationID)
	if result.Error != nil {
		return handleDatabaseError(result.Error)
	}
	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrApplicationNotFound)
	}
	return nil
}

