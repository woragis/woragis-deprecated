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
