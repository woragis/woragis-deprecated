package problemsolutions

import (
	"context"

	"github.com/google/uuid"
)

// Service orchestrates problem solution workflows.
type Service interface {
	CreateProblemSolution(ctx context.Context, req CreateProblemSolutionRequest) (*ProblemSolution, error)
	UpdateProblemSolution(ctx context.Context, req UpdateProblemSolutionRequest) (*ProblemSolution, error)
	GetProblemSolution(ctx context.Context, problemSolutionID uuid.UUID, userID uuid.UUID) (*ProblemSolution, error)
	GetProblemSolutionPublic(ctx context.Context, problemSolutionID uuid.UUID) (*ProblemSolution, error)
	ListProblemSolutions(ctx context.Context, userID uuid.UUID) ([]ProblemSolution, error)
	ListFeaturedProblemSolutions(ctx context.Context) ([]ProblemSolution, error)
	DeleteProblemSolution(ctx context.Context, req DeleteProblemSolutionRequest) error
	GetProblemSolutionMatrix(ctx context.Context) ([]ProblemSolutionMatrixEntry, error)
}

type service struct {
	repo Repository
}

var _ Service = (*service)(nil)

// NewService constructs a Service.
func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// Request payloads

type CreateProblemSolutionRequest struct {
	UserID      uuid.UUID
	Problem     string
	Context     string
	Solution    string
	Technologies []string
	Impact      string
	Metrics     *MetricsData
	Featured    bool
}

type UpdateProblemSolutionRequest struct {
	ProblemSolutionID uuid.UUID
	UserID             uuid.UUID
	Problem            *string
	Context            *string
	Solution           *string
	Technologies       []string
	Impact             *string
	Metrics            *MetricsData
	Featured           *bool
}

type DeleteProblemSolutionRequest struct {
	ProblemSolutionID uuid.UUID
	UserID             uuid.UUID
}

// Service implementations

func (s *service) CreateProblemSolution(ctx context.Context, req CreateProblemSolutionRequest) (*ProblemSolution, error) {
	problemSolution, err := NewProblemSolution(req.UserID, req.Problem, req.Context, req.Solution)
	if err != nil {
		return nil, err
	}

	if len(req.Technologies) > 0 {
		problemSolution.SetTechnologies(req.Technologies)
	}

	if req.Impact != "" {
		problemSolution.Impact = req.Impact
	}

	if req.Metrics != nil {
		problemSolution.SetMetrics(req.Metrics)
	}

	problemSolution.SetFeatured(req.Featured)

	if err := s.repo.CreateProblemSolution(ctx, problemSolution); err != nil {
		return nil, err
	}

	return problemSolution, nil
}

func (s *service) UpdateProblemSolution(ctx context.Context, req UpdateProblemSolutionRequest) (*ProblemSolution, error) {
	problemSolution, err := s.repo.GetProblemSolution(ctx, req.ProblemSolutionID, req.UserID)
	if err != nil {
		return nil, err
	}

	problem := ""
	if req.Problem != nil {
		problem = *req.Problem
	}
	context := ""
	if req.Context != nil {
		context = *req.Context
	}
	solution := ""
	if req.Solution != nil {
		solution = *req.Solution
	}
	impact := ""
	if req.Impact != nil {
		impact = *req.Impact
	}

	if err := problemSolution.UpdateDetails(problem, context, solution, impact); err != nil {
		return nil, err
	}

	if req.Technologies != nil {
		problemSolution.SetTechnologies(req.Technologies)
	}

	if req.Metrics != nil {
		problemSolution.SetMetrics(req.Metrics)
	}

	if req.Featured != nil {
		problemSolution.SetFeatured(*req.Featured)
	}

	if err := s.repo.UpdateProblemSolution(ctx, problemSolution); err != nil {
		return nil, err
	}

	return problemSolution, nil
}

func (s *service) GetProblemSolution(ctx context.Context, problemSolutionID uuid.UUID, userID uuid.UUID) (*ProblemSolution, error) {
	return s.repo.GetProblemSolution(ctx, problemSolutionID, userID)
}

func (s *service) GetProblemSolutionPublic(ctx context.Context, problemSolutionID uuid.UUID) (*ProblemSolution, error) {
	return s.repo.GetProblemSolutionPublic(ctx, problemSolutionID)
}

func (s *service) ListProblemSolutions(ctx context.Context, userID uuid.UUID) ([]ProblemSolution, error) {
	return s.repo.ListProblemSolutions(ctx, userID)
}

func (s *service) ListFeaturedProblemSolutions(ctx context.Context) ([]ProblemSolution, error) {
	return s.repo.ListFeaturedProblemSolutions(ctx)
}

func (s *service) DeleteProblemSolution(ctx context.Context, req DeleteProblemSolutionRequest) error {
	return s.repo.DeleteProblemSolution(ctx, req.ProblemSolutionID, req.UserID)
}

func (s *service) GetProblemSolutionMatrix(ctx context.Context) ([]ProblemSolutionMatrixEntry, error) {
	return s.repo.GetProblemSolutionMatrix(ctx)
}

