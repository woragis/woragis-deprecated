package aimlintegrations

import (
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes AI/ML integration endpoints.
type Handler interface {
	CreateAIMLIntegration(c *fiber.Ctx) error
	UpdateAIMLIntegration(c *fiber.Ctx) error
	GetAIMLIntegration(c *fiber.Ctx) error
	GetAIMLIntegrationPublic(c *fiber.Ctx) error
	ListAIMLIntegrations(c *fiber.Ctx) error
	ListFeaturedAIMLIntegrations(c *fiber.Ctx) error
	GetIntegrationsByProject(c *fiber.Ctx) error
	GetIntegrationsByType(c *fiber.Ctx) error
	GetIntegrationsByFramework(c *fiber.Ctx) error
	DeleteAIMLIntegration(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs an AI/ML integration handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type createAIMLIntegrationPayload struct {
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Type            IntegrationType `json:"type"`
	Framework       Framework       `json:"framework"`
	ModelName       string          `json:"modelName,omitempty"`
	ModelVersion    string          `json:"modelVersion,omitempty"`
	UseCase         string          `json:"useCase,omitempty"`
	Impact          string          `json:"impact,omitempty"`
	Technologies    []string        `json:"technologies,omitempty"`
	Architecture    string          `json:"architecture,omitempty"`
	Metrics         string          `json:"metrics,omitempty"`
	ProjectID       string          `json:"projectId,omitempty"`
	CaseStudyID     string          `json:"caseStudyId,omitempty"`
	Featured        bool            `json:"featured,omitempty"`
	DisplayOrder    int             `json:"displayOrder,omitempty"`
	DemoURL         string          `json:"demoUrl,omitempty"`
	DocumentationURL string          `json:"documentationUrl,omitempty"`
	GitHubURL       string          `json:"githubUrl,omitempty"`
}

type updateAIMLIntegrationPayload struct {
	Title           *string          `json:"title,omitempty"`
	Description     *string          `json:"description,omitempty"`
	Type            *IntegrationType `json:"type,omitempty"`
	Framework       *Framework       `json:"framework,omitempty"`
	ModelName       *string          `json:"modelName,omitempty"`
	ModelVersion    *string          `json:"modelVersion,omitempty"`
	UseCase         *string          `json:"useCase,omitempty"`
	Impact          *string          `json:"impact,omitempty"`
	Technologies    []string         `json:"technologies,omitempty"`
	Architecture    *string          `json:"architecture,omitempty"`
	Metrics         *string          `json:"metrics,omitempty"`
	ProjectID       *string          `json:"projectId,omitempty"`
	CaseStudyID     *string          `json:"caseStudyId,omitempty"`
	Featured        *bool            `json:"featured,omitempty"`
	DisplayOrder    *int             `json:"displayOrder,omitempty"`
	DemoURL         *string          `json:"demoUrl,omitempty"`
	DocumentationURL *string         `json:"documentationUrl,omitempty"`
	GitHubURL       *string          `json:"githubUrl,omitempty"`
}

// Handlers

func (h *handler) CreateAIMLIntegration(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	var payload createAIMLIntegrationPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req := CreateAIMLIntegrationRequest{
		Title:           payload.Title,
		Description:     payload.Description,
		Type:            payload.Type,
		Framework:       payload.Framework,
		ModelName:       payload.ModelName,
		ModelVersion:    payload.ModelVersion,
		UseCase:         payload.UseCase,
		Impact:          payload.Impact,
		Technologies:    payload.Technologies,
		Architecture:    payload.Architecture,
		Metrics:         payload.Metrics,
		Featured:        payload.Featured,
		DisplayOrder:    payload.DisplayOrder,
		DemoURL:         payload.DemoURL,
		DocumentationURL: payload.DocumentationURL,
		GitHubURL:       payload.GitHubURL,
	}

	if payload.ProjectID != "" {
		req.ProjectID = &payload.ProjectID
	}
	if payload.CaseStudyID != "" {
		req.CaseStudyID = &payload.CaseStudyID
	}

	integration, err := h.service.CreateAIMLIntegration(c.Context(), userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, integration)
}

func (h *handler) UpdateAIMLIntegration(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	integrationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid integration id",
		})
	}

	var payload updateAIMLIntegrationPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req := UpdateAIMLIntegrationRequest{}

	if payload.Title != nil {
		req.Title = payload.Title
	}
	if payload.Description != nil {
		req.Description = payload.Description
	}
	if payload.Type != nil {
		req.Type = payload.Type
	}
	if payload.Framework != nil {
		req.Framework = payload.Framework
	}
	if payload.ModelName != nil {
		req.ModelName = payload.ModelName
	}
	if payload.ModelVersion != nil {
		req.ModelVersion = payload.ModelVersion
	}
	if payload.UseCase != nil {
		req.UseCase = payload.UseCase
	}
	if payload.Impact != nil {
		req.Impact = payload.Impact
	}
	if payload.Technologies != nil {
		req.Technologies = payload.Technologies
	}
	if payload.Architecture != nil {
		req.Architecture = payload.Architecture
	}
	if payload.Metrics != nil {
		req.Metrics = payload.Metrics
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
	if payload.DemoURL != nil {
		req.DemoURL = payload.DemoURL
	}
	if payload.DocumentationURL != nil {
		req.DocumentationURL = payload.DocumentationURL
	}
	if payload.GitHubURL != nil {
		req.GitHubURL = payload.GitHubURL
	}

	integration, err := h.service.UpdateAIMLIntegration(c.Context(), userID, integrationID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, integration)
}

func (h *handler) GetAIMLIntegration(c *fiber.Ctx) error {
	integrationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid integration id",
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	integration, err := h.service.GetAIMLIntegration(c.Context(), integrationID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, integration)
}

func (h *handler) GetAIMLIntegrationPublic(c *fiber.Ctx) error {
	integrationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid integration id",
		})
	}

	integration, err := h.service.GetAIMLIntegrationPublic(c.Context(), integrationID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, integration)
}

func (h *handler) ListAIMLIntegrations(c *fiber.Ctx) error {
	var userID *uuid.UUID
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		uid := apiKey.UserID
		userID = &uid
	} else if uid, err := authdomain.UserIDFromContext(c); err == nil {
		userID = &uid
	}

	filters := ListAIMLIntegrationFilters{
		UserID: userID,
	}

	// Parse query parameters
	if typeStr := c.Query("type"); typeStr != "" {
		integrationType := IntegrationType(typeStr)
		filters.Type = &integrationType
	}

	if frameworkStr := c.Query("framework"); frameworkStr != "" {
		framework := Framework(frameworkStr)
		filters.Framework = &framework
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

	integrations, err := h.service.ListAIMLIntegrations(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, integrations)
}

func (h *handler) ListFeaturedAIMLIntegrations(c *fiber.Ctx) error {
	integrations, err := h.service.ListFeaturedAIMLIntegrations(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, integrations)
}

func (h *handler) GetIntegrationsByProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid project id",
		})
	}

	integrations, err := h.service.GetIntegrationsByProject(c.Context(), projectID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, integrations)
}

func (h *handler) GetIntegrationsByType(c *fiber.Ctx) error {
	integrationTypeStr := c.Params("type")
	integrationType := IntegrationType(integrationTypeStr)

	if !isValidIntegrationType(integrationType) {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidType, fiber.Map{
			"message": "invalid integration type",
		})
	}

	integrations, err := h.service.GetIntegrationsByType(c.Context(), integrationType)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, integrations)
}

func (h *handler) GetIntegrationsByFramework(c *fiber.Ctx) error {
	frameworkStr := c.Params("framework")
	framework := Framework(frameworkStr)

	if !isValidFramework(framework) {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidFramework, fiber.Map{
			"message": "invalid framework",
		})
	}

	integrations, err := h.service.GetIntegrationsByFramework(c.Context(), framework)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, integrations)
}

func (h *handler) DeleteAIMLIntegration(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	integrationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid integration id",
		})
	}

	if err := h.service.DeleteAIMLIntegration(c.Context(), userID, integrationID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "integration deleted successfully",
	})
}

// Helper functions

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	domainErr, ok := AsDomainError(err)
	if !ok {
		h.logger.Error("unexpected error in AI/ML integration handler", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{
			"message": "internal server error",
		})
	}

	statusCode := fiber.StatusInternalServerError
	switch domainErr.Code {
	case ErrCodeInvalidPayload, ErrCodeInvalidType, ErrCodeInvalidFramework, ErrCodeInvalidTitle:
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

