package logger

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		env  string
	}{
		{
			name: "development environment",
			env:  "development",
		},
		{
			name: "production environment",
			env:  "production",
		},
		{
			name: "empty environment",
			env:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.env)
			if logger == nil {
				t.Error("New() returned nil logger")
			}
		})
	}
}

func TestNewWithConfig(t *testing.T) {
	cfg := LogConfig{
		Env:       "development",
		LogToFile: false,
		LogDir:    "logs",
	}

	logger := NewWithConfig(cfg)
	if logger == nil {
		t.Error("NewWithConfig() returned nil logger")
	}
}

func TestWithTraceID(t *testing.T) {
	ctx := context.Background()
	traceID := "test-trace-id"
	ctxWithTrace := WithTraceID(ctx, traceID)

	retrievedTraceID := GetTraceID(ctxWithTrace)
	if retrievedTraceID != traceID {
		t.Errorf("GetTraceID() = %v, want %v", retrievedTraceID, traceID)
	}
}

func TestGetTraceID(t *testing.T) {
	ctx := context.Background()
	traceID := GetTraceID(ctx)
	// TraceID may be empty if not set, which is valid
	if traceID != "" {
		t.Errorf("GetTraceID() = %v, want empty string for context without trace ID", traceID)
	}

	ctxWithTrace := WithTraceID(ctx, "test-id")
	traceID = GetTraceID(ctxWithTrace)
	if traceID != "test-id" {
		t.Errorf("GetTraceID() = %v, want test-id", traceID)
	}
}
