package monitoring

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository abstracts persistence of monitoring events.
type Repository interface {
	StoreEvent(ctx context.Context, event Event) error
	ListEvents(ctx context.Context, limit int) ([]Event, error)
}

// GormRepository persists events via GORM.
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository constructs a repository backed by GORM database.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// StoreEvent persists the monitoring event.
func (r *GormRepository) StoreEvent(ctx context.Context, event Event) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}

	return r.db.WithContext(ctx).Create(&event).Error
}

// ListEvents returns most recent events up to the limit specified.
func (r *GormRepository) ListEvents(ctx context.Context, limit int) ([]Event, error) {
	if limit <= 0 {
		limit = 20
	}

	var events []Event
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&events).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return events, nil
}
