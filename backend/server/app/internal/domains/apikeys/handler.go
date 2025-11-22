package apikeys

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes API key endpoints.
type Handler interface {
	CreateAPIKey(c *fiber.Ctx) error
	UpdateAPIKey(c *fiber.Ctx) error
	DeleteAPIKey(c *fiber.Ctx) error
	GetAPIKey(c *fiber.Ctx) error
	ListAPIKeys(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs an API key handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type createAPIKeyPayload struct {
	Name      string  `json:"name"`
	ExpiresAt *string `json:"expiresAt,omitempty"` // ISO 8601 format
}

type updateAPIKeyPayload struct {
	Name string `json:"name"`
}

// Handlers

func (h *handler) CreateAPIKey(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var payload createAPIKeyPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var expiresAt *time.Time
	if payload.ExpiresAt != nil && *payload.ExpiresAt != "" {
		parsed, err := time.Parse(time.RFC3339, *payload.ExpiresAt)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
				"message": "invalid expiresAt format, use RFC3339",
			})
		}
		expiresAt = &parsed
	}

	apiKeyWithToken, err := h.service.CreateAPIKey(c.Context(), userID, payload.Name, expiresAt)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toAPIKeyWithTokenResponse(apiKeyWithToken))
}

func (h *handler) UpdateAPIKey(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	apiKeyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateAPIKeyPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	apiKey, err := h.service.UpdateAPIKey(c.Context(), userID, apiKeyID, payload.Name)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toAPIKeyResponse(apiKey))
}

func (h *handler) DeleteAPIKey(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	apiKeyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeleteAPIKey(c.Context(), userID, apiKeyID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "API key deleted"})
}

func (h *handler) GetAPIKey(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	apiKeyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	apiKey, err := h.service.GetAPIKey(c.Context(), userID, apiKeyID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toAPIKeyResponse(apiKey))
}

func (h *handler) ListAPIKeys(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	apiKeys, err := h.service.ListAPIKeys(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]apiKeyResponse, len(apiKeys))
	for i, apiKey := range apiKeys {
		responses[i] = toAPIKeyResponse(&apiKey)
	}

	return response.Success(c, fiber.StatusOK, responses)
}

// Response helpers

type apiKeyResponse struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Prefix     string     `json:"prefix"`
	UserID     string     `json:"userId"`
	LastUsedAt *string    `json:"lastUsedAt,omitempty"`
	ExpiresAt  *string    `json:"expiresAt,omitempty"`
	CreatedAt  string     `json:"createdAt"`
	UpdatedAt  string     `json:"updatedAt"`
}

type apiKeyWithTokenResponse struct {
	apiKeyResponse
	Token string `json:"token"` // Only included on creation
}

func toAPIKeyResponse(apiKey *APIKey) apiKeyResponse {
	resp := apiKeyResponse{
		ID:        apiKey.ID.String(),
		Name:      apiKey.Name,
		Prefix:    apiKey.Prefix,
		UserID:    apiKey.UserID.String(),
		CreatedAt: apiKey.CreatedAt.Format(time.RFC3339),
		UpdatedAt: apiKey.UpdatedAt.Format(time.RFC3339),
	}

	if apiKey.LastUsedAt != nil {
		lastUsed := apiKey.LastUsedAt.Format(time.RFC3339)
		resp.LastUsedAt = &lastUsed
	}

	if apiKey.ExpiresAt != nil {
		expires := apiKey.ExpiresAt.Format(time.RFC3339)
		resp.ExpiresAt = &expires
	}

	return resp
}

func toAPIKeyWithTokenResponse(apiKeyWithToken *APIKeyWithToken) apiKeyWithTokenResponse {
	return apiKeyWithTokenResponse{
		apiKeyResponse: toAPIKeyResponse(&apiKeyWithToken.APIKey),
		Token:          apiKeyWithToken.Token,
	}
}

// Error handling

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	domainErr, ok := AsDomainError(err)
	if !ok {
		h.logger.Error("unexpected error", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, nil)
	}

	statusCode := fiber.StatusInternalServerError
	switch domainErr.Code {
	case ErrCodeInvalidPayload, ErrCodeInvalidName:
		statusCode = fiber.StatusBadRequest
	case ErrCodeNotFound:
		statusCode = fiber.StatusNotFound
	case ErrCodeInvalidToken, ErrCodeExpiredToken:
		statusCode = fiber.StatusUnauthorized
	}

	return response.Error(c, statusCode, domainErr.Code, domainErr.Message)
}

func unauthorizedResponse(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusUnauthorized, 0, "unauthorized")
}

