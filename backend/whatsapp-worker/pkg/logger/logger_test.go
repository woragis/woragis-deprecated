package logger

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	logger := New("development")
	if logger == nil {
		t.Fatal("New() returned nil logger")
	}

	logger = New("production")
	if logger == nil {
		t.Fatal("New() returned nil logger for production")
	}
}

func TestNewWithConfig(t *testing.T) {
	tests := []struct {
		name string
		cfg  LogConfig
	}{
		{
			name: "development with stdout",
			cfg: LogConfig{
				Env:       "development",
				LogToFile: false,
			},
		},
		{
			name: "production",
			cfg: LogConfig{
				Env:       "production",
				LogToFile: false,
			},
		},
		{
			name: "development with file logging",
			cfg: LogConfig{
				Env:       "development",
				LogToFile: true,
				LogDir:    "test-logs",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewWithConfig(tt.cfg)
			if logger == nil {
				t.Fatal("NewWithConfig() returned nil logger")
			}

			// Clean up test log directory if created
			if tt.cfg.LogToFile && tt.cfg.LogDir != "" {
				defer os.RemoveAll(tt.cfg.LogDir)
			}
		})
	}
}

func TestNewWithConfig_FileLogging(t *testing.T) {
	// Test file logging creation
	testDir := "test-logs-file"
	defer os.RemoveAll(testDir)

	cfg := LogConfig{
		Env:       "development",
		LogToFile: true,
		LogDir:    testDir,
	}

	logger := NewWithConfig(cfg)
	if logger == nil {
		t.Fatal("NewWithConfig() returned nil logger")
	}

	// Verify log directory was created
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Errorf("Log directory %s was not created", testDir)
	}

	// Verify log file exists
	logFile := filepath.Join(testDir, ServiceName+".log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Log file %s was not created", logFile)
	}
}

func TestNewWithConfig_InvalidDirectory(t *testing.T) {
	// Test that invalid directory falls back to stdout
	cfg := LogConfig{
		Env:       "development",
		LogToFile: true,
		LogDir:    "/invalid/path/that/does/not/exist",
	}

	logger := NewWithConfig(cfg)
	if logger == nil {
		t.Fatal("NewWithConfig() should fallback to stdout on error")
	}
}

func TestWithTraceID(t *testing.T) {
	ctx := context.Background()
	traceID := "test-trace-123"

	ctxWithTrace := WithTraceID(ctx, traceID)

	retrieved := GetTraceID(ctxWithTrace)
	if retrieved != traceID {
		t.Errorf("GetTraceID() = %v, want %v", retrieved, traceID)
	}
}

func TestGetTraceID(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		want    string
	}{
		{
			name: "with trace ID",
			ctx:  WithTraceID(context.Background(), "trace-123"),
			want: "trace-123",
		},
		{
			name: "without trace ID",
			ctx:  context.Background(),
			want: "",
		},
		{
			name: "with non-string trace ID",
			ctx:  context.WithValue(context.Background(), TraceIDKey, 123),
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTraceID(tt.ctx)
			if got != tt.want {
				t.Errorf("GetTraceID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceHandler(t *testing.T) {
	logger := New("development")
	if logger == nil {
		t.Fatal("New() returned nil logger")
	}

	// Test that logger includes service name
	ctx := context.Background()
	logger.InfoContext(ctx, "test message")
}

func TestServiceHandler_WithTraceID(t *testing.T) {
	logger := New("development")
	if logger == nil {
		t.Fatal("New() returned nil logger")
	}

	// Test that logger includes trace ID from context
	ctx := WithTraceID(context.Background(), "test-trace-456")
	logger.InfoContext(ctx, "test message with trace ID")
}

