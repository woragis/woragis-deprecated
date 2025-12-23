package validation

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

// RequestSizeLimitMiddleware limits request body size
func RequestSizeLimitMiddleware(maxSize int64) fiber.Handler {
	return func(c *fiber.Ctx) error {
		contentLength := int64(c.Request().Header.ContentLength())
		if contentLength > maxSize {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error": "Request body too large",
			})
		}
		return c.Next()
	}
}

// RateLimitMiddleware creates a rate limiter
func RateLimitMiddleware(maxRequests int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use user ID if available, otherwise use IP
			userID := c.Locals("user_id")
			if userID != nil {
				return "user:" + fmt.Sprintf("%v", userID)
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
			})
		},
	})
}

// InputSanitizationMiddleware sanitizes input
func InputSanitizationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Sanitize query parameters
		// Note: c.Queries() returns map[string]string, so we sanitize values
		queries := c.Queries()
		for key, value := range queries {
			sanitized := SanitizeString(value)
			if sanitized != value {
				// Update query parameter if sanitization changed it
				c.Request().URI().QueryArgs().Set(key, sanitized)
			}
		}
		
		return c.Next()
	}
}
