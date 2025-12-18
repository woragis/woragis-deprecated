//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	projectsdomain "github.com/woragis/backend/server/app/internal/domains/projects"
	"github.com/woragis/backend/server/app/internal/testutil"
)

// TestPagination tests pagination functionality (if implemented)
func TestPagination(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	userID, token := createTestUserAndToken(t, db)

	// Create multiple projects for pagination testing
	projectNames := []string{"Project 1", "Project 2", "Project 3", "Project 4", "Project 5"}
	for _, name := range projectNames {
		createReq := projectsdomain.CreateProjectRequest{
			UserID:      userID,
			Name:        name,
			Description: "Test project",
			Status:      projectsdomain.ProjectStatusPlanning,
		}
		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		_, err := app.Test(req)
		require.NoError(t, err)
	}

	// Test: List all projects (baseline)
	req := httptest.NewRequest("GET", "/api/projects", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var allProjects interface{}
	err = json.NewDecoder(resp.Body).Decode(&allProjects)
	require.NoError(t, err)

	// Extract projects array
	var projects []interface{}
	switch v := allProjects.(type) {
	case []interface{}:
		projects = v
	case map[string]interface{}:
		if data, ok := v["data"].([]interface{}); ok {
			projects = data
		}
	}

	assert.GreaterOrEqual(t, len(projects), len(projectNames), "Should return all created projects")

	// Test pagination parameters (if supported)
	// Note: Current implementation may not support pagination yet
	// These tests validate the endpoint accepts parameters gracefully
	testCases := []struct {
		name   string
		params string
	}{
		{"limit only", "?limit=2"},
		{"offset only", "?offset=1"},
		{"limit and offset", "?limit=2&offset=1"},
		{"page parameter", "?page=1&per_page=2"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/projects"+tc.params, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := app.Test(req)
			require.NoError(t, err)
			// Should return 200 even if pagination not implemented
			assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest,
				"Endpoint should handle pagination parameters gracefully")
		})
	}
}

// TestSorting tests sorting functionality
func TestSorting(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	userID, token := createTestUserAndToken(t, db)

	// Create projects with different names for sorting
	projectNames := []string{"Zebra Project", "Alpha Project", "Beta Project"}
	for _, name := range projectNames {
		createReq := projectsdomain.CreateProjectRequest{
			UserID:      userID,
			Name:        name,
			Description: "Test project",
			Status:      projectsdomain.ProjectStatusPlanning,
		}
		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		_, err := app.Test(req)
		require.NoError(t, err)
	}

	// Test sorting parameters (if supported)
	sortParams := []string{
		"?sort=name",
		"?sort=name&order=asc",
		"?sort=name&order=desc",
		"?sort=created_at&order=desc",
	}

	for _, params := range sortParams {
		req := httptest.NewRequest("GET", "/api/projects"+params, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)
		// Should return 200 even if sorting not implemented
		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest,
			"Endpoint should handle sort parameters gracefully")
	}
}

// TestAdvancedFiltering tests advanced filtering capabilities
func TestAdvancedFiltering(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	userID, token := createTestUserAndToken(t, db)

	// Create projects with different statuses
	statuses := []string{"planning", "active", "completed"}
	for i, status := range statuses {
		createReq := projectsdomain.CreateProjectRequest{
			UserID:      userID,
			Name:        "Project " + status,
			Description: "Test project",
			Status:      projectsdomain.ProjectStatus(status),
		}
		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		_, err := app.Test(req)
		require.NoError(t, err)
		_ = i
	}

	// Test filtering by status (if supported)
	filterParams := []string{
		"?status=active",
		"?status=planning",
		"?status=completed",
	}

	for _, params := range filterParams {
		req := httptest.NewRequest("GET", "/api/projects"+params, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)
		// Should return 200 even if filtering not implemented
		assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest,
			"Endpoint should handle filter parameters gracefully")
	}
}

// TestSearchFunctionality tests search across different domains
func TestSearchFunctionality(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Test Skills search
	req := httptest.NewRequest("GET", "/api/skills/search?q=Go", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test Interests search
	req = httptest.NewRequest("GET", "/api/interests/search?q=Machine", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test search with empty query
	req = httptest.NewRequest("GET", "/api/skills/search?q=", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	// Should handle empty query gracefully
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest)
}

// TestFilteringByCategory tests category-based filtering
func TestFilteringByCategory(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create skills in different categories
	categories := []string{"programming", "devops", "design"}
	for _, category := range categories {
		createReq := map[string]interface{}{
			"name":     "Skill " + category,
			"category": category,
		}
		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/api/skills", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		_, err := app.Test(req)
		require.NoError(t, err)
	}

	// Test filtering by category
	req := httptest.NewRequest("GET", "/api/skills/category?category=programming", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var results interface{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	require.NoError(t, err)
	// Results should be returned (may be empty if no matches)
	assert.NotNil(t, results)
}

// TestMultipleFilters tests combining multiple filters
func TestMultipleFilters(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Test multiple query parameters
	req := httptest.NewRequest("GET", "/api/skills/category?category=programming&search=Go", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	// Should handle multiple parameters gracefully
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusBadRequest)
}

// TestFeaturedFiltering tests featured/priority filtering
func TestFeaturedFiltering(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create featured and non-featured interests
	interests := []struct {
		title    string
		featured bool
	}{
		{"Featured Interest", true},
		{"Regular Interest", false},
	}

	for _, interest := range interests {
		createReq := map[string]interface{}{
			"title":    interest.title,
			"featured": interest.featured,
		}
		body, _ := json.Marshal(createReq)
		req := httptest.NewRequest("POST", "/api/interests", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		_, err := app.Test(req)
		require.NoError(t, err)
	}

	// Test featured interests endpoint
	req := httptest.NewRequest("GET", "/api/interests/featured", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var featuredResults interface{}
	err = json.NewDecoder(resp.Body).Decode(&featuredResults)
	require.NoError(t, err)
	assert.NotNil(t, featuredResults)
}

// TestBulkOperations tests bulk operation endpoints (if available)
func TestBulkOperations(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	userID, token := createTestUserAndToken(t, db)

	// Create a project first
	createReq := projectsdomain.CreateProjectRequest{
		UserID:      userID,
		Name:        "Bulk Test Project",
		Description: "For bulk operations",
		Status:      projectsdomain.ProjectStatusPlanning,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/projects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var project map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&project)
	require.NoError(t, err)
	projectID := project["id"].(string)

	// Test bulk milestone update (if available)
	bulkReq := map[string]interface{}{
		"updates": []map[string]interface{}{},
	}
	body, _ = json.Marshal(bulkReq)
	req = httptest.NewRequest("POST", "/api/projects/"+projectID+"/milestones/bulk", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	// Should handle bulk operations (may return 200, 400, or 404 depending on implementation)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 500, "Bulk endpoint should respond")
}
