package ideas

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes HTTP endpoints for ideas canvas.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler constructs a new handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type createIdeaPayload struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PosX        float64 `json:"pos_x"`
	PosY        float64 `json:"pos_y"`
	Color       string  `json:"color"`
	ProjectID   string  `json:"project_id"`
}

type updateIdeaPayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
	ProjectID   string `json:"project_id"`
}

type moveIdeaPayload struct {
	PosX float64 `json:"pos_x"`
	PosY float64 `json:"pos_y"`
}

type bulkMoveItem struct {
	IdeaID string  `json:"idea_id"`
	PosX   float64 `json:"pos_x"`
	PosY   float64 `json:"pos_y"`
}

type bulkMovePayload struct {
	Items []bulkMoveItem `json:"items"`
}

type bulkUpdateItem struct {
	IdeaID      string `json:"idea_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
	ProjectID   string `json:"project_id"`
}

type bulkUpdatePayload struct {
	Items []bulkUpdateItem `json:"items"`
}

type bulkIDsPayload struct {
	IdeaIDs []string `json:"idea_ids"`
}

type createLinkPayload struct {
	SourceIdeaID  string  `json:"source_idea_id"`
	TargetIdeaID  string  `json:"target_idea_id"`
	Relation      string  `json:"relation"`
	Weight        float64 `json:"weight"`
	Bidirectional bool    `json:"bidirectional"`
}

type collaboratorPayload struct {
	OwnerID        string `json:"owner_id"`
	CollaboratorID string `json:"collaborator_id"`
	Role           string `json:"role"`
}

// PostIdea handles idea creation.
func (h *Handler) PostIdea(c *fiber.Ctx) error {
	var payload createIdeaPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateCreateIdeaPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	var projectID *uuid.UUID
	if payload.ProjectID != "" {
		id, err := uuid.Parse(payload.ProjectID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		projectID = &id
	}

	idea, err := h.service.CreateIdea(c.Context(), CreateIdeaRequest{
		UserID:      userID,
		Title:       payload.Title,
		Description: payload.Description,
		PosX:        payload.PosX,
		PosY:        payload.PosY,
		Color:       payload.Color,
		ProjectID:   projectID,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, idea)
}

// PatchIdea handles metadata updates.
func (h *Handler) PatchIdea(c *fiber.Ctx) error {
	ideaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateIdeaPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateUpdateIdeaPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	var projectID *uuid.UUID
	if payload.ProjectID != "" {
		id, err := uuid.Parse(payload.ProjectID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		projectID = &id
	}

	idea, err := h.service.UpdateIdea(c.Context(), UpdateIdeaRequest{
		ActorID:     actorID,
		IdeaID:      ideaID,
		Title:       payload.Title,
		Description: payload.Description,
		Color:       payload.Color,
		ProjectID:   projectID,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, idea)
}

// PatchIdeaPosition handles moving idea nodes.
func (h *Handler) PatchIdeaPosition(c *fiber.Ctx) error {
	ideaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload moveIdeaPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateMoveIdeaPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	idea, err := h.service.MoveIdea(c.Context(), MoveIdeaRequest{
		ActorID: actorID,
		IdeaID:  ideaID,
		PosX:    payload.PosX,
		PosY:    payload.PosY,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, idea)
}

// PostBulkMove handles bulk movement of ideas.
func (h *Handler) PostBulkMove(c *fiber.Ctx) error {
	var payload bulkMovePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateBulkMovePayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	moves := make([]IdeaPositionUpdate, 0, len(payload.Items))
	for _, item := range payload.Items {
		id, err := uuid.Parse(item.IdeaID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		moves = append(moves, IdeaPositionUpdate{
			IdeaID: id,
			PosX:   item.PosX,
			PosY:   item.PosY,
		})
	}

	if err := h.service.BulkMoveIdeas(c.Context(), BulkMoveIdeasRequest{
		ActorID: actorID,
		Moves:   moves,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "moved"})
}

// PostBulkUpdate handles bulk metadata updates.
func (h *Handler) PostBulkUpdate(c *fiber.Ctx) error {
	var payload bulkUpdatePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateBulkUpdatePayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	updates := make([]IdeaDetailUpdate, 0, len(payload.Items))
	for _, item := range payload.Items {
		id, err := uuid.Parse(item.IdeaID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		var projectID *uuid.UUID
		if item.ProjectID != "" {
			pid, err := uuid.Parse(item.ProjectID)
			if err != nil {
				return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
			}
			projectID = &pid
		}
		updates = append(updates, IdeaDetailUpdate{
			IdeaID:      id,
			Title:       item.Title,
			Description: item.Description,
			Color:       item.Color,
			ProjectID:   projectID,
		})
	}

	if err := h.service.BulkUpdateIdeas(c.Context(), BulkUpdateIdeasRequest{
		ActorID: actorID,
		Updates: updates,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "updated"})
}

// PostBulkDelete handles soft deletion of ideas.
func (h *Handler) PostBulkDelete(c *fiber.Ctx) error {
	var payload bulkIDsPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateBulkIDsPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}
	ids, err := parseUUIDs(payload.IdeaIDs)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeleteIdeas(c.Context(), BulkIDsRequest{
		ActorID: actorID,
		IDs:     ids,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "deleted"})
}

// PostBulkRestore handles restoring soft-deleted ideas.
func (h *Handler) PostBulkRestore(c *fiber.Ctx) error {
	var payload bulkIDsPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateBulkIDsPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}
	ids, err := parseUUIDs(payload.IdeaIDs)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.RestoreIdeas(c.Context(), BulkIDsRequest{
		ActorID: actorID,
		IDs:     ids,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "restored"})
}

// PostLink creates a relationship between ideas.
func (h *Handler) PostLink(c *fiber.Ctx) error {
	var payload createLinkPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateCreateLinkPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	sourceID, err := uuid.Parse(payload.SourceIdeaID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}
	targetID, err := uuid.Parse(payload.TargetIdeaID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	link, err := h.service.CreateLink(c.Context(), CreateLinkRequest{
		ActorID:       actorID,
		SourceIdeaID:  sourceID,
		TargetIdeaID:  targetID,
		Relation:      payload.Relation,
		Weight:        payload.Weight,
		Bidirectional: payload.Bidirectional,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, link)
}

// ListIdeas returns ideas filtered by owner.
func (h *Handler) ListIdeas(c *fiber.Ctx) error {
	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	var ownerID uuid.UUID
	if ownerParam := c.Query("owner_id"); ownerParam != "" {
		ownerID, err = uuid.Parse(ownerParam)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	ideas, err := h.service.ListIdeas(c.Context(), ListIdeasRequest{
		ActorID: actorID,
		OwnerID: ownerID,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, ideas)
}

// GetIdeaBySlug returns an idea resolved by slug.
func (h *Handler) GetIdeaBySlug(c *fiber.Ctx) error {
	slug := strings.TrimSpace(c.Params("slug"))
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	idea, err := h.service.GetIdeaBySlug(c.Context(), actorID, slug)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, idea)
}

// GetIdeaVersions returns version history for an idea.
func (h *Handler) GetIdeaVersions(c *fiber.Ctx) error {
	ideaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	limit := 20
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsed, err := strconv.Atoi(limitParam); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	versions, err := h.service.ListVersions(c.Context(), ListVersionsRequest{
		ActorID: actorID,
		IdeaID:  ideaID,
		Limit:   limit,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, versions)
}

// ListLinks returns links for a user, optionally filtered.
func (h *Handler) ListLinks(c *fiber.Ctx) error {
	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	var ownerID uuid.UUID
	if ownerParam := c.Query("owner_id"); ownerParam != "" {
		ownerID, err = uuid.Parse(ownerParam)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	var ideaID uuid.UUID
	if ideaParam := c.Query("idea_id"); ideaParam != "" {
		ideaID, err = uuid.Parse(ideaParam)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	var minWeight *float64
	if minParam := c.Query("min_weight"); minParam != "" {
		if parsed, err := strconv.ParseFloat(minParam, 64); err == nil {
			minWeight = &parsed
		} else {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	var maxWeight *float64
	if maxParam := c.Query("max_weight"); maxParam != "" {
		if parsed, err := strconv.ParseFloat(maxParam, 64); err == nil {
			maxWeight = &parsed
		} else {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	var bidirectional *bool
	if bidParam := c.Query("bidirectional"); bidParam != "" {
		if parsed, err := strconv.ParseBool(bidParam); err == nil {
			bidirectional = &parsed
		} else {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	links, err := h.service.ListLinks(c.Context(), ListLinksRequest{
		ActorID:       actorID,
		OwnerID:       ownerID,
		IdeaID:        ideaID,
		Relation:      c.Query("relation"),
		Search:        c.Query("search"),
		MinWeight:     minWeight,
		MaxWeight:     maxWeight,
		Bidirectional: bidirectional,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, links)
}

// PostCollaborator grants access to a collaborator.
func (h *Handler) PostCollaborator(c *fiber.Ctx) error {
	var payload collaboratorPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateCollaboratorPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	ownerID := actorID
	if payload.OwnerID != "" {
		ownerID, err = uuid.Parse(payload.OwnerID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	collaboratorID, err := uuid.Parse(payload.CollaboratorID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	entry, err := h.service.AddCollaborator(c.Context(), CollaboratorRequest{
		ActorID:        actorID,
		OwnerID:        ownerID,
		CollaboratorID: collaboratorID,
		Role:           payload.Role,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, entry)
}

// DeleteCollaborator revokes collaborator access.
func (h *Handler) DeleteCollaborator(c *fiber.Ctx) error {
	collabID, err := uuid.Parse(c.Params("collaborator_id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	ownerID := actorID
	if ownerParam := c.Query("owner_id"); ownerParam != "" {
		ownerID, err = uuid.Parse(ownerParam)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	if err := h.service.RemoveCollaborator(c.Context(), CollaboratorRequest{
		ActorID:        actorID,
		OwnerID:        ownerID,
		CollaboratorID: collabID,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "removed"})
}

// ListCollaborators returns collaborators for a board owner.
func (h *Handler) ListCollaborators(c *fiber.Ctx) error {
	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	ownerID := actorID
	if ownerParam := c.Query("owner_id"); ownerParam != "" {
		ownerID, err = uuid.Parse(ownerParam)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	collaborators, err := h.service.ListCollaborators(c.Context(), actorID, ownerID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, collaborators)
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromError(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, fiber.Map{"message": domainErr.Message})
	}

	h.logError("ideas: unexpected error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func statusFromError(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidTitle, ErrCodeInvalidRelation, ErrCodeInvalidCollaborator, ErrCodeCollaboratorConflict:
		return fiber.StatusBadRequest
	case ErrCodeNotFound:
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
	}
}

func (h *Handler) logWarn(message string) {
	if h.logger != nil {
		h.logger.Warn(message)
	}
}

func (h *Handler) logError(message string, err error) {
	if h.logger != nil {
		h.logger.Error(message, slog.Any("error", err))
	}
}

func parseUUIDs(values []string) ([]uuid.UUID, error) {
	result := make([]uuid.UUID, 0, len(values))
	for _, raw := range values {
		id, err := uuid.Parse(raw)
		if err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, nil
}

// IdeaNode handlers

type createIdeaNodePayload struct {
	IdeaID      string  `json:"idea_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PosX        float64 `json:"pos_x"`
	PosY        float64 `json:"pos_y"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Color       string  `json:"color"`
	Type        string  `json:"type"`
}

type updateIdeaNodePayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Type        string `json:"type"`
}

type moveIdeaNodePayload struct {
	PosX float64 `json:"pos_x"`
	PosY float64 `json:"pos_y"`
}

type resizeIdeaNodePayload struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type createIdeaNodeConnectionPayload struct {
	IdeaID       string `json:"idea_id"`
	SourceNodeID string `json:"source_node_id"`
	TargetNodeID string `json:"target_node_id"`
	Direction    string `json:"direction"`
	Label        string `json:"label"`
}

// PostIdeaNode handles node creation within an idea canvas.
func (h *Handler) PostIdeaNode(c *fiber.Ctx) error {
	var payload createIdeaNodePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	ideaID, err := uuid.Parse(payload.IdeaID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	node, err := h.service.CreateIdeaNode(c.Context(), CreateIdeaNodeRequest{
		ActorID:     actorID,
		IdeaID:      ideaID,
		Title:       payload.Title,
		Description: payload.Description,
		PosX:        payload.PosX,
		PosY:        payload.PosY,
		Width:       payload.Width,
		Height:      payload.Height,
		Color:       payload.Color,
		Type:        payload.Type,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, node)
}

// GetIdeaNodes returns all nodes for an idea.
func (h *Handler) GetIdeaNodes(c *fiber.Ctx) error {
	ideaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	nodes, err := h.service.ListIdeaNodes(c.Context(), actorID, ideaID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, nodes)
}

// PatchIdeaNode handles node metadata updates.
func (h *Handler) PatchIdeaNode(c *fiber.Ctx) error {
	nodeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateIdeaNodePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	node, err := h.service.UpdateIdeaNode(c.Context(), UpdateIdeaNodeRequest{
		ActorID:     actorID,
		NodeID:      nodeID,
		Title:       payload.Title,
		Description: payload.Description,
		Color:       payload.Color,
		Type:        payload.Type,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, node)
}

// PatchIdeaNodePosition handles node position updates.
func (h *Handler) PatchIdeaNodePosition(c *fiber.Ctx) error {
	nodeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload moveIdeaNodePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	node, err := h.service.MoveIdeaNode(c.Context(), MoveIdeaNodeRequest{
		ActorID: actorID,
		NodeID:  nodeID,
		PosX:    payload.PosX,
		PosY:    payload.PosY,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, node)
}

// PatchIdeaNodeResize handles node dimension updates.
func (h *Handler) PatchIdeaNodeResize(c *fiber.Ctx) error {
	nodeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload resizeIdeaNodePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	node, err := h.service.ResizeIdeaNode(c.Context(), ResizeIdeaNodeRequest{
		ActorID: actorID,
		NodeID:  nodeID,
		Width:   payload.Width,
		Height:  payload.Height,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, node)
}

// DeleteIdeaNode handles node deletion.
func (h *Handler) DeleteIdeaNode(c *fiber.Ctx) error {
	nodeID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeleteIdeaNode(c.Context(), actorID, nodeID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "deleted"})
}

// PostIdeaNodeConnection handles connection creation between nodes.
func (h *Handler) PostIdeaNodeConnection(c *fiber.Ctx) error {
	var payload createIdeaNodeConnectionPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	ideaID, err := uuid.Parse(payload.IdeaID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	sourceNodeID, err := uuid.Parse(payload.SourceNodeID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	targetNodeID, err := uuid.Parse(payload.TargetNodeID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	conn, err := h.service.CreateIdeaNodeConnection(c.Context(), CreateIdeaNodeConnectionRequest{
		ActorID:      actorID,
		IdeaID:       ideaID,
		SourceNodeID: sourceNodeID,
		TargetNodeID: targetNodeID,
		Direction:    ConnectionDirection(payload.Direction),
		Label:        payload.Label,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, conn)
}

// GetIdeaNodeConnections returns all connections for an idea's canvas.
func (h *Handler) GetIdeaNodeConnections(c *fiber.Ctx) error {
	ideaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	connections, err := h.service.ListIdeaNodeConnections(c.Context(), actorID, ideaID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, connections)
}

// DeleteIdeaNodeConnection handles connection deletion.
func (h *Handler) DeleteIdeaNodeConnection(c *fiber.Ctx) error {
	connID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeleteIdeaNodeConnection(c.Context(), actorID, connID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "deleted"})
}

// Document handlers

type createDocumentPayload struct {
	IdeaID  string `json:"idea_id"`
	NodeID  string `json:"node_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updateDocumentPayload struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// PostDocument handles document creation.
func (h *Handler) PostDocument(c *fiber.Ctx) error {
	var payload createDocumentPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	ideaID, err := uuid.Parse(payload.IdeaID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var nodeID *uuid.UUID
	if payload.NodeID != "" {
		id, err := uuid.Parse(payload.NodeID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		nodeID = &id
	}

	doc, err := h.service.CreateDocument(c.Context(), CreateDocumentRequest{
		ActorID: actorID,
		IdeaID:  ideaID,
		NodeID:  nodeID,
		Title:   payload.Title,
		Content: payload.Content,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, doc)
}

// GetDocuments returns all documents for an idea, optionally filtered by node.
func (h *Handler) GetDocuments(c *fiber.Ctx) error {
	ideaID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	var nodeID *uuid.UUID
	if nodeIDParam := c.Query("node_id"); nodeIDParam != "" {
		id, err := uuid.Parse(nodeIDParam)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		nodeID = &id
	}

	docs, err := h.service.ListDocuments(c.Context(), actorID, ideaID, nodeID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, docs)
}

// PatchDocument handles document updates.
func (h *Handler) PatchDocument(c *fiber.Ctx) error {
	docID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateDocumentPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	doc, err := h.service.UpdateDocument(c.Context(), UpdateDocumentRequest{
		ActorID:   actorID,
		DocumentID: docID,
		Title:     payload.Title,
		Content:   payload.Content,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, doc)
}

// DeleteDocument handles document deletion.
func (h *Handler) DeleteDocument(c *fiber.Ctx) error {
	docID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	actorID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeleteDocument(c.Context(), actorID, docID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "deleted"})
}
