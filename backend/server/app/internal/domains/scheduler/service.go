package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"

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
	Email       string
	PhoneNumber string
	Active      *bool
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
	if req.Email != "" {
		schedule.Email = req.Email
	}
	if req.PhoneNumber != "" {
		schedule.PhoneNumber = req.PhoneNumber
	}
	if req.Active != nil {
		schedule.Active = *req.Active
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

	if s.reports == nil {
		return fmt.Errorf("reports service not configured")
	}

	summary, err := s.reports.GenerateSummary(ctx, schedule.UserID)
	if err != nil {
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
		if s.logger != nil {
			s.logger.Error("scheduler: dispatch summary failed", slog.String("schedule_id", schedule.ID.String()), slog.Any("error", err))
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
	hour, minute, err := parseTime(schedule.TimeOfDay)
	if err != nil {
		return time.Time{}, err
	}

	next := time.Date(refLocal.Year(), refLocal.Month(), refLocal.Day(), hour, minute, 0, 0, loc)
	if !next.After(refLocal) {
		next = next.Add(24 * time.Hour)
	}

	if schedule.Frequency == "weekly" {
		targetWeekday, err := parseWeekday(schedule.Weekday)
		if err != nil {
			return time.Time{}, err
		}
		for next.Weekday() != targetWeekday || !next.After(refLocal) {
			next = next.Add(24 * time.Hour)
		}
	}

	return next.UTC(), nil
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
