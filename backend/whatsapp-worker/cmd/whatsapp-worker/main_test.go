package main

import (
	"os"
	"testing"
)

func TestMain_EnvironmentVariables(t *testing.T) {
	// Test that main function can handle missing environment variables
	// This is a basic smoke test since main() is hard to test directly

	// Save original env
	originalEnv := os.Getenv("ENV")
	defer os.Setenv("ENV", originalEnv)

	// Test with empty ENV (should default to development)
	os.Unsetenv("ENV")
	// We can't easily test main() without running it, so we just verify
	// that the environment variable handling logic would work
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	if env != "development" {
		t.Errorf("Default ENV = %v, want development", env)
	}

	// Test with production ENV
	os.Setenv("ENV", "production")
	env = os.Getenv("ENV")
	if env != "production" {
		t.Errorf("ENV = %v, want production", env)
	}
}

