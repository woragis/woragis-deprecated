package chats

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	langchainservice "github.com/woragis/backend/server/app/internal/services/langchain"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes HTTP endpoints for chat conversations.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler constructs a handler instance.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type createConversationPayload struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IdeaID      string `json:"idea_id"`
	ProjectID   string `json:"project_id"`
}

type appendMessagePayload struct {
	UserID        string  `json:"user_id"`
	Role          string  `json:"role"`
	Content       string  `json:"content"`
	GenerateReply bool    `json:"generate_reply"`
	Provider      string  `json:"provider"`
	Model         string  `json:"model"`
	MaxTokens     int     `json:"max_tokens"`
	Temperature   float64 `json:"temperature"`
}

type conversationResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	IdeaID      *string `json:"idea_id"`
	ProjectID   *string `json:"project_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type messageResponse struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	Role           string `json:"role"`
	Content        string `json:"content"`
	CreatedAt      string `json:"created_at"`
}

// CreateConversation handles POST /chats/conversations.
func (h *Handler) CreateConversation(c *fiber.Ctx) error {
	var payload createConversationPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var ideaID *uuid.UUID
	if payload.IdeaID != "" {
		parsed, err := uuid.Parse(payload.IdeaID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		ideaID = &parsed
	}

	var projectID *uuid.UUID
	if payload.ProjectID != "" {
		parsed, err := uuid.Parse(payload.ProjectID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		projectID = &parsed
	}

	conversation, err := h.service.CreateConversation(c.Context(), CreateConversationRequest{
		UserID:      userID,
		Title:       payload.Title,
		Description: payload.Description,
		IdeaID:      ideaID,
		ProjectID:   projectID,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toConversationResponse(conversation))
}

// ListConversations handles GET /chats/conversations.
func (h *Handler) ListConversations(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	conversations, err := h.service.ListConversations(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]conversationResponse, 0, len(conversations))
	for _, conv := range conversations {
		convCopy := conv
		resp = append(resp, toConversationResponse(&convCopy))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// AppendMessage handles POST /chats/conversations/:id/messages.
func (h *Handler) AppendMessage(c *fiber.Ctx) error {
	convID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload appendMessagePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	messages, err := h.service.AppendMessage(c.Context(), AppendMessageRequest{
		ConversationID: convID,
		UserID:         userID,
		Role:           payload.Role,
		Content:        payload.Content,
		GenerateReply:  payload.GenerateReply,
		Provider:       langchainservice.ModelProvider(strings.ToLower(payload.Provider)),
		Model:          payload.Model,
		MaxTokens:      payload.MaxTokens,
		Temperature:    payload.Temperature,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toMessageResponses(messages))
}

// ListMessages handles GET /chats/conversations/:id/messages.
func (h *Handler) ListMessages(c *fiber.Ctx) error {
	convID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	messages, err := h.service.ListMessages(c.Context(), convID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toMessageResponses(messages))
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromError(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, nil)
	}

	h.logError("chats: unexpected error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func statusFromError(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidTitle, ErrCodeInvalidRole, ErrCodeInvalidContent:
		return fiber.StatusBadRequest
	case ErrCodeNotFound:
		return fiber.StatusNotFound
	case ErrCodeLLMFailure:
		return fiber.StatusBadGateway
	case ErrCodeConversationAccessDenied:
		return fiber.StatusForbidden
	default:
		return fiber.StatusInternalServerError
	}
}

func toConversationResponse(conv *Conversation) conversationResponse {
	var ideaID *string
	if conv.IdeaID != nil {
		str := conv.IdeaID.String()
		ideaID = &str
	}

	var projectID *string
	if conv.ProjectID != nil {
		str := conv.ProjectID.String()
		projectID = &str
	}

	return conversationResponse{
		ID:          conv.ID.String(),
		UserID:      conv.UserID.String(),
		Title:       conv.Title,
		Description: conv.Description,
		IdeaID:      ideaID,
		ProjectID:   projectID,
		CreatedAt:   conv.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   conv.UpdatedAt.Format(time.RFC3339),
	}
}

func toMessageResponses(messages []Message) []messageResponse {
	resp := make([]messageResponse, 0, len(messages))
	for _, msg := range messages {
		resp = append(resp, messageResponse{
			ID:             msg.ID.String(),
			ConversationID: msg.ConversationID.String(),
			Role:           msg.Role,
			Content:        msg.Content,
			CreatedAt:      msg.CreatedAt.Format(time.RFC3339),
		})
	}
	return resp
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
