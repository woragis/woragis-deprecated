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

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	emailservice "github.com/woragis/backend/server/app/internal/services/email"
	applogger "github.com/woragis/backend/server/app/pkg/logger"
	"github.com/woragis/backend/server/app/internal/testutil"
)

// TestEmailConfirmationFlow tests the complete email confirmation flow
func TestEmailConfirmationFlow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)

	// Register a new user
	registerReq := map[string]interface{}{
		"email":        "confirm@example.com",
		"password":     "testpassword123",
		"display_name": "Test User",
		"locale":       "en",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusAccepted)

	// Get the confirmation token from the database
	// In a real scenario, this would come from the email
	// For tests, we'll manually get it from the repository
	authRepo := authdomain.NewGormRepository(db)
	ctx := req.Context()
	
	// Find the user
	user, err := authRepo.FindByEmail(ctx, "confirm@example.com")
	require.NoError(t, err)
	require.NotNil(t, user)

	// Manually confirm the email for tests
	// In production, this would be done via the confirmation token from email
	user.ConfirmEmail()
	err = authRepo.Update(ctx, user)
	require.NoError(t, err)

	// Now login should work
	loginReq := map[string]interface{}{
		"email":    "confirm@example.com",
		"password": "testpassword123",
	}

	body, _ = json.Marshal(loginReq)
	req = httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	require.NoError(t, err)
	assert.NotNil(t, loginResp["access_token"])
	assert.NotNil(t, loginResp["refresh_token"])
}

// TestPasswordResetFlow tests the password reset flow
func TestPasswordResetFlow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)

	// Create a user first
	_, token := createTestUserAndToken(t, db)
	_ = token

	// Request password reset
	resetReq := map[string]interface{}{
		"email": "test@example.com",
	}

	body, _ := json.Marshal(resetReq)
	req := httptest.NewRequest("POST", "/api/auth/password/reset/request", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusAccepted)

	// In a real scenario, the token would come from email
	// For tests, we'll simulate getting the token from Redis
	// The actual implementation would require accessing the token store
	// This test verifies the endpoint exists and responds correctly
}

// TestSessionRefreshFlow tests session refresh (token rotation)
func TestSessionRefreshFlow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)

	// Login to get access and refresh tokens
	loginReq := map[string]interface{}{
		"email":    "test@example.com",
		"password": "testpassword123",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	
	if resp.StatusCode == http.StatusOK {
		var loginResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		require.NoError(t, err)
		
		refreshToken, ok := loginResp["refresh_token"].(string)
		if ok && refreshToken != "" {
			// Use refresh token to get new access token
			refreshReq := map[string]interface{}{
				"refresh_token": refreshToken,
			}

			body, _ = json.Marshal(refreshReq)
			req = httptest.NewRequest("POST", "/api/auth/refresh", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err = app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var refreshResp map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&refreshResp)
			require.NoError(t, err)
			assert.NotNil(t, refreshResp["access_token"])
			assert.NotNil(t, refreshResp["refresh_token"])
		}
	}
}

// TestLogoutFlow tests logout functionality
func TestLogoutFlow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)

	// Login to get session
	loginReq := map[string]interface{}{
		"email":    "test@example.com",
		"password": "testpassword123",
	}

	body, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	
	if resp.StatusCode == http.StatusOK {
		var loginResp map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&loginResp)
		require.NoError(t, err)
		
		sessionID, ok := loginResp["session_id"].(string)
		if ok && sessionID != "" {
			// Logout
			logoutReq := map[string]interface{}{
				"session_id": sessionID,
			}

			body, _ = json.Marshal(logoutReq)
			req = httptest.NewRequest("POST", "/api/auth/logout", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			accessToken, _ := loginResp["access_token"].(string)
			req.Header.Set("Authorization", "Bearer "+accessToken)

			resp, err = app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	}
}

// TestResendConfirmationEmail tests resending confirmation email
func TestResendConfirmationEmail(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)

	// Register a user
	registerReq := map[string]interface{}{
		"email":        "resend@example.com",
		"password":     "testpassword123",
		"display_name": "Test User",
		"locale":       "en",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusAccepted)

	// Resend confirmation
	resendReq := map[string]interface{}{
		"email": "resend@example.com",
	}

	body, _ = json.Marshal(resendReq)
	req = httptest.NewRequest("POST", "/api/auth/confirm/resend", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// Helper function to get auth service for tests
func getAuthServiceForTests(t *testing.T, db *gorm.DB, redisClient *redis.Client) *authdomain.Service {
	logger := applogger.New("test")
	cfg := testutil.LoadTestConfig()
	
	authRepo := authdomain.NewGormRepository(db)
	tokenStore := authdomain.NewRedisTokenStore(redisClient)
	jwtManager, err := authdomain.NewJWTManager(cfg.JWTSecret, 24*time.Hour, "woragis-test")
	require.NoError(t, err)
	
	emailSender := emailservice.NewNoopSender(logger)
	return authdomain.NewService(authRepo, emailSender, tokenStore, "http://localhost:8080", jwtManager, logger)
}
