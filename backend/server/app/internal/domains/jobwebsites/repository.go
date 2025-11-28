package jobwebsites

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for job websites.
type Repository interface {
	CreateJobWebsite(ctx context.Context, website *JobWebsite) error
	UpdateJobWebsite(ctx context.Context, website *JobWebsite) error
	GetJobWebsite(ctx context.Context, websiteID uuid.UUID) (*JobWebsite, error)
	GetJobWebsiteByName(ctx context.Context, name string) (*JobWebsite, error)
	ListJobWebsites(ctx context.Context, enabledOnly bool) ([]JobWebsite, error)
	DeleteJobWebsite(ctx context.Context, websiteID uuid.UUID) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateJobWebsite(ctx context.Context, website *JobWebsite) error {
	if err := website.Validate(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(website).Error
}

func (r *gormRepository) UpdateJobWebsite(ctx context.Context, website *JobWebsite) error {
	if err := website.Validate(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(website).Error
}

func (r *gormRepository) GetJobWebsite(ctx context.Context, websiteID uuid.UUID) (*JobWebsite, error) {
	var website JobWebsite
	if err := r.db.WithContext(ctx).Where("id = ?", websiteID).First(&website).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrWebsiteNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &website, nil
}

func (r *gormRepository) GetJobWebsiteByName(ctx context.Context, name string) (*JobWebsite, error) {
	var website JobWebsite
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&website).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrWebsiteNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &website, nil
}

func (r *gormRepository) ListJobWebsites(ctx context.Context, enabledOnly bool) ([]JobWebsite, error) {
	var websites []JobWebsite
	query := r.db.WithContext(ctx).Model(&JobWebsite{})

	if enabledOnly {
		query = query.Where("enabled = ?", true)
	}

	query = query.Order("name ASC")

	if err := query.Find(&websites).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return websites, nil
}

func (r *gormRepository) DeleteJobWebsite(ctx context.Context, websiteID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&JobWebsite{}, websiteID)
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrWebsiteNotFound)
	}
	return nil
}

