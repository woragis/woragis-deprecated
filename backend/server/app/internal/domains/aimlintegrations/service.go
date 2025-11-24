package aimlintegrations

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service orchestrates AI/ML integration workflows.
type Service interface {
	CreateAIMLIntegration(ctx context.Context, userID uuid.UUID, req CreateAIMLIntegrationRequest) (*AIMLIntegration, error)
	UpdateAIMLIntegration(ctx context.Context, userID, integrationID uuid.UUID, req UpdateAIMLIntegrationRequest) (*AIMLIntegration, error)
	GetAIMLIntegration(ctx context.Context, integrationID uuid.UUID, userID uuid.UUID) (*AIMLIntegration, error)
	GetAIMLIntegrationPublic(ctx context.Context, integrationID uuid.UUID) (*AIMLIntegration, error)
	ListAIMLIntegrations(ctx context.Context, filters ListAIMLIntegrationFilters) ([]AIMLIntegration, error)
	ListFeaturedAIMLIntegrations(ctx context.Context) ([]AIMLIntegration, error)
	GetIntegrationsByProject(ctx context.Context, projectID uuid.UUID) ([]AIMLIntegration, error)
	GetIntegrationsByType(ctx context.Context, integrationType IntegrationType) ([]AIMLIntegration, error)
	GetIntegrationsByFramework(ctx context.Context, framework Framework) ([]AIMLIntegration, error)
	DeleteAIMLIntegration(ctx context.Context, userID, integrationID uuid.UUID) error
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

type CreateAIMLIntegrationRequest struct {
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Type            IntegrationType `json:"type"`
	Framework       Framework       `json:"framework"`
	ModelName       string          `json:"modelName,omitempty"`
	ModelVersion    string          `json:"modelVersion,omitempty"`
	UseCase         string          `json:"useCase,omitempty"`
	Impact          string          `json:"impact,omitempty"`
	Technologies    []string        `json:"technologies,omitempty"`
	Architecture    string          `json:"architecture,omitempty"`
	Metrics         string          `json:"metrics,omitempty"`
	ProjectID       *string         `json:"projectId,omitempty"`
	CaseStudyID     *string         `json:"caseStudyId,omitempty"`
	Featured        bool            `json:"featured,omitempty"`
	DisplayOrder    int             `json:"displayOrder,omitempty"`
	DemoURL         string          `json:"demoUrl,omitempty"`
	DocumentationURL string          `json:"documentationUrl,omitempty"`
	GitHubURL       string          `json:"githubUrl,omitempty"`
}

type UpdateAIMLIntegrationRequest struct {
	Title           *string          `json:"title,omitempty"`
	Description     *string          `json:"description,omitempty"`
	Type            *IntegrationType  `json:"type,omitempty"`
	Framework       *Framework       `json:"framework,omitempty"`
	ModelName       *string          `json:"modelName,omitempty"`
	ModelVersion    *string          `json:"modelVersion,omitempty"`
	UseCase         *string          `json:"useCase,omitempty"`
	Impact          *string          `json:"impact,omitempty"`
	Technologies    []string         `json:"technologies,omitempty"`
	Architecture    *string          `json:"architecture,omitempty"`
	Metrics         *string          `json:"metrics,omitempty"`
	ProjectID      *string           `json:"projectId,omitempty"`
	CaseStudyID     *string          `json:"caseStudyId,omitempty"`
	Featured        *bool            `json:"featured,omitempty"`
	DisplayOrder    *int             `json:"displayOrder,omitempty"`
	DemoURL         *string          `json:"demoUrl,omitempty"`
	DocumentationURL *string         `json:"documentationUrl,omitempty"`
	GitHubURL       *string          `json:"githubUrl,omitempty"`
}

type ListAIMLIntegrationFilters struct {
	UserID      *uuid.UUID
	Type        *IntegrationType
	Framework   *Framework
	ProjectID   *uuid.UUID
	Featured    *bool
	Limit       int
	Offset      int
	OrderBy     string
	Order       string
}

// Service methods

func (s *service) CreateAIMLIntegration(ctx context.Context, userID uuid.UUID, req CreateAIMLIntegrationRequest) (*AIMLIntegration, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	integration, err := NewAIMLIntegration(userID, req.Title, req.Description, req.Type, req.Framework)
	if err != nil {
		return nil, err
	}

	if req.ModelName != "" {
		integration.SetModelInfo(req.ModelName, req.ModelVersion)
	}
	if req.UseCase != "" || req.Impact != "" || req.Architecture != "" || req.Metrics != "" {
		integration.UpdateDetails(req.Title, req.Description, req.UseCase, req.Impact, req.Architecture, req.Metrics)
	}
	if len(req.Technologies) > 0 {
		integration.SetTechnologies(req.Technologies)
	}
	if req.ProjectID != nil {
		projectID, err := uuid.Parse(*req.ProjectID)
		if err != nil {
			return nil, NewDomainError(ErrCodeInvalidPayload, "invalid project id format")
		}
		integration.SetProjectLink(projectID)
	}
	if req.CaseStudyID != nil {
		caseStudyID, err := uuid.Parse(*req.CaseStudyID)
		if err != nil {
			return nil, NewDomainError(ErrCodeInvalidPayload, "invalid case study id format")
		}
		integration.SetCaseStudyLink(caseStudyID)
	}
	integration.Featured = req.Featured
	integration.DisplayOrder = req.DisplayOrder
	if req.DemoURL != "" || req.DocumentationURL != "" || req.GitHubURL != "" {
		integration.SetURLs(req.DemoURL, req.DocumentationURL, req.GitHubURL)
	}

	if err := integration.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateAIMLIntegration(ctx, integration); err != nil {
		return nil, err
	}

	return integration, nil
}

func (s *service) UpdateAIMLIntegration(ctx context.Context, userID, integrationID uuid.UUID, req UpdateAIMLIntegrationRequest) (*AIMLIntegration, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if integrationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyIntegrationID)
	}

	// Get existing integration
	integration, err := s.repo.GetAIMLIntegration(ctx, integrationID, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	title := integration.Title
	if req.Title != nil {
		title = *req.Title
	}
	description := integration.Description
	if req.Description != nil {
		description = *req.Description
	}
	useCase := integration.UseCase
	if req.UseCase != nil {
		useCase = *req.UseCase
	}
	impact := integration.Impact
	if req.Impact != nil {
		impact = *req.Impact
	}
	architecture := integration.Architecture
	if req.Architecture != nil {
		architecture = *req.Architecture
	}
	metrics := integration.Metrics
	if req.Metrics != nil {
		metrics = *req.Metrics
	}

	integration.UpdateDetails(title, description, useCase, impact, architecture, metrics)

	if req.Type != nil {
		integration.Type = *req.Type
	}
	if req.Framework != nil {
		integration.Framework = *req.Framework
	}
	if req.ModelName != nil || req.ModelVersion != nil {
		modelName := integration.ModelName
		modelVersion := integration.ModelVersion
		if req.ModelName != nil {
			modelName = *req.ModelName
		}
		if req.ModelVersion != nil {
			modelVersion = *req.ModelVersion
		}
		integration.SetModelInfo(modelName, modelVersion)
	}
	if req.Technologies != nil {
		integration.SetTechnologies(req.Technologies)
	}
	if req.ProjectID != nil {
		if *req.ProjectID == "" {
			integration.ClearProjectLink()
		} else {
			projectID, err := uuid.Parse(*req.ProjectID)
			if err != nil {
				return nil, NewDomainError(ErrCodeInvalidPayload, "invalid project id format")
			}
			integration.SetProjectLink(projectID)
		}
	}
	if req.CaseStudyID != nil {
		if *req.CaseStudyID == "" {
			integration.ClearCaseStudyLink()
		} else {
			caseStudyID, err := uuid.Parse(*req.CaseStudyID)
			if err != nil {
				return nil, NewDomainError(ErrCodeInvalidPayload, "invalid case study id format")
			}
			integration.SetCaseStudyLink(caseStudyID)
		}
	}
	if req.Featured != nil {
		integration.SetFeatured(*req.Featured)
	}
	if req.DisplayOrder != nil {
		integration.SetDisplayOrder(*req.DisplayOrder)
	}
	if req.DemoURL != nil || req.DocumentationURL != nil || req.GitHubURL != nil {
		demoURL := integration.DemoURL
		docURL := integration.DocumentationURL
		githubURL := integration.GitHubURL
		if req.DemoURL != nil {
			demoURL = *req.DemoURL
		}
		if req.DocumentationURL != nil {
			docURL = *req.DocumentationURL
		}
		if req.GitHubURL != nil {
			githubURL = *req.GitHubURL
		}
		integration.SetURLs(demoURL, docURL, githubURL)
	}

	if err := integration.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateAIMLIntegration(ctx, integration); err != nil {
		return nil, err
	}

	return integration, nil
}

func (s *service) GetAIMLIntegration(ctx context.Context, integrationID uuid.UUID, userID uuid.UUID) (*AIMLIntegration, error) {
	if integrationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyIntegrationID)
	}

	return s.repo.GetAIMLIntegration(ctx, integrationID, userID)
}

func (s *service) GetAIMLIntegrationPublic(ctx context.Context, integrationID uuid.UUID) (*AIMLIntegration, error) {
	if integrationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyIntegrationID)
	}

	return s.repo.GetAIMLIntegrationPublic(ctx, integrationID)
}

func (s *service) ListAIMLIntegrations(ctx context.Context, filters ListAIMLIntegrationFilters) ([]AIMLIntegration, error) {
	repoFilters := AIMLIntegrationFilters{
		UserID:    filters.UserID,
		Type:      filters.Type,
		Framework: filters.Framework,
		ProjectID: filters.ProjectID,
		Featured:  filters.Featured,
		Limit:     filters.Limit,
		Offset:    filters.Offset,
		OrderBy:   filters.OrderBy,
		Order:     filters.Order,
	}

	return s.repo.ListAIMLIntegrations(ctx, repoFilters)
}

func (s *service) ListFeaturedAIMLIntegrations(ctx context.Context) ([]AIMLIntegration, error) {
	return s.repo.ListFeaturedAIMLIntegrations(ctx)
}

func (s *service) GetIntegrationsByProject(ctx context.Context, projectID uuid.UUID) ([]AIMLIntegration, error) {
	if projectID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, "project id cannot be empty")
	}

	return s.repo.GetIntegrationsByProject(ctx, projectID)
}

func (s *service) GetIntegrationsByType(ctx context.Context, integrationType IntegrationType) ([]AIMLIntegration, error) {
	return s.repo.GetIntegrationsByType(ctx, integrationType)
}

func (s *service) GetIntegrationsByFramework(ctx context.Context, framework Framework) ([]AIMLIntegration, error) {
	return s.repo.GetIntegrationsByFramework(ctx, framework)
}

func (s *service) DeleteAIMLIntegration(ctx context.Context, userID, integrationID uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if integrationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyIntegrationID)
	}

	return s.repo.DeleteAIMLIntegration(ctx, integrationID, userID)
}

