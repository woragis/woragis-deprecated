package scheduler

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository provides persistence operations for schedules.
type Repository interface {
	Create(ctx context.Context, schedule *Schedule) error
	Update(ctx context.Context, schedule *Schedule) error
	Get(ctx context.Context, id, userID uuid.UUID) (*Schedule, error)
	List(ctx context.Context, userID uuid.UUID) ([]Schedule, error)
	ListDue(ctx context.Context, now time.Time) ([]Schedule, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) Create(ctx context.Context, schedule *Schedule) error {
	if err := schedule.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(schedule).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) Update(ctx context.Context, schedule *Schedule) error {
	if err := schedule.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(schedule).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	return nil
}

func (r *gormRepository) Get(ctx context.Context, id, userID uuid.UUID) (*Schedule, error) {
	var schedule Schedule
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&schedule).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrScheduleNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &schedule, nil
}

func (r *gormRepository) List(ctx context.Context, userID uuid.UUID) ([]Schedule, error) {
	var schedules []Schedule
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("next_run asc").
		Find(&schedules).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return schedules, nil
}

func (r *gormRepository) ListDue(ctx context.Context, now time.Time) ([]Schedule, error) {
	var schedules []Schedule
	if err := r.db.WithContext(ctx).
		Where("active = ? AND next_run <= ?", true, now).
		Find(&schedules).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return schedules, nil
}
