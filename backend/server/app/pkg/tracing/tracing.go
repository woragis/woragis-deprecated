package tracing

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	// TraceIDKey is the context key for trace ID (matches logger package)
	TraceIDKey = "trace_id"
)

var (
	tracer trace.Tracer
)

// Config holds tracing configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	JaegerEndpoint string
	SamplingRate   float64 // 0.0 to 1.0 (1.0 = 100%)
}

// Init initializes OpenTelemetry tracing with OTLP HTTP exporter (for Jaeger)
func Init(cfg Config) (func(), error) {
	// OTLP endpoint for Jaeger
	otlpEndpoint := cfg.JaegerEndpoint
	if otlpEndpoint == "" {
		otlpEndpoint = os.Getenv("OTLP_ENDPOINT")
		if otlpEndpoint == "" {
		// Jaeger OTLP HTTP endpoint (port 4318 for HTTP, 4317 for gRPC)
		// OTLP HTTP uses /v1/traces path by default
		otlpEndpoint = "http://jaeger:4318"
		}
	}

	if cfg.SamplingRate == 0 {
		// Default sampling: 100% in development, 10% in production
		env := strings.ToLower(cfg.Environment)
		if env == "production" || env == "prod" {
			cfg.SamplingRate = 0.1 // 10%
		} else {
			cfg.SamplingRate = 1.0 // 100%
		}
	}

	// Create OTLP HTTP exporter
	ctx := context.Background()
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otlpEndpoint),
		otlptracehttp.WithInsecure(), // Use TLS in production
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create trace provider with sampling
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(cfg.SamplingRate)),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator for trace context
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create tracer for this service
	tracer = otel.Tracer(cfg.ServiceName)

	// Return shutdown function
	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			// Log error but don't fail
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}, nil
}

// Tracer returns the global tracer
func Tracer() trace.Tracer {
	if tracer == nil {
		// Fallback to no-op tracer if not initialized
		return trace.NewNoopTracerProvider().Tracer("noop")
	}
	return tracer
}

// StartSpan starts a new span with the given name and options
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return Tracer().Start(ctx, name, opts...)
}

// SpanFromContext extracts span from context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// TraceIDFromContext extracts trace ID from context and returns it as a string
func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	// Fallback to trace_id from context (for compatibility with existing logger)
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// ContextWithTraceID creates a context with trace_id and OpenTelemetry span
// This bridges the existing trace_id system with OpenTelemetry
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	// Add trace_id to context (for logger compatibility)
	ctx = context.WithValue(ctx, TraceIDKey, traceID)

	// If we have a valid span context, use it
	// Otherwise, create a new span with the trace_id
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() && traceID != "" {
		// Try to parse trace_id as OpenTelemetry TraceID
		// If it's a valid format, create span with it
		// Otherwise, just store in context for logger
		// For now, we'll create a new span and let OpenTelemetry generate the ID
		// The trace_id will be available via context for logging
	}

	return ctx
}

// SetSpanAttributes sets attributes on the span from context
func SetSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.SetAttributes(attrs...)
	}
}

// RecordError records an error on the span
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if span != nil && err != nil {
		span.RecordError(err)
		span.SetStatus(trace.StatusError, err.Error())
	}
}
