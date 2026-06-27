package testhelpers

import (
	"time"

	"github.com/google/uuid"

	domains "woragis-files-service/internal/domains"
)

// CreateTestUser creates a test user entity
func CreateTestUser() *domains.User {
	now := time.Now()
	return &domains.User{
		ID:         uuid.New(),
		Username:   "testuser",
		Email:      "test@example.com",
		Password:   "$2a$10$abcdefghijklmnopqrstuvwxyz1234567890", // bcrypt hash
		Role:       "user",
		IsActive:   true,
		IsVerified: false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// CreateTestUserWithID creates a test user with specific ID
func CreateTestUserWithID(id uuid.UUID) *domains.User {
	user := CreateTestUser()
	user.ID = id
	return user
}

// CreateTestProfile creates a test profile
func CreateTestProfile(userID uuid.UUID) *domains.Profile {
	now := time.Now()
	dob := now.AddDate(-25, 0, 0)
	return &domains.Profile{
		UserID:      userID,
		Avatar:      "https://example.com/avatar.jpg",
		Bio:         "Test bio",
		DateOfBirth: &dob,
		Gender:      "male",
		Phone:       "+1234567890",
		Location:    "Test City",
		Website:     "https://example.com",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// CreateTestSession creates a test session
func CreateTestSession(userID uuid.UUID) *domains.Session {
	now := time.Now()
	return &domains.Session{
		ID:           uuid.New(),
		UserID:       userID,
		RefreshToken: "test-refresh-token-" + uuid.New().String(),
		UserAgent:    "test-agent",
		IPAddress:    "127.0.0.1",
		IsActive:     true,
		ExpiresAt:    now.Add(24 * time.Hour),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// CreateTestRegisterRequest creates a test register request
func CreateTestRegisterRequest() *domains.RegisterRequest {
	return &domains.RegisterRequest{
		Username:  "newuser",
		Email:     "newuser@example.com",
		Password:  "SecurePass123!",
		FirstName: "New",
		LastName:  "User",
	}
}

// CreateTestLoginRequest creates a test login request
func CreateTestLoginRequest() *domains.LoginRequest {
	return &domains.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
}
