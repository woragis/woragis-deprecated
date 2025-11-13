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

type requestResetPayload struct {
	Email string `json:"email"`
}

type resetConfirmPayload struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Token     string    `json:"token,omitempty"`
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

	token, err := h.service.IssueToken(c.Context(), user)
	if err != nil {
		h.logError("auth: token issuance failed", err)
		return response.Error(c, fiber.StatusInternalServerError, ErrCodeTokenIssuanceFailure, nil)
	}

	return response.Success(c, fiber.StatusCreated, userResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Token:     token,
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

	token, err := h.service.IssueToken(c.Context(), user)
	if err != nil {
		h.logError("auth: token issuance failed", err)
		return response.Error(c, fiber.StatusInternalServerError, ErrCodeTokenIssuanceFailure, nil)
	}

	return response.Success(c, fiber.StatusOK, userResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Token:     token,
	})
}

// RequestPasswordReset handles POST /auth/password/reset/request.
func (h *Handler) RequestPasswordReset(c *fiber.Ctx) error {
	var payload requestResetPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.RequestPasswordReset(c.Context(), PasswordResetRequest{Email: payload.Email}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "email dispatched"})
}

// ConfirmPasswordReset handles POST /auth/password/reset/confirm.
func (h *Handler) ConfirmPasswordReset(c *fiber.Ctx) error {
	var payload resetConfirmPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.ResetPassword(c.Context(), PasswordResetConfirmRequest{Token: payload.Token, Password: payload.Password}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "password updated"})
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
	case ErrCodeInvalidToken:
		return fiber.StatusUnauthorized
	case ErrCodeTokenIssuanceFailure:
		return fiber.StatusInternalServerError
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
