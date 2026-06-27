package platforms

import (
	"encoding/json"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"woragis-social-media-service/pkg/middleware"
	"woragis-social-media-service/pkg/response"
)

// Handler exposes platform configuration endpoints.
type Handler interface {
	ListConfigs(c *fiber.Ctx) error
	GetConfig(c *fiber.Ctx) error
	GetConfigByName(c *fiber.Ctx) error
	UpdateConfig(c *fiber.Ctx) error
	GetOptimalTimes(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a platform handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type updateConfigPayload struct {
	DisplayName      *string                        `json:"displayName,omitempty"`
	PostingFrequency *int                           `json:"postingFrequency,omitempty"`
	BestDays         []string                       `json:"bestDays,omitempty"`
	BestTimes        []string                       `json:"bestTimes,omitempty"`
	SupportedFormats []ContentFormat `json:"supportedFormats,omitempty"`
	IsActive         *bool                          `json:"isActive,omitempty"`
}

// Handlers

func (h *handler) ListConfigs(c *fiber.Ctx) error {
	activeOnly := c.Query("activeOnly", "false") == "true"

	configs, err := h.service.ListConfigs(c.Context(), activeOnly)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]configResponse, len(configs))
	for i := range configs {
		responses[i] = toConfigResponse(&configs[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) GetConfig(c *fiber.Ctx) error {
	configID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	config, err := h.service.GetConfig(c.Context(), configID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toConfigResponse(config))
}

func (h *handler) GetConfigByName(c *fiber.Ctx) error {
	platformName := Platform(c.Params("name"))
	if platformName == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	config, err := h.service.GetConfigByName(c.Context(), platformName)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toConfigResponse(config))
}

func (h *handler) UpdateConfig(c *fiber.Ctx) error {
	_, err := middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	configID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateConfigPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	config, err := h.service.UpdateConfig(c.Context(), UpdateConfigRequest{
		ConfigID:         configID,
		DisplayName:      payload.DisplayName,
		PostingFrequency: payload.PostingFrequency,
		BestDays:         payload.BestDays,
		BestTimes:        payload.BestTimes,
		SupportedFormats: payload.SupportedFormats,
		IsActive:         payload.IsActive,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toConfigResponse(config))
}

func (h *handler) GetOptimalTimes(c *fiber.Ctx) error {
	platformName := Platform(c.Params("name"))
	if platformName == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	optimalTimes, err := h.service.GetOptimalTimes(c.Context(), platformName)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, optimalTimes)
}

// Response helpers

type configResponse struct {
	ID               string                        `json:"id"`
	Name             Platform      `json:"name"`
	DisplayName      string        `json:"displayName"`
	PostingFrequency *int          `json:"postingFrequency,omitempty"`
	BestDays         []string      `json:"bestDays,omitempty"`
	BestTimes        []string      `json:"bestTimes,omitempty"`
	SupportedFormats []ContentFormat `json:"supportedFormats"`
	IsActive         bool                          `json:"isActive"`
	CreatedAt        string                        `json:"createdAt"`
	UpdatedAt        string                        `json:"updatedAt"`
}

func toConfigResponse(config *PlatformConfig) configResponse {
	response := configResponse{
		ID:               config.ID.String(),
		Name:             config.Name,
		DisplayName:      config.DisplayName,
		PostingFrequency: config.PostingFrequency,
		IsActive:         config.IsActive,
		CreatedAt:        config.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        config.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Parse supported formats
	var formats []ContentFormat
	if len(config.SupportedFormats) > 0 {
		_ = json.Unmarshal(config.SupportedFormats, &formats)
	}
	response.SupportedFormats = formats

	// Parse best days
	if len(config.BestDays) > 0 {
		var days []string
		if err := json.Unmarshal(config.BestDays, &days); err == nil {
			response.BestDays = days
		}
	}

	// Parse best times
	if len(config.BestTimes) > 0 {
		var times []string
		if err := json.Unmarshal(config.BestTimes, &times); err == nil {
			response.BestTimes = times
		}
	}

	return response
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
	case ErrCodeInvalidPayload, ErrCodeInvalidPlatform:
		statusCode = fiber.StatusBadRequest
	case ErrCodeNotFound:
		statusCode = fiber.StatusNotFound
	case ErrCodeConflict:
		statusCode = fiber.StatusConflict
	case ErrCodeRepositoryFailure:
		statusCode = fiber.StatusInternalServerError
	}

	return response.Error(c, statusCode, domainErr.Code, fiber.Map{"message": domainErr.Message})
}

func unauthorizedResponse(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "unauthorized"})
}
