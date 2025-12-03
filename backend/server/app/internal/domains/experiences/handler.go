package experiences

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

// Handler exposes experience endpoints.
type Handler interface {
	CreateExperience(c *fiber.Ctx) error
	UpdateExperience(c *fiber.Ctx) error
	GetExperience(c *fiber.Ctx) error
	ListExperiences(c *fiber.Ctx) error
	DeleteExperience(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs an experience handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type createExperiencePayload struct {
	Company      string          `json:"company"`
	Position     string          `json:"position"`
	PeriodStart  *time.Time      `json:"periodStart,omitempty"`
	PeriodEnd    *time.Time      `json:"periodEnd,omitempty"`
	PeriodText   string          `json:"periodText,omitempty"`
	Location     string          `json:"location,omitempty"`
	Description  string          `json:"description,omitempty"`
	Type         ExperienceType  `json:"type,omitempty"`
	CompanyURL   string          `json:"companyUrl,omitempty"`
	LinkedInURL  string          `json:"linkedinUrl,omitempty"`
	DisplayOrder int             `json:"displayOrder,omitempty"`
	IsCurrent    bool            `json:"isCurrent,omitempty"`
	Technologies []string        `json:"technologies,omitempty"`
	Projects     []ProjectInput  `json:"projects,omitempty"`
	Achievements []AchievementInput `json:"achievements,omitempty"`
}

type updateExperiencePayload struct {
	Company      *string         `json:"company,omitempty"`
	Position     *string         `json:"position,omitempty"`
	PeriodStart  *time.Time      `json:"periodStart,omitempty"`
	PeriodEnd    *time.Time      `json:"periodEnd,omitempty"`
	PeriodText   *string         `json:"periodText,omitempty"`
	Location     *string         `json:"location,omitempty"`
	Description  *string         `json:"description,omitempty"`
	Type         *ExperienceType `json:"type,omitempty"`
	CompanyURL   *string         `json:"companyUrl,omitempty"`
	LinkedInURL  *string         `json:"linkedinUrl,omitempty"`
	DisplayOrder *int            `json:"displayOrder,omitempty"`
	IsCurrent    *bool           `json:"isCurrent,omitempty"`
	Technologies []string        `json:"technologies,omitempty"`
	Projects     []ProjectInput  `json:"projects,omitempty"`
	Achievements []AchievementInput `json:"achievements,omitempty"`
}

// Handlers

func (h *handler) CreateExperience(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	var payload createExperiencePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	experience, err := h.service.CreateExperience(c.Context(), userID, CreateExperienceRequest{
		Company:      payload.Company,
		Position:     payload.Position,
		PeriodStart:  payload.PeriodStart,
		PeriodEnd:    payload.PeriodEnd,
		PeriodText:   payload.PeriodText,
		Location:     payload.Location,
		Description:  payload.Description,
		Type:         payload.Type,
		CompanyURL:   payload.CompanyURL,
		LinkedInURL:  payload.LinkedInURL,
		DisplayOrder: payload.DisplayOrder,
		IsCurrent:    payload.IsCurrent,
		Technologies: payload.Technologies,
		Projects:     payload.Projects,
		Achievements: payload.Achievements,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, experience)
}

func (h *handler) UpdateExperience(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	experienceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid experience id",
		})
	}

	var payload updateExperiencePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	experience, err := h.service.UpdateExperience(c.Context(), userID, experienceID, UpdateExperienceRequest{
		Company:      payload.Company,
		Position:     payload.Position,
		PeriodStart:  payload.PeriodStart,
		PeriodEnd:    payload.PeriodEnd,
		PeriodText:   payload.PeriodText,
		Location:     payload.Location,
		Description:  payload.Description,
		Type:         payload.Type,
		CompanyURL:   payload.CompanyURL,
		LinkedInURL:  payload.LinkedInURL,
		DisplayOrder: payload.DisplayOrder,
		IsCurrent:    payload.IsCurrent,
		Technologies: payload.Technologies,
		Projects:     payload.Projects,
		Achievements: payload.Achievements,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, experience)
}

func (h *handler) GetExperience(c *fiber.Ctx) error {
	experienceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid experience id",
		})
	}

	experience, err := h.service.GetExperience(c.Context(), experienceID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, experience)
}

func (h *handler) ListExperiences(c *fiber.Ctx) error {
	// For GET requests, try to get userID from API key context first, then JWT
	var userID *uuid.UUID
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		uid := apiKey.UserID
		userID = &uid
	} else if uid, err := authdomain.UserIDFromContext(c); err == nil {
		userID = &uid
	}

	filters := ListExperiencesFilters{
		UserID: userID,
	}

	// Parse query parameters
	if typeStr := c.Query("type"); typeStr != "" {
		expType := ExperienceType(typeStr)
		filters.Type = &expType
	}

	if isCurrentStr := c.Query("isCurrent"); isCurrentStr != "" {
		if isCurrent, err := strconv.ParseBool(isCurrentStr); err == nil {
			filters.IsCurrent = &isCurrent
		}
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

	experiences, err := h.service.ListExperiences(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, experiences)
}

func (h *handler) DeleteExperience(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	experienceID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid experience id",
		})
	}

	if err := h.service.DeleteExperience(c.Context(), userID, experienceID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "experience deleted successfully",
	})
}

// Error handling

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	domainErr, ok := AsDomainError(err)
	if !ok {
		h.logger.Error("unexpected error in experience handler", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{
			"message": "internal server error",
		})
	}

	statusCode := fiber.StatusInternalServerError
	switch domainErr.Code {
	case ErrCodeInvalidPayload, ErrCodeInvalidCompany, ErrCodeInvalidPosition, ErrCodeInvalidType:
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

