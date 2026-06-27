package logger

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	apptracing "woragis-skills-service/pkg/tracing"
)

// RequestIDMiddleware generates and adds a trace_id (request ID) to each request.
// The trace_id is added to the context and response headers for distributed tracing.
// This works with OpenTelemetry tracing - if a trace ID exists from OpenTelemetry,
// it will be used; otherwise, a new one is generated.
func RequestIDMiddleware(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// First, try to get trace ID from OpenTelemetry span context
		// This ensures compatibility with distributed tracing
		traceID := apptracing.TraceIDFromContext(ctx)

		// If no trace ID from OpenTelemetry, check header
		if traceID == "" {
			traceID = c.Get("X-Trace-ID")
		}

		// If still no trace ID, generate new one
		if traceID == "" {
			traceID = uuid.New().String()
		}

		// Add trace_id to response header
		c.Set("X-Trace-ID", traceID)

		// Add trace_id to context for logging (preserves OpenTelemetry context)
		ctx = WithTraceID(ctx, traceID)

		// Store trace_id in locals for easy access
		c.Locals("trace_id", traceID)

		// Update request context
		c.SetUserContext(ctx)

		return c.Next()
	}
}

// RequestLoggerMiddleware logs HTTP requests with trace_id.
// This should be used after RequestIDMiddleware to ensure trace_id is available.
func RequestLoggerMiddleware(logger *slog.Logger) fiber.Handler {
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

		// Build log attributes
		attrs := []slog.Attr{
			slog.String("method", c.Method()),
			slog.String("path", c.Path()),
			slog.Int("status", c.Response().StatusCode()),
			slog.Duration("duration", duration),
		}

		if traceID != "" {
			attrs = append(attrs, slog.String("trace_id", traceID))
		}

		// Log based on status code
		ctx := c.UserContext()
		if traceID != "" {
			ctx = WithTraceID(ctx, traceID)
		}

		// Log request
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, "http request",
				attrs...,
			)
		} else if c.Response().StatusCode() >= 500 {
			logger.LogAttrs(ctx, slog.LevelError, "http request",
				attrs...,
			)
		} else if c.Response().StatusCode() >= 400 {
			logger.LogAttrs(ctx, slog.LevelWarn, "http request",
				attrs...,
			)
		} else {
			logger.LogAttrs(ctx, slog.LevelInfo, "http request",
				attrs...,
			)
		}

		return err
	}
}

