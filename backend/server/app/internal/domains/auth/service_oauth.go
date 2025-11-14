package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

// OAuthFlowMode indicates the intent of the OAuth flow.
type OAuthFlowMode string

const (
	// OAuthModeLogin is used for signing users in or creating new accounts.
	OAuthModeLogin OAuthFlowMode = "login"
	// OAuthModeLink links a provider to an already authenticated user.
	OAuthModeLink OAuthFlowMode = "link"
)

// OAuthProviderSettings configures the runtime behaviour for an OAuth provider.
type OAuthProviderSettings struct {
	Name        string
	Config      *oauth2.Config
	UserInfoURL string
}

type oauthProviderConfig struct {
	name        string
	config      *oauth2.Config
	userInfoURL string
}

type oauthStateData struct {
	Provider          OAuthProvider
	Mode              OAuthFlowMode
	RedirectOrigin    string
	UserID            *uuid.UUID
	DeviceFingerprint string
	DeviceName        string
	CreatedAt         time.Time
}

// OAuthStartOptions configures the generation of an OAuth authorisation URL.
type OAuthStartOptions struct {
	Mode              OAuthFlowMode
	UserID            *uuid.UUID
	RedirectOrigin    string
	DeviceFingerprint string
	DeviceName        string
}

// OAuthCallbackInput carries callback data from an external provider.
type OAuthCallbackInput struct {
	Provider  OAuthProvider
	State     string
	Code      string
	IPAddress string
	UserAgent string
}

// OAuthCallbackResult summarises the outcome of an OAuth callback.
type OAuthCallbackResult struct {
	Provider       OAuthProvider
	Mode           OAuthFlowMode
	RedirectOrigin string
	Success        bool
	Login          *LoginResponse
	Message        string
}

// OAuthProviderInfo describes a configured provider exposed to clients.
type OAuthProviderInfo struct {
	ID   OAuthProvider `json:"id"`
	Name string        `json:"name"`
}

// ConfigureOAuthProviders updates the service with the supplied provider settings.
func (s *Service) ConfigureOAuthProviders(configs map[OAuthProvider]OAuthProviderSettings) {
	s.oauthMu.Lock()
	defer s.oauthMu.Unlock()

	s.oauthProviders = make(map[OAuthProvider]*oauthProviderConfig, len(configs))

	for provider, cfg := range configs {
		if cfg.Config == nil || cfg.UserInfoURL == "" {
			continue
		}

		clone := *cfg.Config
		if len(clone.Scopes) == 0 {
			clone.Scopes = []string{"openid", "email"}
		}

		s.oauthProviders[provider] = &oauthProviderConfig{
			name:        cfg.Name,
			config:      &clone,
			userInfoURL: cfg.UserInfoURL,
		}
	}
}

// ListOAuthProviders returns configured OAuth providers.
func (s *Service) ListOAuthProviders() []OAuthProviderInfo {
	s.oauthMu.Lock()
	defer s.oauthMu.Unlock()

	providers := make([]OAuthProviderInfo, 0, len(s.oauthProviders))
	for id, cfg := range s.oauthProviders {
		providers = append(providers, OAuthProviderInfo{
			ID:   id,
			Name: cfg.name,
		})
	}
	return providers
}

// BeginOAuth generates an authorisation URL and state token.
func (s *Service) BeginOAuth(ctx context.Context, provider OAuthProvider, opts OAuthStartOptions) (string, string, error) {
	if opts.Mode == "" {
		opts.Mode = OAuthModeLogin
	}
	if strings.TrimSpace(opts.RedirectOrigin) == "" {
		return "", "", NewDomainError(ErrCodeInvalidPayload, "auth: redirect origin is required")
	}

	s.oauthMu.Lock()
	cfg, ok := s.oauthProviders[provider]
	if !ok {
		s.oauthMu.Unlock()
		return "", "", NewDomainError(ErrCodeInvalidPayload, ErrEmptyOAuthProvider)
	}
	s.pruneExpiredStatesLocked(time.Now().UTC())

	state, err := generateSecureToken(32)
	if err != nil {
		s.oauthMu.Unlock()
		return "", "", err
	}

	entry := oauthStateData{
		Provider:          provider,
		Mode:              opts.Mode,
		RedirectOrigin:    opts.RedirectOrigin,
		UserID:            opts.UserID,
		DeviceFingerprint: strings.TrimSpace(opts.DeviceFingerprint),
		DeviceName:        strings.TrimSpace(opts.DeviceName),
		CreatedAt:         time.Now().UTC(),
	}

	if entry.Mode == OAuthModeLink && entry.UserID == nil {
		s.oauthMu.Unlock()
		return "", "", NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	s.oauthStates[state] = entry
	s.oauthMu.Unlock()

	authURL := cfg.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return state, authURL, nil
}

// CompleteOAuth finalises the OAuth flow after the provider callback.
func (s *Service) CompleteOAuth(ctx context.Context, input OAuthCallbackInput) (*OAuthCallbackResult, error) {
	result := &OAuthCallbackResult{
		Provider: input.Provider,
	}

	entry, err := s.consumeOAuthState(input.State)
	if err != nil {
		return result, err
	}

	result.Mode = entry.Mode
	result.RedirectOrigin = entry.RedirectOrigin

	cfg := s.lookupProviderConfig(entry.Provider)
	if cfg == nil {
		return result, NewDomainError(ErrCodeInvalidPayload, ErrEmptyOAuthProvider)
	}

	token, err := cfg.config.Exchange(ctx, input.Code)
	if err != nil {
		return result, fmt.Errorf("oauth exchange failed: %w", err)
	}

	userInfo, err := s.fetchOAuthUserInfo(ctx, entry.Provider, cfg, token)
	if err != nil {
		return result, err
	}

	scopeList := scopesFromToken(token, cfg)
	var expiresAt *time.Time
	if !token.Expiry.IsZero() {
		exp := token.Expiry.UTC()
		expiresAt = &exp
	}

	switch entry.Mode {
	case OAuthModeLink:
		if entry.UserID == nil {
			return result, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
		}
		if err := s.linkOAuthIdentity(ctx, *entry.UserID, entry.Provider, userInfo, token, expiresAt, scopeList); err != nil {
			return result, err
		}
		result.Success = true
		result.Message = "provider linked"
	case OAuthModeLogin:
		loginResp, err := s.loginWithOAuth(ctx, entry.Provider, userInfo, token, expiresAt, scopeList, entry.DeviceFingerprint, entry.DeviceName, input.IPAddress, input.UserAgent)
		if err != nil {
			return result, err
		}
		result.Login = loginResp
		result.Success = true
	default:
		return result, NewDomainError(ErrCodeInvalidPayload, "auth: unsupported oauth flow mode")
	}

	return result, nil
}

func (s *Service) consumeOAuthState(state string) (oauthStateData, error) {
	if strings.TrimSpace(state) == "" {
		return oauthStateData{}, NewDomainError(ErrCodeInvalidToken, ErrInvalidOAuthState)
	}

	s.oauthMu.Lock()
	defer s.oauthMu.Unlock()

	entry, ok := s.oauthStates[state]
	if !ok {
		return oauthStateData{}, NewDomainError(ErrCodeInvalidToken, ErrInvalidOAuthState)
	}
	delete(s.oauthStates, state)
	return entry, nil
}

func (s *Service) pruneExpiredStatesLocked(now time.Time) {
	for key, entry := range s.oauthStates {
		if now.Sub(entry.CreatedAt) > s.oauthStateTTL {
			delete(s.oauthStates, key)
		}
	}
}

func (s *Service) lookupProviderConfig(provider OAuthProvider) *oauthProviderConfig {
	s.oauthMu.Lock()
	defer s.oauthMu.Unlock()
	return s.oauthProviders[provider]
}

type oauthUserInfo struct {
	ID    string
	Email string
}

func (s *Service) fetchOAuthUserInfo(ctx context.Context, provider OAuthProvider, cfg *oauthProviderConfig, token *oauth2.Token) (*oauthUserInfo, error) {
	if token == nil || token.AccessToken == "" {
		return nil, NewDomainError(ErrCodeInvalidToken, ErrInvalidOAuthState)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("oauth userinfo failed: status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	info, err := extractOAuthUser(provider, payload)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

func extractOAuthUser(provider OAuthProvider, payload map[string]any) (oauthUserInfo, error) {
	stringVal := func(key string) string {
		if val, ok := payload[key]; ok {
			switch typed := val.(type) {
			case string:
				return strings.TrimSpace(typed)
			case float64:
				return fmt.Sprintf("%.0f", typed)
			}
		}
		return ""
	}

	var id string
	var email string

	switch provider {
	case OAuthProviderGoogle:
		id = stringVal("id")
		if id == "" {
			id = stringVal("sub")
		}
		email = stringVal("email")
	case OAuthProviderGithub:
		id = stringVal("id")
		email = stringVal("email")
		if email == "" {
			email = stringVal("login")
		}
	case OAuthProviderMicrosoft:
		id = stringVal("id")
		if id == "" {
			id = stringVal("sub")
		}
		email = stringVal("mail")
		if email == "" {
			email = stringVal("userPrincipalName")
		}
	default:
		return oauthUserInfo{}, NewDomainError(ErrCodeInvalidPayload, "auth: unsupported provider")
	}

	if id == "" {
		return oauthUserInfo{}, NewDomainError(ErrCodeInvalidPayload, ErrEmptyOAuthSubject)
	}
	if email == "" {
		return oauthUserInfo{}, NewDomainError(ErrCodeInvalidEmail, ErrEmptyEmail)
	}

	return oauthUserInfo{
		ID:    id,
		Email: strings.ToLower(email),
	}, nil
}

func scopesFromToken(token *oauth2.Token, cfg *oauthProviderConfig) []string {
	if token != nil {
		if raw := token.Extra("scope"); raw != nil {
			switch value := raw.(type) {
			case string:
				if strings.TrimSpace(value) != "" {
					return strings.Fields(value)
				}
			case []string:
				if len(value) > 0 {
					return value
				}
			}
		}
	}
	return cfg.config.Scopes
}

func (s *Service) loginWithOAuth(ctx context.Context, provider OAuthProvider, userInfo *oauthUserInfo, token *oauth2.Token, expiresAt *time.Time, scopes []string, deviceFingerprint, deviceName, ip, userAgent string) (*LoginResponse, error) {
	if userInfo == nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyOAuthSubject)
	}

	account, err := s.repo.FindOAuthAccount(ctx, provider, userInfo.ID)
	if err != nil {
		return nil, err
	}

	var user *User
	isNewAccount := false

	if account != nil {
		user, err = s.repo.FindByID(ctx, account.UserID)
		if err != nil {
			return nil, err
		}
	} else {
		user, err = s.repo.FindByEmail(ctx, userInfo.Email)
		if err != nil {
			if domainErr, ok := AsDomainError(err); ok && domainErr.Code == ErrCodeUserNotFound {
				user, err = s.createUserFromOAuth(ctx, userInfo.Email)
				if err != nil {
					return nil, err
				}
				isNewAccount = true
			} else {
				return nil, err
			}
		}
		if user.EmailConfirmedAt == nil {
			user.ConfirmEmail()
			if err := s.repo.Update(ctx, user); err != nil {
				return nil, err
			}
		}

		account, err = NewOAuthAccount(user.ID, provider, userInfo.ID)
		if err != nil {
			return nil, err
		}
	}

	accessToken := ""
	refreshToken := ""
	if token != nil {
		accessToken = token.AccessToken
		refreshToken = token.RefreshToken
	}

	account.UpdateTokens(accessToken, refreshToken, expiresAt, scopes)
	if err := s.repo.UpsertOAuthAccount(ctx, account); err != nil {
		return nil, err
	}
	_ = s.recordAudit(ctx, &user.ID, AuditActionOAuthLinked, map[string]any{"provider": provider}, ip, userAgent)

	if isNewAccount && s.monitor != nil {
		s.monitor.RecordUserRegistration(ctx, user.ID)
	}

	if deviceName == "" {
		deviceName = fmt.Sprintf("%s OAuth", strings.ToUpper(string(provider)))
	}

	metadata := map[string]any{
		"provider": provider,
	}

	return s.issueLoginArtifacts(ctx, user, deviceFingerprint, deviceName, ip, userAgent, metadata)
}

func (s *Service) linkOAuthIdentity(ctx context.Context, userID uuid.UUID, provider OAuthProvider, userInfo *oauthUserInfo, token *oauth2.Token, expiresAt *time.Time, scopes []string) error {
	account, err := NewOAuthAccount(userID, provider, userInfo.ID)
	if err != nil {
		if domainErr, ok := AsDomainError(err); ok && domainErr.Code == ErrCodeInvalidPayload {
			return domainErr
		}
		return err
	}

	refresh := ""
	if token != nil {
		refresh = token.RefreshToken
	}

	account.UpdateTokens(token.AccessToken, refresh, expiresAt, scopes)

	if err := s.repo.UpsertOAuthAccount(ctx, account); err != nil {
		return err
	}

	_ = s.recordAudit(ctx, &userID, AuditActionOAuthLinked, map[string]any{"provider": provider}, "", "")
	return nil
}

func (s *Service) createUserFromOAuth(ctx context.Context, email string) (*User, error) {
	passwordSeed, err := generateSecureToken(32)
	if err != nil {
		return nil, err
	}

	hash, err := hashPassword(passwordSeed)
	if err != nil {
		return nil, err
	}

	user, err := NewUser(email, hash)
	if err != nil {
		return nil, err
	}

	user.ConfirmEmail()

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	_ = s.recordAudit(ctx, &user.ID, AuditActionUserRegistered, map[string]any{
		"source": "oauth",
	}, "", "")

	return user, nil
}
