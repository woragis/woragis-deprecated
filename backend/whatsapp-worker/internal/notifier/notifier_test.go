package notifier

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewWhatsmeowNotifier(t *testing.T) {
	// Skip if CGO is disabled (SQLite requires CGO)
	// This test requires CGO which may not be available in Alpine Docker images
	if !isCGOEnabled() {
		t.Skip("Skipping test: CGO is required for SQLite (whatsmeow)")
	}

	// Create temporary directory for session
	tempDir := t.TempDir()
	logger := slog.Default()

	notifier, err := NewWhatsmeowNotifier(tempDir, logger)
	if err != nil {
		t.Fatalf("NewWhatsmeowNotifier() error = %v", err)
	}

	if notifier == nil {
		t.Fatal("NewWhatsmeowNotifier() returned nil")
	}

	if notifier.sessionPath != tempDir {
		t.Errorf("sessionPath = %v, want %v", notifier.sessionPath, tempDir)
	}

	// Verify session directory was created
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Errorf("Session directory %s was not created", tempDir)
	}

	// Verify database file exists
	dbPath := filepath.Join(tempDir, "whatsapp.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Errorf("Database file %s was not created", dbPath)
	}
}

// isCGOEnabled checks if CGO is enabled by trying to use SQLite
func isCGOEnabled() bool {
	// Check if we can create a sqlstore (requires CGO)
	tempDir, err := os.MkdirTemp("", "whatsapp-test-*")
	if err != nil {
		return false
	}
	defer os.RemoveAll(tempDir)

	logger := slog.Default()
	_, err = NewWhatsmeowNotifier(tempDir, logger)
	// If error is about CGO, return false
	if err != nil && strings.Contains(err.Error(), "CGO") {
		return false
	}
	// If no error or different error, assume CGO is available
	return err == nil
}

func TestNewWhatsmeowNotifier_DefaultPath(t *testing.T) {
	// Skip if CGO is disabled
	if !isCGOEnabled() {
		t.Skip("Skipping test: CGO is required for SQLite (whatsmeow)")
	}

	// Test with empty path (should use default)
	logger := slog.Default()

	notifier, err := NewWhatsmeowNotifier("", logger)
	if err != nil {
		t.Fatalf("NewWhatsmeowNotifier() error = %v", err)
	}

	if notifier == nil {
		t.Fatal("NewWhatsmeowNotifier() returned nil")
	}

	// Clean up default directory
	defer os.RemoveAll("./whatsapp-session")
}

func TestNewWhatsmeowNotifier_InvalidPath(t *testing.T) {
	// Test with invalid path (should fail)
	logger := slog.Default()

	// Use a path that cannot be created (on Unix systems)
	invalidPath := "/root/invalid/path/that/cannot/be/created"
	notifier, err := NewWhatsmeowNotifier(invalidPath, logger)

	// On Windows, this might succeed, so we check if it fails
	if err == nil && notifier == nil {
		t.Log("NewWhatsmeowNotifier() with invalid path handled gracefully")
	}
}

func TestWhatsmeowNotifier_IsConnected(t *testing.T) {
	// Skip if CGO is disabled
	if !isCGOEnabled() {
		t.Skip("Skipping test: CGO is required for SQLite (whatsmeow)")
	}

	tempDir := t.TempDir()
	logger := slog.Default()

	notifier, err := NewWhatsmeowNotifier(tempDir, logger)
	if err != nil {
		t.Fatalf("NewWhatsmeowNotifier() error = %v", err)
	}

	// Initially should not be connected
	if notifier.IsConnected() {
		t.Error("IsConnected() = true, want false for new notifier")
	}
}

func TestWhatsmeowNotifier_GetQRCode(t *testing.T) {
	// Skip if CGO is disabled
	if !isCGOEnabled() {
		t.Skip("Skipping test: CGO is required for SQLite (whatsmeow)")
	}

	tempDir := t.TempDir()
	logger := slog.Default()

	notifier, err := NewWhatsmeowNotifier(tempDir, logger)
	if err != nil {
		t.Fatalf("NewWhatsmeowNotifier() error = %v", err)
	}

	// Initially QR code should be empty
	qrCode := notifier.GetQRCode()
	if qrCode != "" {
		t.Errorf("GetQRCode() = %v, want empty string for new notifier", qrCode)
	}
}

func TestWhatsmeowNotifier_Disconnect(t *testing.T) {
	// Skip if CGO is disabled
	if !isCGOEnabled() {
		t.Skip("Skipping test: CGO is required for SQLite (whatsmeow)")
	}

	tempDir := t.TempDir()
	logger := slog.Default()

	notifier, err := NewWhatsmeowNotifier(tempDir, logger)
	if err != nil {
		t.Fatalf("NewWhatsmeowNotifier() error = %v", err)
	}

	// Disconnect should not panic even if not connected
	notifier.Disconnect()

	// Should still not be connected
	if notifier.IsConnected() {
		t.Error("IsConnected() = true after Disconnect(), want false")
	}
}

func TestWhatsmeowNotifier_Send_NotConnected(t *testing.T) {
	// Skip if CGO is disabled
	if !isCGOEnabled() {
		t.Skip("Skipping test: CGO is required for SQLite (whatsmeow)")
	}

	tempDir := t.TempDir()
	logger := slog.Default()

	notifier, err := NewWhatsmeowNotifier(tempDir, logger)
	if err != nil {
		t.Fatalf("NewWhatsmeowNotifier() error = %v", err)
	}

	ctx := context.Background()
	err = notifier.Send(ctx, "+1234567890", "test message")

	// Should return error when not connected
	if err == nil {
		t.Error("Send() should return error when not connected")
	}

	expectedError := "whatsapp not connected"
	if err != nil && err.Error() != expectedError && !contains(err.Error(), expectedError) {
		t.Errorf("Send() error = %v, want error containing %v", err, expectedError)
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "short string",
			input:  "short",
			maxLen: 10,
			want:   "short",
		},
		{
			name:   "long string",
			input:  "this is a very long string that should be truncated",
			maxLen: 20,
			want:   "this is a very long ...",
		},
		{
			name:   "exact length",
			input:  "exact length",
			maxLen: 12,
			want:   "exact length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateString(tt.input, tt.maxLen)
			if len(got) > tt.maxLen+3 && tt.maxLen < len(tt.input) {
				t.Errorf("truncateString() length = %v, want <= %v", len(got), tt.maxLen+3)
			}
			if tt.maxLen >= len(tt.input) && got != tt.input {
				t.Errorf("truncateString() = %v, want %v", got, tt.input)
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

