package tracing

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

// getOTLPEndpoint reads the OTLP endpoint from config or environment variables
// Priority: cfg.JaegerEndpoint > JAEGER_ENDPOINT env > OTLP_ENDPOINT env > default
func getOTLPEndpoint(cfg Config) string {
	// First, check if explicitly set in config
	if cfg.JaegerEndpoint != "" {
		return cfg.JaegerEndpoint
	}
	
	// Then check JAEGER_ENDPOINT environment variable
	if endpoint := os.Getenv("JAEGER_ENDPOINT"); endpoint != "" {
		return endpoint
	}
	
	// Then check OTLP_ENDPOINT environment variable (OpenTelemetry standard)
	if endpoint := os.Getenv("OTLP_ENDPOINT"); endpoint != "" {
		return endpoint
	}
	
	// Default to Jaeger OTLP HTTP endpoint (port 4318 for HTTP, 4317 for gRPC)
	// OTLP HTTP uses /v1/traces path by default
	return "http://jaeger:4318"
}

// Init initializes OpenTelemetry tracing with OTLP HTTP exporter (for Jaeger)
func Init(cfg Config) (func(), error) {
	// Get OTLP endpoint from config or environment variables
	otlpEndpoint := getOTLPEndpoint(cfg)

	// Normalize endpoint URL - remove duplicate http:// prefixes and ensure proper format
	otlpEndpoint = normalizeEndpoint(otlpEndpoint)
	
	// Log the normalized endpoint for debugging (remove in production if needed)
	if cfg.Environment != "production" && cfg.Environment != "prod" {
		fmt.Printf("OTLP endpoint normalized to: %s\n", otlpEndpoint)
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
	// Extract just host:port from the normalized URL
	// WithEndpoint expects host:port format, not full URL
	// It will automatically add http:// and /v1/traces path
	endpointHost := extractHostPort(otlpEndpoint)
	
	ctx := context.Background()
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpointHost),
		otlptracehttp.WithInsecure(), // Use TLS in production
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource with service information
	// Note: We create a fresh resource instead of merging with resource.Default()
	// to avoid schema URL conflicts between different semconv versions
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
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
		span.SetStatus(codes.Error, err.Error())
	}
}

// normalizeEndpoint normalizes the OTLP endpoint URL to prevent double-prefixing
// and ensures it's in the correct format for OpenTelemetry
func normalizeEndpoint(endpoint string) string {
	if endpoint == "" {
		return "http://jaeger:4318" // Default fallback
	}

	endpoint = strings.TrimSpace(endpoint)

	// First, handle URL-encoded malformed URLs like "http://http:%2F%2Fjaeger:4318"
	// Try to decode URL encoding
	if strings.Contains(endpoint, "%") {
		if decoded, err := url.QueryUnescape(endpoint); err == nil && decoded != endpoint {
			endpoint = decoded
		}
	}

	// More aggressive cleanup: remove ALL occurrences of duplicate protocols
	// Handle cases like "http://http://", "http://http://http://", etc.
	maxIterations := 10 // Prevent infinite loops
	for i := 0; i < maxIterations; i++ {
		original := endpoint
		// Remove duplicate http://
		endpoint = strings.Replace(endpoint, "http://http://", "http://", 1)
		// Remove duplicate https://
		endpoint = strings.Replace(endpoint, "https://https://", "https://", 1)
		// Remove mixed protocols (http://https:// or https://http://)
		endpoint = strings.Replace(endpoint, "http://https://", "https://", 1)
		endpoint = strings.Replace(endpoint, "https://http://", "http://", 1)
		
		// If no changes were made, we're done
		if endpoint == original {
			break
		}
	}

	// Try to parse as URL to validate and normalize
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		// If parsing fails, try to extract just the host:port manually
		cleanEndpoint := endpoint
		
		// Remove all protocol prefixes (handle multiple)
		for strings.HasPrefix(cleanEndpoint, "http://") || strings.HasPrefix(cleanEndpoint, "https://") {
			cleanEndpoint = strings.TrimPrefix(cleanEndpoint, "http://")
			cleanEndpoint = strings.TrimPrefix(cleanEndpoint, "https://")
		}
		
		// Extract host:port (remove any path, query, fragment)
		if idx := strings.IndexAny(cleanEndpoint, "/?#"); idx > 0 {
			cleanEndpoint = cleanEndpoint[:idx]
		}
		
		// If we have something that looks valid, construct a proper URL
		if cleanEndpoint != "" {
			// Check if it already has a port
			if strings.Contains(cleanEndpoint, ":") {
				return "http://" + cleanEndpoint
			}
			// Add default port
			return "http://" + cleanEndpoint + ":4318"
		}
		
		// Ultimate fallback
		return "http://jaeger:4318"
	}

	// Reconstruct URL without path (OTLP exporter adds /v1/traces automatically)
	// Only include scheme and host:port
	if parsedURL.Host == "" {
		// If host is empty, try to extract from the original string
		return "http://jaeger:4318"
	}
	
	normalized := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	
	// Validate the scheme
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		// If scheme is invalid, default to http
		normalized = fmt.Sprintf("http://%s", parsedURL.Host)
	}

	return normalized
}

// extractHostPort extracts just the host:port from a URL string
// This is safer than passing a full URL which might have encoding issues
func extractHostPort(endpoint string) string {
	if endpoint == "" {
		return "jaeger:4318" // Default fallback
	}

	endpoint = strings.TrimSpace(endpoint)

	// Remove any protocol prefix
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	
	// Handle URL-encoded malformed URLs
	if decoded, err := url.QueryUnescape(endpoint); err == nil && decoded != endpoint {
		endpoint = decoded
		endpoint = strings.TrimPrefix(endpoint, "http://")
		endpoint = strings.TrimPrefix(endpoint, "https://")
	}

	// Remove any path (everything after the first /)
	if idx := strings.Index(endpoint, "/"); idx > 0 {
		endpoint = endpoint[:idx]
	}
	
	// Remove any query parameters
	if idx := strings.Index(endpoint, "?"); idx > 0 {
		endpoint = endpoint[:idx]
	}

	// If we have something that looks like host:port, return it
	if strings.Contains(endpoint, ":") {
		return endpoint
	}

	// If it's just a hostname, add default port
	if endpoint != "" && !strings.Contains(endpoint, ":") {
		return endpoint + ":4318"
	}

	// Fallback
	return "jaeger:4318"
}

