package auth

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"

	"github.com/woragis/backend/server/app/internal/monitoring"
	emailservice "github.com/woragis/backend/server/app/internal/services/email"
	emailtemplates "github.com/woragis/backend/server/app/internal/services/email/templates"
)

const (
	defaultPasswordResetTTL     = 30 * time.Minute
	defaultEmailConfirmationTTL = 48 * time.Hour
	defaultSessionTTL           = 30 * 24 * time.Hour
	defaultOAuthStateTTL        = 10 * time.Minute
	refreshTokenEntropyBytes    = 48
	minPasswordLength           = 8
	maxPasswordLength           = 72
)

// Service exposes use-cases for the auth domain.
type Service struct {
	repo             Repository
	emailSender      emailservice.Sender
	tokenStore       TokenStore
	monitor          monitoring.Tracker
	logger           *slog.Logger
	publicURL        string
	passwordResetTTL time.Duration
	confirmationTTL  time.Duration
	sessionTTL       time.Duration
	jwtManager       *JWTManager
	renderer         *emailtemplates.Renderer
	httpClient       *http.Client
	oauthProviders   map[OAuthProvider]*oauthProviderConfig
	oauthStates      map[string]oauthStateData
	oauthMu          sync.Mutex
	oauthStateTTL    time.Duration
}

// NewService wires a Service with its collaborators.
func NewService(
	repo Repository,
	emailSender emailservice.Sender,
	tokenStore TokenStore,
	monitor monitoring.Tracker,
	publicURL string,
	jwtManager *JWTManager,
	logger *slog.Logger,
) *Service {
	if publicURL == "" {
		publicURL = "http://localhost:8080"
	}

	return &Service{
		repo:             repo,
		emailSender:      emailSender,
		tokenStore:       tokenStore,
		monitor:          monitor,
		logger:           logger,
		publicURL:        strings.TrimRight(publicURL, "/"),
		passwordResetTTL: defaultPasswordResetTTL,
		confirmationTTL:  defaultEmailConfirmationTTL,
		sessionTTL:       defaultSessionTTL,
		jwtManager:       jwtManager,
		renderer:         emailtemplates.NewRenderer("en"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		oauthProviders: make(map[OAuthProvider]*oauthProviderConfig),
		oauthStates:    make(map[string]oauthStateData),
		oauthStateTTL:  defaultOAuthStateTTL,
	}
}

// RegisterRequest transports user registration data.
type RegisterRequest struct {
	Email       string
	Password    string
	Locale      string
	DisplayName string
}

// ConfirmEmailRequest contains a confirmation token payload.
type ConfirmEmailRequest struct {
	Token string
}

// ResendConfirmationRequest triggers confirmation re-dispatch.
type ResendConfirmationRequest struct {
	Email string
}

// LoginRequest includes credentials and device metadata.
type LoginRequest struct {
	Email             string
	Password          string
	DeviceFingerprint string
	DeviceName        string
	IPAddress         string
	UserAgent         string
	MFACode           string
}

// LoginResponse returns authentication artefacts.
type LoginResponse struct {
	User         *User
	AccessToken  string
	RefreshToken string
	SessionID    uuid.UUID
}

// LogoutRequest revokes a particular session.
type LogoutRequest struct {
	SessionID uuid.UUID
}

// RefreshSessionRequest rotates tokens.
type RefreshSessionRequest struct {
	RefreshToken string
	IPAddress    string
	UserAgent    string
}

// RefreshSessionResponse is returned after successful rotation.
type RefreshSessionResponse struct {
	AccessToken  string
	RefreshToken string
	SessionID    uuid.UUID
}

// PasswordResetRequest is used to initiate a reset flow.
type PasswordResetRequest struct {
	Email string
}

// PasswordResetConfirmRequest finalises password change.
type PasswordResetConfirmRequest struct {
	Token    string
	Password string
}

// EnableMFARequest starts MFA enrolment.
type EnableMFARequest struct {
	UserID uuid.UUID
	Issuer string
	Label  string
	Code   string
}

// EnableMFAResponse returns generated secrets.
type EnableMFAResponse struct {
	Secret          string
	BackupCodes     []string
	ProvisioningURI string
}

// VerifyMFACodeRequest validates a TOTP code.
type VerifyMFACodeRequest struct {
	UserID uuid.UUID
	Code   string
}

// DisableMFARequest cancels MFA enrolment.
type DisableMFARequest struct {
	UserID uuid.UUID
}

// OAuthLinkRequest links an external identity.
type OAuthLinkRequest struct {
	UserID         uuid.UUID
	Provider       OAuthProvider
	ProviderUserID string
	AccessToken    string
	RefreshToken   string
	ExpiresAt      *time.Time
	Scopes         []string
}

// BulkUserUpdateRequest enables mass administrative updates.
type BulkUserUpdateRequest struct {
	UserIDs      []uuid.UUID
	SetRole      *string
	DisableMFA   bool
	ConfirmEmail bool
}

// Register registers a new user account and dispatches confirmation e-mail.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	return s.RegisterUser(ctx, req)
}

// RegisterUser registers a new account enforcing validation and confirmation.
func (s *Service) RegisterUser(ctx context.Context, req RegisterRequest) (*User, error) {
	email := normalizeEmail(req.Email)
	if email == "" {
		return nil, NewDomainError(ErrCodeInvalidEmail, ErrEmptyEmail)
	}

	if err := validatePassword(req.Password); err != nil {
		return nil, err
	}

	if existing, err := s.repo.FindByEmail(ctx, email); err == nil && existing != nil {
		return nil, NewDomainError(ErrCodeEmailAlreadyExists, ErrUserAlreadyExists)
	} else if err != nil {
		if domainErr, ok := AsDomainError(err); ok {
			if domainErr.Code != ErrCodeUserNotFound {
				return nil, domainErr
			}
		} else {
			return nil, err
		}
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := NewUser(email, hash)
	if err != nil {
		return nil, err
	}

	if locale := strings.TrimSpace(req.Locale); locale != "" {
		user.PreferredLocale = locale
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	if s.monitor != nil {
		s.monitor.RecordUserRegistration(ctx, user.ID)
	}
	_ = s.recordAudit(ctx, &user.ID, AuditActionUserRegistered, nil, "", "")

	if err := s.dispatchConfirmationEmail(ctx, user); err != nil && s.logger != nil {
		s.logger.Warn("auth: confirmation email dispatch failed", slog.String("email", user.Email), slog.Any("error", err))
	}

	if err := s.sendWelcomeEmail(ctx, user); err != nil && s.logger != nil {
		s.logger.Warn("auth: welcome email dispatch failed", slog.String("email", user.Email), slog.Any("error", err))
	}

	return user, nil
}

// ConfirmEmail marks the user email as verified using the provided token.
func (s *Service) ConfirmEmail(ctx context.Context, req ConfirmEmailRequest) (*User, error) {
	token := strings.TrimSpace(req.Token)
	if token == "" {
		return nil, NewDomainError(ErrCodeInvalidToken, ErrEmptyEmailToken)
	}

	emailToken, err := s.repo.FindEmailToken(ctx, EmailTokenTypeConfirmation, hashToken(token))
	if err != nil {
		return nil, err
	}

	if emailToken.IsExpired(time.Now()) {
		return nil, NewDomainError(ErrCodeInvalidToken, ErrInvalidResetToken)
	}

	user, err := s.repo.FindByID(ctx, emailToken.UserID)
	if err != nil {
		return nil, err
	}

	if user.EmailConfirmedAt != nil {
		return nil, NewDomainError(ErrCodeEmailAlreadyExists, ErrEmailAlreadyConfirmed)
	}

	user.ConfirmEmail()
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	emailToken.Consume()
	_ = s.repo.UpdateEmailToken(ctx, emailToken)
	_ = s.repo.DeleteEmailTokensByUser(ctx, user.ID, EmailTokenTypeConfirmation)

	_ = s.recordAudit(ctx, &user.ID, AuditActionEmailConfirmed, nil, "", "")

	return user, nil
}

// ResendConfirmation issues a new confirmation email if the user isn't verified.
func (s *Service) ResendConfirmation(ctx context.Context, req ResendConfirmationRequest) error {
	email := normalizeEmail(req.Email)
	if email == "" {
		return NewDomainError(ErrCodeInvalidEmail, ErrEmptyEmail)
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}

	if user.EmailConfirmedAt != nil {
		return NewDomainError(ErrCodeEmailAlreadyExists, ErrEmailAlreadyConfirmed)
	}

	return s.dispatchConfirmationEmail(ctx, user)
}

// Authenticate maintains backward compatibility with legacy login flow.
func (s *Service) Authenticate(ctx context.Context, email, password string) (*User, error) {
	resp, err := s.Login(ctx, LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

// Login authenticates credentials, enforces MFA, and creates a session.
func (s *Service) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	email := normalizeEmail(req.Email)
	if email == "" || strings.TrimSpace(req.Password) == "" {
		return nil, NewDomainError(ErrCodeInvalidCredentials, ErrCredentialsMismatch)
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user.EmailConfirmedAt == nil {
		return nil, NewDomainError(ErrCodeEmailNotConfirmed, ErrEmailNotConfirmed)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		_ = s.recordAudit(ctx, &user.ID, AuditActionLoginFailed, map[string]any{"email": email}, req.IPAddress, req.UserAgent)
		return nil, NewDomainError(ErrCodeInvalidCredentials, ErrCredentialsMismatch)
	}

	if user.MFAEnabled {
		if err := s.validateMFACode(ctx, user.ID, req.MFACode); err != nil {
			return nil, err
		}
	}

	resp, err := s.issueLoginArtifacts(ctx, user, req.DeviceFingerprint, req.DeviceName, req.IPAddress, req.UserAgent, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Logout revokes an active session.
func (s *Service) Logout(ctx context.Context, req LogoutRequest) error {
	if req.SessionID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptySessionID)
	}

	if err := s.repo.RevokeSession(ctx, req.SessionID); err != nil {
		return err
	}

	_ = s.recordAudit(ctx, nil, AuditActionSessionRevoked, map[string]any{"session_id": req.SessionID.String()}, "", "")
	return nil
}

// RefreshSession rotates refresh token and issues a new access token.
func (s *Service) RefreshSession(ctx context.Context, req RefreshSessionRequest) (*RefreshSessionResponse, error) {
	rawToken := strings.TrimSpace(req.RefreshToken)
	if rawToken == "" {
		return nil, NewDomainError(ErrCodeInvalidToken, ErrEmptyRefreshTokenHash)
	}

	session, err := s.repo.FindSessionByRefreshHash(ctx, hashToken(rawToken))
	if err != nil {
		return nil, err
	}

	if !session.IsActive(time.Now()) {
		return nil, NewDomainError(ErrCodeSessionExpired, ErrSessionExpired)
	}

	user, err := s.repo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := generateSecureToken(refreshTokenEntropyBytes)
	if err != nil {
		return nil, NewDomainError(ErrCodeTokenIssuanceFailure, ErrUnableToIssueToken)
	}

	session.RefreshTokenHash = hashToken(newRefreshToken)
	session.ExpiresAt = time.Now().Add(s.sessionTTL)
	session.Touch(req.IPAddress, req.UserAgent)

	if err := s.repo.UpdateSession(ctx, session); err != nil {
		return nil, err
	}

	accessToken, err := s.issueAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &RefreshSessionResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		SessionID:    session.ID,
	}, nil
}

// RequestPasswordReset sends a reset link to the supplied email.
func (s *Service) RequestPasswordReset(ctx context.Context, req PasswordResetRequest) error {
	email := normalizeEmail(req.Email)
	if email == "" {
		return NewDomainError(ErrCodeInvalidEmail, ErrEmptyEmail)
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		if s.logger != nil {
			s.logger.Info("auth: password reset requested for unknown email", slog.String("email", email))
		}
		return nil
	}

	if s.tokenStore == nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrInvalidResetToken)
	}

	token, err := s.tokenStore.CreateToken(ctx, user.ID, s.passwordResetTTL)
	if err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrInvalidResetToken)
	}

	if err := s.dispatchPasswordResetEmail(ctx, user, token); err != nil {
		return err
	}

	_ = s.recordAudit(ctx, &user.ID, AuditActionPasswordResetRequested, nil, "", "")
	return nil
}

// ResetPassword maintains backward compatibility with legacy handler.
func (s *Service) ResetPassword(ctx context.Context, req PasswordResetConfirmRequest) error {
	return s.CompletePasswordReset(ctx, req)
}

// CompletePasswordReset validates the token and updates the password hash.
func (s *Service) CompletePasswordReset(ctx context.Context, req PasswordResetConfirmRequest) error {
	if err := validatePassword(req.Password); err != nil {
		return err
	}

	if s.tokenStore == nil {
		return NewDomainError(ErrCodeInvalidToken, ErrInvalidResetToken)
	}

	userID, err := s.tokenStore.ValidateToken(ctx, req.Token)
	if err != nil {
		return NewDomainError(ErrCodeInvalidToken, ErrInvalidResetToken)
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePasswordHash(ctx, userID, hash); err != nil {
		return err
	}

	_ = s.tokenStore.InvalidateToken(ctx, req.Token)

	user, err := s.repo.FindByID(ctx, userID)
	if err == nil {
		_ = s.sendPasswordChangedNotice(ctx, user)
		_ = s.recordAudit(ctx, &user.ID, AuditActionPasswordResetCompleted, nil, "", "")
	}

	return nil
}

// IssueToken generates a signed JWT for the provided user.
func (s *Service) IssueToken(ctx context.Context, user *User) (string, error) {
	if user == nil {
		return "", NewDomainError(ErrCodeInvalidPayload, ErrNilUser)
	}
	return s.issueAccessToken(user.ID, user.Email)
}

// ListActiveSessions returns active sessions for a user.
func (s *Service) ListActiveSessions(ctx context.Context, userID uuid.UUID) ([]Session, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	return s.repo.ListActiveSessions(ctx, userID)
}

// RevokeOtherSessions revokes all sessions except the provided one.
func (s *Service) RevokeOtherSessions(ctx context.Context, userID, exclude uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	return s.repo.RevokeUserSessions(ctx, userID, exclude)
}

// EnableMFA generates a new TOTP secret and optionally activates immediately.
func (s *Service) EnableMFA(ctx context.Context, req EnableMFARequest) (*EnableMFAResponse, error) {
	if req.UserID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	user, err := s.repo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      pickString(req.Issuer, "Woragis"),
		AccountName: pickString(req.Label, user.Email),
	})
	if err != nil {
		return nil, NewDomainError(ErrCodeMFAGenerationFailure, ErrUnableToGenerateMFASecret)
	}

	token, err := NewMFAToken(user.ID, MFATypeTOTP, key.Issuer(), key.AccountName())
	if err != nil {
		return nil, err
	}
	token.Secret = key.Secret()

	if err := s.repo.CreateMFAToken(ctx, token); err != nil {
		return nil, err
	}

	if code := strings.TrimSpace(req.Code); code != "" {
		if err := s.activateMFA(ctx, user, token, code); err != nil {
			return nil, err
		}
	}

	return &EnableMFAResponse{
		Secret:          token.Secret,
		BackupCodes:     token.BackupCodes,
		ProvisioningURI: key.URL(),
	}, nil
}

// VerifyMFACode finalises MFA enrolment for a user.
func (s *Service) VerifyMFACode(ctx context.Context, req VerifyMFACodeRequest) error {
	if req.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if strings.TrimSpace(req.Code) == "" {
		return NewDomainError(ErrCodeMFAInvalid, ErrMFAInvalidCode)
	}

	user, err := s.repo.FindByID(ctx, req.UserID)
	if err != nil {
		return err
	}

	token, err := s.repo.FindActiveMFAToken(ctx, user.ID, MFATypeTOTP)
	if err != nil {
		return err
	}

	return s.activateMFA(ctx, user, token, req.Code)
}

// DisableMFA revokes MFA enrolment for a user.
func (s *Service) DisableMFA(ctx context.Context, req DisableMFARequest) error {
	if req.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	token, err := s.repo.FindActiveMFAToken(ctx, req.UserID, MFATypeTOTP)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteMFAToken(ctx, token.ID); err != nil {
		return err
	}

	user, err := s.repo.FindByID(ctx, req.UserID)
	if err != nil {
		return err
	}
	user.MFAEnabled = false
	user.MFAMethod = ""
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	_ = s.recordAudit(ctx, &user.ID, AuditActionMFADisabled, nil, "", "")

	return nil
}

// LinkOAuthAccount associates an external identity provider.
func (s *Service) LinkOAuthAccount(ctx context.Context, req OAuthLinkRequest) error {
	account, err := NewOAuthAccount(req.UserID, req.Provider, req.ProviderUserID)
	if err != nil {
		return err
	}
	account.UpdateTokens(req.AccessToken, req.RefreshToken, req.ExpiresAt, req.Scopes)

	if err := s.repo.UpsertOAuthAccount(ctx, account); err != nil {
		return err
	}

	_ = s.recordAudit(ctx, &req.UserID, AuditActionOAuthLinked, map[string]any{"provider": req.Provider}, "", "")
	return nil
}

// UnlinkOAuthAccount removes an external provider association.
func (s *Service) UnlinkOAuthAccount(ctx context.Context, userID uuid.UUID, provider OAuthProvider) error {
	if err := s.repo.DeleteOAuthAccount(ctx, userID, provider); err != nil {
		return err
	}
	_ = s.recordAudit(ctx, &userID, AuditActionOAuthUnlinked, map[string]any{"provider": provider}, "", "")
	return nil
}

// ListOAuthAccounts returns the linked providers for a user.
func (s *Service) ListOAuthAccounts(ctx context.Context, userID uuid.UUID) ([]OAuthAccount, error) {
	return s.repo.ListOAuthAccounts(ctx, userID)
}

// BulkUpdateUsers applies administrative updates.
func (s *Service) BulkUpdateUsers(ctx context.Context, req BulkUserUpdateRequest) error {
	if len(req.UserIDs) == 0 {
		return nil
	}

	updates := make(map[string]any)
	if req.SetRole != nil {
		updates["role"] = strings.TrimSpace(*req.SetRole)
	}
	if req.DisableMFA {
		updates["mfa_enabled"] = false
		updates["mfa_method"] = ""
	}
	if req.ConfirmEmail {
		now := time.Now().UTC()
		updates["email_confirmed_at"] = now
	}

	if len(updates) == 0 {
		return nil
	}

	return s.repo.BulkUpdateUserStatus(ctx, req.UserIDs, updates)
}

// ListAuditLogs retrieves audit entries.
func (s *Service) ListAuditLogs(ctx context.Context, userID uuid.UUID, limit int) ([]AuditLog, error) {
	return s.repo.ListAuditLogs(ctx, userID, limit)
}

// Helper methods ---------------------------------------------------------

func (s *Service) issueLoginArtifacts(ctx context.Context, user *User, deviceFingerprint, deviceName, ip, userAgent string, metadata map[string]any) (*LoginResponse, error) {
	if user == nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrNilUser)
	}

	device, err := s.ensureDevice(ctx, user.ID, deviceFingerprint, deviceName)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateSecureToken(refreshTokenEntropyBytes)
	if err != nil {
		return nil, NewDomainError(ErrCodeTokenIssuanceFailure, ErrUnableToIssueToken)
	}

	session, err := NewSession(user.ID, device.ID, hashToken(refreshToken), userAgent, ip, s.sessionTTL)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	user.MarkLogin()
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	accessToken, err := s.issueAccessToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	if err := s.sendSessionAlert(ctx, user, session, device); err != nil && s.logger != nil {
		s.logger.Warn("auth: session alert send failed", slog.Any("error", err), slog.String("user_id", user.ID.String()))
	}

	auditMetadata := map[string]any{
		"session_id": session.ID.String(),
	}
	for key, value := range metadata {
		auditMetadata[key] = value
	}

	_ = s.recordAudit(ctx, &user.ID, AuditActionLoginSucceeded, auditMetadata, ip, userAgent)

	return &LoginResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID:    session.ID,
	}, nil
}

func (s *Service) dispatchConfirmationEmail(ctx context.Context, user *User) error {
	if user == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilUser)
	}

	if err := s.repo.DeleteEmailTokensByUser(ctx, user.ID, EmailTokenTypeConfirmation); err != nil {
		return err
	}

	rawToken, err := generateSecureToken(refreshTokenEntropyBytes)
	if err != nil {
		return NewDomainError(ErrCodeTokenIssuanceFailure, ErrUnableToSendEmail)
	}

	emailToken, err := NewEmailToken(user.ID, EmailTokenTypeConfirmation, rawToken, s.confirmationTTL)
	if err != nil {
		return err
	}

	if err := s.repo.CreateEmailToken(ctx, emailToken); err != nil {
		return err
	}

	confirmationURL := fmt.Sprintf("%s/auth/confirm-email?token=%s", s.publicURL, rawToken)
	data := map[string]any{
		"DisplayName":     pickString(user.Email, "there"),
		"ConfirmationURL": confirmationURL,
		"ConfirmationTTL": fmt.Sprintf("%.0f minutes", s.confirmationTTL.Minutes()),
		"SupportURL":      s.publicURL + "/support",
		"Preheader":       "Confirm your Woragis account",
	}

	return s.sendTemplatedEmail(ctx, emailtemplates.TemplateConfirmation, user.Email, user.PreferredLocale, data)
}

func (s *Service) dispatchPasswordResetEmail(ctx context.Context, user *User, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.publicURL, token)
	data := map[string]any{
		"DisplayName": pickString(user.Email, "there"),
		"ResetURL":    resetURL,
		"ResetTTL":    fmt.Sprintf("%.0f minutes", s.passwordResetTTL.Minutes()),
		"SupportURL":  s.publicURL + "/support",
		"Preheader":   "Reset your Woragis password",
	}

	return s.sendTemplatedEmail(ctx, emailtemplates.TemplatePasswordReset, user.Email, user.PreferredLocale, data)
}

func (s *Service) sendPasswordChangedNotice(ctx context.Context, user *User) error {
	data := map[string]any{
		"DisplayName":       pickString(user.Email, "there"),
		"DeviceName":        "Your account",
		"ApproxLocation":    "Unknown location",
		"IPAddress":         "",
		"SignedInAt":        time.Now().Format(time.RFC1123),
		"ManageSessionsURL": s.publicURL + "/settings/security",
		"SupportURL":        s.publicURL + "/support",
		"Preheader":         "Your Woragis password has changed",
	}

	return s.sendTemplatedEmail(ctx, emailtemplates.TemplateSessionAlert, user.Email, user.PreferredLocale, data)
}

func (s *Service) sendSessionAlert(ctx context.Context, user *User, session *Session, device *Device) error {
	var approxLocation string
	if session.Location != nil {
		if city, ok := session.Location["city"]; ok {
			approxLocation = fmt.Sprintf("%v", city)
		}
	}
	if approxLocation == "" {
		approxLocation = "Unknown location"
	}

	data := map[string]any{
		"DisplayName":       pickString(user.Email, "there"),
		"DeviceName":        device.Name,
		"ApproxLocation":    approxLocation,
		"IPAddress":         session.IP,
		"SignedInAt":        session.LastSeenAt.Format(time.RFC1123),
		"ManageSessionsURL": s.publicURL + "/settings/security",
		"SupportURL":        s.publicURL + "/support",
		"Preheader":         "New Woragis sign-in detected",
	}

	return s.sendTemplatedEmail(ctx, emailtemplates.TemplateSessionAlert, user.Email, user.PreferredLocale, data)
}

func (s *Service) sendWelcomeEmail(ctx context.Context, user *User) error {
	data := map[string]any{
		"DisplayName":       pickString(user.Email, "there"),
		"GettingStartedURL": s.publicURL,
		"SupportURL":        s.publicURL + "/support",
		"Preheader":         "Welcome to Woragis",
	}

	return s.sendTemplatedEmail(ctx, emailtemplates.TemplateWelcome, user.Email, user.PreferredLocale, data)
}

func (s *Service) sendTemplatedEmail(ctx context.Context, templateID emailtemplates.TemplateID, to, locale string, data map[string]any) error {
	if s.emailSender == nil {
		return nil
	}

	data = withDefaultEmailData(data, s.publicURL)

	output, err := s.renderer.Render(emailtemplates.RenderInput{
		Template: templateID,
		Locale:   locale,
		Data:     data,
	})
	if err != nil {
		return err
	}

	message := emailservice.Message{
		To:       to,
		Subject:  output.Subject,
		TextBody: output.Text,
		HTMLBody: output.HTML,
	}

	return s.emailSender.Send(ctx, message)
}

func (s *Service) ensureDevice(ctx context.Context, userID uuid.UUID, fingerprint, name string) (*Device, error) {
	fp := strings.TrimSpace(fingerprint)
	if fp == "" {
		token, err := generateSecureToken(24)
		if err != nil {
			return nil, err
		}
		fp = token
	}

	device, err := NewDevice(userID, name, fp)
	if err != nil {
		return nil, err
	}

	return s.repo.UpsertDevice(ctx, device)
}

func (s *Service) validateMFACode(ctx context.Context, userID uuid.UUID, code string) error {
	code = strings.TrimSpace(code)
	if code == "" {
		return NewDomainError(ErrCodeMFARequired, ErrMFAEnrollmentRequired)
	}

	token, err := s.repo.FindActiveMFAToken(ctx, userID, MFATypeTOTP)
	if err != nil {
		return err
	}

	if !totp.Validate(code, token.Secret) {
		return NewDomainError(ErrCodeMFAInvalid, ErrMFAInvalidCode)
	}

	token.Touch()
	_ = s.repo.UpdateMFAToken(ctx, token)

	return nil
}

func (s *Service) activateMFA(ctx context.Context, user *User, token *MFAToken, code string) error {
	if !totp.Validate(code, token.Secret) {
		return NewDomainError(ErrCodeMFAInvalid, ErrMFAInvalidCode)
	}

	token.Activate()
	if err := s.repo.UpdateMFAToken(ctx, token); err != nil {
		return err
	}

	user.MFAEnabled = true
	user.MFAMethod = string(MFATypeTOTP)
	if err := s.repo.Update(ctx, user); err != nil {
		return err
	}

	_ = s.recordAudit(ctx, &user.ID, AuditActionMFAEnabled, nil, "", "")
	return nil
}

func (s *Service) issueAccessToken(userID uuid.UUID, email string) (string, error) {
	if s.jwtManager == nil {
		return "", NewDomainError(ErrCodeTokenIssuanceFailure, ErrUnableToIssueToken)
	}

	return s.jwtManager.Generate(userID, email)
}

func (s *Service) recordAudit(ctx context.Context, userID *uuid.UUID, action AuditAction, metadata map[string]any, ip, userAgent string) error {
	entry, err := NewAuditLog(userID, action, metadata, ip, userAgent)
	if err != nil {
		return err
	}
	if err := s.repo.InsertAuditLog(ctx, entry); err != nil && s.logger != nil {
		s.logger.Warn("auth: audit log insert failed", slog.Any("error", err))
	}
	return nil
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}

func validatePassword(password string) error {
	if len(password) < minPasswordLength {
		return NewDomainError(ErrCodeInvalidPassword, ErrPasswordTooShort)
	}
	if len(password) > maxPasswordLength {
		return NewDomainError(ErrCodeInvalidPassword, ErrPasswordTooLong)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", NewDomainError(ErrCodePasswordHashFailure, ErrHashPassword)
	}
	return string(hash), nil
}

func generateSecureToken(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(buf), nil
}

func withDefaultEmailData(data map[string]any, baseURL string) map[string]any {
	if data == nil {
		data = make(map[string]any)
	}
	if _, ok := data["SupportURL"]; !ok {
		data["SupportURL"] = baseURL + "/support"
	}
	return data
}

func pickString(values ...string) string {
	for _, value := range values {
		if v := strings.TrimSpace(value); v != "" {
			return v
		}
	}
	return ""
}

// ErrUserNotFoundErr standardises missing user comparison.
