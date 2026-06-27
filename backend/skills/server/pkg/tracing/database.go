package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// WithDatabaseSpan creates a span for a database operation
// This should be used to wrap database operations in repositories
func WithDatabaseSpan(ctx context.Context, operation, table string, fn func() error) error {
	ctx, span := StartSpan(ctx, "db."+operation,
		trace.WithAttributes(
			attribute.String("db.operation", operation),
			attribute.String("db.table", table),
		),
	)
	defer span.End()

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	// Record duration as attribute
	span.SetAttributes(attribute.Float64("db.duration_ms", float64(duration.Nanoseconds())/1e6))

	if err != nil {
		RecordError(ctx, err)
		span.SetAttributes(attribute.String("db.error", err.Error()))
		return err
	}

	span.SetStatus(codes.Ok, "OK")
	return nil
}
