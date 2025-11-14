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
	ID          uuid.UUID                           `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID                           `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	ReportType  string                              `gorm:"column:report_type;size:64;not null" json:"reportType"`
	AgentAlias  string                              `gorm:"column:agent_alias;size:64;not null" json:"agentAlias"`
	Frequency   string                              `gorm:"column:frequency;size:32;not null" json:"frequency"`  // daily, weekly, custom
	Weekday     string                              `gorm:"column:weekday;size:16" json:"weekday"`               // monday ... sunday
	TimeOfDay   string                              `gorm:"column:time_of_day;size:8;not null" json:"timeOfDay"` // HH:MM (24h)
	Timezone    string                              `gorm:"column:timezone;size:64;not null" json:"timezone"`
	RRule       string                              `gorm:"column:rrule;type:text" json:"rrule"`
	Priority    int                                 `gorm:"column:priority;default:0" json:"priority"`
	Email       string                              `gorm:"column:email;size:255" json:"email"`
	PhoneNumber string                              `gorm:"column:phone_number;size:64" json:"phoneNumber"`
	Channels    datatypes.JSONType[map[string]bool] `gorm:"column:channels;type:jsonb" json:"channels"`
	Active      bool                                `gorm:"column:active;default:true" json:"active"`
	Paused      bool                                `gorm:"column:paused;default:false" json:"paused"`
	NextRun     time.Time                           `gorm:"column:next_run;index" json:"nextRun"`
	LastRun     *time.Time                          `gorm:"column:last_run" json:"lastRun,omitempty"`
	CreatedAt   time.Time                           `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time                           `gorm:"column:updated_at" json:"updatedAt"`
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
	ID           uuid.UUID                          `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID                          `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	ScheduleID   uuid.UUID                          `gorm:"column:schedule_id;type:uuid;index;not null" json:"scheduleId"`
	Status       string                             `gorm:"column:status;size:32;index" json:"status"`
	Output       string                             `gorm:"column:output;size:255" json:"output"`
	ErrorMessage string                             `gorm:"column:error_message;size:255" json:"errorMessage"`
	Metadata     datatypes.JSONType[map[string]any] `gorm:"column:metadata;type:jsonb" json:"metadata"`
	StartedAt    *time.Time                         `gorm:"column:started_at" json:"startedAt,omitempty"`
	CompletedAt  *time.Time                         `gorm:"column:completed_at" json:"completedAt,omitempty"`
	CreatedAt    time.Time                          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time                          `gorm:"column:updated_at" json:"updatedAt"`
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
