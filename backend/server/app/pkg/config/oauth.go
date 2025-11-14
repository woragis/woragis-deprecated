package config

import (
	"os"
	"strings"
)

// OAuthProviderConfig describes the runtime configuration for a single OAuth provider.
type OAuthProviderConfig struct {
	Name         string
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// OAuthConfig aggregates the configured OAuth providers keyed by provider identifier.
type OAuthConfig struct {
	Providers map[string]OAuthProviderConfig
}

// LoadOAuthConfig reads OAuth provider configuration from environment variables.
//
// The following environment variables are recognised:
//
//   - GOOGLE_OAUTH_CLIENT_ID / GOOGLE_OAUTH_CLIENT_SECRET
//
//   - GOOGLE_OAUTH_REDIRECT_URL (defaults to {publicURL}/api/auth/oauth/callback/google)
//
//   - GOOGLE_OAUTH_SCOPES (comma or space separated; defaults to "openid email profile")
//
//   - GITHUB_OAUTH_CLIENT_ID / GITHUB_OAUTH_CLIENT_SECRET
//
//   - GITHUB_OAUTH_REDIRECT_URL (defaults to {publicURL}/api/auth/oauth/callback/github)
//
//   - GITHUB_OAUTH_SCOPES (comma or space separated; defaults to "read:user user:email")
//
//   - MICROSOFT_OAUTH_CLIENT_ID / MICROSOFT_OAUTH_CLIENT_SECRET
//
//   - MICROSOFT_OAUTH_REDIRECT_URL (defaults to {publicURL}/api/auth/oauth/callback/microsoft)
//
//   - MICROSOFT_OAUTH_SCOPES (comma or space separated; defaults to "User.Read email openid profile")
//
// Providers are considered enabled when both client ID and client secret are present.
func LoadOAuthConfig(publicURL string) OAuthConfig {
	publicURL = strings.TrimRight(publicURL, "/")

	cfg := OAuthConfig{
		Providers: make(map[string]OAuthProviderConfig),
	}

	addIfEnabled := func(id, name, defaultScopes, redirectPath string, overrides map[string]string) {
		clientID := strings.TrimSpace(os.Getenv(overrides["client_id"]))
		clientSecret := strings.TrimSpace(os.Getenv(overrides["client_secret"]))
		if clientID == "" || clientSecret == "" {
			return
		}

		redirectURL := strings.TrimSpace(os.Getenv(overrides["redirect"]))
		if redirectURL == "" {
			redirectURL = publicURL + redirectPath
		}

		scopeValue := strings.TrimSpace(os.Getenv(overrides["scopes"]))
		if scopeValue == "" {
			scopeValue = defaultScopes
		}

		cfg.Providers[id] = OAuthProviderConfig{
			Name:         name,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       splitScopes(scopeValue),
		}
	}

	addIfEnabled("google", "Google", "openid email profile", "/api/auth/oauth/callback/google", map[string]string{
		"client_id":     "GOOGLE_OAUTH_CLIENT_ID",
		"client_secret": "GOOGLE_OAUTH_CLIENT_SECRET",
		"redirect":      "GOOGLE_OAUTH_REDIRECT_URL",
		"scopes":        "GOOGLE_OAUTH_SCOPES",
	})

	addIfEnabled("github", "GitHub", "read:user user:email", "/api/auth/oauth/callback/github", map[string]string{
		"client_id":     "GITHUB_OAUTH_CLIENT_ID",
		"client_secret": "GITHUB_OAUTH_CLIENT_SECRET",
		"redirect":      "GITHUB_OAUTH_REDIRECT_URL",
		"scopes":        "GITHUB_OAUTH_SCOPES",
	})

	addIfEnabled("microsoft", "Microsoft", "User.Read email openid profile", "/api/auth/oauth/callback/microsoft", map[string]string{
		"client_id":     "MICROSOFT_OAUTH_CLIENT_ID",
		"client_secret": "MICROSOFT_OAUTH_CLIENT_SECRET",
		"redirect":      "MICROSOFT_OAUTH_REDIRECT_URL",
		"scopes":        "MICROSOFT_OAUTH_SCOPES",
	})

	return cfg
}

func splitScopes(raw string) []string {
	if raw == "" {
		return nil
	}

	raw = strings.ReplaceAll(raw, ",", " ")
	fields := strings.Fields(raw)
	seen := make(map[string]struct{}, len(fields))
	result := make([]string, 0, len(fields))
	for _, scope := range fields {
		scope = strings.TrimSpace(scope)
		if scope == "" {
			continue
		}
		if _, exists := seen[scope]; exists {
			continue
		}
		seen[scope] = struct{}{}
		result = append(result, scope)
	}
	return result
}
