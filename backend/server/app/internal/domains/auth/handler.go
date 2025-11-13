package auth

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

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
	Email       string `json:"email"`
	Password    string `json:"password"`
	Locale      string `json:"locale"`
	DisplayName string `json:"display_name"`
}

type confirmPayload struct {
	Token string `json:"token"`
}

type resendPayload struct {
	Email string `json:"email"`
}

type loginPayload struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	DeviceFingerprint string `json:"device_fingerprint"`
	DeviceName        string `json:"device_name"`
	MFACode           string `json:"mfa_code"`
	UserAgent         string `json:"user_agent"`
	IPAddress         string `json:"ip_address"`
}

type refreshPayload struct {
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent"`
	IPAddress    string `json:"ip_address"`
}

type logoutPayload struct {
	SessionID string `json:"session_id"`
}

type requestResetPayload struct {
	Email string `json:"email"`
}

type resetConfirmPayload struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type revokeSessionsPayload struct {
	KeepSessionID string `json:"keep_session_id"`
}

type enableMFAPayload struct {
	Issuer string `json:"issuer"`
	Label  string `json:"label"`
	Code   string `json:"code"`
}

type verifyMFAPayload struct {
	Code string `json:"code"`
}

type userResponse struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	CreatedAt       time.Time `json:"created_at"`
	EmailConfirmed  bool      `json:"email_confirmed"`
	MFAEnabled      bool      `json:"mfa_enabled"`
	PreferredLocale string    `json:"preferred_locale,omitempty"`
	Role            string    `json:"role,omitempty"`
}

type loginResponse struct {
	User         userResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	SessionID    string       `json:"session_id"`
}

type sessionResponse struct {
	ID         string    `json:"id"`
	DeviceID   string    `json:"device_id"`
	UserAgent  string    `json:"user_agent"`
	IP         string    `json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	LastSeenAt time.Time `json:"last_seen_at"`
	IsRevoked  bool      `json:"is_revoked"`
}

type enableMFAResponse struct {
	Secret          string   `json:"secret"`
	BackupCodes     []string `json:"backup_codes"`
	ProvisioningURI string   `json:"provisioning_uri"`
}

// Register handles user registration requests.
func (h *Handler) Register(c *fiber.Ctx) error {
	var payload registerPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	user, err := h.service.Register(c.Context(), RegisterRequest{
		Email:       payload.Email,
		Password:    payload.Password,
		Locale:      payload.Locale,
		DisplayName: payload.DisplayName,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusAccepted, fiber.Map{
		"status": "confirmation_required",
		"user":   toUserResponse(user),
	})
}

// ConfirmEmail validates email confirmation tokens.
func (h *Handler) ConfirmEmail(c *fiber.Ctx) error {
	var payload confirmPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	user, err := h.service.ConfirmEmail(c.Context(), ConfirmEmailRequest{Token: payload.Token})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"user": toUserResponse(user),
	})
}

// ResendConfirmation triggers a new confirmation email.
func (h *Handler) ResendConfirmation(c *fiber.Ctx) error {
	var payload resendPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.ResendConfirmation(c.Context(), ResendConfirmationRequest{Email: payload.Email}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "confirmation_email_resent"})
}

// Login handles authentication requests.
func (h *Handler) Login(c *fiber.Ctx) error {
	var payload loginPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	ip := payload.IPAddress
	if ip == "" {
		ip = c.IP()
	}
	userAgent := payload.UserAgent
	if userAgent == "" {
		userAgent = c.Get("User-Agent")
	}

	resp, err := h.service.Login(c.Context(), LoginRequest{
		Email:             payload.Email,
		Password:          payload.Password,
		DeviceFingerprint: payload.DeviceFingerprint,
		DeviceName:        payload.DeviceName,
		MFACode:           payload.MFACode,
		IPAddress:         ip,
		UserAgent:         userAgent,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, loginResponse{
		User:         toUserResponse(resp.User),
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		SessionID:    resp.SessionID.String(),
	})
}

// RefreshSession rotates session tokens.
func (h *Handler) RefreshSession(c *fiber.Ctx) error {
	var payload refreshPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	ip := payload.IPAddress
	if ip == "" {
		ip = c.IP()
	}
	userAgent := payload.UserAgent
	if userAgent == "" {
		userAgent = c.Get("User-Agent")
	}

	resp, err := h.service.RefreshSession(c.Context(), RefreshSessionRequest{
		RefreshToken: payload.RefreshToken,
		IPAddress:    ip,
		UserAgent:    userAgent,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"access_token":  resp.AccessToken,
		"refresh_token": resp.RefreshToken,
		"session_id":    resp.SessionID.String(),
	})
}

// Logout revokes a session.
func (h *Handler) Logout(c *fiber.Ctx) error {
	var payload logoutPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	sessionID, err := uuid.Parse(payload.SessionID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "invalid session_id"})
	}

	if err := h.service.Logout(c.Context(), LogoutRequest{SessionID: sessionID}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "session_revoked"})
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

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "reset_email_dispatched"})
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

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "password_updated"})
}

// ListSessions returns active sessions for the current user.
func (h *Handler) ListSessions(c *fiber.Ctx) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	sessions, err := h.service.ListActiveSessions(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]sessionResponse, 0, len(sessions))
	for _, session := range sessions {
		responses = append(responses, sessionResponse{
			ID:         session.ID.String(),
			DeviceID:   session.DeviceID.String(),
			UserAgent:  session.UserAgent,
			IP:         session.IP,
			CreatedAt:  session.CreatedAt,
			ExpiresAt:  session.ExpiresAt,
			LastSeenAt: session.LastSeenAt,
			IsRevoked:  session.RevokedAt != nil,
		})
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"sessions": responses})
}

// RevokeOtherSessions revokes all sessions except the specified one.
func (h *Handler) RevokeOtherSessions(c *fiber.Ctx) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	var payload revokeSessionsPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var keep uuid.UUID
	if payload.KeepSessionID != "" {
		if keep, err = uuid.Parse(payload.KeepSessionID); err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "invalid keep_session_id"})
		}
	}

	if err := h.service.RevokeOtherSessions(c.Context(), userID, keep); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "sessions_revoked"})
}

// EnableMFA starts MFA enrolment for the current user.
func (h *Handler) EnableMFA(c *fiber.Ctx) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	var payload enableMFAPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	resp, err := h.service.EnableMFA(c.Context(), EnableMFARequest{
		UserID: userID,
		Issuer: payload.Issuer,
		Label:  payload.Label,
		Code:   payload.Code,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, enableMFAResponse{
		Secret:          resp.Secret,
		BackupCodes:     resp.BackupCodes,
		ProvisioningURI: resp.ProvisioningURI,
	})
}

// VerifyMFA finalises TOTP enrolment.
func (h *Handler) VerifyMFA(c *fiber.Ctx) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	var payload verifyMFAPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.VerifyMFACode(c.Context(), VerifyMFACodeRequest{
		UserID: userID,
		Code:   payload.Code,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "mfa_enabled"})
}

// DisableMFA revokes MFA enrolment.
func (h *Handler) DisableMFA(c *fiber.Ctx) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	if err := h.service.DisableMFA(c.Context(), DisableMFARequest{UserID: userID}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "mfa_disabled"})
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromErrorCode(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, fiber.Map{"message": domainErr.Message})
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
	case ErrCodeInvalidToken, ErrCodeSessionRevoked, ErrCodeSessionExpired, ErrCodeMFARequired, ErrCodeMFAInvalid:
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

func toUserResponse(user *User) userResponse {
	if user == nil {
		return userResponse{}
	}

	return userResponse{
		ID:              user.ID.String(),
		Email:           user.Email,
		CreatedAt:       user.CreatedAt,
		EmailConfirmed:  user.EmailConfirmedAt != nil,
		MFAEnabled:      user.MFAEnabled,
		PreferredLocale: user.PreferredLocale,
		Role:            user.Role,
	}
}

func currentUserID(c *fiber.Ctx) (uuid.UUID, error) {
	raw := c.Locals(ContextUserIDKey)
	if raw == nil {
		return uuid.Nil, response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidToken, fiber.Map{"message": "unauthenticated"})
	}

	userID, err := uuid.Parse(raw.(string))
	if err != nil {
		return uuid.Nil, response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidToken, fiber.Map{"message": "invalid subject"})
	}

	return userID, nil
}
