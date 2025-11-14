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
	ActorID     uuid.UUID
	IdeaID      uuid.UUID
	Title       string
	Description string
	Color       string
	ProjectID   *uuid.UUID
}

// MoveIdeaRequest updates canvas coordinates.
type MoveIdeaRequest struct {
	ActorID uuid.UUID
	IdeaID  uuid.UUID
	PosX    float64
	PosY    float64
}

// CreateLinkRequest links two ideas.
type CreateLinkRequest struct {
	ActorID       uuid.UUID
	SourceIdeaID  uuid.UUID
	TargetIdeaID  uuid.UUID
	Relation      string
	Weight        float64
	Bidirectional bool
}

// BulkMoveIdeasRequest captures move operations.
type BulkMoveIdeasRequest struct {
	ActorID uuid.UUID
	Moves   []IdeaPositionUpdate
}

// BulkUpdateIdeasRequest captures metadata operations.
type BulkUpdateIdeasRequest struct {
	ActorID uuid.UUID
	Updates []IdeaDetailUpdate
}

// BulkIDsRequest captures list of idea ids.
type BulkIDsRequest struct {
	ActorID uuid.UUID
	IDs     []uuid.UUID
}

// ListIdeasRequest controls idea listing.
type ListIdeasRequest struct {
	ActorID uuid.UUID
	OwnerID uuid.UUID
}

// ListLinksRequest controls relationship listing.
type ListLinksRequest struct {
	ActorID       uuid.UUID
	OwnerID       uuid.UUID
	IdeaID        uuid.UUID
	Relation      string
	Search        string
	MinWeight     *float64
	MaxWeight     *float64
	Bidirectional *bool
}

// ListVersionsRequest controls version history.
type ListVersionsRequest struct {
	ActorID uuid.UUID
	IdeaID  uuid.UUID
	Limit   int
}

// CollaboratorRequest describes collaborator actions.
type CollaboratorRequest struct {
	ActorID        uuid.UUID
	OwnerID        uuid.UUID
	CollaboratorID uuid.UUID
	Role           string
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

	s.recordVersion(ctx, idea, req.UserID, ChangeTypeCreated)
	return idea, nil
}

// UpdateIdea updates metadata of an existing idea.
func (s *Service) UpdateIdea(ctx context.Context, req UpdateIdeaRequest) (*Idea, error) {
	idea, err := s.loadIdea(ctx, req.IdeaID, req.ActorID)
	if err != nil {
		return nil, err
	}

	if err := idea.UpdateDetails(req.Title, req.Description, req.Color); err != nil {
		return nil, err
	}

	if req.ProjectID != nil {
		idea.ProjectID = req.ProjectID
		idea.Touch()
	}

	if err := s.repo.UpdateIdea(ctx, idea); err != nil {
		return nil, err
	}

	s.recordVersion(ctx, idea, req.ActorID, ChangeTypeEdited)
	return idea, nil
}

// MoveIdea updates an idea position.
func (s *Service) MoveIdea(ctx context.Context, req MoveIdeaRequest) (*Idea, error) {
	idea, err := s.loadIdea(ctx, req.IdeaID, req.ActorID)
	if err != nil {
		return nil, err
	}

	idea.Move(req.PosX, req.PosY)

	if err := s.repo.UpdateIdea(ctx, idea); err != nil {
		return nil, err
	}

	s.recordVersion(ctx, idea, req.ActorID, ChangeTypeMoved)
	return idea, nil
}

// BulkMoveIdeas applies coordinate updates.
func (s *Service) BulkMoveIdeas(ctx context.Context, req BulkMoveIdeasRequest) error {
	if len(req.Moves) == 0 {
		return nil
	}

	var ownerID uuid.UUID
	for _, move := range req.Moves {
		idea, err := s.loadIdea(ctx, move.IdeaID, req.ActorID)
		if err != nil {
			return err
		}
		if ownerID == uuid.Nil {
			ownerID = idea.UserID
		}
		if ownerID != idea.UserID {
			return NewDomainError(ErrCodeInvalidPayload, "ideas: bulk move requires a single board owner")
		}
	}

	if err := s.repo.BulkMoveIdeas(ctx, ownerID, req.Moves); err != nil {
		return err
	}

	for _, move := range req.Moves {
		if idea, err := s.repo.GetIdeaByID(ctx, move.IdeaID); err == nil {
			s.recordVersion(ctx, idea, req.ActorID, ChangeTypeMoved)
		}
	}

	return nil
}

// BulkUpdateIdeas applies metadata updates.
func (s *Service) BulkUpdateIdeas(ctx context.Context, req BulkUpdateIdeasRequest) error {
	if len(req.Updates) == 0 {
		return nil
	}

	var ownerID uuid.UUID
	for _, upd := range req.Updates {
		idea, err := s.loadIdea(ctx, upd.IdeaID, req.ActorID)
		if err != nil {
			return err
		}
		if ownerID == uuid.Nil {
			ownerID = idea.UserID
		}
		if ownerID != idea.UserID {
			return NewDomainError(ErrCodeInvalidPayload, "ideas: bulk update requires a single board owner")
		}
	}

	if err := s.repo.BulkUpdateDetails(ctx, ownerID, req.Updates); err != nil {
		return err
	}

	for _, upd := range req.Updates {
		if idea, err := s.repo.GetIdeaByID(ctx, upd.IdeaID); err == nil {
			s.recordVersion(ctx, idea, req.ActorID, ChangeTypeEdited)
		}
	}

	return nil
}

// DeleteIdeas soft deletes ideas in bulk.
func (s *Service) DeleteIdeas(ctx context.Context, req BulkIDsRequest) error {
	if len(req.IDs) == 0 {
		return nil
	}

	var ownerID uuid.UUID
	for _, id := range req.IDs {
		idea, err := s.loadIdea(ctx, id, req.ActorID)
		if err != nil {
			return err
		}
		if ownerID == uuid.Nil {
			ownerID = idea.UserID
		}
		if ownerID != idea.UserID {
			return NewDomainError(ErrCodeInvalidPayload, "ideas: bulk delete requires a single board owner")
		}
	}

	return s.repo.DeleteIdeas(ctx, ownerID, req.IDs)
}

// RestoreIdeas clears soft delete flags.
func (s *Service) RestoreIdeas(ctx context.Context, req BulkIDsRequest) error {
	if len(req.IDs) == 0 {
		return nil
	}

	var ownerID uuid.UUID
	for _, id := range req.IDs {
		idea, err := s.repo.GetIdeaByID(ctx, id)
		if err != nil {
			return err
		}
		if ownerID == uuid.Nil {
			ownerID = idea.UserID
		}
		if ownerID != idea.UserID {
			return NewDomainError(ErrCodeInvalidPayload, "ideas: bulk restore requires a single board owner")
		}
		if err := s.ensureAccess(ctx, idea.UserID, req.ActorID); err != nil {
			return err
		}
	}

	return s.repo.RestoreIdeas(ctx, ownerID, req.IDs)
}

// CreateLink links two ideas.
func (s *Service) CreateLink(ctx context.Context, req CreateLinkRequest) (*IdeaLink, error) {
	source, err := s.loadIdea(ctx, req.SourceIdeaID, req.ActorID)
	if err != nil {
		return nil, err
	}
	target, err := s.loadIdea(ctx, req.TargetIdeaID, req.ActorID)
	if err != nil {
		return nil, err
	}
	if source.UserID != target.UserID {
		return nil, NewDomainError(ErrCodeInvalidRelation, "ideas: links require ideas from the same board")
	}

	link, err := NewIdeaLink(source.UserID, source.ID, target.ID, req.Relation, req.Weight, req.Bidirectional)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateLink(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

// ListIdeas returns ideas for an owner with access control.
func (s *Service) ListIdeas(ctx context.Context, req ListIdeasRequest) ([]Idea, error) {
	ownerID := req.OwnerID
	if ownerID == uuid.Nil {
		ownerID = req.ActorID
	}
	if err := s.ensureAccess(ctx, ownerID, req.ActorID); err != nil {
		return nil, err
	}
	return s.repo.ListIdeas(ctx, ownerID)
}

// ListLinks returns links with filters.
func (s *Service) ListLinks(ctx context.Context, req ListLinksRequest) ([]IdeaLink, error) {
	ownerID := req.OwnerID
	if ownerID == uuid.Nil {
		ownerID = req.ActorID
	}
	if err := s.ensureAccess(ctx, ownerID, req.ActorID); err != nil {
		return nil, err
	}

	filters := LinkFilters{
		UserID:           ownerID,
		IdeaID:           req.IdeaID,
		Relation:         req.Relation,
		RelationContains: req.Search,
		MinWeight:        req.MinWeight,
		MaxWeight:        req.MaxWeight,
		Bidirectional:    req.Bidirectional,
	}

	return s.repo.ListLinks(ctx, filters)
}

// ListVersions returns version history.
func (s *Service) ListVersions(ctx context.Context, req ListVersionsRequest) ([]IdeaVersion, error) {
	idea, err := s.loadIdea(ctx, req.IdeaID, req.ActorID)
	if err != nil {
		return nil, err
	}
	return s.repo.ListVersions(ctx, idea.ID, idea.UserID, req.Limit)
}

// AddCollaborator grants canvas access.
func (s *Service) AddCollaborator(ctx context.Context, req CollaboratorRequest) (*IdeaCollaborator, error) {
	if req.ActorID != req.OwnerID {
		if err := s.ensureAccess(ctx, req.OwnerID, req.ActorID); err != nil {
			return nil, err
		}
	}

	entry, err := NewIdeaCollaborator(req.OwnerID, req.CollaboratorID, req.Role)
	if err != nil {
		return nil, err
	}

	if err := s.repo.AddCollaborator(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// RemoveCollaborator revokes canvas access.
func (s *Service) RemoveCollaborator(ctx context.Context, req CollaboratorRequest) error {
	if req.ActorID != req.OwnerID {
		if err := s.ensureAccess(ctx, req.OwnerID, req.ActorID); err != nil {
			return err
		}
	}
	return s.repo.RemoveCollaborator(ctx, req.OwnerID, req.CollaboratorID)
}

// ListCollaborators returns collaborators for an owner.
func (s *Service) ListCollaborators(ctx context.Context, actorID, ownerID uuid.UUID) ([]IdeaCollaborator, error) {
	if ownerID == uuid.Nil {
		ownerID = actorID
	}
	if ownerID != actorID {
		if err := s.ensureAccess(ctx, ownerID, actorID); err != nil {
			return nil, err
		}
	}
	return s.repo.ListCollaborators(ctx, ownerID)
}

func (s *Service) loadIdea(ctx context.Context, ideaID, actorID uuid.UUID) (*Idea, error) {
	if ideaID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyIdeaID)
	}
	if actorID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	idea, err := s.repo.GetIdea(ctx, ideaID, actorID)
	if err == nil {
		return idea, nil
	}
	if domainErr, ok := AsDomainError(err); ok && domainErr.Code == ErrCodeNotFound {
		idea, fetchErr := s.repo.GetIdeaByID(ctx, ideaID)
		if fetchErr != nil {
			return nil, fetchErr
		}
		if err := s.ensureAccess(ctx, idea.UserID, actorID); err != nil {
			return nil, err
		}
		return idea, nil
	}
	return nil, err
}

func (s *Service) ensureAccess(ctx context.Context, ownerID, actorID uuid.UUID) error {
	if ownerID == uuid.Nil || actorID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if ownerID == actorID {
		return nil
	}
	ok, err := s.repo.HasCollaborator(ctx, ownerID, actorID)
	if err != nil {
		return err
	}
	if !ok {
		return NewDomainError(ErrCodeInvalidCollaborator, ErrCollaboratorUnauthorized)
	}
	return nil
}

func (s *Service) recordVersion(ctx context.Context, idea *Idea, editorID uuid.UUID, changeType string) {
	if idea == nil {
		return
	}
	version := NewIdeaVersion(idea, editorID, changeType)
	if err := s.repo.CreateVersion(ctx, version); err != nil && s.logger != nil {
		s.logger.Warn("ideas: failed to persist version history", slog.Any("error", err), slog.String("idea_id", idea.ID.String()))
	}
}
