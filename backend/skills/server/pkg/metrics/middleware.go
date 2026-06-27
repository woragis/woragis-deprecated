package metrics

import (
	"strconv"
	"strings"
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

		// Normalize endpoint
		endpoint := normalizeEndpoint(c.Path())

		// Get status code as string
		status := strconv.Itoa(c.Response().StatusCode())

		// Record metrics
		RecordHTTPRequest(c.Method(), endpoint, status, duration)

		return err
	}
}

// normalizeEndpoint normalizes the endpoint path for metrics
// Replaces path parameters with placeholders (e.g., /users/123 -> /users/:id)
func normalizeEndpoint(path string) string {
	// For now, just return the path as-is
	// In a more sophisticated implementation, you might want to:
	// - Replace UUIDs with :id
	// - Replace numeric IDs with :id
	// - Normalize query parameters
	// This helps reduce cardinality in metrics
	
	// Simple normalization: remove trailing slashes
	path = strings.TrimSuffix(path, "/")
	if path == "" {
		path = "/"
	}
	
	return path
}

