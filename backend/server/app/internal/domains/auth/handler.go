package auth

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes HTTP endpoints for the auth domain.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler constructs a new Handler instance.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type registerPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Register handles user registration requests.
func (h *Handler) Register(c *fiber.Ctx) error {
	var payload registerPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	user, err := h.service.Register(c.Context(), RegisterRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, userResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

// Login handles authentication requests.
func (h *Handler) Login(c *fiber.Ctx) error {
	var payload loginPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	user, err := h.service.Authenticate(c.Context(), payload.Email, payload.Password)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, userResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromErrorCode(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, nil)
	}

	h.logError("auth: unexpected error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func statusFromErrorCode(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidEmail, ErrCodeInvalidPassword:
		return fiber.StatusBadRequest
	case ErrCodeEmailAlreadyExists:
		return fiber.StatusConflict
	case ErrCodeInvalidCredentials:
		return fiber.StatusUnauthorized
	case ErrCodeUserNotFound:
		return fiber.StatusNotFound
	case ErrCodePasswordHashFailure, ErrCodeRepositoryFailure, ErrCodeEmailDispatchFailure:
		return fiber.StatusInternalServerError
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
