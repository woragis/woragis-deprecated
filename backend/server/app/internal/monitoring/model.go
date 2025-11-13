package monitoring

import (
	"time"

	"github.com/google/uuid"
)

// Event represents a persisted monitoring event (stored only in production).
type Event struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Type      string    `gorm:"size:120;index;not null"`
	Reference string    `gorm:"size:120"`
	Payload   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
