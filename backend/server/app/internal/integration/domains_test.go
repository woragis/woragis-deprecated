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

	"github.com/woragis/backend/server/app/internal/testutil"
)

// TestSkillsAPI tests the Skills API endpoints
func TestSkillsAPI(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create a skill
	createReq := map[string]interface{}{
		"name":        "Go",
		"description": "Go programming language",
		"category":    "programming",
		"icon":        "go-icon",
		"color":       "#00ADD8",
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/skills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var skill map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&skill)
	require.NoError(t, err)
	assert.NotNil(t, skill["id"])
	assert.Equal(t, "Go", skill["name"])
	skillID := skill["id"].(string)

	// Get skill by ID
	req = httptest.NewRequest("GET", "/api/skills/"+skillID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// List skills
	req = httptest.NewRequest("GET", "/api/skills", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Search skills
	req = httptest.NewRequest("GET", "/api/skills/search?q=Go", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// List skills by category
	req = httptest.NewRequest("GET", "/api/skills/category?category=programming", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Update skill
	updateReq := map[string]interface{}{
		"description": "Updated description",
		"proficiencyLevel": "intermediate",
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest("PATCH", "/api/skills/"+skillID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete skill
	req = httptest.NewRequest("DELETE", "/api/skills/"+skillID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify deletion
	req = httptest.NewRequest("GET", "/api/skills/"+skillID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// TestPostsAPI tests the Posts API endpoints
func TestPostsAPI(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create a post
	createReq := map[string]interface{}{
		"title":   "Test Post",
		"content": "This is a test post content",
		"excerpt": "Test excerpt",
		"status":  "draft",
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var post map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&post)
	require.NoError(t, err)
	assert.NotNil(t, post["id"])
	assert.Equal(t, "Test Post", post["title"])
	postID := post["id"].(string)

	// Get post by ID
	req = httptest.NewRequest("GET", "/api/posts/"+postID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// List posts
	req = httptest.NewRequest("GET", "/api/posts", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Update post
	updateReq := map[string]interface{}{
		"title":   "Updated Post Title",
		"status":  "published",
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest("PATCH", "/api/posts/"+postID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete post
	req = httptest.NewRequest("DELETE", "/api/posts/"+postID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestInterestsAPI tests the Interests API endpoints
func TestInterestsAPI(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create an interest
	createReq := map[string]interface{}{
		"title":       "Machine Learning",
		"description": "Interest in ML and AI",
		"featured":    true,
		"fullWidth":   false,
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/interests", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var interest map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&interest)
	require.NoError(t, err)
	assert.NotNil(t, interest["id"])
	assert.Equal(t, "Machine Learning", interest["title"])
	interestID := interest["id"].(string)

	// Get interest by ID
	req = httptest.NewRequest("GET", "/api/interests/"+interestID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// List interests
	req = httptest.NewRequest("GET", "/api/interests", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// List featured interests
	req = httptest.NewRequest("GET", "/api/interests/featured", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Search interests
	req = httptest.NewRequest("GET", "/api/interests/search?q=Machine", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Update interest
	updateReq := map[string]interface{}{
		"description": "Updated description",
		"featured":    false,
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest("PATCH", "/api/interests/"+interestID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete interest
	req = httptest.NewRequest("DELETE", "/api/interests/"+interestID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestSkillsSearchAndFiltering tests advanced search and filtering features
func TestSkillsSearchAndFiltering(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create multiple skills with different categories
	skills := []map[string]interface{}{
		{"name": "Go", "category": "programming", "description": "Go language"},
		{"name": "Python", "category": "programming", "description": "Python language"},
		{"name": "Docker", "category": "devops", "description": "Containerization"},
		{"name": "Kubernetes", "category": "devops", "description": "Orchestration"},
	}

	skillIDs := make([]string, 0)
	for _, skillData := range skills {
		body, _ := json.Marshal(skillData)
		req := httptest.NewRequest("POST", "/api/skills", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var skill map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&skill)
		require.NoError(t, err)
		skillIDs = append(skillIDs, skill["id"].(string))
	}

	// Test search functionality
	req := httptest.NewRequest("GET", "/api/skills/search?q=Go", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var searchResults interface{}
	err = json.NewDecoder(resp.Body).Decode(&searchResults)
	require.NoError(t, err)

	// Test filtering by category
	req = httptest.NewRequest("GET", "/api/skills/category?category=devops", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var categoryResults interface{}
	err = json.NewDecoder(resp.Body).Decode(&categoryResults)
	require.NoError(t, err)

	// Test get all skills with project counts
	req = httptest.NewRequest("GET", "/api/skills/with-counts", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestPostsRelationships tests post relationships (skills, categories, tags)
func TestPostsRelationships(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create a skill first
	skillReq := map[string]interface{}{
		"name":     "Go",
		"category": "programming",
	}
	body, _ := json.Marshal(skillReq)
	req := httptest.NewRequest("POST", "/api/skills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	var skill map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&skill)
	require.NoError(t, err)
	skillID := skill["id"].(string)

	// Create a category
	categoryReq := map[string]interface{}{
		"name": "Technology",
	}
	body, _ = json.Marshal(categoryReq)
	req = httptest.NewRequest("POST", "/api/posts/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	var category map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&category)
	require.NoError(t, err)
	categoryID := category["id"].(string)

	// Create a post with relationships
	postReq := map[string]interface{}{
		"title":       "Test Post with Relationships",
		"content":     "Content",
		"skillIds":    []string{skillID},
		"categoryIds": []string{categoryID},
		"tagNames":    []string{"test", "go"},
	}
	body, _ = json.Marshal(postReq)
	req = httptest.NewRequest("POST", "/api/posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var post map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&post)
	require.NoError(t, err)
	postID := post["id"].(string)

	// Get post skills
	req = httptest.NewRequest("GET", "/api/posts/"+postID+"/skills", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get post categories
	req = httptest.NewRequest("GET", "/api/posts/"+postID+"/categories", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get post tags
	req = httptest.NewRequest("GET", "/api/posts/"+postID+"/tags", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
