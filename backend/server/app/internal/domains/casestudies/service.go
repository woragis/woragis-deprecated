package casestudies

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates case study workflows.
type Service interface {
	CreateCaseStudy(ctx context.Context, userID uuid.UUID, req CreateCaseStudyRequest) (*CaseStudy, error)
	UpdateCaseStudy(ctx context.Context, userID, caseStudyID uuid.UUID, req UpdateCaseStudyRequest) (*CaseStudy, error)
	GetCaseStudy(ctx context.Context, caseStudyID uuid.UUID) (*CaseStudy, error)
	GetCaseStudyByProjectSlug(ctx context.Context, projectSlug string) (*CaseStudy, error)
	GetCaseStudyByProjectID(ctx context.Context, projectID uuid.UUID) (*CaseStudy, error)
	ListCaseStudies(ctx context.Context, filters ListCaseStudiesFilters) ([]CaseStudy, error)
	DeleteCaseStudy(ctx context.Context, userID, caseStudyID uuid.UUID) error
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

var _ Service = (*service)(nil)

// NewService constructs a Service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// Request payloads

type CreateCaseStudyRequest struct {
	ProjectID      uuid.UUID           `json:"projectId"`
	ProjectSlug    string              `json:"projectSlug"`
	Title          string              `json:"title"`
	Problem        string              `json:"problem"`
	Context        string              `json:"context"`
	Solution       string              `json:"solution"`
	Approach       []string            `json:"approach,omitempty"`
	Architecture   *ArchitectureData   `json:"architecture,omitempty"`
	Metrics        *MetricsData        `json:"metrics,omitempty"`
	LessonsLearned []string            `json:"lessonsLearned,omitempty"`
	Technologies   []string            `json:"technologies,omitempty"`
	Featured       bool                `json:"featured,omitempty"`
}

type UpdateCaseStudyRequest struct {
	Title          *string             `json:"title,omitempty"`
	Problem        *string             `json:"problem,omitempty"`
	Context        *string             `json:"context,omitempty"`
	Solution       *string             `json:"solution,omitempty"`
	Approach       []string            `json:"approach,omitempty"`
	Architecture   *ArchitectureData   `json:"architecture,omitempty"`
	Metrics        *MetricsData        `json:"metrics,omitempty"`
	LessonsLearned []string            `json:"lessonsLearned,omitempty"`
	Technologies   []string            `json:"technologies,omitempty"`
	Featured       *bool               `json:"featured,omitempty"`
}

type ListCaseStudiesFilters struct {
	UserID      *uuid.UUID
	ProjectID   *uuid.UUID
	ProjectSlug *string
	Featured    *bool
	Limit       int
	Offset      int
	OrderBy     string
	Order       string
}

// Service methods

func (s *service) CreateCaseStudy(ctx context.Context, userID uuid.UUID, req CreateCaseStudyRequest) (*CaseStudy, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	caseStudy := NewCaseStudy(userID, req.ProjectID, req.ProjectSlug, req.Title, req.Problem, req.Context, req.Solution)
	
	if req.Approach != nil {
		caseStudy.Approach = JSONArray(req.Approach)
	}
	if req.Architecture != nil {
		caseStudy.Architecture = req.Architecture
	}
	if req.Metrics != nil {
		caseStudy.Metrics = req.Metrics
	}
	if req.LessonsLearned != nil {
		caseStudy.LessonsLearned = JSONArray(req.LessonsLearned)
	}
	if req.Technologies != nil {
		caseStudy.Technologies = JSONArray(req.Technologies)
	}
	caseStudy.Featured = req.Featured

	if err := caseStudy.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateCaseStudy(ctx, caseStudy); err != nil {
		return nil, err
	}

	return caseStudy, nil
}

func (s *service) UpdateCaseStudy(ctx context.Context, userID, caseStudyID uuid.UUID, req UpdateCaseStudyRequest) (*CaseStudy, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if caseStudyID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyCaseStudyID)
	}

	// Get existing case study
	caseStudy, err := s.repo.GetCaseStudy(ctx, caseStudyID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if caseStudy.UserID != userID {
		return nil, NewDomainError(ErrCodeUnauthorized, ErrUnauthorized)
	}

	// Update fields
	if req.Title != nil {
		caseStudy.Title = *req.Title
	}
	if req.Problem != nil {
		caseStudy.Problem = *req.Problem
	}
	if req.Context != nil {
		caseStudy.Context = *req.Context
	}
	if req.Solution != nil {
		caseStudy.Solution = *req.Solution
	}
	if req.Approach != nil {
		caseStudy.Approach = JSONArray(req.Approach)
	}
	if req.Architecture != nil {
		caseStudy.Architecture = req.Architecture
	}
	if req.Metrics != nil {
		caseStudy.Metrics = req.Metrics
	}
	if req.LessonsLearned != nil {
		caseStudy.LessonsLearned = JSONArray(req.LessonsLearned)
	}
	if req.Technologies != nil {
		caseStudy.Technologies = JSONArray(req.Technologies)
	}
	if req.Featured != nil {
		caseStudy.Featured = *req.Featured
	}

	caseStudy.UpdatedAt = time.Now().UTC()

	if err := caseStudy.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateCaseStudy(ctx, caseStudy); err != nil {
		return nil, err
	}

	return caseStudy, nil
}

func (s *service) GetCaseStudy(ctx context.Context, caseStudyID uuid.UUID) (*CaseStudy, error) {
	if caseStudyID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyCaseStudyID)
	}

	return s.repo.GetCaseStudy(ctx, caseStudyID)
}

func (s *service) GetCaseStudyByProjectSlug(ctx context.Context, projectSlug string) (*CaseStudy, error) {
	if projectSlug == "" {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectSlug)
	}

	return s.repo.GetCaseStudyByProjectSlug(ctx, projectSlug)
}

func (s *service) GetCaseStudyByProjectID(ctx context.Context, projectID uuid.UUID) (*CaseStudy, error) {
	if projectID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}

	return s.repo.GetCaseStudyByProjectID(ctx, projectID)
}

func (s *service) ListCaseStudies(ctx context.Context, filters ListCaseStudiesFilters) ([]CaseStudy, error) {
	repoFilters := CaseStudyFilters{
		UserID:      filters.UserID,
		ProjectID:   filters.ProjectID,
		ProjectSlug: filters.ProjectSlug,
		Featured:    filters.Featured,
		Limit:       filters.Limit,
		Offset:      filters.Offset,
		OrderBy:     filters.OrderBy,
		Order:       filters.Order,
	}

	return s.repo.ListCaseStudies(ctx, repoFilters)
}

func (s *service) DeleteCaseStudy(ctx context.Context, userID, caseStudyID uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if caseStudyID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCaseStudyID)
	}

	// Verify ownership
	caseStudy, err := s.repo.GetCaseStudy(ctx, caseStudyID)
	if err != nil {
		return err
	}

	if caseStudy.UserID != userID {
		return NewDomainError(ErrCodeUnauthorized, ErrUnauthorized)
	}

	return s.repo.DeleteCaseStudy(ctx, caseStudyID)
}

