//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"woragis-files-service/internal/config"
	"woragis-files-service/internal/domains"
	testhelpers "woragis-files-service/internal/testing"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
)

func setupIntegrationTest(t *testing.T) (*gorm.DB, *redis.Client, *domains.Repository, *domains.Service, *domains.Handler, *fiber.App) {
	cfg := testhelpers.LoadTestConfig()

	// Setup database
	db, err := testhelpers.SetupTestDB(cfg.DatabaseURL)
	require.NoError(t, err)

	// Setup Redis
	redisClient, err := testhelpers.SetupTestRedis(cfg.RedisURL)
	require.NoError(t, err)

	// Run migrations
	require.NoError(t, domains.MigrateAuthTables(db))

	// Setup JWT manager
	jwtManager := testhelpers.SetupTestJWTManager(cfg.JWTSecret, redisClient)

	// Setup repository, service, handler
	repo := domains.NewRepository(db)
	service := domains.NewService(repo, jwtManager, 10)
	handler := domains.NewHandler(service)

	// Setup Fiber app
	appConfig := &config.Config{
		AppName: "woragis-files-service-test",
		Env:     "test",
		Port:    "8080",
	}
	app := config.CreateFiberApp(appConfig)

	// Setup routes
	api := app.Group("/api/v1")
	domains.SetupRoutes(api, handler, jwtManager)

	return db, redisClient, repo, service, handler, app
}

func teardownIntegrationTest(t *testing.T, db *gorm.DB, redisClient *redis.Client) {
	if db != nil {
		testhelpers.CleanupTestDB(db)
	}
	if redisClient != nil {
		testhelpers.CleanupTestRedis(redisClient)
	}
}

func TestIntegration_RegisterAndLogin(t *testing.T) {
	db, redis, _, _, _, app := setupIntegrationTest(t)
	defer teardownIntegrationTest(t, db, redis)

	// Register user
	registerReq := domains.RegisterRequest{
		Username:  "testuser" + uuid.New().String()[:8],
		Email:     "test" + uuid.New().String()[:8] + "@example.com",
		Password:  "SecurePass123!",
		FirstName: "Test",
		LastName:  "User",
	}

	registerBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var registerResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&registerResp)
	assert.NotEmpty(t, registerResp["data"].(map[string]interface{})["access_token"])

	// Login with registered user
	loginReq := domains.LoginRequest{
		Email:    registerReq.Email,
		Password: registerReq.Password,
	}

	loginBody, _ := json.Marshal(loginReq)
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&loginResp)
	assert.NotEmpty(t, loginResp["data"].(map[string]interface{})["access_token"])
}

func TestIntegration_RefreshToken(t *testing.T) {
	db, redis, _, service, _, app := setupIntegrationTest(t)
	defer teardownIntegrationTest(t, db, redis)

	// Register and login
	registerReq := domains.RegisterRequest{
		Username:  "testuser" + uuid.New().String()[:8],
		Email:     "test" + uuid.New().String()[:8] + "@example.com",
		Password:  "SecurePass123!",
		FirstName: "Test",
		LastName:  "User",
	}

	authResp, err := service.register(&registerReq)
	require.NoError(t, err)

	// Refresh token
	refreshReq := map[string]string{
		"refresh_token": authResp.RefreshToken,
	}

	refreshBody, _ := json.Marshal(refreshReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(refreshBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var refreshResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&refreshResp)
	assert.NotEmpty(t, refreshResp["data"].(map[string]interface{})["access_token"])
}

func TestIntegration_Logout(t *testing.T) {
	db, redis, _, service, _, app := setupIntegrationTest(t)
	defer teardownIntegrationTest(t, db, redis)

	// Register and login
	registerReq := domains.RegisterRequest{
		Username:  "testuser" + uuid.New().String()[:8],
		Email:     "test" + uuid.New().String()[:8] + "@example.com",
		Password:  "SecurePass123!",
		FirstName: "Test",
		LastName:  "User",
	}

	authResp, err := service.register(&registerReq)
	require.NoError(t, err)

	// Logout
	logoutReq := map[string]string{
		"refresh_token": authResp.RefreshToken,
	}

	logoutBody, _ := json.Marshal(logoutReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/logout", bytes.NewBuffer(logoutBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

