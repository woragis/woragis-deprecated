package technicalwritings

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes technical writing endpoints.
type Handler interface {
	CreateTechnicalWriting(c *fiber.Ctx) error
	UpdateTechnicalWriting(c *fiber.Ctx) error
	GetTechnicalWriting(c *fiber.Ctx) error
	GetTechnicalWritingPublic(c *fiber.Ctx) error
	ListTechnicalWritings(c *fiber.Ctx) error
	ListFeaturedTechnicalWritings(c *fiber.Ctx) error
	GetWritingsByProject(c *fiber.Ctx) error
	GetWritingsByType(c *fiber.Ctx) error
	GetWritingsByPlatform(c *fiber.Ctx) error
	SearchTechnicalWritings(c *fiber.Ctx) error
	DeleteTechnicalWriting(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a technical writing handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type createTechnicalWritingPayload struct {
	Title         string             `json:"title"`
	Description   string             `json:"description"`
	Type          WritingType        `json:"type"`
	Platform      PublicationPlatform `json:"platform"`
	URL           string             `json:"url"`
	Content       string             `json:"content,omitempty"`
	Excerpt       string             `json:"excerpt,omitempty"`
	CanonicalURL  string             `json:"canonicalUrl,omitempty"`
	CoverImageURL string             `json:"coverImageUrl,omitempty"`
	PublishedAt   string             `json:"publishedAt,omitempty"` // ISO 8601 date string
	ReadingTime   int                `json:"readingTime,omitempty"`
	Topics        []string           `json:"topics,omitempty"`
	Technologies  []string           `json:"technologies,omitempty"`
	Views         *int               `json:"views,omitempty"`
	Likes         *int               `json:"likes,omitempty"`
	Shares        *int               `json:"shares,omitempty"`
	Comments      *int               `json:"comments,omitempty"`
	ProjectID     string              `json:"projectId,omitempty"`
	CaseStudyID   string              `json:"caseStudyId,omitempty"`
	Featured      bool                `json:"featured,omitempty"`
	DisplayOrder  int                 `json:"displayOrder,omitempty"`
}

type updateTechnicalWritingPayload struct {
	Title         *string             `json:"title,omitempty"`
	Description   *string             `json:"description,omitempty"`
	Type          *WritingType        `json:"type,omitempty"`
	Platform      *PublicationPlatform `json:"platform,omitempty"`
	URL           *string             `json:"url,omitempty"`
	Content       *string             `json:"content,omitempty"`
	Excerpt       *string             `json:"excerpt,omitempty"`
	CanonicalURL  *string             `json:"canonicalUrl,omitempty"`
	CoverImageURL *string             `json:"coverImageUrl,omitempty"`
	PublishedAt   *string             `json:"publishedAt,omitempty"`
	ReadingTime   *int                `json:"readingTime,omitempty"`
	Topics        []string            `json:"topics,omitempty"`
	Technologies  []string            `json:"technologies,omitempty"`
	Views         *int                `json:"views,omitempty"`
	Likes         *int                `json:"likes,omitempty"`
	Shares        *int                `json:"shares,omitempty"`
	Comments      *int                `json:"comments,omitempty"`
	ProjectID     *string             `json:"projectId,omitempty"`
	CaseStudyID   *string             `json:"caseStudyId,omitempty"`
	Featured      *bool               `json:"featured,omitempty"`
	DisplayOrder  *int                `json:"displayOrder,omitempty"`
}

// Handlers

func (h *handler) CreateTechnicalWriting(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	var payload createTechnicalWritingPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req := CreateTechnicalWritingRequest{
		Title:         payload.Title,
		Description:   payload.Description,
		Type:          payload.Type,
		Platform:      payload.Platform,
		URL:           payload.URL,
		Content:       payload.Content,
		Excerpt:       payload.Excerpt,
		CanonicalURL:  payload.CanonicalURL,
		CoverImageURL: payload.CoverImageURL,
		ReadingTime:   payload.ReadingTime,
		Topics:        payload.Topics,
		Technologies:  payload.Technologies,
		Views:         payload.Views,
		Likes:         payload.Likes,
		Shares:        payload.Shares,
		Comments:      payload.Comments,
		Featured:      payload.Featured,
		DisplayOrder:  payload.DisplayOrder,
	}

	// Parse published date
	if payload.PublishedAt != "" {
		publishedAt, err := parseDate(payload.PublishedAt)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
				"message": "invalid published date format",
			})
		}
		req.PublishedAt = &publishedAt
	}

	if payload.ProjectID != "" {
		req.ProjectID = &payload.ProjectID
	}
	if payload.CaseStudyID != "" {
		req.CaseStudyID = &payload.CaseStudyID
	}

	writing, err := h.service.CreateTechnicalWriting(c.Context(), userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, writing)
}

func (h *handler) UpdateTechnicalWriting(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	writingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid writing id",
		})
	}

	var payload updateTechnicalWritingPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req := UpdateTechnicalWritingRequest{}

	if payload.Title != nil {
		req.Title = payload.Title
	}
	if payload.Description != nil {
		req.Description = payload.Description
	}
	if payload.Type != nil {
		req.Type = payload.Type
	}
	if payload.Platform != nil {
		req.Platform = payload.Platform
	}
	if payload.URL != nil {
		req.URL = payload.URL
	}
	if payload.Content != nil {
		req.Content = payload.Content
	}
	if payload.Excerpt != nil {
		req.Excerpt = payload.Excerpt
	}
	if payload.CanonicalURL != nil {
		req.CanonicalURL = payload.CanonicalURL
	}
	if payload.CoverImageURL != nil {
		req.CoverImageURL = payload.CoverImageURL
	}
	if payload.PublishedAt != nil {
		if *payload.PublishedAt != "" {
			publishedAt, err := parseDate(*payload.PublishedAt)
			if err != nil {
				return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
					"message": "invalid published date format",
				})
			}
			req.PublishedAt = &publishedAt
		} else {
			req.PublishedAt = nil
		}
	}
	if payload.ReadingTime != nil {
		req.ReadingTime = payload.ReadingTime
	}
	if payload.Topics != nil {
		req.Topics = payload.Topics
	}
	if payload.Technologies != nil {
		req.Technologies = payload.Technologies
	}
	if payload.Views != nil {
		req.Views = payload.Views
	}
	if payload.Likes != nil {
		req.Likes = payload.Likes
	}
	if payload.Shares != nil {
		req.Shares = payload.Shares
	}
	if payload.Comments != nil {
		req.Comments = payload.Comments
	}
	if payload.ProjectID != nil {
		req.ProjectID = payload.ProjectID
	}
	if payload.CaseStudyID != nil {
		req.CaseStudyID = payload.CaseStudyID
	}
	if payload.Featured != nil {
		req.Featured = payload.Featured
	}
	if payload.DisplayOrder != nil {
		req.DisplayOrder = payload.DisplayOrder
	}

	writing, err := h.service.UpdateTechnicalWriting(c.Context(), userID, writingID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, writing)
}

func (h *handler) GetTechnicalWriting(c *fiber.Ctx) error {
	writingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid writing id",
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	writing, err := h.service.GetTechnicalWriting(c.Context(), writingID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, writing)
}

func (h *handler) GetTechnicalWritingPublic(c *fiber.Ctx) error {
	writingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid writing id",
		})
	}

	writing, err := h.service.GetTechnicalWritingPublic(c.Context(), writingID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, writing)
}

func (h *handler) ListTechnicalWritings(c *fiber.Ctx) error {
	var userID *uuid.UUID
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		uid := apiKey.UserID
		userID = &uid
	} else if uid, err := authdomain.UserIDFromContext(c); err == nil {
		userID = &uid
	}

	filters := ListTechnicalWritingFilters{
		UserID: userID,
	}

	// Parse query parameters
	if typeStr := c.Query("type"); typeStr != "" {
		writingType := WritingType(typeStr)
		filters.Type = &writingType
	}

	if platformStr := c.Query("platform"); platformStr != "" {
		platform := PublicationPlatform(platformStr)
		filters.Platform = &platform
	}

	if projectIDStr := c.Query("projectId"); projectIDStr != "" {
		projectID, err := uuid.Parse(projectIDStr)
		if err == nil {
			filters.ProjectID = &projectID
		}
	}

	if featuredStr := c.Query("featured"); featuredStr != "" {
		featured := featuredStr == "true"
		filters.Featured = &featured
	}

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

	if orderBy := c.Query("orderBy"); orderBy != "" {
		filters.OrderBy = orderBy
	}

	if order := c.Query("order"); order != "" {
		filters.Order = order
	}

	writings, err := h.service.ListTechnicalWritings(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, writings)
}

func (h *handler) ListFeaturedTechnicalWritings(c *fiber.Ctx) error {
	writings, err := h.service.ListFeaturedTechnicalWritings(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, writings)
}

func (h *handler) GetWritingsByProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid project id",
		})
	}

	writings, err := h.service.GetWritingsByProject(c.Context(), projectID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, writings)
}

func (h *handler) GetWritingsByType(c *fiber.Ctx) error {
	writingTypeStr := c.Params("type")
	writingType := WritingType(writingTypeStr)

	if !isValidWritingType(writingType) {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidType, fiber.Map{
			"message": "invalid writing type",
		})
	}

	writings, err := h.service.GetWritingsByType(c.Context(), writingType)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, writings)
}

func (h *handler) GetWritingsByPlatform(c *fiber.Ctx) error {
	platformStr := c.Params("platform")
	platform := PublicationPlatform(platformStr)

	if !isValidPlatform(platform) {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPlatform, fiber.Map{
			"message": "invalid platform",
		})
	}

	writings, err := h.service.GetWritingsByPlatform(c.Context(), platform)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, writings)
}

func (h *handler) SearchTechnicalWritings(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "search query is required",
		})
	}

	writings, err := h.service.SearchTechnicalWritings(c.Context(), query)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, writings)
}

func (h *handler) DeleteTechnicalWriting(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	writingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid writing id",
		})
	}

	if err := h.service.DeleteTechnicalWriting(c.Context(), userID, writingID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "writing deleted successfully",
	})
}

// Helper functions

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	domainErr, ok := AsDomainError(err)
	if !ok {
		h.logger.Error("unexpected error in technical writing handler", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{
			"message": "internal server error",
		})
	}

	statusCode := fiber.StatusInternalServerError
	switch domainErr.Code {
	case ErrCodeInvalidPayload, ErrCodeInvalidType, ErrCodeInvalidPlatform, ErrCodeInvalidTitle:
		statusCode = fiber.StatusBadRequest
	case ErrCodeNotFound:
		statusCode = fiber.StatusNotFound
	case ErrCodeUnauthorized:
		statusCode = fiber.StatusUnauthorized
	case ErrCodeConflict:
		statusCode = fiber.StatusConflict
	}

	return response.Error(c, statusCode, domainErr.Code, fiber.Map{
		"message": domainErr.Message,
	})
}

func parseDate(dateStr string) (time.Time, error) {
	// Try ISO 8601 format first (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t, nil
	}
	// Try RFC3339 format
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t, nil
	}
	return time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid date format")
}

