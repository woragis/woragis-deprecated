package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/teambition/rrule-go"
	"gorm.io/datatypes"

	reportsdomain "github.com/woragis/backend/server/app/internal/domains/reports"
)

// Service orchestrates schedule management and execution.
type Service struct {
	repo            Repository
	reports         *reportsdomain.Service
	logger          *slog.Logger
	defaultEmail    bool
	defaultWhatsApp bool
}

// NewService constructs a Service.
func NewService(repo Repository, reports *reportsdomain.Service, logger *slog.Logger) *Service {
	return &Service{
		repo:    repo,
		reports: reports,
		logger:  logger,
	}
}

// CreateRequest encapsulates schedule creation data.
type CreateRequest struct {
	UserID      uuid.UUID
	ReportType  string
	AgentAlias  string
	Frequency   string
	Weekday     string
	TimeOfDay   string
	Timezone    string
	RRule       string
	Priority    int
	Channels    map[string]bool
	Email       string
	PhoneNumber string
}

// UpdateRequest updates schedule metadata.
type UpdateRequest struct {
	UserID      uuid.UUID
	ScheduleID  uuid.UUID
	ReportType  string
	AgentAlias  string
	Frequency   string
	Weekday     string
	TimeOfDay   string
	Timezone    string
	RRule       string
	Priority    *int
	Channels    map[string]bool
	Email       string
	PhoneNumber string
	Active      *bool
	Paused      *bool
}

// ListRunsRequest filters execution runs.
type ListRunsRequest struct {
	UserID     uuid.UUID
	ScheduleID uuid.UUID
	Status     string
	Limit      int
	Offset     int
}

// Create creates a new schedule and computes its next run.
func (s *Service) Create(ctx context.Context, req CreateRequest) (*Schedule, error) {
	schedule, err := NewSchedule(
		req.UserID,
		req.ReportType,
		req.AgentAlias,
		req.Frequency,
		req.Weekday,
		req.TimeOfDay,
		req.Timezone,
		req.RRule,
		req.Priority,
		normaliseChannels(req.Channels, req.Email, req.PhoneNumber),
		req.Email,
		req.PhoneNumber,
	)
	if err != nil {
		return nil, err
	}

	nextRun, err := computeNextRun(schedule, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	schedule.SetNextRun(nextRun)

	if err := s.repo.Create(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

// Update modifies an existing schedule.
func (s *Service) Update(ctx context.Context, req UpdateRequest) (*Schedule, error) {
	schedule, err := s.repo.Get(ctx, req.ScheduleID, req.UserID)
	if err != nil {
		return nil, err
	}

	if req.ReportType != "" {
		schedule.ReportType = req.ReportType
	}
	if req.AgentAlias != "" {
		schedule.AgentAlias = req.AgentAlias
	}
	if req.Frequency != "" {
		schedule.Frequency = strings.ToLower(req.Frequency)
	}
	if req.Weekday != "" {
		schedule.Weekday = strings.ToLower(req.Weekday)
	}
	if req.TimeOfDay != "" {
		schedule.TimeOfDay = req.TimeOfDay
	}
	if req.Timezone != "" {
		schedule.Timezone = req.Timezone
	}
	if req.RRule != "" {
		schedule.RRule = req.RRule
	}
	if req.Priority != nil {
		schedule.Priority = *req.Priority
	}
	if req.Channels != nil {
		schedule.Channels = datatypes.JSONType[map[string]bool]{Data: normaliseChannels(req.Channels, schedule.Email, schedule.PhoneNumber)}
	}
	if req.Email != "" {
		schedule.Email = req.Email
	}
	if req.PhoneNumber != "" {
		schedule.PhoneNumber = req.PhoneNumber
	}
	if req.Active != nil {
		schedule.Active = *req.Active
	}
	if req.Paused != nil {
		if *req.Paused {
			schedule.Pause()
		} else {
			schedule.Resume()
		}
	}

	if err := schedule.Validate(); err != nil {
		return nil, err
	}

	nextRun, err := computeNextRun(schedule, time.Now().UTC())
	if err != nil {
		return nil, err
	}
	schedule.SetNextRun(nextRun)

	if err := s.repo.Update(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

// Delete removes a schedule and its execution runs.
func (s *Service) Delete(ctx context.Context, scheduleID, userID uuid.UUID) error {
	return s.repo.Delete(ctx, scheduleID, userID)
}

// List returns schedules for the user.
func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]Schedule, error) {
	return s.repo.List(ctx, userID)
}

// ListDue returns schedules due for execution.
func (s *Service) ListDue(ctx context.Context, now time.Time) ([]Schedule, error) {
	return s.repo.ListDue(ctx, now)
}

// Execute triggers a schedule and computes its next run.
func (s *Service) Execute(ctx context.Context, schedule *Schedule) error {
	if !schedule.Active {
		return nil
	}
	if schedule.Paused {
		return nil
	}

	if s.reports == nil {
		return fmt.Errorf("reports service not configured")
	}

	run := NewExecutionRun(schedule.UserID, schedule.ID, RunStatusPending, nil)
	if err := s.repo.InsertRun(ctx, run); err != nil && s.logger != nil {
		s.logger.Warn("scheduler: failed to insert run", slog.Any("error", err))
	}
	run.MarkStarted()
	if err := s.repo.UpdateRun(ctx, run); err != nil && s.logger != nil {
		s.logger.Warn("scheduler: failed to update run", slog.Any("error", err))
	}

	summary, err := s.reports.GenerateSummary(ctx, schedule.UserID)
	if err != nil {
		run.MarkFailed(err)
		_ = s.repo.UpdateRun(ctx, run)
		return err
	}

	opts := reportsdomain.DispatchOptions{
		SendEmail:    schedule.Email != "",
		EmailAddress: schedule.Email,
		SendWhatsApp: schedule.PhoneNumber != "",
		PhoneNumber:  schedule.PhoneNumber,
		AgentAlias:   schedule.AgentAlias,
	}

	if err := s.reports.DispatchSummary(ctx, summary, opts); err != nil {
		run.MarkFailed(err)
		_ = s.repo.UpdateRun(ctx, run)
		if s.logger != nil {
			s.logger.Error("scheduler: dispatch summary failed", slog.String("schedule_id", schedule.ID.String()), slog.Any("error", err))
		}
	} else {
		run.MarkCompleted("dispatched")
		if err := s.repo.UpdateRun(ctx, run); err != nil && s.logger != nil {
			s.logger.Warn("scheduler: failed to finalize run", slog.Any("error", err))
		}
	}

	nextRun, err := computeNextRun(schedule, time.Now().UTC().Add(time.Minute))
	if err != nil {
		return err
	}

	schedule.MarkExecuted(nextRun)
	return s.repo.Update(ctx, schedule)
}

// computeNextRun returns the next execution time after reference.
func computeNextRun(schedule *Schedule, reference time.Time) (time.Time, error) {
	loc, err := time.LoadLocation(schedule.Timezone)
	if err != nil {
		loc = time.UTC
	}

	refLocal := reference.In(loc)

	switch schedule.Frequency {
	case "daily":
		return nextDailyRun(schedule.TimeOfDay, refLocal)
	case "weekly":
		return nextWeeklyRun(schedule.TimeOfDay, schedule.Weekday, refLocal)
	case "custom":
		return nextRRuleRun(schedule, refLocal)
	default:
		return time.Time{}, NewDomainError(ErrCodeInvalidFrequency, ErrUnsupportedFrequency)
	}
}

func parseTime(value string) (int, int, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid time of day: %s", value)
	}

	var hour, minute int
	if _, err := fmt.Sscanf(value, "%d:%d", &hour, &minute); err != nil {
		return 0, 0, fmt.Errorf("invalid time of day: %s", value)
	}

	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, 0, fmt.Errorf("invalid time of day: %s", value)
	}

	return hour, minute, nil
}

func parseWeekday(value string) (time.Weekday, error) {
	switch strings.ToLower(value) {
	case "sunday", "sun":
		return time.Sunday, nil
	case "monday", "mon":
		return time.Monday, nil
	case "tuesday", "tue":
		return time.Tuesday, nil
	case "wednesday", "wed":
		return time.Wednesday, nil
	case "thursday", "thu":
		return time.Thursday, nil
	case "friday", "fri":
		return time.Friday, nil
	case "saturday", "sat":
		return time.Saturday, nil
	default:
		return time.Sunday, fmt.Errorf("invalid weekday: %s", value)
	}
}

func nextDailyRun(timeOfDay string, reference time.Time) (time.Time, error) {
	hour, minute, err := parseTime(timeOfDay)
	if err != nil {
		return time.Time{}, err
	}

	next := time.Date(reference.Year(), reference.Month(), reference.Day(), hour, minute, 0, 0, reference.Location())
	if !next.After(reference) {
		next = next.Add(24 * time.Hour)
	}
	return next.UTC(), nil
}

func nextWeeklyRun(timeOfDay, weekday string, reference time.Time) (time.Time, error) {
	hour, minute, err := parseTime(timeOfDay)
	if err != nil {
		return time.Time{}, err
	}
	targetWeekday, err := parseWeekday(weekday)
	if err != nil {
		return time.Time{}, err
	}

	next := time.Date(reference.Year(), reference.Month(), reference.Day(), hour, minute, 0, 0, reference.Location())
	for !next.After(reference) || next.Weekday() != targetWeekday {
		next = next.Add(24 * time.Hour)
	}
	return next.UTC(), nil
}

func nextRRuleRun(schedule *Schedule, reference time.Time) (time.Time, error) {
	if schedule.RRule == "" {
		return time.Time{}, NewDomainError(ErrCodeInvalidFrequency, ErrRRuleRequired)
	}
	rule, err := rrule.StrToRRule(schedule.RRule)
	if err != nil {
		return time.Time{}, NewDomainError(ErrCodeInvalidFrequency, ErrRRuleRequired)
	}
	next := rule.After(reference, false)
	if next.IsZero() {
		return time.Time{}, NewDomainError(ErrCodeInvalidFrequency, ErrUnableToComputeNextRun)
	}
	return next.UTC(), nil
}

// BulkSetActive toggles active flag in bulk.
func (s *Service) BulkSetActive(ctx context.Context, userID uuid.UUID, ids []uuid.UUID, active bool) error {
	if len(ids) == 0 {
		return nil
	}
	updates := map[string]any{
		"active":     active,
		"updated_at": time.Now().UTC(),
	}
	if !active {
		updates["paused"] = false
	}
	return s.repo.BulkUpdateState(ctx, userID, ids, updates)
}

// BulkPause toggles paused flag in bulk.
func (s *Service) BulkPause(ctx context.Context, userID uuid.UUID, ids []uuid.UUID, paused bool) error {
	if len(ids) == 0 {
		return nil
	}
	updates := map[string]any{
		"paused":     paused,
		"updated_at": time.Now().UTC(),
	}
	return s.repo.BulkUpdateState(ctx, userID, ids, updates)
}

// ListRuns returns execution history.
func (s *Service) ListRuns(ctx context.Context, req ListRunsRequest) ([]ExecutionRun, error) {
	return s.repo.ListRuns(ctx, req.UserID, RunFilters{
		ScheduleID: req.ScheduleID,
		Status:     req.Status,
		Limit:      req.Limit,
		Offset:     req.Offset,
	})
}

func normaliseChannels(ch map[string]bool, email, phone string) map[string]bool {
	if ch == nil {
		ch = make(map[string]bool)
	}
	if email != "" {
		ch["email"] = true
	}
	if phone != "" {
		ch["whatsapp"] = true
	}
	result := make(map[string]bool, len(ch))
	for key, value := range ch {
		if value {
			result[strings.ToLower(strings.TrimSpace(key))] = true
		}
	}
	return result
}
