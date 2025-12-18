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

	var testimonial map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&testimonial)
	require.NoError(t, err)
	assert.NotNil(t, testimonial["id"])
	testimonialID := testimonial["id"].(string)

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

	var experience map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&experience)
	require.NoError(t, err)
	assert.NotNil(t, experience["id"])
	experienceID := experience["id"].(string)

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

	var certification map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&certification)
	require.NoError(t, err)
	assert.NotNil(t, certification["id"])
	certificationID := certification["id"].(string)

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

	var caseStudy map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&caseStudy)
	require.NoError(t, err)
	assert.NotNil(t, caseStudy["id"])
	caseStudyID := caseStudy["id"].(string)

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
// Note: This test is skipped if social media posts routes are not set up in test app
func TestSocialMediaPostsAPI(t *testing.T) {
	t.Skip("Social Media Posts requires subdomain handlers setup - skipping for now")
	// TODO: Set up social media posts with all subdomain handlers for full testing
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
	var skill map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&skill)
	require.NoError(t, err)
	skillID := skill["id"].(string)

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
