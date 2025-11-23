package testimonials

import (
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes testimonial endpoints.
type Handler interface {
	CreateTestimonial(c *fiber.Ctx) error
	UpdateTestimonial(c *fiber.Ctx) error
	GetTestimonial(c *fiber.Ctx) error
	ListTestimonials(c *fiber.Ctx) error
	DeleteTestimonial(c *fiber.Ctx) error
	ApproveTestimonial(c *fiber.Ctx) error
	RejectTestimonial(c *fiber.Ctx) error
	HideTestimonial(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a testimonial handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type createTestimonialPayload struct {
	AuthorName    string `json:"authorName"`
	AuthorRole    string `json:"authorRole,omitempty"`
	AuthorCompany string `json:"authorCompany,omitempty"`
	AuthorPhoto   string `json:"authorPhoto,omitempty"`
	Content       string `json:"content"`
	Rating        *int   `json:"rating,omitempty"`
	LinkedInURL   string `json:"linkedinUrl,omitempty"`
	DisplayOrder  int    `json:"displayOrder,omitempty"`
}

type updateTestimonialPayload struct {
	AuthorName    *string            `json:"authorName,omitempty"`
	AuthorRole    *string            `json:"authorRole,omitempty"`
	AuthorCompany *string            `json:"authorCompany,omitempty"`
	AuthorPhoto   *string            `json:"authorPhoto,omitempty"`
	Content       *string            `json:"content,omitempty"`
	Rating        *int               `json:"rating,omitempty"`
	LinkedInURL   *string            `json:"linkedinUrl,omitempty"`
	Status        *TestimonialStatus `json:"status,omitempty"`
	DisplayOrder  *int               `json:"displayOrder,omitempty"`
}

// Handlers

func (h *handler) CreateTestimonial(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	var payload createTestimonialPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	testimonial, err := h.service.CreateTestimonial(c.Context(), userID, CreateTestimonialRequest{
		AuthorName:    payload.AuthorName,
		AuthorRole:    payload.AuthorRole,
		AuthorCompany: payload.AuthorCompany,
		AuthorPhoto:   payload.AuthorPhoto,
		Content:       payload.Content,
		Rating:        payload.Rating,
		LinkedInURL:   payload.LinkedInURL,
		DisplayOrder:  payload.DisplayOrder,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, testimonial)
}

func (h *handler) UpdateTestimonial(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	testimonialID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid testimonial id",
		})
	}

	var payload updateTestimonialPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	testimonial, err := h.service.UpdateTestimonial(c.Context(), userID, testimonialID, UpdateTestimonialRequest{
		AuthorName:    payload.AuthorName,
		AuthorRole:    payload.AuthorRole,
		AuthorCompany: payload.AuthorCompany,
		AuthorPhoto:   payload.AuthorPhoto,
		Content:       payload.Content,
		Rating:        payload.Rating,
		LinkedInURL:   payload.LinkedInURL,
		Status:        payload.Status,
		DisplayOrder:  payload.DisplayOrder,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, testimonial)
}

func (h *handler) GetTestimonial(c *fiber.Ctx) error {
	testimonialID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid testimonial id",
		})
	}

	testimonial, err := h.service.GetTestimonial(c.Context(), testimonialID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, testimonial)
}

func (h *handler) ListTestimonials(c *fiber.Ctx) error {
	// For GET requests, try to get userID from API key context first, then JWT
	var userID *uuid.UUID
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		uid := apiKey.UserID
		userID = &uid
	} else if uid, err := authdomain.UserIDFromContext(c); err == nil {
		userID = &uid
	}

	filters := ListTestimonialsFilters{
		UserID: userID,
	}

	// Parse query parameters
	if statusStr := c.Query("status"); statusStr != "" {
		status := TestimonialStatus(statusStr)
		filters.Status = &status
	} else {
		// Default to approved for public access
		approved := TestimonialStatusApproved
		filters.Status = &approved
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

	testimonials, err := h.service.ListTestimonials(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, testimonials)
}

func (h *handler) DeleteTestimonial(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	testimonialID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid testimonial id",
		})
	}

	if err := h.service.DeleteTestimonial(c.Context(), userID, testimonialID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "testimonial deleted successfully",
	})
}

func (h *handler) ApproveTestimonial(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	testimonialID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid testimonial id",
		})
	}

	testimonial, err := h.service.ApproveTestimonial(c.Context(), userID, testimonialID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, testimonial)
}

func (h *handler) RejectTestimonial(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	testimonialID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid testimonial id",
		})
	}

	testimonial, err := h.service.RejectTestimonial(c.Context(), userID, testimonialID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, testimonial)
}

func (h *handler) HideTestimonial(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	testimonialID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid testimonial id",
		})
	}

	testimonial, err := h.service.HideTestimonial(c.Context(), userID, testimonialID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, testimonial)
}

// Error handling

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	domainErr, ok := AsDomainError(err)
	if !ok {
		h.logger.Error("unexpected error in testimonial handler", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{
			"message": "internal server error",
		})
	}

	statusCode := fiber.StatusInternalServerError
	switch domainErr.Code {
	case ErrCodeInvalidPayload, ErrCodeInvalidAuthorName, ErrCodeInvalidContent, ErrCodeInvalidRating, ErrCodeInvalidStatus:
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

