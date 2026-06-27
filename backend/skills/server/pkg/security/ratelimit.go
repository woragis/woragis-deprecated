package security

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimitMiddleware creates a rate limiter middleware
// maxRequests: maximum number of requests allowed
// window: time window for the rate limit
func RateLimitMiddleware(maxRequests int, window time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        maxRequests,
		Expiration: window,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use user ID if available (from JWT), otherwise use IP
			userID := c.Locals("user_id")
			if userID != nil {
				return fmt.Sprintf("user:%v", userID)
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded",
				"message": fmt.Sprintf("Maximum %d requests per %v exceeded", maxRequests, window),
			})
		},
		// Add rate limit headers
		Next: func(c *fiber.Ctx) bool {
			// Skip rate limiting for health checks and metrics
			return c.Path() == "/healthz" || c.Path() == "/metrics"
		},
	})
}

