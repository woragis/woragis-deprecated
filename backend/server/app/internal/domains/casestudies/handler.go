package casestudies

import (
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes case study endpoints.
type Handler interface {
	CreateCaseStudy(c *fiber.Ctx) error
	UpdateCaseStudy(c *fiber.Ctx) error
	GetCaseStudy(c *fiber.Ctx) error
	GetCaseStudyByProjectSlug(c *fiber.Ctx) error
	ListCaseStudies(c *fiber.Ctx) error
	DeleteCaseStudy(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a case study handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type createCaseStudyPayload struct {
	ProjectID      uuid.UUID           `json:"projectId"`
	ProjectSlug    string              `json:"projectSlug"`
	Title          string              `json:"title"`
	Problem        string              `json:"problem"`
	Context        string              `json:"context"`
	Solution       string              `json:"solution"`
	Approach       []string            `json:"approach,omitempty"`
	Architecture   *ArchitectureData   `json:"architecture,omitempty"`
	Metrics        *MetricsData        `json:"metrics,omitempty"`
	LessonsLearned []string            `json:"lessonsLearned,omitempty"`
	Technologies   []string            `json:"technologies,omitempty"`
	Featured       bool                `json:"featured,omitempty"`
}

type updateCaseStudyPayload struct {
	Title          *string             `json:"title,omitempty"`
	Problem        *string             `json:"problem,omitempty"`
	Context        *string             `json:"context,omitempty"`
	Solution       *string             `json:"solution,omitempty"`
	Approach       []string            `json:"approach,omitempty"`
	Architecture   *ArchitectureData   `json:"architecture,omitempty"`
	Metrics        *MetricsData        `json:"metrics,omitempty"`
	LessonsLearned []string            `json:"lessonsLearned,omitempty"`
	Technologies   []string            `json:"technologies,omitempty"`
	Featured       *bool               `json:"featured,omitempty"`
}

// Handlers

func (h *handler) CreateCaseStudy(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	var payload createCaseStudyPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	caseStudy, err := h.service.CreateCaseStudy(c.Context(), userID, CreateCaseStudyRequest{
		ProjectID:      payload.ProjectID,
		ProjectSlug:    payload.ProjectSlug,
		Title:          payload.Title,
		Problem:        payload.Problem,
		Context:        payload.Context,
		Solution:       payload.Solution,
		Approach:       payload.Approach,
		Architecture:   payload.Architecture,
		Metrics:        payload.Metrics,
		LessonsLearned: payload.LessonsLearned,
		Technologies:   payload.Technologies,
		Featured:        payload.Featured,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, caseStudy)
}

func (h *handler) UpdateCaseStudy(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	caseStudyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid case study id",
		})
	}

	var payload updateCaseStudyPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	caseStudy, err := h.service.UpdateCaseStudy(c.Context(), userID, caseStudyID, UpdateCaseStudyRequest{
		Title:          payload.Title,
		Problem:        payload.Problem,
		Context:        payload.Context,
		Solution:       payload.Solution,
		Approach:       payload.Approach,
		Architecture:   payload.Architecture,
		Metrics:        payload.Metrics,
		LessonsLearned: payload.LessonsLearned,
		Technologies:   payload.Technologies,
		Featured:       payload.Featured,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, caseStudy)
}

func (h *handler) GetCaseStudy(c *fiber.Ctx) error {
	caseStudyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid case study id",
		})
	}

	caseStudy, err := h.service.GetCaseStudy(c.Context(), caseStudyID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, caseStudy)
}

func (h *handler) GetCaseStudyByProjectSlug(c *fiber.Ctx) error {
	projectSlug := c.Params("projectSlug")
	if projectSlug == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "project slug is required",
		})
	}

	caseStudy, err := h.service.GetCaseStudyByProjectSlug(c.Context(), projectSlug)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, caseStudy)
}

func (h *handler) ListCaseStudies(c *fiber.Ctx) error {
	// For GET requests, try to get userID from API key context first, then JWT
	var userID *uuid.UUID
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		uid := apiKey.UserID
		userID = &uid
	} else if uid, err := authdomain.UserIDFromContext(c); err == nil {
		userID = &uid
	}

	filters := ListCaseStudiesFilters{
		UserID: userID,
	}

	// Parse query parameters
	if projectIDStr := c.Query("projectId"); projectIDStr != "" {
		if projectID, err := uuid.Parse(projectIDStr); err == nil {
			filters.ProjectID = &projectID
		}
	}

	if projectSlug := c.Query("projectSlug"); projectSlug != "" {
		filters.ProjectSlug = &projectSlug
	}

	if featuredStr := c.Query("featured"); featuredStr != "" {
		if featured, err := strconv.ParseBool(featuredStr); err == nil {
			filters.Featured = &featured
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

	caseStudies, err := h.service.ListCaseStudies(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, caseStudies)
}

func (h *handler) DeleteCaseStudy(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	caseStudyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid case study id",
		})
	}

	if err := h.service.DeleteCaseStudy(c.Context(), userID, caseStudyID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "case study deleted successfully",
	})
}

// Error handling

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	domainErr, ok := AsDomainError(err)
	if !ok {
		h.logger.Error("unexpected error in case study handler", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{
			"message": "internal server error",
		})
	}

	statusCode := fiber.StatusInternalServerError
	switch domainErr.Code {
	case ErrCodeInvalidPayload, ErrCodeInvalidTitle:
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

