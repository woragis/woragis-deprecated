package userpreferences

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes HTTP endpoints for user preferences.
type Handler struct {
	service Service
	logger  *slog.Logger
}

// NewHandler constructs a handler instance.
func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type updatePreferencesPayload struct {
	DefaultLanguage string `json:"defaultLanguage"`
	DefaultCurrency string `json:"defaultCurrency"`
}

type preferencesResponse struct {
	ID              string `json:"id"`
	UserID          string `json:"userId"`
	DefaultLanguage string `json:"defaultLanguage"`
	DefaultCurrency string `json:"defaultCurrency"`
	DefaultWebsite  string `json:"defaultWebsite"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// GetPreferences handles GET /user-preferences.
func (h *Handler) GetPreferences(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	prefs, err := h.service.GetOrCreatePreferences(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toPreferencesResponse(prefs))
}

// UpdatePreferences handles PATCH /user-preferences.
func (h *Handler) UpdatePreferences(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	var payload updatePreferencesPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	prefs, err := h.service.UpdatePreferences(c.Context(), userID, payload.DefaultLanguage, payload.DefaultCurrency)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toPreferencesResponse(prefs))
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromError(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, nil)
	}

	h.logError("userpreferences: unexpected error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func statusFromError(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidLanguage, ErrCodeInvalidCurrency:
		return fiber.StatusBadRequest
	case ErrCodeNotFound:
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
	}
}

func toPreferencesResponse(prefs *UserPreferences) preferencesResponse {
	return preferencesResponse{
		ID:              prefs.ID.String(),
		UserID:          prefs.UserID.String(),
		DefaultLanguage: prefs.DefaultLanguage,
		DefaultCurrency: prefs.DefaultCurrency,
		DefaultWebsite:  prefs.DefaultWebsite,
		CreatedAt:       prefs.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       prefs.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
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

