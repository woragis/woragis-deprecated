package auth

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/woragis/backend/server/app/pkg/response"
)

// NewAuthMiddleware produces a Fiber middleware that enforces JWT authentication.
func NewAuthMiddleware(manager *JWTManager, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if manager == nil {
			return response.Error(c, fiber.StatusInternalServerError, ErrCodeTokenIssuanceFailure, fiber.Map{
				"message": "authentication is not configured",
			})
		}

		authHeader := strings.TrimSpace(c.Get("Authorization"))
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidToken, fiber.Map{
				"message": "missing Authorization header",
			})
		}

		const prefix = "bearer "
		if len(authHeader) < len(prefix) || !strings.EqualFold(authHeader[:len(prefix)], prefix) {
			return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidToken, fiber.Map{
				"message": "invalid Authorization header format",
			})
		}

		rawToken := strings.TrimSpace(authHeader[len(prefix):])
		claims, err := manager.Verify(rawToken)
		if err != nil {
			if logger != nil {
				logger.Warn("auth: token verification failed", slog.Any("error", err))
			}

			var message string
			if errors.Is(err, jwt.ErrTokenExpired) {
				message = "token expired"
			} else {
				message = "invalid token"
			}

			return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidToken, fiber.Map{
				"message": message,
			})
		}

		c.Locals(ContextUserIDKey, claims.Subject)
		c.Locals(ContextUserEmailKey, claims.Email)

		return c.Next()
	}
}
