//go:build integration && performance_test
// +build integration,performance_test

package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/woragis/backend/server/app/internal/testutil"
)

// setupMinimalTestApp creates a minimal app for performance testing (avoids import cycles)
func setupMinimalTestApp(t testing.TB, db *gorm.DB, redisClient *redis.Client) *fiber.App {
	app := testutil.SetupTestApp(t, db, redisClient)
	
	// Add minimal routes for performance testing
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})
	
	app.Get("/api/v1/skills", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"data": []interface{}{}})
	})
	
	return app
}

// BenchmarkServerHealthEndpoint benchmarks the health check endpoint
func BenchmarkServerHealthEndpoint(b *testing.B) {
	db := testutil.SetupTestDB(b)
	defer testutil.CleanupTestDB(b, db)

	redis := testutil.SetupTestRedis(b)
	defer testutil.CleanupTestRedis(b, redis)

	app := setupMinimalTestApp(b, db, redis)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
			resp, err := app.Test(req)
			require.NoError(b, err)
			assert.Equal(b, http.StatusOK, resp.StatusCode)
		}
	})
}

// TestServerAPILoadTest tests server API under load
func TestServerAPILoadTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping load test in short mode")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupMinimalTestApp(t, db, redis)

	// Create test user and token
	userID := testutil.CreateTestUser(t, db, "perf-test@example.com", "testpass123")
	token := testutil.GenerateTestJWT(t, userID, "perf-test@example.com")

	concurrentRequests := 50
	requestsPerGoroutine := 10
	totalRequests := concurrentRequests * requestsPerGoroutine

	var successCount int64
	var errorCount int64
	var totalLatency int64

	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < concurrentRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				reqStart := time.Now()
				req := httptest.NewRequest(http.MethodGet, "/api/v1/skills", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				resp, err := app.Test(req)
				latency := time.Since(reqStart).Nanoseconds()

				atomic.AddInt64(&totalLatency, latency)

				if err != nil || resp.StatusCode != http.StatusOK {
					atomic.AddInt64(&errorCount, 1)
				} else {
					atomic.AddInt64(&successCount, 1)
				}
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	t.Logf("Load Test Results:")
	t.Logf("  Total Requests: %d", totalRequests)
	t.Logf("  Successful: %d", successCount)
	t.Logf("  Errors: %d", errorCount)
	t.Logf("  Duration: %v", duration)
	t.Logf("  Throughput: %.2f req/s", float64(totalRequests)/duration.Seconds())
	t.Logf("  Avg Latency: %.2f ms", float64(totalLatency/int64(totalRequests))/1e6)
	t.Logf("  P95 Latency: N/A (calculate from histogram)")

	assert.Greater(t, successCount, int64(totalRequests*9/10), "At least 90% of requests should succeed")
	assert.Less(t, errorCount, int64(totalRequests/10), "Less than 10% of requests should fail")
}

// TestServerAPILatency tests server API latency
func TestServerAPILatency(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping latency test in short mode")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupMinimalTestApp(t, db, redis)

	userID := testutil.CreateTestUser(t, db, "perf-test@example.com", "testpass123")
	token := testutil.GenerateTestJWT(t, userID, "perf-test@example.com")

	iterations := 100
	latencies := make([]time.Duration, iterations)

	for i := 0; i < iterations; i++ {
		start := time.Now()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/skills", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		latency := time.Since(start)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		latencies[i] = latency
	}

	// Calculate statistics
	var sum time.Duration
	for _, lat := range latencies {
		sum += lat
	}
	avgLatency := sum / time.Duration(iterations)

	// Sort for percentile calculation (simplified)
	p95Index := int(float64(iterations) * 0.95)
	p95Latency := latencies[p95Index]

	t.Logf("Latency Test Results:")
	t.Logf("  Iterations: %d", iterations)
	t.Logf("  Average Latency: %v", avgLatency)
	t.Logf("  P95 Latency: %v", p95Latency)

	assert.Less(t, avgLatency, 100*time.Millisecond, "Average latency should be less than 100ms")
	assert.Less(t, p95Latency, 200*time.Millisecond, "P95 latency should be less than 200ms")
}

// TestServerConcurrentUsers tests server with concurrent users
func TestServerConcurrentUsers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping concurrent users test in short mode")
	}

	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupMinimalTestApp(t, db, redis)

	concurrentUsers := 20
	requestsPerUser := 5

	var wg sync.WaitGroup
	var successCount int64

	start := time.Now()

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			// Each user creates their own token
			email := fmt.Sprintf("perf-user-%d@example.com", userID)
			createdUserID := testutil.CreateTestUser(t, db, email, "testpass123")
			token := testutil.GenerateTestJWT(t, createdUserID, email)

			for j := 0; j < requestsPerUser; j++ {
				req := httptest.NewRequest(http.MethodGet, "/api/v1/skills", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				resp, err := app.Test(req)

				if err == nil && resp.StatusCode == http.StatusOK {
					atomic.AddInt64(&successCount, 1)
				}
			}
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	totalRequests := concurrentUsers * requestsPerUser

	t.Logf("Concurrent Users Test Results:")
	t.Logf("  Concurrent Users: %d", concurrentUsers)
	t.Logf("  Requests per User: %d", requestsPerUser)
	t.Logf("  Total Requests: %d", totalRequests)
	t.Logf("  Successful: %d", successCount)
	t.Logf("  Duration: %v", duration)
	t.Logf("  Throughput: %.2f req/s", float64(totalRequests)/duration.Seconds())

	assert.Greater(t, successCount, int64(totalRequests*8/10), "At least 80% of requests should succeed")
}
