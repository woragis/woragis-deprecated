package skills

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Register_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/register", handler.Register)

	req := RegisterRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "SecurePass123!",
		FirstName: "Test",
		LastName:  "User",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	authResp := &AuthResponse{
		User: &User{
			ID:       uuid.New(),
			Username: req.Username,
			Email:    req.Email,
		},
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    1234567890,
	}

	mockService.On("register", &req).Return(authResp, nil)

	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestHandler_Register_InvalidRequest(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/register", handler.Register)

	// Invalid JSON
	httpReq := httptest.NewRequest("POST", "/register", bytes.NewBufferString("invalid json"))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestHandler_Login_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/login", handler.Login)

	req := LoginRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("User-Agent", "test-agent")

	authResp := &AuthResponse{
		User: &User{
			ID:    uuid.New(),
			Email: req.Email,
		},
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresAt:    1234567890,
	}

	mockService.On("login", &req, "test-agent", mock.AnythingOfType("string")).Return(authResp, nil)

	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/login", handler.Login)

	req := LoginRequest{
		Email:    "test@example.com",
		Password: "WrongPassword",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	mockService.On("login", &req, "", mock.AnythingOfType("string")).Return(nil, ErrInvalidCredentials)

	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestHandler_RefreshToken_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/refresh", handler.RefreshToken)

	req := map[string]string{
		"refresh_token": "test-refresh-token",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/refresh", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	authResp := &AuthResponse{
		User: &User{
			ID: uuid.New(),
		},
		AccessToken:  "new-access-token",
		RefreshToken: "new-refresh-token",
		ExpiresAt:    1234567890,
	}

	mockService.On("refreshToken", "test-refresh-token").Return(authResp, nil)

	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestHandler_Logout_Success(t *testing.T) {
	mockService := new(MockService)
	handler := NewHandler(mockService)

	app := fiber.New()
	app.Post("/logout", handler.Logout)

	req := map[string]string{
		"refresh_token": "test-refresh-token",
	}

	reqBody, _ := json.Marshal(req)
	httpReq := httptest.NewRequest("POST", "/logout", bytes.NewBuffer(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	mockService.On("logout", "test-refresh-token").Return(nil)

	resp, err := app.Test(httpReq)

	require.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

// MockService is a mock implementation of Service
type MockService struct {
	mock.Mock
}

func (m *MockService) register(req *RegisterRequest) (*AuthResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AuthResponse), args.Error(1)
}

func (m *MockService) login(req *LoginRequest, userAgent, ipAddress string) (*AuthResponse, error) {
	args := m.Called(req, userAgent, ipAddress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AuthResponse), args.Error(1)
}

func (m *MockService) refreshToken(refreshToken string) (*AuthResponse, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*AuthResponse), args.Error(1)
}

func (m *MockService) logout(refreshToken string) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

func (m *MockService) logoutAll(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockService) getUserProfile(userID uuid.UUID) (*Profile, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Profile), args.Error(1)
}

func (m *MockService) updateUserProfile(userID uuid.UUID, req *ProfileUpdateRequest) (*Profile, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Profile), args.Error(1)
}

func (m *MockService) changePassword(userID uuid.UUID, req *PasswordChangeRequest) error {
	args := m.Called(userID, req)
	return args.Error(0)
}

func (m *MockService) createVerificationToken(userID uuid.UUID, tokenType string) (*VerificationToken, error) {
	args := m.Called(userID, tokenType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*VerificationToken), args.Error(1)
}

func (m *MockService) verifyEmail(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockService) getUserByID(userID uuid.UUID) (*User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockService) listUsers(page, limit int, search string) ([]User, int64, error) {
	args := m.Called(page, limit, search)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]User), args.Get(1).(int64), args.Error(2)
}

func (m *MockService) cleanupExpiredSessions() error {
	args := m.Called()
	return args.Error(0)
}

