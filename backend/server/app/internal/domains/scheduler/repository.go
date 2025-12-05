package scheduler

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository provides persistence operations for schedules.
type Repository interface {
	Create(ctx context.Context, schedule *Schedule) error
	Update(ctx context.Context, schedule *Schedule) error
	Delete(ctx context.Context, id, userID uuid.UUID) error
	Get(ctx context.Context, id, userID uuid.UUID) (*Schedule, error)
	List(ctx context.Context, userID uuid.UUID) ([]Schedule, error)
	ListDue(ctx context.Context, now time.Time) ([]Schedule, error)
	BulkUpdateState(ctx context.Context, userID uuid.UUID, ids []uuid.UUID, updates map[string]any) error
	InsertRun(ctx context.Context, run *ExecutionRun) error
	UpdateRun(ctx context.Context, run *ExecutionRun) error
	ListRuns(ctx context.Context, userID uuid.UUID, filters RunFilters) ([]ExecutionRun, error)
}

// RunFilters filters execution history.
type RunFilters struct {
	ScheduleID uuid.UUID
	Status     string
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

func (r *gormRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	// Verify schedule exists and belongs to user
	if _, err := r.Get(ctx, id, userID); err != nil {
		return err
	}

	// Delete execution runs first
	if err := r.db.WithContext(ctx).Where("schedule_id = ?", id).Delete(&ExecutionRun{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
	}

	// Delete the schedule
	if err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&Schedule{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
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
		Where("active = ? AND paused = ? AND next_run <= ?", true, false, now).
		Find(&schedules).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return schedules, nil
}

func (r *gormRepository) BulkUpdateState(ctx context.Context, userID uuid.UUID, ids []uuid.UUID, updates map[string]any) error {
	if len(ids) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).
		Model(&Schedule{}).
		Where("user_id = ? AND id IN ?", userID, ids).
		Updates(updates).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) InsertRun(ctx context.Context, run *ExecutionRun) error {
	if err := r.db.WithContext(ctx).Create(run).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateRun(ctx context.Context, run *ExecutionRun) error {
	if err := r.db.WithContext(ctx).Save(run).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) ListRuns(ctx context.Context, userID uuid.UUID, filters RunFilters) ([]ExecutionRun, error) {
	var runs []ExecutionRun

	db := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if filters.ScheduleID != uuid.Nil {
		db = db.Where("schedule_id = ?", filters.ScheduleID)
	}
	if strings.TrimSpace(filters.Status) != "" {
		db = db.Where("LOWER(status) = ?", strings.ToLower(filters.Status))
	}
	if filters.Limit > 0 {
		db = db.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		db = db.Offset(filters.Offset)
	}

	if err := db.Order("created_at DESC").Find(&runs).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return runs, nil
}
