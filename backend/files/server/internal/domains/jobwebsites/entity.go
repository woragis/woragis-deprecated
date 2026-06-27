package jobwebsites

import (
	"time"

	"github.com/google/uuid"
)

// JobWebsite represents a job website configuration.
type JobWebsite struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name         string    `gorm:"column:name;size:50;not null;uniqueIndex" json:"name"` // "linkedin", "glassdoor", etc.
	DisplayName  string    `gorm:"column:display_name;size:255;not null" json:"displayName"`
	DailyLimit   int       `gorm:"column:daily_limit;not null;default:50" json:"dailyLimit"`
	CurrentCount int       `gorm:"column:current_count;not null;default:0" json:"currentCount"`
	LastReset    time.Time `gorm:"column:last_reset;not null" json:"lastReset"`
	Enabled      bool      `gorm:"column:enabled;not null;default:true" json:"enabled"`
	BaseURL      string    `gorm:"column:base_url;size:512" json:"baseUrl"`
	LoginURL     string    `gorm:"column:login_url;size:512" json:"loginUrl"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for JobWebsite.
func (JobWebsite) TableName() string {
	return "job_websites"
}

// NewJobWebsite creates a new job website entity.
func NewJobWebsite(name, displayName, baseURL, loginURL string, dailyLimit int) (*JobWebsite, error) {
	website := &JobWebsite{
		ID:          uuid.New(),
		Name:         name,
		DisplayName:  displayName,
		DailyLimit:   dailyLimit,
		CurrentCount: 0,
		LastReset:    time.Now().UTC(),
		Enabled:      true,
		BaseURL:      baseURL,
		LoginURL:     loginURL,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	return website, website.Validate()
}

// Validate ensures job website invariants hold.
func (j *JobWebsite) Validate() error {
	if j.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyWebsiteID)
	}
	if j.Name == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyWebsiteName)
	}
	if j.DisplayName == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDisplayName)
	}
	if j.DailyLimit < 0 {
		return NewDomainError(ErrCodeInvalidPayload, ErrInvalidDailyLimit)
	}
	if j.CurrentCount < 0 {
		return NewDomainError(ErrCodeInvalidPayload, ErrInvalidCurrentCount)
	}
	return nil
}

// IncrementCount increments the current count by 1.
func (j *JobWebsite) IncrementCount() {
	j.CurrentCount++
	j.UpdatedAt = time.Now().UTC()
}

// ResetCount resets the current count to 0 and updates last reset time.
func (j *JobWebsite) ResetCount() {
	j.CurrentCount = 0
	j.LastReset = time.Now().UTC()
	j.UpdatedAt = time.Now().UTC()
}

// IsLimitReached returns true if the daily limit has been reached.
func (j *JobWebsite) IsLimitReached() bool {
	return j.CurrentCount >= j.DailyLimit
}

// ShouldReset checks if the counter should be reset (new day).
func (j *JobWebsite) ShouldReset() bool {
	now := time.Now().UTC()
	lastReset := j.LastReset.UTC()
	
	// Reset if it's a different day
	return now.Year() != lastReset.Year() || 
		   now.YearDay() != lastReset.YearDay()
}

// UpdateDailyLimit updates the daily limit.
func (j *JobWebsite) UpdateDailyLimit(limit int) error {
	if limit < 0 {
		return NewDomainError(ErrCodeInvalidPayload, ErrInvalidDailyLimit)
	}
	j.DailyLimit = limit
	j.UpdatedAt = time.Now().UTC()
	return nil
}

// Enable or disable the website.
func (j *JobWebsite) SetEnabled(enabled bool) {
	j.Enabled = enabled
	j.UpdatedAt = time.Now().UTC()
}

