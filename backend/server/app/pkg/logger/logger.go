package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a slog.Logger configured for the supplied environment.
func New(env string) *slog.Logger {
	var handler slog.Handler

	switch strings.ToLower(env) {
	case "production", "prod":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(handler)
}
