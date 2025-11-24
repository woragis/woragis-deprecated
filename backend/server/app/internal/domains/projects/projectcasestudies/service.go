package projectcasestudies

import (
	"context"

	"github.com/google/uuid"
	projectsdomain "github.com/woragis/backend/server/app/internal/domains/projects"
)

// Service orchestrates project case study workflows.
type Service interface {
	CreateCaseStudy(ctx context.Context, req CreateCaseStudyRequest) (*ProjectCaseStudy, error)
	UpdateCaseStudy(ctx context.Context, req UpdateCaseStudyRequest) (*ProjectCaseStudy, error)
	GetCaseStudy(ctx context.Context, caseStudyID uuid.UUID, userID uuid.UUID) (*ProjectCaseStudy, error)
	GetCaseStudyByProjectID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*ProjectCaseStudy, error)
	GetPublicCaseStudyByProjectID(ctx context.Context, projectID uuid.UUID) (*ProjectCaseStudy, error)
	ListCaseStudies(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectCaseStudy, error)
	DeleteCaseStudy(ctx context.Context, req DeleteCaseStudyRequest) error
}

type service struct {
	repo        Repository
	projectRepo projectsdomain.Repository
}

var _ Service = (*service)(nil)

// NewService constructs a Service.
func NewService(repo Repository, projectRepo projectsdomain.Repository) Service {
	return &service{
		repo:        repo,
		projectRepo: projectRepo,
	}
}

// Request payloads

type CreateCaseStudyRequest struct {
	ProjectID     uuid.UUID
	UserID        uuid.UUID
	Title         string
	Description   string
	Challenge     string
	Solution      string
	Technologies  []string
	Architecture  string
	Metrics       *MetricsData
	Tradeoffs     *TradeoffsData
	LessonsLearned []string
}

type UpdateCaseStudyRequest struct {
	CaseStudyID   uuid.UUID
	UserID        uuid.UUID
	Title         *string
	Description   *string
	Challenge     *string
	Solution      *string
	Technologies  []string
	Architecture  *string
	Metrics       *MetricsData
	Tradeoffs     *TradeoffsData
	LessonsLearned []string
}

type DeleteCaseStudyRequest struct {
	CaseStudyID uuid.UUID
	UserID      uuid.UUID
}

// Service implementations

func (s *service) CreateCaseStudy(ctx context.Context, req CreateCaseStudyRequest) (*ProjectCaseStudy, error) {
	// Verify project exists and user has access
	if _, err := s.projectRepo.GetProject(ctx, req.ProjectID, req.UserID); err != nil {
		return nil, err
	}

	caseStudy, err := NewProjectCaseStudy(
		req.ProjectID,
		req.Title,
		req.Description,
		req.Challenge,
		req.Solution,
		req.Architecture,
	)
	if err != nil {
		return nil, err
	}

	if len(req.Technologies) > 0 {
		caseStudy.SetTechnologies(req.Technologies)
	}

	if req.Metrics != nil {
		caseStudy.SetMetrics(req.Metrics)
	}

	if req.Tradeoffs != nil {
		caseStudy.SetTradeoffs(req.Tradeoffs)
	}

	if len(req.LessonsLearned) > 0 {
		caseStudy.SetLessonsLearned(req.LessonsLearned)
	}

	if err := s.repo.CreateCaseStudy(ctx, caseStudy); err != nil {
		return nil, err
	}

	return caseStudy, nil
}

func (s *service) UpdateCaseStudy(ctx context.Context, req UpdateCaseStudyRequest) (*ProjectCaseStudy, error) {
	caseStudy, err := s.repo.GetCaseStudy(ctx, req.CaseStudyID, req.UserID)
	if err != nil {
		return nil, err
	}

	title := ""
	if req.Title != nil {
		title = *req.Title
	}
	description := ""
	if req.Description != nil {
		description = *req.Description
	}
	challenge := ""
	if req.Challenge != nil {
		challenge = *req.Challenge
	}
	solution := ""
	if req.Solution != nil {
		solution = *req.Solution
	}
	architecture := ""
	if req.Architecture != nil {
		architecture = *req.Architecture
	}

	if err := caseStudy.UpdateDetails(title, description, challenge, solution, architecture); err != nil {
		return nil, err
	}

	if req.Technologies != nil {
		caseStudy.SetTechnologies(req.Technologies)
	}

	if req.Metrics != nil {
		caseStudy.SetMetrics(req.Metrics)
	}

	if req.Tradeoffs != nil {
		caseStudy.SetTradeoffs(req.Tradeoffs)
	}

	if req.LessonsLearned != nil {
		caseStudy.SetLessonsLearned(req.LessonsLearned)
	}

	if err := s.repo.UpdateCaseStudy(ctx, caseStudy); err != nil {
		return nil, err
	}

	return caseStudy, nil
}

func (s *service) GetCaseStudy(ctx context.Context, caseStudyID uuid.UUID, userID uuid.UUID) (*ProjectCaseStudy, error) {
	return s.repo.GetCaseStudy(ctx, caseStudyID, userID)
}

func (s *service) GetCaseStudyByProjectID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*ProjectCaseStudy, error) {
	return s.repo.GetCaseStudyByProjectID(ctx, projectID, userID)
}

func (s *service) GetPublicCaseStudyByProjectID(ctx context.Context, projectID uuid.UUID) (*ProjectCaseStudy, error) {
	return s.repo.GetCaseStudyByProjectIDPublic(ctx, projectID)
}

func (s *service) ListCaseStudies(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectCaseStudy, error) {
	return s.repo.ListCaseStudies(ctx, projectID, userID)
}

func (s *service) DeleteCaseStudy(ctx context.Context, req DeleteCaseStudyRequest) error {
	return s.repo.DeleteCaseStudy(ctx, req.CaseStudyID, req.UserID)
}
