package config

import (
	"strconv"
	"strings"
)

// CORSConfig captures cross-origin settings for the HTTP layer
type CORSConfig struct {
	Enabled          bool
	AllowedOrigins   string
	AllowedMethods   string
	AllowedHeaders   string
	ExposedHeaders   string
	AllowCredentials bool
	MaxAge           int
}

// LoadCORSConfig reads CORS-related environment variables
func LoadCORSConfig() *CORSConfig {
	enabled := strings.ToLower(getEnv("CORS_ENABLED", "true"))
	allowCredentials := strings.ToLower(getEnv("CORS_ALLOW_CREDENTIALS", "true"))

	maxAge, err := strconv.Atoi(getEnv("CORS_MAX_AGE", "86400"))
	if err != nil {
		maxAge = 86400
	}

	defaultOrigins := "http://localhost:3000,http://127.0.0.1:3000,http://localhost:5173,http://127.0.0.1:5173"

	return &CORSConfig{
		Enabled:          enabled != "false" && enabled != "0",
		AllowedOrigins:   sanitizeCSV(getEnv("CORS_ALLOWED_ORIGINS", defaultOrigins)),
		AllowedMethods:   sanitizeCSV(getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")),
		AllowedHeaders:   sanitizeCSV(getEnv("CORS_ALLOWED_HEADERS", "Authorization,Content-Type,X-Requested-With")),
		ExposedHeaders:   sanitizeCSV(getEnv("CORS_EXPOSED_HEADERS", "")),
		AllowCredentials: allowCredentials == "true" || allowCredentials == "1" || allowCredentials == "yes",
		MaxAge:           maxAge,
	}
}

func sanitizeCSV(value string) string {
	parts := strings.Split(value, ",")
	clean := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			clean = append(clean, trimmed)
		}
	}

	if len(clean) == 0 {
		return ""
	}

	return strings.Join(clean, ",")
}
