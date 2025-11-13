package scheduler

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Schedule represents an automated report dispatch configuration.
type Schedule struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;index;not null"`
	ReportType  string    `gorm:"size:64;not null"`
	AgentAlias  string    `gorm:"size:64;not null"`
	Frequency   string    `gorm:"size:32;not null"` // daily, weekly
	Weekday     string    `gorm:"size:16"`          // monday ... sunday
	TimeOfDay   string    `gorm:"size:8;not null"`  // HH:MM (24h)
	Timezone    string    `gorm:"size:64;not null"`
	Email       string    `gorm:"size:255"`
	PhoneNumber string    `gorm:"size:64"`
	Active      bool      `gorm:"default:true"`
	NextRun     time.Time `gorm:"index"`
	LastRun     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewSchedule constructs a schedule.
func NewSchedule(
	userID uuid.UUID,
	reportType,
	agentAlias,
	frequency,
	weekday,
	timeOfDay,
	timezone,
	email,
	phone string,
) (*Schedule, error) {
	schedule := &Schedule{
		ID:          uuid.New(),
		UserID:      userID,
		ReportType:  strings.TrimSpace(reportType),
		AgentAlias:  strings.TrimSpace(agentAlias),
		Frequency:   strings.ToLower(strings.TrimSpace(frequency)),
		Weekday:     strings.ToLower(strings.TrimSpace(weekday)),
		TimeOfDay:   strings.TrimSpace(timeOfDay),
		Timezone:    strings.TrimSpace(timezone),
		Email:       strings.TrimSpace(email),
		PhoneNumber: strings.TrimSpace(phone),
		Active:      true,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	return schedule, schedule.Validate()
}

// Validate checks schedule invariants.
func (s *Schedule) Validate() error {
	if s == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilSchedule)
	}

	if s.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyScheduleID)
	}

	if s.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if s.ReportType == "" {
		return NewDomainError(ErrCodeInvalidReport, ErrEmptyReportType)
	}

	if s.AgentAlias == "" {
		return NewDomainError(ErrCodeInvalidAgent, ErrEmptyAgentAlias)
	}

	if s.Frequency != "daily" && s.Frequency != "weekly" {
		return NewDomainError(ErrCodeInvalidFrequency, ErrUnsupportedFrequency)
	}

	if s.Frequency == "weekly" && s.Weekday == "" {
		return NewDomainError(ErrCodeInvalidFrequency, ErrWeekdayRequired)
	}

	if s.TimeOfDay == "" {
		return NewDomainError(ErrCodeInvalidFrequency, ErrTimeRequired)
	}

	if s.Timezone == "" {
		s.Timezone = "UTC"
	}

	return nil
}

// SetNextRun sets the next run timestamp.
func (s *Schedule) SetNextRun(next time.Time) {
	s.NextRun = next
	s.UpdatedAt = time.Now().UTC()
}

// MarkExecuted updates the last run timestamp.
func (s *Schedule) MarkExecuted(next time.Time) {
	now := time.Now().UTC()
	s.LastRun = &now
	s.NextRun = next
	s.UpdatedAt = now
}
