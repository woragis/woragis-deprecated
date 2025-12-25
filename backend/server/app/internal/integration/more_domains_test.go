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

// TestTestimonialsAPI tests the Testimonials API endpoints
func TestTestimonialsAPI(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create a testimonial
	createReq := map[string]interface{}{
		"author_name":  "John Doe",
		"author_role":  "CEO",
		"author_company": "Acme Corp",
		"content":      "Great work!",
		"rating":       5,
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/testimonials", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var testimonialResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&testimonialResponse)
	require.NoError(t, err)
	
	// Response is wrapped in "data" field
	testimonialData, ok := testimonialResponse["data"].(map[string]interface{})
	require.True(t, ok, "response should have data field")
	assert.NotNil(t, testimonialData["id"])
	testimonialID := testimonialData["id"].(string)

	// Get testimonial by ID
	req = httptest.NewRequest("GET", "/api/testimonials/"+testimonialID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// List testimonials
	req = httptest.NewRequest("GET", "/api/testimonials", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Update testimonial
	updateReq := map[string]interface{}{
		"content": "Updated testimonial content",
		"rating":  4,
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest("PATCH", "/api/testimonials/"+testimonialID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete testimonial
	req = httptest.NewRequest("DELETE", "/api/testimonials/"+testimonialID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestExperiencesAPI tests the Experiences API endpoints
func TestExperiencesAPI(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create an experience
	createReq := map[string]interface{}{
		"title":       "Software Engineer",
		"company":     "Tech Corp",
		"description": "Worked on backend systems",
		"start_date":  "2020-01-01",
		"end_date":    "2022-12-31",
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/experiences", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var experienceResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&experienceResponse)
	require.NoError(t, err)
	
	// Response is wrapped in "data" field
	experienceData, ok := experienceResponse["data"].(map[string]interface{})
	require.True(t, ok, "response should have data field")
	assert.NotNil(t, experienceData["id"])
	experienceID := experienceData["id"].(string)

	// Get experience by ID
	req = httptest.NewRequest("GET", "/api/experiences/"+experienceID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// List experiences
	req = httptest.NewRequest("GET", "/api/experiences", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Update experience
	updateReq := map[string]interface{}{
		"title": "Senior Software Engineer",
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest("PATCH", "/api/experiences/"+experienceID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete experience
	req = httptest.NewRequest("DELETE", "/api/experiences/"+experienceID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestCertificationsAPI tests the Certifications API endpoints
func TestCertificationsAPI(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create a certification
	createReq := map[string]interface{}{
		"name":         "AWS Certified Solutions Architect",
		"issuer":       "Amazon Web Services",
		"issue_date":   "2023-01-15",
		"expiry_date":  "2026-01-15",
		"credential_id": "AWS-12345",
		"featured":     true,
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/certifications", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var certificationResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&certificationResponse)
	require.NoError(t, err)
	
	// Response is wrapped in "data" field
	certificationData, ok := certificationResponse["data"].(map[string]interface{})
	require.True(t, ok, "response should have data field")
	assert.NotNil(t, certificationData["id"])
	certificationID := certificationData["id"].(string)

	// Get certification by ID
	req = httptest.NewRequest("GET", "/api/certifications/"+certificationID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// List certifications
	req = httptest.NewRequest("GET", "/api/certifications", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// List featured certifications
	req = httptest.NewRequest("GET", "/api/certifications/featured", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Update certification
	updateReq := map[string]interface{}{
		"featured": false,
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest("PATCH", "/api/certifications/"+certificationID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete certification
	req = httptest.NewRequest("DELETE", "/api/certifications/"+certificationID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestCaseStudiesAPI tests the Case Studies API endpoints
func TestCaseStudiesAPI(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create a case study
	createReq := map[string]interface{}{
		"title":       "E-commerce Platform Redesign",
		"description": "Complete redesign of e-commerce platform",
		"challenge":   "Legacy system performance issues",
		"solution":    "Modern microservices architecture",
		"results":     "50% performance improvement",
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/case-studies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var caseStudyResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&caseStudyResponse)
	require.NoError(t, err)
	
	// Response is wrapped in "data" field
	caseStudyData, ok := caseStudyResponse["data"].(map[string]interface{})
	require.True(t, ok, "response should have data field")
	assert.NotNil(t, caseStudyData["id"])
	caseStudyID := caseStudyData["id"].(string)

	// Get case study by ID
	req = httptest.NewRequest("GET", "/api/case-studies/"+caseStudyID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// List case studies
	req = httptest.NewRequest("GET", "/api/case-studies", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Update case study
	updateReq := map[string]interface{}{
		"title": "Updated Case Study Title",
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest("PATCH", "/api/case-studies/"+caseStudyID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete case study
	req = httptest.NewRequest("DELETE", "/api/case-studies/"+caseStudyID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestSocialMediaPostsAPI tests the Social Media Posts API endpoints
func TestSocialMediaPostsAPI(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create a social media post
	createReq := map[string]interface{}{
		"platform": "linkedin",
		"format":   "post",
		"title":    "Test Social Media Post",
		"content":  "This is a test social media post content",
	}

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/social-media-posts", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var postResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&postResponse)
	require.NoError(t, err)
	
	// Response is wrapped in "data" field
	postData, ok := postResponse["data"].(map[string]interface{})
	require.True(t, ok, "response should have data field")
	assert.NotNil(t, postData["id"])
	assert.Equal(t, "linkedin", postData["platform"])
	assert.Equal(t, "post", postData["format"])
	postID := postData["id"].(string)

	// Get post by ID
	req = httptest.NewRequest("GET", "/api/social-media-posts/"+postID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// List posts
	req = httptest.NewRequest("GET", "/api/social-media-posts", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Update post
	updateReq := map[string]interface{}{
		"title": "Updated Social Media Post Title",
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest("PATCH", "/api/social-media-posts/"+postID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Update post status
	statusReq := map[string]interface{}{
		"status": "ready",
	}
	body, _ = json.Marshal(statusReq)
	req = httptest.NewRequest("PATCH", "/api/social-media-posts/"+postID+"/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Delete post
	req = httptest.NewRequest("DELETE", "/api/social-media-posts/"+postID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestCertificationsRelationships tests certification relationships (skills, entities)
func TestCertificationsRelationships(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	redis := testutil.SetupTestRedis(t)
	defer testutil.CleanupTestRedis(t, redis)

	app := setupTestAppWithRoutes(t, db, redis)
	_, token := createTestUserAndToken(t, db)

	// Create a skill first
	skillReq := map[string]interface{}{
		"name":     "AWS",
		"category": "cloud",
	}
	body, _ := json.Marshal(skillReq)
	req := httptest.NewRequest("POST", "/api/skills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	var skillResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&skillResponse)
	require.NoError(t, err)
	
	// Response is wrapped in "data" field
	skillData, ok := skillResponse["data"].(map[string]interface{})
	require.True(t, ok, "response should have data field")
	skillID := skillData["id"].(string)

	// Create a certification
	certReq := map[string]interface{}{
		"name":   "AWS Certified",
		"issuer": "AWS",
	}
	body, _ = json.Marshal(certReq)
	req = httptest.NewRequest("POST", "/api/certifications", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	var certification map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&certification)
	require.NoError(t, err)
	certificationID := certification["id"].(string)

	// Add skill to certification
	req = httptest.NewRequest("POST", "/api/certifications/"+certificationID+"/skills/"+skillID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated)

	// Get certification skills
	req = httptest.NewRequest("GET", "/api/certifications/"+certificationID+"/skills", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
