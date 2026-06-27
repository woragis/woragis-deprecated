package interviewstages

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for interview stages.
type Repository interface {
	CreateStage(ctx context.Context, stage *InterviewStage) error
	UpdateStage(ctx context.Context, stage *InterviewStage) error
	GetStage(ctx context.Context, stageID uuid.UUID) (*InterviewStage, error)
	ListStages(ctx context.Context, filters StageFilters) ([]InterviewStage, error)
	DeleteStage(ctx context.Context, stageID uuid.UUID) error
	GetStagesByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]InterviewStage, error)
}

// StageFilters represents filtering options for listing interview stages.
type StageFilters struct {
	JobApplicationID *uuid.UUID
	StageType        *StageType
	Outcome          *StageOutcome
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

func (r *gormRepository) CreateStage(ctx context.Context, stage *InterviewStage) error {
	if err := stage.Validate(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(stage).Error
}

func (r *gormRepository) UpdateStage(ctx context.Context, stage *InterviewStage) error {
	if err := stage.Validate(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(stage).Error
}

func (r *gormRepository) GetStage(ctx context.Context, stageID uuid.UUID) (*InterviewStage, error) {
	var stage InterviewStage
	if err := r.db.WithContext(ctx).Where("id = ?", stageID).First(&stage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrStageNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &stage, nil
}

func (r *gormRepository) ListStages(ctx context.Context, filters StageFilters) ([]InterviewStage, error) {
	var stages []InterviewStage
	query := r.db.WithContext(ctx).Model(&InterviewStage{})

	if filters.JobApplicationID != nil {
		query = query.Where("job_application_id = ?", *filters.JobApplicationID)
	}
	if filters.StageType != nil {
		query = query.Where("stage_type = ?", *filters.StageType)
	}
	if filters.Outcome != nil {
		query = query.Where("outcome = ?", *filters.Outcome)
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	query = query.Order("scheduled_date ASC NULLS LAST, created_at ASC")

	if err := query.Find(&stages).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return stages, nil
}

func (r *gormRepository) DeleteStage(ctx context.Context, stageID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&InterviewStage{}, stageID)
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrStageNotFound)
	}
	return nil
}

func (r *gormRepository) GetStagesByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]InterviewStage, error) {
	return r.ListStages(ctx, StageFilters{
		JobApplicationID: &applicationID,
	})
}

