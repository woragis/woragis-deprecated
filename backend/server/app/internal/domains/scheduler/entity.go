package scheduler

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const (
	RunStatusPending   = "pending"
	RunStatusRunning   = "running"
	RunStatusCompleted = "completed"
	RunStatusFailed    = "failed"
)

// Schedule represents an automated report dispatch configuration.
type Schedule struct {
	ID          uuid.UUID                           `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID                           `gorm:"type:uuid;index;not null"`
	ReportType  string                              `gorm:"size:64;not null"`
	AgentAlias  string                              `gorm:"size:64;not null"`
	Frequency   string                              `gorm:"size:32;not null"` // daily, weekly, custom
	Weekday     string                              `gorm:"size:16"`          // monday ... sunday
	TimeOfDay   string                              `gorm:"size:8;not null"`  // HH:MM (24h)
	Timezone    string                              `gorm:"size:64;not null"`
	RRule       string                              `gorm:"type:text"`
	Priority    int                                 `gorm:"default:0"`
	Email       string                              `gorm:"size:255"`
	PhoneNumber string                              `gorm:"size:64"`
	Channels    datatypes.JSONType[map[string]bool] `gorm:"type:jsonb"`
	Active      bool                                `gorm:"default:true"`
	Paused      bool                                `gorm:"default:false"`
	NextRun     time.Time                           `gorm:"index"`
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
	rrule string,
	priority int,
	channels map[string]bool,
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
		RRule:       strings.TrimSpace(rrule),
		Priority:    priority,
		Channels:    datatypes.JSONType[map[string]bool]{Data: channels},
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

	switch s.Frequency {
	case "daily", "weekly", "custom":
	default:
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

	if s.Frequency == "custom" && s.RRule == "" {
		return NewDomainError(ErrCodeInvalidFrequency, ErrRRuleRequired)
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

// Pause marks the schedule as paused.
func (s *Schedule) Pause() {
	s.Paused = true
	s.UpdatedAt = time.Now().UTC()
}

// Resume clears the paused flag.
func (s *Schedule) Resume() {
	s.Paused = false
	s.UpdatedAt = time.Now().UTC()
}

// ExecutionRun keeps history of schedule executions.
type ExecutionRun struct {
	ID           uuid.UUID                          `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID                          `gorm:"type:uuid;index;not null"`
	ScheduleID   uuid.UUID                          `gorm:"type:uuid;index;not null"`
	Status       string                             `gorm:"size:32;index"`
	Output       string                             `gorm:"size:255"`
	ErrorMessage string                             `gorm:"size:255"`
	Metadata     datatypes.JSONType[map[string]any] `gorm:"type:jsonb"`
	StartedAt    *time.Time
	CompletedAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewExecutionRun constructs a new run entry.
func NewExecutionRun(userID, scheduleID uuid.UUID, status string, metadata map[string]any) *ExecutionRun {
	now := time.Now().UTC()
	return &ExecutionRun{
		ID:         uuid.New(),
		UserID:     userID,
		ScheduleID: scheduleID,
		Status:     strings.ToLower(strings.TrimSpace(status)),
		Metadata:   datatypes.JSONType[map[string]any]{Data: metadata},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// MarkStarted sets start timestamp.
func (r *ExecutionRun) MarkStarted() {
	now := time.Now().UTC()
	r.Status = RunStatusRunning
	r.StartedAt = &now
	r.UpdatedAt = now
}

// MarkCompleted records successful completion.
func (r *ExecutionRun) MarkCompleted(output string) {
	now := time.Now().UTC()
	r.Status = RunStatusCompleted
	r.Output = strings.TrimSpace(output)
	r.CompletedAt = &now
	r.UpdatedAt = now
}

// MarkFailed records failure.
func (r *ExecutionRun) MarkFailed(err error) {
	now := time.Now().UTC()
	r.Status = RunStatusFailed
	if err != nil {
		r.ErrorMessage = strings.TrimSpace(err.Error())
	}
	r.CompletedAt = &now
	r.UpdatedAt = now
}
