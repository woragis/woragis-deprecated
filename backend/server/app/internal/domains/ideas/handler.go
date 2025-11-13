package ideas

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

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
	UserID      string  `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	PosX        float64 `json:"pos_x"`
	PosY        float64 `json:"pos_y"`
	Color       string  `json:"color"`
	ProjectID   string  `json:"project_id"`
}

type updateIdeaPayload struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Color       string `json:"color"`
	ProjectID   string `json:"project_id"`
}

type moveIdeaPayload struct {
	UserID string  `json:"user_id"`
	PosX   float64 `json:"pos_x"`
	PosY   float64 `json:"pos_y"`
}

type createLinkPayload struct {
	UserID        string  `json:"user_id"`
	SourceIdeaID  string  `json:"source_idea_id"`
	TargetIdeaID  string  `json:"target_idea_id"`
	Relation      string  `json:"relation"`
	Weight        float64 `json:"weight"`
	Bidirectional bool    `json:"bidirectional"`
}

// PostIdea handles idea creation.
func (h *Handler) PostIdea(c *fiber.Ctx) error {
	var payload createIdeaPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
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
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateIdeaPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var projectID *uuid.UUID
	if payload.ProjectID != "" {
		parsed, err := uuid.Parse(payload.ProjectID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		projectID = &parsed
	}

	idea, err := h.service.UpdateIdea(c.Context(), UpdateIdeaRequest{
		UserID:      userID,
		IdeaID:      id,
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
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload moveIdeaPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	idea, err := h.service.MoveIdea(c.Context(), MoveIdeaRequest{
		UserID: userID,
		IdeaID: id,
		PosX:   payload.PosX,
		PosY:   payload.PosY,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, idea)
}

// PostLink creates a relationship between ideas.
func (h *Handler) PostLink(c *fiber.Ctx) error {
	var payload createLinkPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
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
		UserID:        userID,
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

// ListIdeas returns all ideas for a user.
func (h *Handler) ListIdeas(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	ideas, err := h.service.ListIdeas(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, ideas)
}

// ListLinks returns links for a user, optionally filtered by idea.
func (h *Handler) ListLinks(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var ideaID uuid.UUID
	if ideaParam := c.Query("idea_id"); ideaParam != "" {
		ideaID, err = uuid.Parse(ideaParam)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	links, err := h.service.ListLinks(c.Context(), userID, ideaID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, links)
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromError(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, nil)
	}

	h.logError("ideas: unexpected error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func statusFromError(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidTitle, ErrCodeInvalidRelation:
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
