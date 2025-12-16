package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"
)

const (
	// ServiceName is the name of the service for structured logging
	ServiceName = "server"
	// TraceIDKey is the context key for trace ID
	TraceIDKey = "trace_id"
)

// New creates a slog.Logger configured for the supplied environment.
// The logger automatically includes service name and supports trace_id from context.
func New(env string) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		AddSource: false, // Set to true if you want source file/line info
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Ensure timestamp is in ISO 8601 format
			if a.Key == slog.TimeKey {
				return slog.String("timestamp", a.Value.Time().Format(time.RFC3339Nano))
			}
			// Ensure level is lowercase
			if a.Key == slog.LevelKey {
				return slog.String("level", a.Value.String())
			}
			return a
		},
	}

	switch strings.ToLower(env) {
	case "production", "prod":
		opts.Level = slog.LevelInfo
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	// Wrap handler to add service name and trace_id
	handler = &serviceHandler{
		Handler: handler,
		service: ServiceName,
	}

	return slog.New(handler)
}

// serviceHandler wraps a slog.Handler to automatically add service name and trace_id
type serviceHandler struct {
	slog.Handler
	service string
}

func (h *serviceHandler) Handle(ctx context.Context, r slog.Record) error {
	// Add service name to all logs
	r.AddAttrs(slog.String("service", h.service))

	// Add trace_id from context if available
	if traceID := ctx.Value(TraceIDKey); traceID != nil {
		if id, ok := traceID.(string); ok && id != "" {
			r.AddAttrs(slog.String("trace_id", id))
		}
	}

	return h.Handler.Handle(ctx, r)
}

// WithTraceID adds a trace_id to the context for distributed tracing
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// GetTraceID retrieves the trace_id from context
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}
