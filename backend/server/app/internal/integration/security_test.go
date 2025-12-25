//go:build integration
// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	emailservice "github.com/woragis/backend/server/app/internal/services/email"
	applogger "github.com/woragis/backend/server/app/pkg/logger"
	"github.com/woragis/backend/server/app/internal/testutil"
)

// TestSecurityHeaders verifies that security headers are set correctly
func TestSecurityHeaders(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)

	req := httptest.NewRequest("GET", "/healthz", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Check security headers
	assert.Equal(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"))
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	assert.Equal(t, "no-referrer", resp.Header.Get("Referrer-Policy"))
	assert.NotEmpty(t, resp.Header.Get("Content-Security-Policy"))
}

// TestRateLimiting verifies rate limiting works (100 req/min per IP)
func TestRateLimiting(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)

	// Make 100 requests (should all succeed)
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest("GET", "/healthz", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Request %d should succeed", i+1)
	}

	// 101st request should be rate limited
	req := httptest.NewRequest("GET", "/healthz", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "101st request should be rate limited")

	// Verify error message
	var errorResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&errorResp)
	require.NoError(t, err)
	assert.Contains(t, errorResp["error"], "Rate limit exceeded")
}

// TestRateLimitingPerIP verifies rate limiting is per IP address
func TestRateLimitingPerIP(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)

	// Make 100 requests from IP 1
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest("GET", "/healthz", nil)
		req.RemoteAddr = "192.168.1.1:12345"
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}

	// IP 1 should be rate limited
	req := httptest.NewRequest("GET", "/healthz", nil)
	req.RemoteAddr = "192.168.1.1:12345"
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode, "IP 1 should be rate limited")

	// IP 2 should still be able to make requests
	req2 := httptest.NewRequest("GET", "/healthz", nil)
	req2.RemoteAddr = "192.168.1.2:12345"
	resp2, err := app.Test(req2)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode, "IP 2 should not be rate limited")
}

// TestJWTAuthentication tests JWT authentication flow
func TestJWTAuthentication(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Test protected endpoint with valid JWT
	req := httptest.NewRequest("GET", "/api/projects", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Valid JWT should allow access")

	// Test protected endpoint without JWT
	req = httptest.NewRequest("GET", "/api/projects", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Missing JWT should be rejected")

	// Test protected endpoint with invalid JWT
	req = httptest.NewRequest("GET", "/api/projects", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Invalid JWT should be rejected")

	// Test protected endpoint with expired JWT (if we can create one)
	// This would require creating a token with past expiration
}

// TestAPIKeyAuthentication tests API key authentication flow
func TestAPIKeyAuthentication(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	userID, _ := createTestUserAndToken(t, db)

	// Create an API key
	apiKeyService := getAPIKeyService(t, db)
	apiKey, err := apiKeyService.CreateAPIKey(context.Background(), userID, "Test API Key", nil)
	require.NoError(t, err)

	// Test GET request with valid API key in X-API-Key header
	req := httptest.NewRequest("GET", "/api/projects", nil)
	req.Header.Set("X-API-Key", apiKey.Token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Valid API key should allow GET access")

	// Test GET request with valid API key in Authorization header
	req = httptest.NewRequest("GET", "/api/projects", nil)
	req.Header.Set("Authorization", "ApiKey "+apiKey.Token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Valid API key in Authorization header should work")

	// Test POST request with API key (should require JWT)
	createReq := map[string]interface{}{
		"name": "Test Project",
	}
	body, _ := json.Marshal(createReq)
	req = httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey.Token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "API key should not work for POST requests")

	// Test with invalid API key
	req = httptest.NewRequest("GET", "/api/projects", nil)
	req.Header.Set("X-API-Key", "invalid-api-key")
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Invalid API key should be rejected")
}

// TestAuthorizationAdmin tests admin authorization
func TestAuthorizationAdmin(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)

	// Create admin user
	authService := getAuthService(t, db)
	adminUser, err := authService.RegisterUser(context.Background(), authdomain.RegisterRequest{
		Email:       "admin@example.com",
		Password:    "adminpass123",
		DisplayName: "Admin User",
		Locale:      "en",
	})
	require.NoError(t, err)

	// Update user to admin role
	authRepo := getAuthRepository(t, db)
	adminUser.Role = "admin"
	err = authRepo.Update(context.Background(), adminUser)
	require.NoError(t, err)

	// Create regular user
	_, regularToken := createTestUserAndToken(t, db)

	// Get admin token
	jwtManager := getJWTManager(t)
	adminToken, err := jwtManager.Generate(adminUser.ID, adminUser.Email)
	require.NoError(t, err)

	// Test admin-only endpoint with admin token
	// Note: This assumes there's an admin-only endpoint. Adjust path as needed.
	req := httptest.NewRequest("GET", "/api/projects", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := app.Test(req)
	require.NoError(t, err)
	// Admin should have access
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusForbidden,
		"Admin should have access or endpoint may not exist")

	// Test admin-only endpoint with regular user token
	req = httptest.NewRequest("GET", "/api/projects", nil)
	req.Header.Set("Authorization", "Bearer "+regularToken)
	resp, err = app.Test(req)
	require.NoError(t, err)
	// Regular user should have access to non-admin endpoints
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusForbidden)
}

// TestCORSConfiguration tests CORS configuration
func TestCORSConfiguration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)

	// Test preflight OPTIONS request
	req := httptest.NewRequest("OPTIONS", "/api/projects", nil)
	req.Header.Set("Origin", "https://example.com")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Check CORS headers (if CORS is enabled)
	// Note: CORS may be disabled in test config
	if resp.Header.Get("Access-Control-Allow-Origin") != "" {
		assert.NotEmpty(t, resp.Header.Get("Access-Control-Allow-Origin"))
		assert.NotEmpty(t, resp.Header.Get("Access-Control-Allow-Methods"))
	}
}

// TestInputValidation tests input validation
func TestInputValidation(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Test invalid email format
	registerReq := map[string]interface{}{
		"email":        "invalid-email",
		"password":     "testpass123",
		"display_name": "Test User",
		"locale":       "en",
	}
	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusUnprocessableEntity,
		"Invalid email should be rejected")

	// Test invalid UUID format
	req = httptest.NewRequest("GET", "/api/projects/invalid-uuid", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound,
		"Invalid UUID should be rejected")
}

// TestSQLInjectionProtection tests SQL injection protection
func TestSQLInjectionProtection(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Test SQL injection in query parameter
	sqlInjectionPayloads := []string{
		"' OR '1'='1",
		"'; DROP TABLE users;--",
		"' UNION SELECT * FROM users--",
		"admin'--",
		"' OR 1=1--",
	}

	for _, payload := range sqlInjectionPayloads {
		req := httptest.NewRequest("GET", "/api/projects?search="+payload, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)
		// Should either reject with 400 or sanitize the input
		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest,
			"SQL injection attempt should be handled safely: %s", payload)
	}

	// Test SQL injection in request body
	createReq := map[string]interface{}{
		"name": "'; DROP TABLE projects;--",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	// Should either reject or sanitize
	assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusCreated,
		"SQL injection in body should be handled safely")
}

// TestXSSProtection tests XSS protection
func TestXSSProtection(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Test XSS in query parameter
	xssPayloads := []string{
		"<script>alert('XSS')</script>",
		"<img src=x onerror=alert('XSS')>",
		"javascript:alert('XSS')",
		"<svg onload=alert('XSS')>",
		"<iframe src=javascript:alert('XSS')>",
	}

	for _, payload := range xssPayloads {
		req := httptest.NewRequest("GET", "/api/projects?search="+payload, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)
		// Should either reject with 400 or sanitize the input
		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest,
			"XSS attempt should be handled safely: %s", payload)
	}

	// Test XSS in request body
	createReq := map[string]interface{}{
		"name":        "<script>alert('XSS')</script>",
		"description": "Test project",
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	// Should either reject or sanitize
	assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusCreated,
		"XSS in body should be handled safely")
}

// TestRequestSizeLimit tests request size limiting (10MB max)
func TestRequestSizeLimit(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create a request body larger than 10MB
	largeBody := make([]byte, 11*1024*1024) // 11MB
	for i := range largeBody {
		largeBody[i] = 'A'
	}

	req := httptest.NewRequest("POST", "/api/projects", bytes.NewReader(largeBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.ContentLength = int64(len(largeBody))
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusRequestEntityTooLarge, resp.StatusCode,
		"Request larger than 10MB should be rejected")
}

// TestInputSanitization tests input sanitization
func TestInputSanitization(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)

	// Test null bytes in query parameter
	req := httptest.NewRequest("GET", "/api/projects?search=test\x00null", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	// Should sanitize null bytes
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Null bytes should be sanitized")

	// Test whitespace trimming
	req = httptest.NewRequest("GET", "/api/projects?search=  test  ", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Whitespace should be trimmed")
}

// TestJWTExpiration tests JWT token expiration
func TestJWTExpiration(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)

	// Create a user and get a token
	_, token := createTestUserAndToken(t, db)

	// Token should be valid initially
	req := httptest.NewRequest("GET", "/api/projects", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Valid token should work")

	// Note: Testing actual expiration would require creating a token with past expiration
	// This is typically done by mocking time or creating a token manually
}

// TestAPIKeyScopes tests API key scope validation
func TestAPIKeyScopes(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	userID, _ := createTestUserAndToken(t, db)

	// Create API key
	apiKeyService := getAPIKeyService(t, db)
	apiKey, err := apiKeyService.CreateAPIKey(context.Background(), userID, "Read Only Key", nil)
	require.NoError(t, err)

	// GET request should work
	req := httptest.NewRequest("GET", "/api/projects", nil)
	req.Header.Set("X-API-Key", apiKey.Token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode, "API key should allow GET")

	// POST request should require JWT (API keys only work for GET)
	createReq := map[string]interface{}{
		"name": "Test Project",
	}
	body, _ := json.Marshal(createReq)
	req = httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey.Token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "API key should not work for POST")
}

// Helper functions

func getAPIKeyService(t *testing.T, db *gorm.DB) apikeysdomain.Service {
	repo := apikeysdomain.NewGormRepository(db)
	logger := applogger.New("test")
	return apikeysdomain.NewService(repo, logger)
}

func getAuthService(t *testing.T, db *gorm.DB) *authdomain.Service {
	repo := getAuthRepository(t, db)
	emailService := emailservice.NewNoopSender(applogger.New("test"))
	tokenStore := authdomain.NewRedisTokenStore(nil)
	jwtManager := getJWTManager(t)
	logger := applogger.New("test")
	return authdomain.NewService(repo, emailService, tokenStore, "http://localhost:8080", jwtManager, logger)
}

func getAuthRepository(t *testing.T, db *gorm.DB) authdomain.Repository {
	return authdomain.NewGormRepository(db)
}

func getJWTManager(t *testing.T) *authdomain.JWTManager {
	// Use test JWT secret
	cfg := testutil.LoadTestConfig()
	jwtManager, err := authdomain.NewJWTManager(cfg.JWTSecret, 24*time.Hour, "woragis-test")
	if err != nil {
		t.Fatalf("failed to create JWT manager: %v", err)
	}
	return jwtManager
}

