package middleware

import (
	"errors"

	"woragis-skills-service/pkg/auth"
	"woragis-skills-service/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// JWTConfig holds the configuration for JWT middleware
type JWTConfig struct {
	JWTManager *auth.JWTManager
}

// JWTMiddleware creates a Fiber JWT authentication middleware
func JWTMiddleware(config JWTConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.UnauthorizedResponse(c, "Authorization header required")
		}

		// Extract token from "Bearer <token>"
		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			return utils.UnauthorizedResponse(c, "Invalid authorization header format")
		}

		// Validate token
		claims, err := config.JWTManager.Validate(token)
		if err != nil {
			switch err {
			case auth.ErrTokenExpired:
				return utils.UnauthorizedResponse(c, "Token has expired")
			case auth.ErrTokenInvalid:
				return utils.UnauthorizedResponse(c, "Invalid token")
			default:
				return utils.UnauthorizedResponse(c, "Token validation failed")
			}
		}

		// Add user information to context
		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)
		c.Locals("userName", claims.Name)

		return c.Next()
	}
}

// RequireRole creates a middleware that checks if user has required role
func RequireRole(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		if userRole == nil {
			return utils.UnauthorizedResponse(c, "User role not found")
		}

		role, ok := userRole.(string)
		if !ok || role != requiredRole {
			return utils.ForbiddenResponse(c, "Insufficient permissions")
		}

		return c.Next()
	}
}

// RequireAdmin creates a middleware that requires admin role
func RequireAdmin() fiber.Handler {
	return RequireRole("admin")
}

// RequireModerator creates a middleware that requires moderator or admin role
func RequireModerator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole")
		if userRole == nil {
			return utils.UnauthorizedResponse(c, "User role not found")
		}

		role, ok := userRole.(string)
		if !ok || (role != "moderator" && role != "admin") {
			return utils.ForbiddenResponse(c, "Moderator or admin access required")
		}

		return c.Next()
	}
}

// OptionalJWTMiddleware creates a middleware that optionally validates JWT
// If token is present and valid, user info is added to context
// If token is missing or invalid, request continues without user info
func OptionalJWTMiddleware(config JWTConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			token, err := auth.ExtractTokenFromHeader(authHeader)
			if err == nil {
				claims, err := config.JWTManager.Validate(token)
				if err == nil {
					// Add user information to context
					c.Locals("userID", claims.UserID)
					c.Locals("userEmail", claims.Email)
					c.Locals("userRole", claims.Role)
					c.Locals("userName", claims.Name)
				}
			}
		}
		return c.Next()
	}
}

// GetUserIDFromFiberContext extracts user ID from Fiber context
func GetUserIDFromFiberContext(c *fiber.Ctx) (uuid.UUID, error) {
	userID := c.Locals("userID")
	if userID == nil {
		return uuid.Nil, errors.New("user not authenticated")
	}

	id, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("invalid user ID in context")
	}

	return id, nil
}

// GetUserRoleFromFiberContext extracts user role from Fiber context
func GetUserRoleFromFiberContext(c *fiber.Ctx) (string, error) {
	userRole := c.Locals("userRole")
	if userRole == nil {
		return "", errors.New("user role not found")
	}

	role, ok := userRole.(string)
	if !ok {
		return "", errors.New("invalid user role in context")
	}

	return role, nil
}

// GetUserEmailFromFiberContext extracts user email from Fiber context
func GetUserEmailFromFiberContext(c *fiber.Ctx) (string, error) {
	userEmail := c.Locals("userEmail")
	if userEmail == nil {
		return "", errors.New("user email not found")
	}

	email, ok := userEmail.(string)
	if !ok {
		return "", errors.New("invalid user email in context")
	}

	return email, nil
}

// GetUserNameFromFiberContext extracts user name from Fiber context
func GetUserNameFromFiberContext(c *fiber.Ctx) (string, error) {
	userName := c.Locals("userName")
	if userName == nil {
		return "", errors.New("user name not found")
	}

	name, ok := userName.(string)
	if !ok {
		return "", errors.New("invalid user name in context")
	}

	return name, nil
}
