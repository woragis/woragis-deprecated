package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"log/slog"
)

// mockRabbitMQConnection is a mock RabbitMQ connection for health checks
type mockRabbitMQConnection struct {
	closed bool
}

func (m *mockRabbitMQConnection) IsClosed() bool {
	return m.closed
}

func TestNewHealthChecker(t *testing.T) {
	logger := slog.Default()
	conn := &mockRabbitMQConnection{closed: false}

	checker := NewHealthChecker(conn, logger)
	if checker == nil {
		t.Error("NewHealthChecker() returned nil")
	}
}

func TestHealthChecker_Check(t *testing.T) {
	tests := []struct {
		name           string
		conn           *mockRabbitMQConnection
		wantStatus     string
		wantRabbitMQOk bool
	}{
		{
			name:           "healthy connection",
			conn:           &mockRabbitMQConnection{closed: false},
			wantStatus:     StatusHealthy,
			wantRabbitMQOk: true,
		},
		{
			name:           "unhealthy connection",
			conn:           &mockRabbitMQConnection{closed: true},
			wantStatus:     StatusUnhealthy,
			wantRabbitMQOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.Default()
			checker := NewHealthChecker(tt.conn, logger)

			result := checker.Check(context.Background())
			if result.Status != tt.wantStatus {
				t.Errorf("HealthChecker.Check() Status = %v, want %v", result.Status, tt.wantStatus)
			}

			// Check RabbitMQ check result
			rabbitmqCheck := findCheck(result.Checks, "rabbitmq")
			if rabbitmqCheck == nil {
				t.Error("HealthChecker.Check() should include rabbitmq check")
			} else {
				if (rabbitmqCheck.Status == "ok") != tt.wantRabbitMQOk {
					t.Errorf("HealthChecker.Check() RabbitMQ status = %v, want ok=%v", rabbitmqCheck.Status, tt.wantRabbitMQOk)
				}
			}
		})
	}
}

func findCheck(checks []CheckResult, name string) *CheckResult {
	for i := range checks {
		if checks[i].Name == name {
			return &checks[i]
		}
	}
	return nil
}

func TestHealthChecker_Handler(t *testing.T) {
	logger := slog.Default()
	conn := &mockRabbitMQConnection{closed: false}
	checker := NewHealthChecker(conn, logger)

	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	handler := checker.Handler()
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("HealthChecker.Handler() status code = %v, want %v", rr.Code, http.StatusOK)
	}

	// Verify response body contains status
	if rr.Body.String() == "" {
		t.Error("HealthChecker.Handler() should return response body")
	}
}

func TestHealthChecker_Caching(t *testing.T) {
	logger := slog.Default()
	conn := &mockRabbitMQConnection{closed: false}
	checker := NewHealthChecker(conn, logger)

	// First check
	result1 := checker.Check(context.Background())

	// Second check immediately after (should use cache)
	result2 := checker.Check(context.Background())

	// Results should be the same (cached)
	if result1.Status != result2.Status {
		t.Error("HealthChecker.Check() should use cache for immediate subsequent calls")
	}

	// Wait for cache to expire (cacheTTL is 5 seconds)
	time.Sleep(6 * time.Second)

	// Third check after cache expiry
	result3 := checker.Check(context.Background())

	// Should still be the same (connection state hasn't changed)
	if result3.Status != result1.Status {
		t.Error("HealthChecker.Check() should return consistent results")
	}
}

func TestHealthChecker_UnhealthyStatus(t *testing.T) {
	logger := slog.Default()
	conn := &mockRabbitMQConnection{closed: true}
	checker := NewHealthChecker(conn, logger)

	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	handler := checker.Handler()
	handler(rr, req)

	// Unhealthy status should return 503
	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("HealthChecker.Handler() status code = %v, want %v for unhealthy", rr.Code, http.StatusServiceUnavailable)
	}
}
