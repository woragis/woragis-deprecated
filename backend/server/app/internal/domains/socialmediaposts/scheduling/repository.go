package scheduling

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for scheduled posts.
type Repository interface {
	CreateSchedule(ctx context.Context, schedule *ScheduledPost) error
	UpdateSchedule(ctx context.Context, schedule *ScheduledPost) error
	GetSchedule(ctx context.Context, scheduleID uuid.UUID) (*ScheduledPost, error)
	GetScheduleBySocialPost(ctx context.Context, socialPostID uuid.UUID) (*ScheduledPost, error)
	GetScheduleForDateRange(ctx context.Context, startDate, endDate time.Time) ([]ScheduledPost, error)
	GetUpcomingPosts(ctx context.Context, limit int) ([]ScheduledPost, error)
	CheckConflict(ctx context.Context, scheduledAt time.Time, excludeScheduleID *uuid.UUID) (bool, error)
	DeleteSchedule(ctx context.Context, scheduleID uuid.UUID) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateSchedule(ctx context.Context, schedule *ScheduledPost) error {
	if err := schedule.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(schedule).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrScheduledPostExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateSchedule(ctx context.Context, schedule *ScheduledPost) error {
	if err := schedule.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(schedule).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetSchedule(ctx context.Context, scheduleID uuid.UUID) (*ScheduledPost, error) {
	var schedule ScheduledPost
	err := r.db.WithContext(ctx).Where("id = ?", scheduleID).First(&schedule).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrScheduledPostNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &schedule, nil
}

func (r *gormRepository) GetScheduleBySocialPost(ctx context.Context, socialPostID uuid.UUID) (*ScheduledPost, error) {
	var schedule ScheduledPost
	err := r.db.WithContext(ctx).
		Where("social_post_id = ? AND status IN ?", socialPostID, []ScheduledPostStatus{
			ScheduledPostStatusPending,
			ScheduledPostStatusScheduled,
		}).
		First(&schedule).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrScheduledPostNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &schedule, nil
}

func (r *gormRepository) GetScheduleForDateRange(ctx context.Context, startDate, endDate time.Time) ([]ScheduledPost, error) {
	var schedules []ScheduledPost
	err := r.db.WithContext(ctx).
		Where("scheduled_at >= ? AND scheduled_at <= ?", startDate, endDate).
		Where("status != ?", ScheduledPostStatusCancelled).
		Order("scheduled_at asc").
		Find(&schedules).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return schedules, nil
}

func (r *gormRepository) GetUpcomingPosts(ctx context.Context, limit int) ([]ScheduledPost, error) {
	var schedules []ScheduledPost
	now := time.Now().UTC()
	query := r.db.WithContext(ctx).
		Where("scheduled_at >= ?", now).
		Where("status IN ?", []ScheduledPostStatus{
			ScheduledPostStatusPending,
			ScheduledPostStatusScheduled,
		}).
		Order("scheduled_at asc").
		Limit(limit)

	if err := query.Find(&schedules).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return schedules, nil
}

func (r *gormRepository) CheckConflict(ctx context.Context, scheduledAt time.Time, excludeScheduleID *uuid.UUID) (bool, error) {
	// Check for conflicts within a 15-minute window
	windowStart := scheduledAt.Add(-15 * time.Minute)
	windowEnd := scheduledAt.Add(15 * time.Minute)

	var count int64
	query := r.db.WithContext(ctx).
		Model(&ScheduledPost{}).
		Where("scheduled_at >= ? AND scheduled_at <= ?", windowStart, windowEnd).
		Where("status IN ?", []ScheduledPostStatus{
			ScheduledPostStatusPending,
			ScheduledPostStatusScheduled,
		})

	if excludeScheduleID != nil {
		query = query.Where("id != ?", *excludeScheduleID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return count > 0, nil
}

func (r *gormRepository) DeleteSchedule(ctx context.Context, scheduleID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&ScheduledPost{}, "id = ?", scheduleID).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

// isUniqueConstraintError checks if the error is a unique constraint violation.
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// Check for PostgreSQL unique constraint violation (code 23505)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}

	// Check for SQLite unique constraint violation
	if err.Error() == "UNIQUE constraint failed" {
		return true
	}

	return false
}
