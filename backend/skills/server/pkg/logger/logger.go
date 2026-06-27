package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	// ServiceName is the name of the service for structured logging
	ServiceName = "woragis-skills-service"
	// TraceIDKey is the context key for trace ID
	TraceIDKey = "trace_id"
	// DefaultLogDir is the default directory for log files in development
	DefaultLogDir = "logs"
)

// LogConfig holds configuration for logger output
type LogConfig struct {
	// Env is the environment (development, production, etc.)
	Env string
	// LogDir is the directory for log files (only used in development)
	// If empty, logs go to stdout
	LogDir string
	// LogToFile enables file logging in development (default: false, uses stdout)
	LogToFile bool
}

// New creates a slog.Logger configured for the supplied environment.
// The logger automatically includes service name and supports trace_id from context.
//
// In production: logs go to stdout (for Kubernetes/log aggregation)
// In development: logs go to stdout by default, or to files if LogToFile is enabled
func New(env string) *slog.Logger {
	return NewWithConfig(LogConfig{
		Env:       env,
		LogToFile: false, // Default to stdout even in development
	})
}

// NewWithConfig creates a logger with custom configuration
func NewWithConfig(cfg LogConfig) *slog.Logger {
	var writer io.Writer = os.Stdout
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

	env := strings.ToLower(cfg.Env)
	isProduction := env == "production" || env == "prod"

	// In production, always use stdout (for Kubernetes/log aggregation)
	// In development, use files if LogToFile is enabled
	if !isProduction && cfg.LogToFile {
		logDir := cfg.LogDir
		if logDir == "" {
			logDir = DefaultLogDir
		}

		// Create log directory if it doesn't exist
		if err := os.MkdirAll(logDir, 0755); err != nil {
			// Fallback to stdout if directory creation fails
			writer = os.Stdout
		} else {
			// Open log file (append mode)
			logFile := filepath.Join(logDir, ServiceName+".log")
			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				// Fallback to stdout if file open fails
				writer = os.Stdout
			} else {
				// Use multi-writer to write to both file and stdout
				writer = io.MultiWriter(file, os.Stdout)
			}
		}
	}

	// Configure handler based on environment
	if isProduction {
		opts.Level = slog.LevelInfo
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(writer, opts)
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

// LogError logs an error with stack trace and context information
// This provides structured error logging for better debugging
func LogError(ctx context.Context, logger *slog.Logger, err error, msg string, attrs ...slog.Attr) {
	if err == nil {
		return
	}

	// Build error attributes
	errorAttrs := []slog.Attr{
		slog.String("error", err.Error()),
		slog.String("error_type", fmt.Sprintf("%T", err)),
	}

	// Add stack trace if available (runtime.Stack)
	stackTrace := getStackTrace()
	if stackTrace != "" {
		errorAttrs = append(errorAttrs, slog.String("stack_trace", stackTrace))
	}

	// Add any additional attributes
	errorAttrs = append(errorAttrs, attrs...)

	// Add trace_id if available
	if traceID := GetTraceID(ctx); traceID != "" {
		errorAttrs = append(errorAttrs, slog.String("trace_id", traceID))
	}

	// Log with error level
	logger.LogAttrs(ctx, slog.LevelError, msg, errorAttrs...)
}

// getStackTrace captures the current stack trace
func getStackTrace() string {
	// Capture stack trace (skip 3 frames: getStackTrace, LogError, caller)
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	if n > 0 {
		// Skip the first few lines (runtime.Stack header)
		lines := strings.Split(string(buf[:n]), "\n")
		if len(lines) > 4 {
			// Return stack trace starting from the actual caller
			return strings.Join(lines[4:], "\n")
		}
		return string(buf[:n])
	}
	return ""
}

