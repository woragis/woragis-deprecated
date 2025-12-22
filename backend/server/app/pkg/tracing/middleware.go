package tracing

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// Middleware creates OpenTelemetry tracing middleware for Fiber
// It extracts trace context from headers and creates spans for each request
// This should be used BEFORE RequestIDMiddleware to ensure trace context is propagated
func Middleware(serviceName string) fiber.Handler {
	propagator := otel.GetTextMapPropagator()
	tracer := otel.Tracer(serviceName)

	return func(c *fiber.Ctx) error {
		// Extract trace context from headers (W3C Trace Context standard)
		ctx := propagator.Extract(c.UserContext(), propagation.HeaderCarrier(c.GetReqHeaders()))

		// Also check for X-Trace-ID header (for compatibility)
		if traceID := c.Get("X-Trace-ID"); traceID != "" {
			// If we have a trace ID in header but no valid span context,
			// we'll let OpenTelemetry create a new span but preserve the trace ID concept
			ctx = ContextWithTraceID(ctx, traceID)
		}

		// Start span for this request
		ctx, span := tracer.Start(
			ctx,
			fmt.Sprintf("%s %s", c.Method(), c.Path()),
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				semconv.HTTPMethodKey.String(c.Method()),
				semconv.HTTPRouteKey.String(c.Path()),
				semconv.HTTPURLKey.String(c.OriginalURL()),
				semconv.HTTPUserAgentKey.String(c.Get("User-Agent")),
				semconv.NetSockPeerAddrKey.String(c.IP()),
			),
		)

		// Store span in context
		c.SetUserContext(ctx)

		// Get trace ID from span and ensure it's in context for logger
		traceID := span.SpanContext().TraceID().String()
		if traceID != "" {
			ctx = ContextWithTraceID(ctx, traceID)
			c.SetUserContext(ctx)
			// Set in response header for propagation
			c.Set("X-Trace-ID", traceID)
			// Also set traceparent header (W3C standard)
			propagator.Inject(ctx, propagation.HeaderCarrier(c.GetRespHeaders()))
		}

		// Process request
		err := c.Next()

		// Set status code
		statusCode := c.Response().StatusCode()
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(statusCode))

		// Set span status based on HTTP status code
		if statusCode >= 500 {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", statusCode))
		} else if statusCode >= 400 {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", statusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}

		// Record error if present
		if err != nil {
			span.RecordError(err)
		}

		// End span
		span.End()

		return err
	}
}
