package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"woragis-files-service/pkg/authservice"
	"woragis-files-service/pkg/utils"
)

// AuthValidationConfig holds configuration for auth validation middleware
type AuthValidationConfig struct {
	AuthServiceClient *authservice.Client
	SkipPaths         []string
}

// DefaultAuthValidationConfig returns default config
func DefaultAuthValidationConfig(authServiceClient *authservice.Client) AuthValidationConfig {
	return AuthValidationConfig{
		AuthServiceClient: authServiceClient,
		SkipPaths: []string{
			"/healthz",
			"/metrics",
		},
	}
}

// AuthValidationMiddleware validates JWT tokens via Auth Service
func AuthValidationMiddleware(config AuthValidationConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Skip validation for certain paths
		for _, path := range config.SkipPaths {
			if strings.HasPrefix(c.Path(), path) {
				return c.Next()
			}
		}

		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse{
				Success: false,
				Message: "missing authorization header",
			})
		}

		// Extract token (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse{
				Success: false,
				Message: "invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token with Auth Service
		response, err := config.AuthServiceClient.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.ErrorResponse{
				Success: false,
				Message: "failed to validate token",
			})
		}

		if !response.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.ErrorResponse{
				Success: false,
				Message: response.Message,
			})
		}

		// Store user info in context
		c.Locals("userID", response.UserID)
		c.Locals("userEmail", response.Email)
		c.Locals("userRole", response.Role)

		return c.Next()
	}
}

// UserIDFromContext extracts user ID from context
func UserIDFromContext(c *fiber.Ctx) (string, error) {
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "user ID not found in context")
	}
	return userID, nil
}

