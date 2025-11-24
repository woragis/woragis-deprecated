package systemdesigns

import (
	"context"

	"github.com/google/uuid"
)

// Service orchestrates system design workflows.
type Service interface {
	CreateSystemDesign(ctx context.Context, req CreateSystemDesignRequest) (*SystemDesign, error)
	UpdateSystemDesign(ctx context.Context, req UpdateSystemDesignRequest) (*SystemDesign, error)
	GetSystemDesign(ctx context.Context, systemDesignID uuid.UUID, userID uuid.UUID) (*SystemDesign, error)
	GetSystemDesignPublic(ctx context.Context, systemDesignID uuid.UUID) (*SystemDesign, error)
	ListSystemDesigns(ctx context.Context, userID uuid.UUID) ([]SystemDesign, error)
	ListFeaturedSystemDesigns(ctx context.Context) ([]SystemDesign, error)
	DeleteSystemDesign(ctx context.Context, req DeleteSystemDesignRequest) error
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

type CreateSystemDesignRequest struct {
	UserID      uuid.UUID
	Title       string
	Description string
	Components  *ComponentsData
	DataFlow    string
	Scalability string
	Reliability string
	Diagram     string
	Featured    bool
}

type UpdateSystemDesignRequest struct {
	SystemDesignID uuid.UUID
	UserID          uuid.UUID
	Title           *string
	Description     *string
	Components      *ComponentsData
	DataFlow        *string
	Scalability     *string
	Reliability     *string
	Diagram         *string
	Featured        *bool
}

type DeleteSystemDesignRequest struct {
	SystemDesignID uuid.UUID
	UserID          uuid.UUID
}

// Service implementations

func (s *service) CreateSystemDesign(ctx context.Context, req CreateSystemDesignRequest) (*SystemDesign, error) {
	systemDesign, err := NewSystemDesign(req.UserID, req.Title, req.Description)
	if err != nil {
		return nil, err
	}

	if req.Components != nil {
		systemDesign.SetComponents(req.Components)
	}

	if req.DataFlow != "" {
		systemDesign.DataFlow = req.DataFlow
	}

	if req.Scalability != "" {
		systemDesign.Scalability = req.Scalability
	}

	if req.Reliability != "" {
		systemDesign.Reliability = req.Reliability
	}

	if req.Diagram != "" {
		systemDesign.Diagram = req.Diagram
	}

	systemDesign.SetFeatured(req.Featured)

	if err := s.repo.CreateSystemDesign(ctx, systemDesign); err != nil {
		return nil, err
	}

	return systemDesign, nil
}

func (s *service) UpdateSystemDesign(ctx context.Context, req UpdateSystemDesignRequest) (*SystemDesign, error) {
	systemDesign, err := s.repo.GetSystemDesign(ctx, req.SystemDesignID, req.UserID)
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
	dataFlow := ""
	if req.DataFlow != nil {
		dataFlow = *req.DataFlow
	}
	scalability := ""
	if req.Scalability != nil {
		scalability = *req.Scalability
	}
	reliability := ""
	if req.Reliability != nil {
		reliability = *req.Reliability
	}
	diagram := ""
	if req.Diagram != nil {
		diagram = *req.Diagram
	}

	if err := systemDesign.UpdateDetails(title, description, dataFlow, scalability, reliability, diagram); err != nil {
		return nil, err
	}

	if req.Components != nil {
		systemDesign.SetComponents(req.Components)
	}

	if req.Featured != nil {
		systemDesign.SetFeatured(*req.Featured)
	}

	if err := s.repo.UpdateSystemDesign(ctx, systemDesign); err != nil {
		return nil, err
	}

	return systemDesign, nil
}

func (s *service) GetSystemDesign(ctx context.Context, systemDesignID uuid.UUID, userID uuid.UUID) (*SystemDesign, error) {
	return s.repo.GetSystemDesign(ctx, systemDesignID, userID)
}

func (s *service) GetSystemDesignPublic(ctx context.Context, systemDesignID uuid.UUID) (*SystemDesign, error) {
	return s.repo.GetSystemDesignPublic(ctx, systemDesignID)
}

func (s *service) ListSystemDesigns(ctx context.Context, userID uuid.UUID) ([]SystemDesign, error) {
	return s.repo.ListSystemDesigns(ctx, userID)
}

func (s *service) ListFeaturedSystemDesigns(ctx context.Context) ([]SystemDesign, error) {
	return s.repo.ListFeaturedSystemDesigns(ctx)
}

func (s *service) DeleteSystemDesign(ctx context.Context, req DeleteSystemDesignRequest) error {
	return s.repo.DeleteSystemDesign(ctx, req.SystemDesignID, req.UserID)
}

