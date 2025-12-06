package interviewstages

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates interview stage workflows.
type Service interface {
	CreateStage(ctx context.Context, jobApplicationID uuid.UUID, stageType StageType) (*InterviewStage, error)
	GetStage(ctx context.Context, stageID uuid.UUID) (*InterviewStage, error)
	ListStages(ctx context.Context, filters StageFilters) ([]InterviewStage, error)
	UpdateStage(ctx context.Context, stageID uuid.UUID, updates UpdateStageRequest) (*InterviewStage, error)
	DeleteStage(ctx context.Context, stageID uuid.UUID) error
	GetStagesByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]InterviewStage, error)
	ScheduleStage(ctx context.Context, stageID uuid.UUID, scheduledDate time.Time) (*InterviewStage, error)
	CompleteStage(ctx context.Context, stageID uuid.UUID, completedDate time.Time, outcome StageOutcome) (*InterviewStage, error)
}

// UpdateStageRequest represents fields that can be updated on a stage.
type UpdateStageRequest struct {
	ScheduledDate    *time.Time
	InterviewerName  *string
	InterviewerEmail *string
	Location         *string
	Notes            *string
	Feedback         *string
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService constructs a Service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) CreateStage(ctx context.Context, jobApplicationID uuid.UUID, stageType StageType) (*InterviewStage, error) {
	stage, err := NewInterviewStage(jobApplicationID, stageType)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateStage(ctx, stage); err != nil {
		return nil, err
	}

	return stage, nil
}

func (s *service) GetStage(ctx context.Context, stageID uuid.UUID) (*InterviewStage, error) {
	return s.repo.GetStage(ctx, stageID)
}

func (s *service) ListStages(ctx context.Context, filters StageFilters) ([]InterviewStage, error) {
	return s.repo.ListStages(ctx, filters)
}

func (s *service) UpdateStage(ctx context.Context, stageID uuid.UUID, updates UpdateStageRequest) (*InterviewStage, error) {
	stage, err := s.repo.GetStage(ctx, stageID)
	if err != nil {
		return nil, err
	}

	if updates.ScheduledDate != nil {
		stage.Schedule(*updates.ScheduledDate)
	}

	if updates.InterviewerName != nil || updates.InterviewerEmail != nil {
		name := ""
		email := ""
		if updates.InterviewerName != nil {
			name = *updates.InterviewerName
		} else {
			name = stage.InterviewerName
		}
		if updates.InterviewerEmail != nil {
			email = *updates.InterviewerEmail
		} else {
			email = stage.InterviewerEmail
		}
		stage.UpdateInterviewerInfo(name, email)
	}

	if updates.Location != nil {
		stage.UpdateLocation(*updates.Location)
	}

	if updates.Notes != nil {
		stage.UpdateNotes(*updates.Notes)
	}

	if updates.Feedback != nil {
		stage.UpdateFeedback(*updates.Feedback)
	}

	if err := s.repo.UpdateStage(ctx, stage); err != nil {
		return nil, err
	}

	return stage, nil
}

func (s *service) DeleteStage(ctx context.Context, stageID uuid.UUID) error {
	return s.repo.DeleteStage(ctx, stageID)
}

func (s *service) GetStagesByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]InterviewStage, error) {
	return s.repo.GetStagesByApplicationID(ctx, applicationID)
}

func (s *service) ScheduleStage(ctx context.Context, stageID uuid.UUID, scheduledDate time.Time) (*InterviewStage, error) {
	stage, err := s.repo.GetStage(ctx, stageID)
	if err != nil {
		return nil, err
	}

	stage.Schedule(scheduledDate)
	if err := s.repo.UpdateStage(ctx, stage); err != nil {
		return nil, err
	}

	return stage, nil
}

func (s *service) CompleteStage(ctx context.Context, stageID uuid.UUID, completedDate time.Time, outcome StageOutcome) (*InterviewStage, error) {
	stage, err := s.repo.GetStage(ctx, stageID)
	if err != nil {
		return nil, err
	}

	if err := stage.Complete(completedDate, outcome); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateStage(ctx, stage); err != nil {
		return nil, err
	}

	return stage, nil
}

