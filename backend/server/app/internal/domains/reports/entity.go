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
	ID          uuid.UUID                          `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID                          `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Name        string                             `gorm:"column:name;size:120;not null" json:"name"`
	Description string                             `gorm:"column:description;size:255" json:"description"`
	Sections    datatypes.JSONType[map[string]any] `gorm:"column:sections;type:jsonb" json:"sections"`
	Filters     datatypes.JSONType[map[string]any] `gorm:"column:filters;type:jsonb" json:"filters"`
	IsFavorite  bool                               `gorm:"column:is_favorite" json:"isFavorite"`
	ArchivedAt  *time.Time                         `gorm:"column:archived_at" json:"archivedAt,omitempty"`
	CreatedAt   time.Time                          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time                          `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt                     `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
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
	ID        uuid.UUID                          `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ReportID  uuid.UUID                          `gorm:"column:report_id;type:uuid;index;not null" json:"reportId"`
	Cron      string                             `gorm:"column:cron;size:120" json:"cron"`
	Frequency string                             `gorm:"column:frequency;size:32" json:"frequency"`
	Timezone  string                             `gorm:"column:timezone;size:64" json:"timezone"`
	NextRun   *time.Time                         `gorm:"column:next_run" json:"nextRun,omitempty"`
	LastRunAt *time.Time                         `gorm:"column:last_run_at" json:"lastRunAt,omitempty"`
	Enabled   bool                               `gorm:"column:enabled" json:"enabled"`
	CreatedAt time.Time                          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time                          `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt                     `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
	Meta      datatypes.JSONType[map[string]any] `gorm:"column:meta;type:jsonb" json:"meta"`
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
	ID        uuid.UUID                          `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ReportID  uuid.UUID                          `gorm:"column:report_id;type:uuid;index;not null" json:"reportId"`
	Channel   string                             `gorm:"column:channel;size:32;not null" json:"channel"`
	Target    string                             `gorm:"column:target;size:255" json:"target"`
	Template  datatypes.JSONType[map[string]any] `gorm:"column:template;type:jsonb" json:"template"`
	Enabled   bool                               `gorm:"column:enabled" json:"enabled"`
	CreatedAt time.Time                          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time                          `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt gorm.DeletedAt                     `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
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
	ID             uuid.UUID                          `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ReportID       uuid.UUID                          `gorm:"column:report_id;type:uuid;index;not null" json:"reportId"`
	RequestedBy    uuid.UUID                          `gorm:"column:requested_by;type:uuid;index" json:"requestedBy"`
	Status         string                             `gorm:"column:status;size:32;index" json:"status"`
	StartedAt      *time.Time                         `gorm:"column:started_at" json:"startedAt,omitempty"`
	CompletedAt    *time.Time                         `gorm:"column:completed_at" json:"completedAt,omitempty"`
	OutputLocation string                             `gorm:"column:output_location;size:255" json:"outputLocation"`
	ErrorMessage   string                             `gorm:"column:error_message;size:255" json:"errorMessage"`
	Metadata       datatypes.JSONType[map[string]any] `gorm:"column:metadata;type:jsonb" json:"metadata"`
	CreatedAt      time.Time                          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time                          `gorm:"column:updated_at" json:"updatedAt"`
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
