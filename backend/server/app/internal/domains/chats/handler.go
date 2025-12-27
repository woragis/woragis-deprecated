package chats

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	langchainservice "github.com/woragis/backend/server/app/internal/services/langchain"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes HTTP endpoints for chat conversations.
type Handler struct {
	service *Service
	logger  *slog.Logger
	stream  *StreamHub
}

// NewHandler constructs a handler instance.
func NewHandler(service *Service, logger *slog.Logger, stream *StreamHub) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
		stream:  stream,
	}
}

type createConversationPayload struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	IdeaID           string `json:"ideaId"`           // Accept camelCase from frontend
	ProjectID        string `json:"projectId"`        // Accept camelCase from frontend
	JobApplicationID string `json:"jobApplicationId"`  // Accept camelCase from frontend
}

type appendMessagePayload struct {
	Role          string  `json:"role"`
	Content       string  `json:"content"`
	GenerateReply bool    `json:"generate_reply"`
	Agent         string  `json:"agent"`
	Provider      string  `json:"provider"`
	Model         string  `json:"model"`
	MaxTokens     int     `json:"max_tokens"`
	Temperature   float64 `json:"temperature"`
}

type bulkUpdatePayload struct {
	ConversationIDs []string `json:"conversation_ids"`
}

type searchQueryParams struct {
	Query            string
	IncludeArchived  bool
	JobApplicationID string
	Limit            int
}

type shareTranscriptPayload struct {
	ExpireAfter *string `json:"expire_after"`
}

type assignmentPayload struct {
	AgentID   string `json:"agent_id"`
	AgentName string `json:"agent_name"`
	Notes     string `json:"notes"`
}

type conversationResponse struct {
	ID                string  `json:"id"`
	UserID            string  `json:"user_id"`
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	IdeaID            *string `json:"idea_id"`
	ProjectID         *string `json:"project_id"`
	JobApplicationID  *string `json:"job_application_id"`
	AssignedAgentID   *string `json:"assigned_agent_id"`
	SharedTranscript  string  `json:"shared_transcript"`
	ArchivedAt        *string `json:"archived_at"`
	DeletedAt         *string `json:"deleted_at"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

type messageResponse struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	Role           string `json:"role"`
	Content        string `json:"content"`
	CreatedAt      string `json:"created_at"`
}

type transcriptResponse struct {
	ID        string  `json:"id"`
	ShareCode string  `json:"share_code"`
	Content   string  `json:"content,omitempty"`
	CreatedAt string  `json:"created_at"`
	ExpiresAt *string `json:"expires_at,omitempty"`
}

type assignmentResponse struct {
	ID           string  `json:"id"`
	AgentID      string  `json:"agent_id"`
	AgentName    string  `json:"agent_name"`
	AssignedAt   string  `json:"assigned_at"`
	UnassignedAt *string `json:"unassigned_at,omitempty"`
	Notes        string  `json:"notes,omitempty"`
}

// CreateConversation handles POST /chats/conversations.
func (h *Handler) CreateConversation(c *fiber.Ctx) error {
	var payload createConversationPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateCreateConversationPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
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

	var jobApplicationID *uuid.UUID
	if payload.JobApplicationID != "" {
		parsed, err := uuid.Parse(payload.JobApplicationID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		jobApplicationID = &parsed
	}

	conversation, err := h.service.CreateConversation(c.Context(), CreateConversationRequest{
		UserID:           userID,
		Title:            payload.Title,
		Description:      payload.Description,
		IdeaID:           ideaID,
		ProjectID:        projectID,
		JobApplicationID: jobApplicationID,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toConversationResponse(conversation))
}

// ListConversations handles GET /chats/conversations.
func (h *Handler) ListConversations(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
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

// GetConversation handles GET /chats/conversations/:id.
func (h *Handler) GetConversation(c *fiber.Ctx) error {
	convID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	conversation, err := h.service.GetConversation(c.Context(), convID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toConversationResponse(conversation))
}

// SearchConversations handles GET /chats/conversations/search.
func (h *Handler) SearchConversations(c *fiber.Ctx) error {
	query := c.Query("q")
	limit := c.QueryInt("limit", 20)
	
	// Validate query parameters
	if err := ValidateSearchQueryParams(query, limit); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	params := searchQueryParams{
		Query:            query,
		JobApplicationID: c.Query("job_application_id"),
	}
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}
	includeArchived := strings.ToLower(c.Query("include_archived")) == "true"
	params.IncludeArchived = includeArchived
	if limit > 0 {
		params.Limit = limit
	}

	var jobApplicationID *uuid.UUID
	if params.JobApplicationID != "" {
		parsed, err := uuid.Parse(params.JobApplicationID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		jobApplicationID = &parsed
		h.logger.Info("SearchConversations",
			"user_id", userID.String(),
			"job_application_id", jobApplicationID.String(),
			"include_archived", includeArchived,
			"query", params.Query)
	} else {
		h.logger.Info("SearchConversations",
			"user_id", userID.String(),
			"job_application_id", "nil",
			"include_archived", includeArchived,
			"query", params.Query)
	}

	conversations, err := h.service.SearchConversations(c.Context(), SearchConversationsRequest{
		UserID:           userID,
		Query:            params.Query,
		IncludeArchived:  params.IncludeArchived,
		JobApplicationID: jobApplicationID,
		Limit:            params.Limit,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	jobAppIDStr := "nil"
	if jobApplicationID != nil {
		jobAppIDStr = jobApplicationID.String()
	}
	h.logger.Info("SearchConversations result",
		"user_id", userID.String(),
		"job_application_id", jobAppIDStr,
		"count", len(conversations))

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

	// Validate payload
	if err := ValidateAppendMessagePayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	messages, err := h.service.AppendMessage(c.Context(), AppendMessageRequest{
		ConversationID: convID,
		UserID:         userID,
		Role:           payload.Role,
		Content:        payload.Content,
		GenerateReply:  payload.GenerateReply,
		Agent:         strings.ToLower(strings.TrimSpace(payload.Agent)),
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

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	messages, err := h.service.ListMessages(c.Context(), convID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toMessageResponses(messages))
}

// ArchiveConversations handles POST /chats/conversations/archive.
func (h *Handler) ArchiveConversations(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	ids, err := h.parseBulkPayload(c)
	if err != nil {
		return err
	}

	if err := h.service.ArchiveConversations(c.Context(), BulkUpdateRequest{
		UserID:          userID,
		ConversationIDs: ids,
	}); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "archived"})
}

// DeleteConversations handles POST /chats/conversations/delete.
func (h *Handler) DeleteConversations(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	ids, err := h.parseBulkPayload(c)
	if err != nil {
		return err
	}

	if err := h.service.DeleteConversations(c.Context(), BulkUpdateRequest{
		UserID:          userID,
		ConversationIDs: ids,
	}); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "deleted"})
}

// RestoreConversations handles POST /chats/conversations/restore.
func (h *Handler) RestoreConversations(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	ids, err := h.parseBulkPayload(c)
	if err != nil {
		return err
	}

	if err := h.service.RestoreConversations(c.Context(), BulkUpdateRequest{
		UserID:          userID,
		ConversationIDs: ids,
	}); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "restored"})
}

// ShareTranscript handles POST /chats/conversations/:id/transcripts.
func (h *Handler) ShareTranscript(c *fiber.Ctx) error {
	convID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload shareTranscriptPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	var expire time.Duration
	if payload.ExpireAfter != nil {
		if parsed, parseErr := time.ParseDuration(*payload.ExpireAfter); parseErr == nil {
			expire = parsed
		}
	}

	transcript, err := h.service.ShareTranscript(c.Context(), ShareTranscriptRequest{
		UserID:         userID,
		ConversationID: convID,
		ExpireAfter:    expire,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toTranscriptResponse(transcript, true))
}

// ListTranscripts handles GET /chats/conversations/:id/transcripts.
func (h *Handler) ListTranscripts(c *fiber.Ctx) error {
	convID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	transcripts, err := h.service.ListTranscripts(c.Context(), userID, convID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]transcriptResponse, 0, len(transcripts))
	for _, t := range transcripts {
		resp = append(resp, toTranscriptResponse(&t, false))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// GetTranscript handles GET /chats/transcripts/:code.
func (h *Handler) GetTranscript(c *fiber.Ctx) error {
	shareCode := c.Params("code")
	transcript, err := h.service.GetTranscript(c.Context(), shareCode)
	if err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, toTranscriptResponse(transcript, true))
}

// AssignConversation handles POST /chats/conversations/:id/assign.
func (h *Handler) AssignConversation(c *fiber.Ctx) error {
	convID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload assignmentPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	agentID, err := uuid.Parse(payload.AgentID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.AssignConversation(c.Context(), AssignConversationRequest{
		UserID:         userID,
		ConversationID: convID,
		AgentID:        agentID,
		AgentName:      payload.AgentName,
		Notes:          payload.Notes,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "assigned"})
}

// UnassignConversation handles POST /chats/conversations/:id/unassign.
func (h *Handler) UnassignConversation(c *fiber.Ctx) error {
	convID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload assignmentPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	if err := h.service.UnassignConversation(c.Context(), userID, convID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "unassigned"})
}

// ListAssignments handles GET /chats/conversations/:id/assignments.
func (h *Handler) ListAssignments(c *fiber.Ctx) error {
	convID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	assignments, err := h.service.ListAssignments(c.Context(), convID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]assignmentResponse, 0, len(assignments))
	for _, assignment := range assignments {
		resp = append(resp, toAssignmentResponse(assignment))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// GetContextPreview handles GET /chats/conversations/:id/context.
func (h *Handler) GetContextPreview(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "authentication required"})
	}

	conversationID := c.Params("id")
	if conversationID == "" {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "conversation ID is required"})
	}

	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, 0, fiber.Map{"message": "invalid conversation ID"})
	}

	conv, err := h.service.GetConversation(c.Context(), convID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Build context with default options
	if h.service.GetContextBuilder() == nil {
		return response.Success(c, fiber.StatusOK, fiber.Map{
			"context": "",
			"options": GetDefaultContextOptions(),
			"message": "Context builder not available",
		})
	}

	contextStr, err := h.service.GetContextBuilder().BuildContext(
		c.Context(),
		userID,
		conv,
		GetDefaultContextOptions(),
	)

	if err != nil {
		h.logger.Error("failed to build context preview", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, fiber.Map{"message": "failed to build context"})
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"context": contextStr,
		"options": GetDefaultContextOptions(),
	})
}

// StreamConversation handles websocket upgrades for streaming events.
func (h *Handler) StreamConversation(conn *websocket.Conn) {
	conversationIDParam := conn.Params("id")
	conversationID, err := uuid.Parse(conversationIDParam)
	if err != nil {
		_ = conn.Close()
		return
	}

	if h.stream == nil {
		if h.logger != nil {
			h.logger.Error("chats: stream hub not configured")
		}
		_ = conn.Close()
		return
	}

	if _, err := userIDFromConn(conn); err != nil {
		_ = conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "unauthorized"))
		_ = conn.Close()
		return
	}

	h.stream.Register(conversationID, conn)
	defer func() {
		h.stream.Unregister(conversationID, conn)
		_ = conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
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

	var jobApplicationID *string
	if conv.JobApplicationID != nil {
		str := conv.JobApplicationID.String()
		jobApplicationID = &str
	}

	var assignedAgentID *string
	if conv.AssignedAgentID != nil {
		str := conv.AssignedAgentID.String()
		assignedAgentID = &str
	}

	var archivedAt *string
	if conv.ArchivedAt != nil {
		str := conv.ArchivedAt.Format(time.RFC3339)
		archivedAt = &str
	}

	var deletedAt *string
	if conv.DeletedAt != nil {
		str := conv.DeletedAt.Format(time.RFC3339)
		deletedAt = &str
	}

	return conversationResponse{
		ID:                conv.ID.String(),
		UserID:            conv.UserID.String(),
		Title:             conv.Title,
		Description:       conv.Description,
		IdeaID:            ideaID,
		ProjectID:         projectID,
		JobApplicationID:  jobApplicationID,
		AssignedAgentID:   assignedAgentID,
		SharedTranscript:  conv.SharedTranscript,
		ArchivedAt:        archivedAt,
		DeletedAt:         deletedAt,
		CreatedAt:         conv.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         conv.UpdatedAt.Format(time.RFC3339),
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

func (h *Handler) parseBulkPayload(c *fiber.Ctx) ([]uuid.UUID, error) {
	var payload bulkUpdatePayload
	if err := c.BodyParser(&payload); err != nil {
		return nil, response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	// Validate payload
	if err := ValidateBulkUpdatePayload(&payload); err != nil {
		return nil, response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, map[string]string{
			"message": err.Error(),
		})
	}

	conversationIDs := make([]uuid.UUID, 0, len(payload.ConversationIDs))
	for _, id := range payload.ConversationIDs {
		parsed, parseErr := uuid.Parse(id)
		if parseErr != nil {
			return nil, response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		conversationIDs = append(conversationIDs, parsed)
	}

	return conversationIDs, nil
}

func userIDFromConn(conn *websocket.Conn) (uuid.UUID, error) {
	raw := conn.Locals(authdomain.ContextUserIDKey)
	idStr, _ := raw.(string)
	if idStr == "" {
		return uuid.Nil, fiber.ErrUnauthorized
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fiber.ErrUnauthorized
	}

	return userID, nil
}

func toTranscriptResponse(transcript *ConversationTranscript, includeContent bool) transcriptResponse {
	if transcript == nil {
		return transcriptResponse{}
	}
	resp := transcriptResponse{
		ID:        transcript.ID.String(),
		ShareCode: transcript.ShareCode,
		CreatedAt: transcript.CreatedAt.Format(time.RFC3339),
	}
	if includeContent {
		resp.Content = transcript.Content
	}
	if transcript.ExpiresAt != nil {
		str := transcript.ExpiresAt.Format(time.RFC3339)
		resp.ExpiresAt = &str
	}
	return resp
}

func toAssignmentResponse(assignment ConversationAssignment) assignmentResponse {
	resp := assignmentResponse{
		ID:         assignment.ID.String(),
		AgentID:    assignment.AgentID.String(),
		AgentName:  assignment.AgentName,
		AssignedAt: assignment.AssignedAt.Format(time.RFC3339),
		Notes:      assignment.Notes,
	}
	if assignment.UnassignedAt != nil {
		str := assignment.UnassignedAt.Format(time.RFC3339)
		resp.UnassignedAt = &str
	}
	return resp
}
