package reports

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ReportDefinition models a saved custom report configuration.
type ReportDefinition struct {
	ID          uuid.UUID                          `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID                          `gorm:"type:uuid;index;not null"`
	Name        string                             `gorm:"size:120;not null"`
	Description string                             `gorm:"size:255"`
	Sections    datatypes.JSONType[map[string]any] `gorm:"type:jsonb"`
	Filters     datatypes.JSONType[map[string]any] `gorm:"type:jsonb"`
	IsFavorite  bool
	ArchivedAt  *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// NewReportDefinition constructs a report definition entity.
func NewReportDefinition(userID uuid.UUID, name, description string, sections, filters datatypes.JSONType[map[string]any], favorite bool) (*ReportDefinition, error) {
	def := &ReportDefinition{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Sections:    sections,
		Filters:     filters,
		IsFavorite:  favorite,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	return def, def.Validate()
}

// Validate ensures required fields are populated.
func (r *ReportDefinition) Validate() error {
	if r == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilReportDefinition)
	}
	if r.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyReportDefinitionID)
	}
	if r.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if strings.TrimSpace(r.Name) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyReportName)
	}
	return nil
}

// Archive marks the report as archived.
func (r *ReportDefinition) Archive() {
	now := time.Now().UTC()
	r.ArchivedAt = &now
	r.UpdatedAt = now
}

// Restore clears the archived flag.
func (r *ReportDefinition) Restore() {
	r.ArchivedAt = nil
	r.UpdatedAt = time.Now().UTC()
}

// ReportSchedule defines automation for a report.
type ReportSchedule struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	ReportID  uuid.UUID `gorm:"type:uuid;index;not null"`
	Cron      string    `gorm:"size:120"`
	Frequency string    `gorm:"size:32"`
	Timezone  string    `gorm:"size:64"`
	NextRun   *time.Time
	LastRunAt *time.Time
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt                     `gorm:"index"`
	Meta      datatypes.JSONType[map[string]any] `gorm:"type:jsonb"`
}

// NewReportSchedule constructs a schedule entity.
func NewReportSchedule(reportID uuid.UUID, cron, frequency, timezone string, nextRun *time.Time, enabled bool, meta datatypes.JSONType[map[string]any]) (*ReportSchedule, error) {
	s := &ReportSchedule{
		ID:        uuid.New(),
		ReportID:  reportID,
		Cron:      strings.TrimSpace(cron),
		Frequency: strings.ToLower(strings.TrimSpace(frequency)),
		Timezone:  strings.TrimSpace(timezone),
		NextRun:   nextRun,
		Enabled:   enabled,
		Meta:      meta,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return s, s.Validate()
}

// Validate ensures schedule invariants.
func (s *ReportSchedule) Validate() error {
	if s == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilReportSchedule)
	}
	if s.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyScheduleID)
	}
	if s.ReportID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyReportDefinitionID)
	}
	return nil
}

// Toggle sets the enabled flag.
func (s *ReportSchedule) Toggle(enabled bool) {
	s.Enabled = enabled
	s.UpdatedAt = time.Now().UTC()
}

// ReportDelivery defines how a report is delivered.
type ReportDelivery struct {
	ID        uuid.UUID                          `gorm:"type:uuid;primaryKey"`
	ReportID  uuid.UUID                          `gorm:"type:uuid;index;not null"`
	Channel   string                             `gorm:"size:32;not null"`
	Target    string                             `gorm:"size:255"`
	Template  datatypes.JSONType[map[string]any] `gorm:"type:jsonb"`
	Enabled   bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// NewReportDelivery constructs a delivery entity.
func NewReportDelivery(reportID uuid.UUID, channel, target string, template datatypes.JSONType[map[string]any], enabled bool) (*ReportDelivery, error) {
	d := &ReportDelivery{
		ID:        uuid.New(),
		ReportID:  reportID,
		Channel:   strings.ToLower(strings.TrimSpace(channel)),
		Target:    strings.TrimSpace(target),
		Template:  template,
		Enabled:   enabled,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return d, d.Validate()
}

// Validate ensures delivery invariants.
func (d *ReportDelivery) Validate() error {
	if d == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilReportDelivery)
	}
	if d.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDeliveryID)
	}
	if d.ReportID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyReportDefinitionID)
	}
	if d.Channel == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDeliveryChannel)
	}
	return nil
}

// Toggle enables or disables the delivery.
func (d *ReportDelivery) Toggle(enabled bool) {
	d.Enabled = enabled
	d.UpdatedAt = time.Now().UTC()
}

// ReportRun tracks regeneration/export requests.
type ReportRun struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	ReportID       uuid.UUID `gorm:"type:uuid;index;not null"`
	RequestedBy    uuid.UUID `gorm:"type:uuid;index"`
	Status         string    `gorm:"size:32;index"`
	StartedAt      *time.Time
	CompletedAt    *time.Time
	OutputLocation string                             `gorm:"size:255"`
	ErrorMessage   string                             `gorm:"size:255"`
	Metadata       datatypes.JSONType[map[string]any] `gorm:"type:jsonb"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewReportRun constructs a pending run entity.
func NewReportRun(reportID, requestedBy uuid.UUID, metadata datatypes.JSONType[map[string]any]) *ReportRun {
	return &ReportRun{
		ID:          uuid.New(),
		ReportID:    reportID,
		RequestedBy: requestedBy,
		Status:      RunStatusPending,
		Metadata:    metadata,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

// MarkStarted updates the run status to in-progress.
func (r *ReportRun) MarkStarted() {
	now := time.Now().UTC()
	r.Status = RunStatusRunning
	r.StartedAt = &now
	r.UpdatedAt = now
}

// MarkCompleted marks the run as completed.
func (r *ReportRun) MarkCompleted(output string) {
	now := time.Now().UTC()
	r.Status = RunStatusCompleted
	r.CompletedAt = &now
	r.OutputLocation = strings.TrimSpace(output)
	r.UpdatedAt = now
}

// MarkFailed marks the run as failed.
func (r *ReportRun) MarkFailed(err error) {
	now := time.Now().UTC()
	r.Status = RunStatusFailed
	r.ErrorMessage = strings.TrimSpace(err.Error())
	r.CompletedAt = &now
	r.UpdatedAt = now
}
