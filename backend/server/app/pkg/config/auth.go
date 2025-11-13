package config

import (
	"fmt"
	"os"
	"time"
)

// AuthConfig holds settings for issuing JWT access tokens.
type AuthConfig struct {
	JWTSecret string
	JWTTTL    time.Duration
}

const (
	defaultJWTSecret = "dev-secret-change-me"
	defaultJWTTTL    = 24 * time.Hour
)

// LoadAuthConfig reads auth-related configuration from the environment.
func LoadAuthConfig() (AuthConfig, error) {
	secret := os.Getenv("AUTH_JWT_SECRET")
	if secret == "" {
		secret = defaultJWTSecret
	}

	ttlStr := os.Getenv("AUTH_JWT_TTL")
	if ttlStr == "" {
		return AuthConfig{
			JWTSecret: secret,
			JWTTTL:    defaultJWTTTL,
		}, nil
	}

	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		return AuthConfig{}, fmt.Errorf("invalid AUTH_JWT_TTL value %q: %w", ttlStr, err)
	}

	if ttl <= 0 {
		return AuthConfig{}, fmt.Errorf("AUTH_JWT_TTL must be positive, got %s", ttl)
	}

	return AuthConfig{
		JWTSecret: secret,
		JWTTTL:    ttl,
	}, nil
}
