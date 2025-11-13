package reports

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes report-related endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler builds a Handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type generateSummaryPayload struct {
	UserID       string `json:"user_id"`
	SendEmail    bool   `json:"send_email"`
	EmailAddress string `json:"email_address"`
	SendWhatsApp bool   `json:"send_whatsapp"`
	PhoneNumber  string `json:"phone_number"`
}

// PostSummary generates a summary and dispatches it.
func (h *Handler) PostSummary(c *fiber.Ctx) error {
	var payload generateSummaryPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	summary, err := h.service.GenerateSummary(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	opts := DispatchOptions{
		SendEmail:    payload.SendEmail,
		EmailAddress: payload.EmailAddress,
		SendWhatsApp: payload.SendWhatsApp,
		PhoneNumber:  payload.PhoneNumber,
	}

	if err := h.service.DispatchSummary(c.Context(), summary, opts); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, summary)
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromError(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, nil)
	}

	h.logError("reports: unexpected error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func statusFromError(code int) int {
	switch code {
	case ErrCodeInvalidPayload:
		return fiber.StatusBadRequest
	case ErrCodeNotFound:
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
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
