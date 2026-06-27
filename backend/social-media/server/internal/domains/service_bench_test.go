package socialmedia

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	authPkg "woragis-social-media-service/pkg/auth"
)

func BenchmarkService_Register(b *testing.B) {
	mockRepo := new(MockRepository)
	jwtManager := authPkg.NewJWTManager("test-secret", "test-issuer", 1*time.Hour, 24*time.Hour)
	service := NewService(mockRepo, jwtManager, 10)

	req := &RegisterRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "SecurePass123!",
		FirstName: "Test",
		LastName:  "User",
	}

	mockRepo.On("getUserByEmail", req.Email).Return(nil, ErrUserNotFound)
	mockRepo.On("createUser", mock.AnythingOfType("*auth.User")).Return(nil)
	mockRepo.On("createProfile", mock.AnythingOfType("*auth.Profile")).Return(nil)
	mockRepo.On("createSession", mock.AnythingOfType("*auth.Session")).Return(nil)
	mockRepo.On("updateLastLogin", mock.AnythingOfType("uuid.UUID")).Return(nil)
	mockRepo.On("getUserByID", mock.AnythingOfType("uuid.UUID")).Return(&User{
		ID:       uuid.New(),
		Email:    req.Email,
		Username: req.Username,
	}, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.register(req)
	}
}

func BenchmarkService_Login(b *testing.B) {
	mockRepo := new(MockRepository)
	jwtManager := authPkg.NewJWTManager("test-secret", "test-issuer", 1*time.Hour, 24*time.Hour)
	service := NewService(mockRepo, jwtManager, 10)

	hashedPassword, _ := authPkg.HashPassword("SecurePass123!", 10)
	user := &User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Password:  hashedPassword,
		IsActive:  true,
		IsVerified: true,
	}

	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	mockRepo.On("getUserByEmail", req.Email).Return(user, nil)
	mockRepo.On("updateLastLogin", user.ID).Return(nil)
	mockRepo.On("createSession", mock.AnythingOfType("*auth.Session")).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.login(req, "test-agent", "127.0.0.1")
	}
}

func BenchmarkService_RefreshToken(b *testing.B) {
	mockRepo := new(MockRepository)
	jwtManager := authPkg.NewJWTManager("test-secret", "test-issuer", 1*time.Hour, 24*time.Hour)
	service := NewService(mockRepo, jwtManager, 10)

	userID := uuid.New()
	_, refreshToken, _ := jwtManager.Generate(userID, "test@example.com", "user", "Test")

	user := &User{
		ID:       userID,
		IsActive: true,
	}

	session := &Session{
		ID:           uuid.New(),
		UserID:       userID,
		RefreshToken: refreshToken,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("getSessionByRefreshToken", refreshToken).Return(session, nil)
	mockRepo.On("getUserByID", userID).Return(user, nil)
	mockRepo.On("updateSession", mock.AnythingOfType("*auth.Session")).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.refreshToken(refreshToken)
	}
}

