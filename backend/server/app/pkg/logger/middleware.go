package logger

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestIDMiddleware generates and adds a trace_id (request ID) to each request.
// The trace_id is added to the context and response headers for distributed tracing.
func RequestIDMiddleware(log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check if trace_id already exists in header (for distributed tracing)
		traceID := c.Get("X-Trace-ID")
		if traceID == "" {
			// Generate new trace_id if not present
			traceID = uuid.New().String()
		}

		// Add trace_id to response header
		c.Set("X-Trace-ID", traceID)

		// Add trace_id to context for logging
		ctx := WithTraceID(c.UserContext(), traceID)
		c.SetUserContext(ctx)

		// Store trace_id in locals for easy access
		c.Locals("trace_id", traceID)

		return c.Next()
	}
}

// RequestLoggerMiddleware logs HTTP requests with structured fields.
// This should be used after RequestIDMiddleware to ensure trace_id is available.
func RequestLoggerMiddleware(log *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get trace_id from context
		traceID := GetTraceID(c.UserContext())
		if traceID == "" {
			// Fallback to locals if not in context
			if id, ok := c.Locals("trace_id").(string); ok {
				traceID = id
			}
		}

		// Get user ID if available (from JWT middleware)
		userID := ""
		// Check auth middleware's user ID key
		if user, ok := c.Locals("auth.user_id").(string); ok {
			userID = user
		} else if userUUID, ok := c.Locals("auth.user_id").(uuid.UUID); ok {
			userID = userUUID.String()
		}

		// Build log attributes
		attrs := []slog.Attr{
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.String("ip", c.IP()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration", duration),
			slog.Int("bytes_in", len(c.Body())),
			slog.Int("bytes_out", len(c.Response().Body())),
		}

		if traceID != "" {
			attrs = append(attrs, slog.String("trace_id", traceID))
		}

		if userID != "" {
			attrs = append(attrs, slog.String("user_id", userID))
		}

		// Add query parameters if present
		if len(c.Queries()) > 0 {
			attrs = append(attrs, slog.String("query", c.OriginalURL()))
		}

		// Log based on status code
		ctx := c.UserContext()
		if traceID != "" {
			ctx = WithTraceID(ctx, traceID)
		}

		status := c.Response().StatusCode()
		switch {
		case status >= 500:
			log.LogAttrs(ctx, slog.LevelError, "http request", attrs...)
		case status >= 400:
			log.LogAttrs(ctx, slog.LevelWarn, "http request", attrs...)
		default:
			log.LogAttrs(ctx, slog.LevelInfo, "http request", attrs...)
		}

		return err
	}
}
