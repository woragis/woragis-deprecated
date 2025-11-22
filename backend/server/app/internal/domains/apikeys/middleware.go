package apikeys

import (
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Context key for API key in Fiber context
const ContextAPIKeyKey = "apikeys.api_key"

// NewAPIKeyMiddleware produces a Fiber middleware that validates API keys for GET requests.
// This middleware should be used before the auth middleware to allow API keys for GET requests.
func NewAPIKeyMiddleware(service Service, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only allow API keys for GET requests
		if c.Method() != fiber.MethodGet {
			return c.Next()
		}

		// Try to extract API key from header
		apiKey := extractAPIKeyFromHeader(c.Get("X-API-Key"))
		if apiKey == "" {
			// Also check Authorization header with "ApiKey" prefix
			apiKey = extractAPIKeyFromAuthHeader(c.Get("Authorization"))
		}

		if apiKey == "" {
			// No API key provided, continue to next middleware (JWT auth)
			return c.Next()
		}

		// Validate the API key
		validatedKey, err := service.ValidateAPIKey(c.Context(), apiKey)
		if err != nil {
			if logger != nil {
				logger.Warn("api key validation failed", slog.Any("error", err))
			}
			// Don't fail here, let JWT auth try
			return c.Next()
		}

		// Store API key in context
		c.Locals(ContextAPIKeyKey, validatedKey)

		// API key is valid, allow the request
		return c.Next()
	}
}

// extractAPIKeyFromHeader extracts API key from X-API-Key header.
func extractAPIKeyFromHeader(header string) string {
	return strings.TrimSpace(header)
}

// extractAPIKeyFromAuthHeader extracts API key from Authorization header with "ApiKey" prefix.
func extractAPIKeyFromAuthHeader(header string) string {
	authHeader := strings.TrimSpace(header)
	if authHeader == "" {
		return ""
	}

	const prefix = "apikey "
	if len(authHeader) < len(prefix) || !strings.EqualFold(authHeader[:len(prefix)], prefix) {
		return ""
	}

	return strings.TrimSpace(authHeader[len(prefix):])
}

// APIKeyFromContext extracts the API key from the Fiber context.
func APIKeyFromContext(c *fiber.Ctx) (*APIKey, bool) {
	apiKey, ok := c.Locals(ContextAPIKeyKey).(*APIKey)
	return apiKey, ok
}

// RequireAPIKeyOrAuth is a middleware that requires either API key (for GET) or JWT auth (for all methods).
func RequireAPIKeyOrAuth(apiKeyService Service, jwtMiddleware fiber.Handler, logger *slog.Logger) fiber.Handler {
	apiKeyMiddleware := NewAPIKeyMiddleware(apiKeyService, logger)

	return func(c *fiber.Ctx) error {
		// For GET requests, try API key first
		if c.Method() == fiber.MethodGet {
			// Run API key middleware
			if err := apiKeyMiddleware(c); err != nil {
				return err
			}

			// Check if API key was set
			if _, hasAPIKey := APIKeyFromContext(c); hasAPIKey {
				// API key is valid, allow request
				return c.Next()
			}
		}

		// For non-GET or if no API key, require JWT auth
		return jwtMiddleware(c)
	}
}

