package projects

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates project workflows.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService constructs a Service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// CreateProjectRequest holds project creation payload.
type CreateProjectRequest struct {
	UserID      uuid.UUID
	Name        string
	Description string
	Status      ProjectStatus
	HealthScore int
	MRR         float64
	CAC         float64
	LTV         float64
	ChurnRate   float64
}

// UpdateStatusRequest updates project stage.
type UpdateStatusRequest struct {
	ProjectID uuid.UUID
	UserID    uuid.UUID
	Status    ProjectStatus
}

// UpdateMetricsRequest updates KPI metrics.
type UpdateMetricsRequest struct {
	ProjectID   uuid.UUID
	UserID      uuid.UUID
	HealthScore int
	MRR         float64
	CAC         float64
	LTV         float64
	ChurnRate   float64
}

// AddMilestoneRequest captures milestone creation data.
type AddMilestoneRequest struct {
	ProjectID   uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
	DueDate     time.Time
}

// ToggleMilestoneRequest toggles milestone completion.
type ToggleMilestoneRequest struct {
	MilestoneID uuid.UUID
	UserID      uuid.UUID
	Completed   bool
}

// CreateProject creates a new project aggregate.
func (s *Service) CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error) {
	if req.Status == "" {
		req.Status = ProjectStatusIdea
	}

	project, err := NewProject(req.UserID, req.Name, req.Description, req.Status, req.HealthScore, req.MRR, req.CAC, req.LTV, req.ChurnRate)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateProject(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// UpdateProjectStatus updates the status of a project.
func (s *Service) UpdateProjectStatus(ctx context.Context, req UpdateStatusRequest) (*Project, error) {
	project, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := project.UpdateStatus(req.Status); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateProject(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// UpdateProjectMetrics updates KPI metrics.
func (s *Service) UpdateProjectMetrics(ctx context.Context, req UpdateMetricsRequest) (*Project, error) {
	project, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := project.UpdateMetrics(req.HealthScore, req.MRR, req.CAC, req.LTV, req.ChurnRate); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateProject(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

// ListProjects returns all user projects.
func (s *Service) ListProjects(ctx context.Context, userID uuid.UUID) ([]Project, error) {
	return s.repo.ListProjects(ctx, userID)
}

// AddMilestone adds a new milestone to a project.
func (s *Service) AddMilestone(ctx context.Context, req AddMilestoneRequest) (*Milestone, error) {
	if _, err := s.repo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return nil, err
	}

	milestone, err := NewMilestone(req.ProjectID, req.Title, req.Description, req.DueDate)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateMilestone(ctx, milestone); err != nil {
		return nil, err
	}

	return milestone, nil
}

// ToggleMilestone updates milestone completion.
func (s *Service) ToggleMilestone(ctx context.Context, req ToggleMilestoneRequest) (*Milestone, error) {
	milestone, err := s.repo.GetMilestone(ctx, req.MilestoneID, req.UserID)
	if err != nil {
		return nil, err
	}

	milestone.MarkCompleted(req.Completed)

	if err := s.repo.UpdateMilestone(ctx, milestone); err != nil {
		return nil, err
	}

	return milestone, nil
}

// ListMilestones returns milestones for a project.
func (s *Service) ListMilestones(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]Milestone, error) {
	return s.repo.ListMilestones(ctx, projectID, userID)
}
