package auth

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"strings"
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

type oauthStartPayload struct {
	Provider          string `json:"provider"`
	Mode              string `json:"mode"`
	RedirectOrigin    string `json:"redirect_origin"`
	DeviceFingerprint string `json:"device_fingerprint"`
	DeviceName        string `json:"device_name"`
}

type oauthProviderResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type oauthAccountResponse struct {
	Provider  string    `json:"provider"`
	LinkedAt  time.Time `json:"linked_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Scopes    []string  `json:"scopes"`
}

type oauthCallbackMessage struct {
	Type     string `json:"type"`
	Provider string `json:"provider"`
	Mode     string `json:"mode"`
	Success  bool   `json:"success"`
	Message  string `json:"message,omitempty"`
	Payload  any    `json:"payload,omitempty"`
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

	return response.Success(c, fiber.StatusOK, toLoginResponsePayload(resp))
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

// ListOAuthProviders returns configured OAuth providers.
func (h *Handler) ListOAuthProviders(c *fiber.Ctx) error {
	providers := h.service.ListOAuthProviders()
	result := make([]oauthProviderResponse, 0, len(providers))
	for _, provider := range providers {
		result = append(result, oauthProviderResponse{
			ID:   string(provider.ID),
			Name: provider.Name,
		})
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"providers": result})
}

// StartOAuth initialises an OAuth flow and returns an authorisation URL.
func (h *Handler) StartOAuth(c *fiber.Ctx) error {
	var payload oauthStartPayload
	if err := c.BodyParser(&payload); err != nil {
		h.logError(ErrMalformedPayload, err)
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "invalid payload"})
	}

	provider := parseOAuthProvider(payload.Provider)
	if provider == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "unknown provider"})
	}

	mode := strings.ToLower(strings.TrimSpace(payload.Mode))
	if mode == "" {
		mode = string(OAuthModeLogin)
	}

	flowMode := OAuthFlowMode(mode)
	if flowMode != OAuthModeLogin && flowMode != OAuthModeLink {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "invalid oauth mode"})
	}

	options := OAuthStartOptions{
		Mode:              flowMode,
		RedirectOrigin:    strings.TrimSpace(payload.RedirectOrigin),
		DeviceFingerprint: strings.TrimSpace(payload.DeviceFingerprint),
		DeviceName:        strings.TrimSpace(payload.DeviceName),
	}

	if flowMode == OAuthModeLink {
		userID, err := currentUserID(c)
		if err != nil {
			return err
		}
		options.UserID = &userID
	}

	state, url, err := h.service.BeginOAuth(c.Context(), provider, options)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"provider":          string(provider),
		"mode":              flowMode,
		"authorization_url": url,
		"state":             state,
	})
}

// OAuthCallback completes the OAuth flow when the provider redirects back.
func (h *Handler) OAuthCallback(c *fiber.Ctx) error {
	provider := parseOAuthProvider(c.Params("provider"))
	redirectOrigin := strings.TrimSpace(c.Query("redirect_origin"))

	message := oauthCallbackMessage{
		Type:     "oauth:result",
		Provider: string(provider),
		Mode:     string(OAuthModeLogin),
		Success:  false,
	}

	if provider == "" {
		message.Message = "Unknown OAuth provider."
		return h.renderOAuthCallback(c, redirectOrigin, message)
	}

	if errParam := strings.TrimSpace(c.Query("error")); errParam != "" {
		message.Message = c.Query("error_description", errParam)
		return h.renderOAuthCallback(c, redirectOrigin, message)
	}

	state := c.Query("state")
	code := c.Query("code")
	if strings.TrimSpace(state) == "" || strings.TrimSpace(code) == "" {
		message.Message = "Missing OAuth parameters."
		return h.renderOAuthCallback(c, redirectOrigin, message)
	}

	result, err := h.service.CompleteOAuth(c.Context(), OAuthCallbackInput{
		Provider:  provider,
		State:     state,
		Code:      code,
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	})

	if result != nil {
		if redirectOrigin == "" {
			redirectOrigin = result.RedirectOrigin
		}
		message.Mode = string(result.Mode)
		message.Success = result.Success && err == nil
		if result.Message != "" {
			message.Message = result.Message
		}
		if result.Login != nil && result.Mode == OAuthModeLogin && message.Success {
			payload := toLoginResponsePayload(result.Login)
			message.Payload = payload
		} else if result.Success && result.Mode == OAuthModeLink {
			message.Payload = fiber.Map{"status": "linked"}
		}
	}

	if err != nil {
		if domainErr, ok := AsDomainError(err); ok {
			message.Message = domainErr.Message
		} else if message.Message == "" {
			message.Message = "Unable to complete OAuth flow."
		}
		message.Success = false
	}

	if message.Message == "" && !message.Success {
		message.Message = "OAuth flow did not complete successfully."
	}

	return h.renderOAuthCallback(c, redirectOrigin, message)
}

// ListOAuthAccounts enumerates linked OAuth providers for the current user.
func (h *Handler) ListOAuthAccounts(c *fiber.Ctx) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	accounts, err := h.service.ListOAuthAccounts(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	result := make([]oauthAccountResponse, 0, len(accounts))
	for _, account := range accounts {
		result = append(result, oauthAccountResponse{
			Provider:  string(account.Provider),
			LinkedAt:  account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
			Scopes:    splitScopeString(account.Scopes),
		})
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"accounts": result})
}

// UnlinkOAuthAccount removes an OAuth association for the authenticated user.
func (h *Handler) UnlinkOAuthAccount(c *fiber.Ctx) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}

	provider := parseOAuthProvider(c.Params("provider"))
	if provider == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "unknown provider"})
	}

	if err := h.service.UnlinkOAuthAccount(c.Context(), userID, provider); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "provider_unlinked"})
}

func (h *Handler) renderOAuthCallback(c *fiber.Ctx, origin string, payload oauthCallbackMessage) error {
	if origin == "" {
		origin = h.service.publicURL
	}

	payload.Provider = strings.TrimSpace(payload.Provider)
	payload.Mode = strings.TrimSpace(payload.Mode)

	data, err := json.Marshal(payload)
	if err != nil {
		h.logError("auth: oauth callback marshal failed", err)
		data = []byte(`{"type":"oauth:result","success":false,"message":"internal error"}`)
	}

	escapedPayload := template.JSEscapeString(string(data))
	escapedOrigin := template.JSEscapeString(origin)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8" />
	<title>OAuth Completed</title>
</head>
<body>
<script>
(function() {
	const data = JSON.parse("%s");
	const targetOrigin = "%s";
	try {
		if (window.opener && !window.opener.closed) {
			window.opener.postMessage(data, targetOrigin);
		}
	} catch (err) {
		console.error('oauth callback error', err);
	}
	window.close();
	setTimeout(function () { window.close(); }, 500);
})();
</script>
<p>Authentication flow complete. You can close this window.</p>
</body>
</html>`, escapedPayload, escapedOrigin)

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(html)
}

func splitScopeString(scopes string) []string {
	fields := strings.Fields(scopes)
	if len(fields) == 0 {
		return []string{}
	}
	return fields
}

func parseOAuthProvider(value string) OAuthProvider {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}
	return OAuthProvider(value)
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
	case ErrCodeInvalidToken, ErrCodeSessionRevoked, ErrCodeSessionExpired, ErrCodeMFARequired, ErrCodeMFAInvalid, ErrCodeOAuthStateInvalid:
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

func toLoginResponsePayload(resp *LoginResponse) loginResponse {
	if resp == nil {
		return loginResponse{}
	}

	return loginResponse{
		User:         toUserResponse(resp.User),
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		SessionID:    resp.SessionID.String(),
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
