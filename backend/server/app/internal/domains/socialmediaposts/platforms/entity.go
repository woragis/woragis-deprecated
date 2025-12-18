package platforms

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	
	"github.com/woragis/backend/server/app/internal/domains/socialmediaposts"
)

// PlatformConfig represents configuration for a social media platform.
type PlatformConfig struct {
	ID               uuid.UUID                    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name             socialmediaposts.Platform     `gorm:"column:name;type:varchar(20);not null;uniqueIndex" json:"name"`
	DisplayName      string                       `gorm:"column:display_name;type:varchar(100);not null" json:"displayName"`
	PostingFrequency *int                         `gorm:"column:posting_frequency;type:integer" json:"postingFrequency,omitempty"` // Posts per week
	BestDays         datatypes.JSON               `gorm:"column:best_days;type:jsonb" json:"bestDays,omitempty"`                    // Array of day names
	BestTimes        datatypes.JSON               `gorm:"column:best_times;type:jsonb" json:"bestTimes,omitempty"`                  // Array of time ranges
	SupportedFormats datatypes.JSON               `gorm:"column:supported_formats;type:jsonb" json:"supportedFormats"`            // Array of ContentFormat
	IsActive         bool                         `gorm:"column:is_active;not null;default:true;index" json:"isActive"`
	Metadata         datatypes.JSON               `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	CreatedAt        time.Time                    `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt        time.Time                    `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for PlatformConfig.
func (PlatformConfig) TableName() string {
	return "platform_configs"
}

// NewPlatformConfig creates a new platform configuration.
func NewPlatformConfig(name socialmediaposts.Platform, displayName string, supportedFormats []socialmediaposts.ContentFormat) (*PlatformConfig, error) {
	config := &PlatformConfig{
		ID:               uuid.New(),
		Name:             name,
		DisplayName:      strings.TrimSpace(displayName),
		SupportedFormats: datatypes.JSON(mustMarshalJSON(supportedFormats)),
		IsActive:         true,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	return config, config.Validate()
}

// Validate ensures platform config invariants hold.
func (p *PlatformConfig) Validate() error {
	if p == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilPlatformConfig)
	}

	if p.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyPlatformConfigID)
	}

	if !isValidPlatform(p.Name) {
		return NewDomainError(ErrCodeInvalidPlatform, ErrUnsupportedPlatform)
	}

	if strings.TrimSpace(p.DisplayName) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDisplayName)
	}

	return nil
}

// UpdatePostingFrequency sets the posting frequency.
func (p *PlatformConfig) UpdatePostingFrequency(frequency int) {
	if frequency > 0 {
		p.PostingFrequency = &frequency
		p.UpdatedAt = time.Now().UTC()
	}
}

// SetBestDays sets the best days for posting.
func (p *PlatformConfig) SetBestDays(days []string) {
	p.BestDays = datatypes.JSON(mustMarshalJSON(days))
	p.UpdatedAt = time.Now().UTC()
}

// SetBestTimes sets the best times for posting.
func (p *PlatformConfig) SetBestTimes(times []string) {
	p.BestTimes = datatypes.JSON(mustMarshalJSON(times))
	p.UpdatedAt = time.Now().UTC()
}

// SetSupportedFormats updates the supported content formats.
func (p *PlatformConfig) SetSupportedFormats(formats []socialmediaposts.ContentFormat) {
	p.SupportedFormats = datatypes.JSON(mustMarshalJSON(formats))
	p.UpdatedAt = time.Now().UTC()
}

// SetActive sets whether the platform is active.
func (p *PlatformConfig) SetActive(active bool) {
	p.IsActive = active
	p.UpdatedAt = time.Now().UTC()
}

// Validation helpers

func isValidPlatform(p socialmediaposts.Platform) bool {
	switch p {
	case socialmediaposts.PlatformLinkedIn, socialmediaposts.PlatformTwitter,
		socialmediaposts.PlatformInstagram, socialmediaposts.PlatformMedium,
		socialmediaposts.PlatformSubstack, socialmediaposts.PlatformValete,
		socialmediaposts.PlatformWebsite:
		return true
	}
	return false
}

// Helper function to marshal JSON (simple wrapper)
func mustMarshalJSON(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}
