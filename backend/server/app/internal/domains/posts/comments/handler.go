package comments

import (
	"errors"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes comment endpoints.
type Handler interface {
	CreateComment(c *fiber.Ctx) error
	UpdateComment(c *fiber.Ctx) error
	GetComment(c *fiber.Ctx) error
	DeleteComment(c *fiber.Ctx) error
	ListComments(c *fiber.Ctx) error
	ApproveComment(c *fiber.Ctx) error
	RejectComment(c *fiber.Ctx) error
	MarkCommentAsSpam(c *fiber.Ctx) error
	GetCommentCount(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a comment handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type createCommentPayload struct {
	PostID      uuid.UUID  `json:"postId"`
	Content     string     `json:"content"`
	AuthorName  string     `json:"authorName"`
	AuthorEmail string     `json:"authorEmail,omitempty"`
	AuthorURL   string     `json:"authorUrl,omitempty"`
	ParentID    *uuid.UUID `json:"parentId,omitempty"`
}

type updateCommentPayload struct {
	Content *string `json:"content,omitempty"`
}

// Handlers

func (h *handler) CreateComment(c *fiber.Ctx) error {
	var userID *uuid.UUID
	if uid, err := authdomain.UserIDFromContext(c); err == nil {
		userID = &uid
	}

	var payload createCommentPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	comment, err := h.service.CreateComment(c.Context(), CreateCommentRequest{
		PostID:      payload.PostID,
		Content:     payload.Content,
		AuthorName:  payload.AuthorName,
		AuthorEmail: payload.AuthorEmail,
		AuthorURL:   payload.AuthorURL,
		ParentID:    payload.ParentID,
		UserID:      userID,
		IPAddress:   c.IP(),
		UserAgent:   c.Get("User-Agent"),
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toCommentResponse(comment))
}

func (h *handler) UpdateComment(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateCommentPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	comment, err := h.service.UpdateComment(c.Context(), userID, commentID, UpdateCommentRequest{
		Content: payload.Content,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toCommentResponse(comment))
}

func (h *handler) GetComment(c *fiber.Ctx) error {
	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	comment, err := h.service.GetComment(c.Context(), commentID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toCommentResponse(comment))
}

func (h *handler) DeleteComment(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeleteComment(c.Context(), userID, commentID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "comment deleted"})
}

func (h *handler) ListComments(c *fiber.Ctx) error {
	filters := CommentFilters{}

	// Post ID is required
	postIDStr := c.Query("postId")
	if postIDStr == "" {
		// Try to get from URL params if nested under posts
		postIDStr = c.Params("postId")
	}
	if postIDStr != "" {
		if postID, err := uuid.Parse(postIDStr); err == nil {
			filters.PostID = &postID
		}
	}

	// User ID if authenticated
	if userID, err := authdomain.UserIDFromContext(c); err == nil {
		filters.UserID = &userID
	}

	// Parent ID for nested comments
	if parentIDStr := c.Query("parentId"); parentIDStr != "" {
		if parentID, err := uuid.Parse(parentIDStr); err == nil {
			filters.ParentID = &parentID
		}
	} else if c.Query("topLevel") == "true" {
		// Top-level comments only
		nilUUID := uuid.Nil
		filters.ParentID = &nilUUID
	}

	// Status filter (default to approved for public)
	if statusStr := c.Query("status"); statusStr != "" {
		status := CommentStatus(statusStr)
		filters.Status = &status
	} else {
		// Default to approved for public access
		approved := CommentStatusApproved
		filters.Status = &approved
	}

	// Pagination
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	// Ordering
	if orderBy := c.Query("orderBy"); orderBy != "" {
		filters.OrderBy = orderBy
	} else {
		filters.OrderBy = "created_at"
	}

	if order := c.Query("order"); order != "" {
		filters.Order = order
	} else {
		filters.Order = "asc"
	}

	comments, err := h.service.ListComments(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]commentResponse, len(comments))
	for i := range comments {
		responses[i] = toCommentResponse(&comments[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) ApproveComment(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.ApproveComment(c.Context(), commentID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "comment approved"})
}

func (h *handler) RejectComment(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.RejectComment(c.Context(), commentID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "comment rejected"})
}

func (h *handler) MarkCommentAsSpam(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	commentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.MarkCommentAsSpam(c.Context(), commentID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "comment marked as spam"})
}

func (h *handler) GetCommentCount(c *fiber.Ctx) error {
	postIDStr := c.Query("postId")
	if postIDStr == "" {
		postIDStr = c.Params("postId")
	}
	if postIDStr == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "postId is required",
		})
	}

	postID, err := uuid.Parse(postIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var status *CommentStatus
	if statusStr := c.Query("status"); statusStr != "" {
		s := CommentStatus(statusStr)
		status = &s
	}

	count, err := h.service.GetCommentCount(c.Context(), postID, status)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"count": count})
}

// Response helpers

type commentResponse struct {
	ID          string     `json:"id"`
	PostID      string     `json:"postId"`
	UserID      *string    `json:"userId,omitempty"`
	ParentID    *string    `json:"parentId,omitempty"`
	AuthorName  string     `json:"authorName"`
	AuthorEmail string     `json:"authorEmail,omitempty"`
	AuthorURL   string     `json:"authorUrl,omitempty"`
	Content     string     `json:"content"`
	Status      CommentStatus `json:"status"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
}

func toCommentResponse(comment *Comment) commentResponse {
	resp := commentResponse{
		ID:         comment.ID.String(),
		PostID:     comment.PostID.String(),
		AuthorName: comment.AuthorName,
		Content:    comment.Content,
		Status:     comment.Status,
		CreatedAt: comment.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: comment.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if comment.UserID != nil {
		userIDStr := comment.UserID.String()
		resp.UserID = &userIDStr
	}

	if comment.ParentID != nil {
		parentIDStr := comment.ParentID.String()
		resp.ParentID = &parentIDStr
	}

	if comment.AuthorEmail != "" {
		resp.AuthorEmail = comment.AuthorEmail
	}

	if comment.AuthorURL != "" {
		resp.AuthorURL = comment.AuthorURL
	}

	return resp
}

// Error handling

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		statusCode := fiber.StatusInternalServerError
		switch domainErr.Code {
		case ErrCodeCommentNotFound:
			statusCode = fiber.StatusNotFound
		case ErrCodeInvalidPayload, ErrCodeInvalidContent, ErrCodeInvalidName:
			statusCode = fiber.StatusBadRequest
		case ErrCodeUnauthorized:
			statusCode = fiber.StatusUnauthorized
		}
		return response.Error(c, statusCode, domainErr.Code, fiber.Map{
			"message": domainErr.Message,
		})
	}

	h.logger.Error("unexpected error in comment handler", "error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{
		"message": "internal server error",
	})
}

// AsDomainError checks if an error is a domain error.
func AsDomainError(err error) (*DomainError, bool) {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr, true
	}
	return nil, false
}

