package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"github.com/woragis/backend/server/app/pkg/response"
)

// NewAuthMiddleware produces a Fiber middleware that enforces JWT authentication.
func NewAuthMiddleware(manager *JWTManager, logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if logger != nil {
			logger.Info("JWT Auth middleware called", slog.String("method", c.Method()), slog.String("path", c.Path()))
		}
		if manager == nil {
			return response.Error(c, fiber.StatusInternalServerError, ErrCodeTokenIssuanceFailure, fiber.Map{
				"message": "authentication is not configured",
			})
		}

		rawToken := extractTokenFromHeader(c.Get("Authorization"))
		if rawToken == "" {
			var err error
			rawToken, err = extractTokenFromCookie(c)
			if err != nil && logger != nil {
				logger.Warn("auth: failed to parse auth cookie", slog.Any("error", err))
			}
		}

		// Fallback for websocket/query-string auth (e.g., ws://.../stream?token=...)
		if rawToken == "" {
			rawToken = strings.TrimSpace(c.Query("token", ""))
		}

		if rawToken == "" {
			return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidToken, fiber.Map{
				"message": "missing credentials",
			})
		}

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

func extractTokenFromHeader(header string) string {
	authHeader := strings.TrimSpace(header)
	if authHeader == "" {
		return ""
	}

	const prefix = "bearer "
	if len(authHeader) < len(prefix) || !strings.EqualFold(authHeader[:len(prefix)], prefix) {
		return ""
	}

	return strings.TrimSpace(authHeader[len(prefix):])
}

func extractTokenFromCookie(c *fiber.Ctx) (string, error) {
	cookieValue := strings.TrimSpace(c.Cookies("woragis_auth"))
	if cookieValue == "" {
		return "", nil
	}

	decoded, err := url.QueryUnescape(cookieValue)
	if err != nil {
		return "", err
	}

	var payload struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal([]byte(decoded), &payload); err != nil {
		return "", err
	}

	return strings.TrimSpace(payload.Token), nil
}
