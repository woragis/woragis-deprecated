package config

import (
	"fmt"
	"os"
)

// AuthConfig holds settings for authentication
type AuthConfig struct {
	JWTSecret             string
	JWTExpireHours        int
	JWTRefreshExpireHours int
	AESKey                string
	HashSalt              string
	BCryptCost            int
}

const (
	defaultJWTSecret             = "dev-secret-change-me"
	defaultJWTExpireHours        = 24
	defaultJWTRefreshExpireHours = 168 // 7 days
	defaultBCryptCost            = 12
)

// LoadAuthConfig reads auth-related configuration from the environment
func LoadAuthConfig() (*AuthConfig, error) {
	// Try AUTH_JWT_SECRET first (matches docker-compose), fallback to JWT_SECRET for backward compatibility
	jwtSecret := getEnv("AUTH_JWT_SECRET", "")
	if jwtSecret == "" {
		jwtSecret = getEnv("JWT_SECRET", defaultJWTSecret)
	}
	if jwtSecret == defaultJWTSecret && os.Getenv("ENV") == "production" {
		return nil, fmt.Errorf("AUTH_JWT_SECRET or JWT_SECRET must be set in production")
	}

	jwtExpireHours := getEnvAsInt("JWT_EXPIRE_HOURS", defaultJWTExpireHours)
	jwtRefreshExpireHours := getEnvAsInt("JWT_REFRESH_EXPIRE_HOURS", defaultJWTRefreshExpireHours)
	bcryptCost := getEnvAsInt("BCRYPT_COST", defaultBCryptCost)

	aesKey := getEnvRequired("AES_KEY")
	hashSalt := getEnvRequired("HASH_SALT")

	return &AuthConfig{
		JWTSecret:             jwtSecret,
		JWTExpireHours:        jwtExpireHours,
		JWTRefreshExpireHours: jwtRefreshExpireHours,
		AESKey:                aesKey,
		HashSalt:              hashSalt,
		BCryptCost:            bcryptCost,
	}, nil
}

