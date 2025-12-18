package metrics

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Middleware creates a Fiber middleware that records HTTP request metrics
func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Increment in-flight requests
		IncHTTPRequestsInFlight()
		defer DecHTTPRequestsInFlight()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get method and path
		method := c.Method()
		endpoint := normalizeEndpoint(c.Path())

		// Get status code
		status := strconv.Itoa(c.Response().StatusCode())

		// Record metrics
		RecordHTTPRequest(method, endpoint, status, duration)

		return err
	}
}

// normalizeEndpoint normalizes the endpoint path for metrics
// Replaces UUIDs and IDs with placeholders to reduce cardinality
func normalizeEndpoint(path string) string {
	// This is a simple implementation - you might want to make it more sophisticated
	// For now, we'll just use the path as-is, but in production you'd want to:
	// - Replace UUIDs with :id
	// - Replace numeric IDs with :id
	// - Normalize query parameters
	return path
}
