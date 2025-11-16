package ideas

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

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

// CreateIdeaNodeRequest describes a new node within an idea canvas.
type CreateIdeaNodeRequest struct {
	ActorID     uuid.UUID
	IdeaID      uuid.UUID
	Title       string
	Description string
	PosX        float64
	PosY        float64
	Width       float64
	Height      float64
	Color       string
	Type        string
}

// UpdateIdeaNodeRequest updates node metadata.
type UpdateIdeaNodeRequest struct {
	ActorID     uuid.UUID
	NodeID      uuid.UUID
	Title       string
	Description string
	Color       string
	Type        string
}

// MoveIdeaNodeRequest updates node position.
type MoveIdeaNodeRequest struct {
	ActorID uuid.UUID
	NodeID  uuid.UUID
	PosX    float64
	PosY    float64
}

// ResizeIdeaNodeRequest updates node dimensions.
type ResizeIdeaNodeRequest struct {
	ActorID uuid.UUID
	NodeID  uuid.UUID
	Width   float64
	Height  float64
}

// CreateIdeaNodeConnectionRequest creates a connection between nodes.
type CreateIdeaNodeConnectionRequest struct {
	ActorID      uuid.UUID
	IdeaID       uuid.UUID
	SourceNodeID uuid.UUID
	TargetNodeID uuid.UUID
	Direction    ConnectionDirection
	Label        string
}

// CreateDocumentRequest describes a new document for an idea.
type CreateDocumentRequest struct {
	ActorID  uuid.UUID
	IdeaID   uuid.UUID
	NodeID   *uuid.UUID
	Title    string
	Content  string
}

// UpdateDocumentRequest updates document metadata.
type UpdateDocumentRequest struct {
	ActorID  uuid.UUID
	DocumentID uuid.UUID
	Title    string
	Content  string
}

// CreateIdea creates a new idea node.
func (s *Service) CreateIdea(ctx context.Context, req CreateIdeaRequest) (*Idea, error) {
	idea, err := NewIdea(req.UserID, req.Title, req.Description, req.PosX, req.PosY, req.Color, req.ProjectID)
	if err != nil {
		return nil, err
	}

	if _, err := s.assignIdeaSlug(ctx, idea); err != nil {
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
	ideas, err := s.repo.ListIdeas(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	for idx := range ideas {
		s.ensureIdeaSlug(ctx, &ideas[idx])
	}
	return ideas, nil
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

// GetIdeaBySlug fetches an idea using its slug.
func (s *Service) GetIdeaBySlug(ctx context.Context, userID uuid.UUID, slug string) (*Idea, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if strings.TrimSpace(slug) == "" {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyIdeaSlug)
	}

	idea, err := s.repo.GetIdeaBySlug(ctx, slug, userID)
	if err != nil {
		return nil, err
	}

	s.ensureIdeaSlug(ctx, idea)
	return idea, nil
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
		s.ensureIdeaSlug(ctx, idea)
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
		s.ensureIdeaSlug(ctx, idea)
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

func (s *Service) assignIdeaSlug(ctx context.Context, idea *Idea) (bool, error) {
	if idea == nil {
		return false, NewDomainError(ErrCodeInvalidPayload, ErrNilIdea)
	}

	base := strings.TrimSpace(idea.Slug)
	if base == "" {
		base = generateIdeaSlug(idea.Title)
	}

	slug := base
	for attempt := 0; attempt < 50; attempt++ {
		taken, err := s.repo.IsIdeaSlugTaken(ctx, idea.UserID, slug, idea.ID)
		if err != nil {
			return false, err
		}
		if !taken {
			changed := idea.Slug != slug
			idea.Slug = slug
			return changed, nil
		}
		slug = fmt.Sprintf("%s-%d", base, attempt+2)
	}

	return false, NewDomainError(ErrCodeRepositoryFailure, "ideas: unable to generate unique slug")
}

func (s *Service) ensureIdeaSlug(ctx context.Context, idea *Idea) {
	if idea == nil {
		return
	}
	updated, err := s.assignIdeaSlug(ctx, idea)
	if err != nil {
		if s.logger != nil {
			s.logger.Warn("ideas: unable to ensure slug", slog.Any("error", err), slog.String("idea_id", idea.ID.String()))
		}
		return
	}
	if !updated {
		return
	}
	if err := s.repo.UpdateIdea(ctx, idea); err != nil && s.logger != nil {
		s.logger.Warn("ideas: unable to backfill slug", slog.Any("error", err), slog.String("idea_id", idea.ID.String()))
	}
}

// CreateIdeaNode creates a new node within an idea's canvas.
func (s *Service) CreateIdeaNode(ctx context.Context, req CreateIdeaNodeRequest) (*IdeaNode, error) {
	// Verify the idea exists and actor has access
	idea, err := s.loadIdea(ctx, req.IdeaID, req.ActorID)
	if err != nil {
		return nil, err
	}

	node, err := NewIdeaNode(idea.ID, req.Title, req.Description, req.PosX, req.PosY, req.Width, req.Height, req.Color, req.Type)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateIdeaNode(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

// UpdateIdeaNode updates node metadata.
func (s *Service) UpdateIdeaNode(ctx context.Context, req UpdateIdeaNodeRequest) (*IdeaNode, error) {
	node, err := s.loadIdeaNode(ctx, req.NodeID, req.ActorID)
	if err != nil {
		return nil, err
	}

	if err := node.UpdateDetails(req.Title, req.Description, req.Color, req.Type); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateIdeaNode(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

// MoveIdeaNode updates node position.
func (s *Service) MoveIdeaNode(ctx context.Context, req MoveIdeaNodeRequest) (*IdeaNode, error) {
	node, err := s.loadIdeaNode(ctx, req.NodeID, req.ActorID)
	if err != nil {
		return nil, err
	}

	node.Move(req.PosX, req.PosY)

	if err := s.repo.UpdateIdeaNode(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

// ResizeIdeaNode updates node dimensions.
func (s *Service) ResizeIdeaNode(ctx context.Context, req ResizeIdeaNodeRequest) (*IdeaNode, error) {
	node, err := s.loadIdeaNode(ctx, req.NodeID, req.ActorID)
	if err != nil {
		return nil, err
	}

	node.Resize(req.Width, req.Height)

	if err := s.repo.UpdateIdeaNode(ctx, node); err != nil {
		return nil, err
	}

	return node, nil
}

// ListIdeaNodes returns all nodes for an idea.
func (s *Service) ListIdeaNodes(ctx context.Context, actorID, ideaID uuid.UUID) ([]IdeaNode, error) {
	// Verify the idea exists and actor has access
	_, err := s.loadIdea(ctx, ideaID, actorID)
	if err != nil {
		return nil, err
	}

	nodes, err := s.repo.ListIdeaNodes(ctx, ideaID)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

// DeleteIdeaNode removes a node from an idea's canvas.
func (s *Service) DeleteIdeaNode(ctx context.Context, actorID, nodeID uuid.UUID) error {
	node, err := s.loadIdeaNode(ctx, nodeID, actorID)
	if err != nil {
		return err
	}

	return s.repo.DeleteIdeaNode(ctx, node.ID)
}

// CreateIdeaNodeConnection creates a connection between two nodes.
func (s *Service) CreateIdeaNodeConnection(ctx context.Context, req CreateIdeaNodeConnectionRequest) (*IdeaNodeConnection, error) {
	// Verify the idea exists and actor has access
	idea, err := s.loadIdea(ctx, req.IdeaID, req.ActorID)
	if err != nil {
		return nil, err
	}

	// Verify both nodes exist and belong to this idea
	sourceNode, err := s.repo.GetIdeaNode(ctx, req.SourceNodeID)
	if err != nil {
		return nil, err
	}
	if sourceNode.IdeaID != idea.ID {
		return nil, NewDomainError(ErrCodeInvalidRelation, "ideas: source node does not belong to this idea")
	}

	targetNode, err := s.repo.GetIdeaNode(ctx, req.TargetNodeID)
	if err != nil {
		return nil, err
	}
	if targetNode.IdeaID != idea.ID {
		return nil, NewDomainError(ErrCodeInvalidRelation, "ideas: target node does not belong to this idea")
	}

	conn, err := NewIdeaNodeConnection(idea.ID, sourceNode.ID, targetNode.ID, req.Direction, req.Label)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateIdeaNodeConnection(ctx, conn); err != nil {
		return nil, err
	}

	return conn, nil
}

// ListIdeaNodeConnections returns all connections for an idea's canvas.
func (s *Service) ListIdeaNodeConnections(ctx context.Context, actorID, ideaID uuid.UUID) ([]IdeaNodeConnection, error) {
	// Verify the idea exists and actor has access
	_, err := s.loadIdea(ctx, ideaID, actorID)
	if err != nil {
		return nil, err
	}

	connections, err := s.repo.ListIdeaNodeConnections(ctx, ideaID)
	if err != nil {
		return nil, err
	}

	return connections, nil
}

// DeleteIdeaNodeConnection removes a connection.
func (s *Service) DeleteIdeaNodeConnection(ctx context.Context, actorID, connID uuid.UUID) error {
	conn, err := s.repo.GetIdeaNodeConnection(ctx, connID)
	if err != nil {
		return err
	}

	// Verify the idea exists and actor has access
	_, err = s.loadIdea(ctx, conn.IdeaID, actorID)
	if err != nil {
		return err
	}

	return s.repo.DeleteIdeaNodeConnection(ctx, conn.ID)
}

func (s *Service) loadIdeaNode(ctx context.Context, nodeID, actorID uuid.UUID) (*IdeaNode, error) {
	if nodeID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, "ideas: node id cannot be empty")
	}
	if actorID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	node, err := s.repo.GetIdeaNode(ctx, nodeID)
	if err != nil {
		return nil, err
	}

	// Verify access through the parent idea
	_, err = s.loadIdea(ctx, node.IdeaID, actorID)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// Document service methods

// CreateDocument creates a new document for an idea.
func (s *Service) CreateDocument(ctx context.Context, req CreateDocumentRequest) (*Document, error) {
	// Verify the idea exists and actor has access
	idea, err := s.loadIdea(ctx, req.IdeaID, req.ActorID)
	if err != nil {
		return nil, err
	}

	// If nodeID is provided, verify the node exists and belongs to this idea
	if req.NodeID != nil {
		node, err := s.repo.GetIdeaNode(ctx, *req.NodeID)
		if err != nil {
			return nil, err
		}
		if node.IdeaID != idea.ID {
			return nil, NewDomainError(ErrCodeInvalidRelation, "ideas: node does not belong to this idea")
		}
	}

	doc, err := NewDocument(idea.ID, req.NodeID, req.Title, req.Content)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateDocument(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// UpdateDocument updates document content.
func (s *Service) UpdateDocument(ctx context.Context, req UpdateDocumentRequest) (*Document, error) {
	doc, err := s.loadDocument(ctx, req.DocumentID, req.ActorID)
	if err != nil {
		return nil, err
	}

	if err := doc.UpdateDetails(req.Title, req.Content); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateDocument(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// ListDocuments returns documents for an idea, optionally filtered by node.
func (s *Service) ListDocuments(ctx context.Context, actorID, ideaID uuid.UUID, nodeID *uuid.UUID) ([]Document, error) {
	// Verify the idea exists and actor has access
	idea, err := s.loadIdea(ctx, ideaID, actorID)
	if err != nil {
		return nil, err
	}

	// If nodeID is provided, verify it belongs to this idea
	if nodeID != nil {
		node, err := s.repo.GetIdeaNode(ctx, *nodeID)
		if err != nil {
			return nil, err
		}
		if node.IdeaID != idea.ID {
			return nil, NewDomainError(ErrCodeInvalidRelation, "ideas: node does not belong to this idea")
		}
	}

	docs, err := s.repo.ListDocuments(ctx, ideaID, nodeID)
	if err != nil {
		return nil, err
	}

	return docs, nil
}

// DeleteDocument removes a document.
func (s *Service) DeleteDocument(ctx context.Context, actorID, docID uuid.UUID) error {
	doc, err := s.loadDocument(ctx, docID, actorID)
	if err != nil {
		return err
	}

	return s.repo.DeleteDocument(ctx, doc.ID)
}

func (s *Service) loadDocument(ctx context.Context, docID, actorID uuid.UUID) (*Document, error) {
	if docID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, "ideas: document id cannot be empty")
	}
	if actorID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	doc, err := s.repo.GetDocument(ctx, docID)
	if err != nil {
		return nil, err
	}

	// Verify access through the parent idea
	_, err = s.loadIdea(ctx, doc.IdeaID, actorID)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
