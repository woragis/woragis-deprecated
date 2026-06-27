package tracing

import (
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/gofiber/fiber/v2"
)

// Middleware creates OpenTelemetry tracing middleware for Fiber
// It extracts trace context from headers and creates spans for each request
// This should be used BEFORE RequestIDMiddleware to ensure trace context is propagated
func Middleware(serviceName string) fiber.Handler {
	tracer := otel.Tracer(serviceName)

	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()

		// Extract trace context from headers (W3C Trace Context standard)
		// Fiber uses fasthttp, so we need to extract headers manually
		headers := make(map[string]string)
		c.Request().Header.VisitAll(func(key, value []byte) {
			headers[string(key)] = string(value)
		})
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(headers))

		// Also check for X-Trace-ID header (for compatibility)
		if traceID := c.Get("X-Trace-ID"); traceID != "" {
			// If we have a trace ID in header but no valid span context,
			// we'll let OpenTelemetry create a new span but preserve the trace ID concept
			ctx = ContextWithTraceID(ctx, traceID)
		}

		// Start span for this request
		ctx, span := tracer.Start(
			ctx,
			c.Method()+" "+c.Path(),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethodKey.String(c.Method()),
				semconv.HTTPRouteKey.String(c.Path()),
				semconv.HTTPURLKey.String(c.OriginalURL()),
			),
		)
		defer span.End()

		// Update request context
		c.SetUserContext(ctx)

		// Get trace ID from span and ensure it's in context for logger
		traceID := span.SpanContext().TraceID().String()
		if traceID != "" {
			ctx = ContextWithTraceID(ctx, traceID)
			c.Set("X-Trace-ID", traceID)
			// Also set traceparent header (W3C standard)
			// Format: 00-{trace-id}-{span-id}-{flags}
			spanID := span.SpanContext().SpanID().String()
			traceparent := fmt.Sprintf("00-%s-%s-01", traceID, spanID)
			c.Set("traceparent", traceparent)
		}

		// Process request
		err := c.Next()

		// Record status code
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(c.Response().StatusCode()))

		// Record error if any
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else if c.Response().StatusCode() >= 400 {
			span.SetStatus(codes.Error, "HTTP error")
		} else {
			span.SetStatus(codes.Ok, "OK")
		}

		return err
	}
}

