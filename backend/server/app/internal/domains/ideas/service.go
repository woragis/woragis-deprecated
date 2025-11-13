package ideas

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service orchestrates idea canvas operations.
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService constructs a new service.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// CreateIdeaRequest describes a new idea node.
type CreateIdeaRequest struct {
	UserID      uuid.UUID
	Title       string
	Description string
	PosX        float64
	PosY        float64
	Color       string
	ProjectID   *uuid.UUID
}

// UpdateIdeaRequest updates textual metadata.
type UpdateIdeaRequest struct {
	UserID      uuid.UUID
	IdeaID      uuid.UUID
	Title       string
	Description string
	Color       string
	ProjectID   *uuid.UUID
}

// MoveIdeaRequest updates canvas coordinates.
type MoveIdeaRequest struct {
	UserID uuid.UUID
	IdeaID uuid.UUID
	PosX   float64
	PosY   float64
}

// CreateLinkRequest links two ideas.
type CreateLinkRequest struct {
	UserID        uuid.UUID
	SourceIdeaID  uuid.UUID
	TargetIdeaID  uuid.UUID
	Relation      string
	Weight        float64
	Bidirectional bool
}

// CreateIdea creates a new idea node.
func (s *Service) CreateIdea(ctx context.Context, req CreateIdeaRequest) (*Idea, error) {
	idea, err := NewIdea(req.UserID, req.Title, req.Description, req.PosX, req.PosY, req.Color, req.ProjectID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateIdea(ctx, idea); err != nil {
		return nil, err
	}

	return idea, nil
}

// UpdateIdea updates metadata of an existing idea.
func (s *Service) UpdateIdea(ctx context.Context, req UpdateIdeaRequest) (*Idea, error) {
	idea, err := s.repo.GetIdea(ctx, req.IdeaID, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := idea.UpdateDetails(req.Title, req.Description, req.Color); err != nil {
		return nil, err
	}

	if req.ProjectID != nil {
		idea.ProjectID = req.ProjectID
	}

	if err := s.repo.UpdateIdea(ctx, idea); err != nil {
		return nil, err
	}

	return idea, nil
}

// MoveIdea updates an idea position.
func (s *Service) MoveIdea(ctx context.Context, req MoveIdeaRequest) (*Idea, error) {
	idea, err := s.repo.GetIdea(ctx, req.IdeaID, req.UserID)
	if err != nil {
		return nil, err
	}

	idea.Move(req.PosX, req.PosY)

	if err := s.repo.UpdateIdea(ctx, idea); err != nil {
		return nil, err
	}

	return idea, nil
}

// CreateLink links two ideas.
func (s *Service) CreateLink(ctx context.Context, req CreateLinkRequest) (*IdeaLink, error) {
	// ensure both ideas exist
	if _, err := s.repo.GetIdea(ctx, req.SourceIdeaID, req.UserID); err != nil {
		return nil, err
	}
	if _, err := s.repo.GetIdea(ctx, req.TargetIdeaID, req.UserID); err != nil {
		return nil, err
	}

	link, err := NewIdeaLink(req.UserID, req.SourceIdeaID, req.TargetIdeaID, req.Relation, req.Weight, req.Bidirectional)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateLink(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

// ListIdeas returns all ideas for a user.
func (s *Service) ListIdeas(ctx context.Context, userID uuid.UUID) ([]Idea, error) {
	return s.repo.ListIdeas(ctx, userID)
}

// ListLinks returns links for a user, optionally filtered by idea.
func (s *Service) ListLinks(ctx context.Context, userID uuid.UUID, ideaID uuid.UUID) ([]IdeaLink, error) {
	return s.repo.ListLinks(ctx, userID, ideaID)
}
