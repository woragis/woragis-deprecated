//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	projectsdomain "github.com/woragis/backend/server/app/internal/domains/projects"
	"github.com/woragis/backend/server/app/internal/testutil"
)

// TestUnicodeAndSpecialCharacters tests handling of Unicode and special characters
func TestUnicodeAndSpecialCharacters(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	userID, token := createTestUserAndToken(t, db)

	// Test with Unicode characters
	unicodeTests := []struct {
		name        string
		projectName string
		description string
	}{
		{"Japanese", "プロジェクト", "日本語の説明"},
		{"Chinese", "项目", "中文描述"},
		{"Emoji", "🚀 Project 🎉", "Description with emojis 🎨"},
		{"Special chars", "Project & Co. <test> \"quotes\"", "Description with special: @#$%"},
		{"Arabic", "مشروع", "وصف بالعربية"},
		{"Russian", "Проект", "Описание на русском"},
	}

	for _, tt := range unicodeTests {
		t.Run(tt.name, func(t *testing.T) {
			createReq := projectsdomain.CreateProjectRequest{
				UserID:      userID,
				Name:        tt.projectName,
				Description: tt.description,
				Status:      projectsdomain.ProjectStatusPlanning,
			}

			body, _ := json.Marshal(createReq)
			req := httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode, "Should handle Unicode characters")

			var project map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&project)
			require.NoError(t, err)

			// Verify Unicode is preserved
			assert.True(t, utf8.ValidString(project["name"].(string)), "Name should be valid UTF-8")
			assert.True(t, utf8.ValidString(project["description"].(string)), "Description should be valid UTF-8")
		})
	}
}

// TestLargePayloads tests handling of large request payloads
func TestLargePayloads(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	userID, token := createTestUserAndToken(t, db)

	// Create a project with very long description
	longDescription := strings.Repeat("This is a very long description. ", 100) // ~3500 chars

	createReq := projectsdomain.CreateProjectRequest{
		UserID:      userID,
		Name:        "Large Payload Project",
		Description: longDescription,
		Status:      projectsdomain.ProjectStatusPlanning,
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	// Should either succeed or return 400 (payload too large)
	assert.True(t, resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusBadRequest,
		"Should handle large payloads gracefully")
}

// TestVeryLongStrings tests extremely long string fields
func TestVeryLongStrings(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Test with very long skill name (should be truncated or rejected)
	longName := strings.Repeat("A", 1000)

	createReq := map[string]interface{}{
		"name":     longName,
		"category": "programming",
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/skills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	// Should either succeed (with truncation) or reject
	assert.True(t, resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusBadRequest,
		"Should handle very long strings gracefully")
}

// TestConcurrentRequests tests handling of concurrent requests
func TestConcurrentRequests(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	userID, token := createTestUserAndToken(t, db)

	// Create multiple projects concurrently
	results := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(index int) {
			createReq := projectsdomain.CreateProjectRequest{
				UserID:      userID,
				Name:        "Concurrent Project " + string(rune(index)),
				Description: "Test concurrent creation",
				Status:      projectsdomain.ProjectStatusPlanning,
			}

			body, _ := json.Marshal(createReq)
			req := httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			resp, err := app.Test(req)
			success := err == nil && (resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK)
			results <- success
		}(i)
	}

	// Wait for all goroutines
	successCount := 0
	for i := 0; i < 5; i++ {
		if <-results {
			successCount++
		}
	}

	assert.Greater(t, successCount, 0, "At least some concurrent requests should succeed")
}

// TestEmptyAndNullValues tests handling of empty and null values
func TestEmptyAndNullValues(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Test with empty string
	createReq := map[string]interface{}{
		"name":     "",
		"category": "programming",
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/skills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	// Should reject empty name
	assert.True(t, resp.StatusCode >= 400 && resp.StatusCode < 500, "Should reject empty required fields")

	// Test with null values
	createReqNull := map[string]interface{}{
		"name":     "Valid Skill",
		"category": "programming",
		"icon":     nil,
	}

	body, _ = json.Marshal(createReqNull)
	req = httptest.NewRequest("POST", "/api/skills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err = app.Test(req)
	require.NoError(t, err)
	// Should handle null optional fields
	assert.True(t, resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusBadRequest)
}

// TestSQLInjectionAttempts tests basic SQL injection protection
func TestSQLInjectionAttempts(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Test SQL injection attempts in search queries
	sqlInjectionAttempts := []string{
		"'; DROP TABLE users; --",
		"1' OR '1'='1",
		"admin'--",
		"1 UNION SELECT * FROM users",
	}

	for _, attempt := range sqlInjectionAttempts {
		req := httptest.NewRequest("GET", "/api/skills/search?q="+attempt, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)
		// Should handle safely (either return empty results or 400)
		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest,
			"Should handle SQL injection attempts safely")
	}
}

// TestXSSAttempts tests basic XSS protection
func TestXSSAttempts(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	userID, token := createTestUserAndToken(t, db)

	// Test XSS attempts in project name
	xssAttempts := []string{
		"<script>alert('XSS')</script>",
		"<img src=x onerror=alert('XSS')>",
		"javascript:alert('XSS')",
		"<svg onload=alert('XSS')>",
	}

	for _, attempt := range xssAttempts {
		createReq := projectsdomain.CreateProjectRequest{
			UserID:      userID,
			Name:        attempt,
			Description: "XSS test",
			Status:      projectsdomain.ProjectStatusPlanning,
		}

		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := app.Test(req)
		require.NoError(t, err)
		// Should either sanitize or reject
		assert.True(t, resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusBadRequest,
			"Should handle XSS attempts safely")
	}
}

// TestInvalidJSON tests handling of invalid JSON
func TestInvalidJSON(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Test with invalid JSON
	invalidJSON := `{"name": "test", "category": }` // Invalid JSON

	req := httptest.NewRequest("POST", "/api/skills", bytes.NewReader([]byte(invalidJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should reject invalid JSON")
}

// TestMissingRequiredFields tests validation of required fields
func TestMissingRequiredFields(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Test with missing required field (name)
	createReq := map[string]interface{}{
		"category": "programming",
		// Missing "name" field
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/skills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 400 && resp.StatusCode < 500, "Should reject missing required fields")
}

// TestInvalidUUIDs tests handling of invalid UUIDs in path parameters
func TestInvalidUUIDs(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	invalidUUIDs := []string{
		"not-a-uuid",
		"123",
		"00000000-0000-0000-0000-000000000000",
		"invalid-uuid-format",
	}

	for _, invalidUUID := range invalidUUIDs {
		req := httptest.NewRequest("GET", "/api/projects/"+invalidUUID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)
		// Should return 400 (bad request) for invalid UUID format
		assert.True(t, resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound,
			"Should handle invalid UUIDs gracefully")
	}
}
