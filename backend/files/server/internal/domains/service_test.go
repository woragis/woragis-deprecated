package files

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	authPkg "woragis-files-service/pkg/auth"
)

// MockRepository is a mock implementation of Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) createUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) getUserByID(id uuid.UUID) (*User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) getUserByEmail(email string) (*User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) getUserByUsername(username string) (*User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockRepository) updateUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) deleteUser(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) updateLastLogin(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) verifyUserEmail(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) listUsers(offset, limit int, search string) ([]User, int64, error) {
	args := m.Called(offset, limit, search)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]User), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) createProfile(profile *Profile) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *MockRepository) getProfileByUserID(userID uuid.UUID) (*Profile, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Profile), args.Error(1)
}

func (m *MockRepository) updateProfile(profile *Profile) error {
	args := m.Called(profile)
	return args.Error(0)
}

func (m *MockRepository) deleteProfile(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockRepository) createSession(session *Session) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockRepository) getSessionByRefreshToken(refreshToken string) (*Session, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Session), args.Error(1)
}

func (m *MockRepository) getSessionsByUserID(userID uuid.UUID) ([]Session, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Session), args.Error(1)
}

func (m *MockRepository) updateSession(session *Session) error {
	args := m.Called(session)
	return args.Error(0)
}

func (m *MockRepository) deactivateSession(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) deactivateAllUserSessions(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockRepository) deleteExpiredSessions() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository) createVerificationToken(token *VerificationToken) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockRepository) getVerificationToken(token string) (*VerificationToken, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*VerificationToken), args.Error(1)
}

func (m *MockRepository) markTokenAsUsed(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) deleteExpiredTokens() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository) deleteUsedTokens() error {
	args := m.Called()
	return args.Error(0)
}

func TestService_Register_Success(t *testing.T) {
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

	userID := uuid.New()
	createdUser := &User{
		ID:        userID,
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	// Mock repository calls
	mockRepo.On("getUserByEmail", req.Email).Return(nil, ErrUserNotFound)
	mockRepo.On("createUser", mock.AnythingOfType("*auth.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(0).(*User)
		user.ID = userID
		createdUser = user
	})
	mockRepo.On("createProfile", mock.AnythingOfType("*auth.Profile")).Return(nil)
	mockRepo.On("createSession", mock.AnythingOfType("*auth.Session")).Return(nil)
	mockRepo.On("updateLastLogin", userID).Return(nil)
	mockRepo.On("getUserByID", userID).Return(createdUser, nil)

	response, err := service.register(req)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.NotNil(t, response.User)
	assert.Equal(t, req.Email, response.User.Email)
	assert.Equal(t, req.Username, response.User.Username)

	mockRepo.AssertExpectations(t)
}

func TestService_Register_UserExists(t *testing.T) {
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

	existingUser := &User{
		ID:       uuid.New(),
		Email:    req.Email,
		Username: req.Username,
	}

	// Mock repository - user already exists
	mockRepo.On("getUserByEmail", req.Email).Return(existingUser, nil)

	response, err := service.register(req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, ErrUserAlreadyExists, err)

	mockRepo.AssertExpectations(t)
}

func TestService_Register_WeakPassword(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtManager := authPkg.NewJWTManager("test-secret", "test-issuer", 1*time.Hour, 24*time.Hour)
	service := NewService(mockRepo, jwtManager, 10)

	req := &RegisterRequest{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "123", // Too weak
		FirstName: "Test",
		LastName:  "User",
	}

	response, err := service.register(req)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, ErrPasswordTooWeak, err)
}

func TestService_Login_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtManager := authPkg.NewJWTManager("test-secret", "test-issuer", 1*time.Hour, 24*time.Hour)
	service := NewService(mockRepo, jwtManager, 10)

	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	// Create user with hashed password
	user := &User{
		ID:        uuid.New(),
		Email:     req.Email,
		Username:  "testuser",
		IsActive:  true,
		IsVerified: true,
	}

	// Hash the password properly
	hashedPassword, err := authPkg.HashPassword("SecurePass123!", 10)
	require.NoError(t, err)
	user.Password = hashedPassword

	// Mock repository calls
	mockRepo.On("getUserByEmail", req.Email).Return(user, nil)
	mockRepo.On("updateLastLogin", user.ID).Return(nil)
	mockRepo.On("createSession", mock.AnythingOfType("*auth.Session")).Return(nil)

	response, err := service.login(req, "test-agent", "127.0.0.1")

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.NotNil(t, response.User)

	mockRepo.AssertExpectations(t)
}

func TestService_Login_InvalidCredentials(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtManager := authPkg.NewJWTManager("test-secret", "test-issuer", 1*time.Hour, 24*time.Hour)
	service := NewService(mockRepo, jwtManager, 10)

	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "WrongPassword",
	}

	user := &User{
		ID:        uuid.New(),
		Email:     req.Email,
		Username:  "testuser",
		IsActive:  true,
		IsVerified: true,
	}

	hashedPassword, err := authPkg.HashPassword("SecurePass123!", 10)
	require.NoError(t, err)
	user.Password = hashedPassword

	mockRepo.On("getUserByEmail", req.Email).Return(user, nil)

	response, err := service.login(req, "test-agent", "127.0.0.1")

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, ErrInvalidCredentials, err)

	mockRepo.AssertExpectations(t)
}

func TestService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtManager := authPkg.NewJWTManager("test-secret", "test-issuer", 1*time.Hour, 24*time.Hour)
	service := NewService(mockRepo, jwtManager, 10)

	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	mockRepo.On("getUserByEmail", req.Email).Return(nil, ErrUserNotFound)

	response, err := service.login(req, "test-agent", "127.0.0.1")

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, ErrInvalidCredentials, err)

	mockRepo.AssertExpectations(t)
}

func TestService_RefreshToken_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtManager := authPkg.NewJWTManager("test-secret", "test-issuer", 1*time.Hour, 24*time.Hour)
	service := NewService(mockRepo, jwtManager, 10)

	userID := uuid.New()
	_, refreshToken, err := jwtManager.Generate(userID, "test@example.com", "user", "Test User")
	require.NoError(t, err)

	user := &User{
		ID:       userID,
		Email:    "test@example.com",
		Username: "testuser",
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

	response, err := service.refreshToken(refreshToken)

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)

	mockRepo.AssertExpectations(t)
}

func TestService_RefreshToken_ExpiredSession(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtManager := authPkg.NewJWTManager("test-secret", "test-issuer", 1*time.Hour, 24*time.Hour)
	service := NewService(mockRepo, jwtManager, 10)

	userID := uuid.New()
	_, refreshToken, err := jwtManager.Generate(userID, "test@example.com", "user", "Test User")
	require.NoError(t, err)

	// Expired session
	session := &Session{
		ID:           uuid.New(),
		UserID:       userID,
		RefreshToken: refreshToken,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(-1 * time.Hour), // Expired
	}

	mockRepo.On("getSessionByRefreshToken", refreshToken).Return(session, nil)

	response, err := service.refreshToken(refreshToken)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, ErrSessionExpired, err)

	mockRepo.AssertExpectations(t)
}

func TestService_Logout_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtManager := authPkg.NewJWTManager("test-secret", "test-issuer", 1*time.Hour, 24*time.Hour)
	service := NewService(mockRepo, jwtManager, 10)

	userID := uuid.New()
	_, refreshToken, err := jwtManager.Generate(userID, "test@example.com", "user", "Test User")
	require.NoError(t, err)

	session := &Session{
		ID:           uuid.New(),
		UserID:       userID,
		RefreshToken: refreshToken,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("getSessionByRefreshToken", refreshToken).Return(session, nil)
	mockRepo.On("deactivateSession", session.ID).Return(nil)

	err = service.logout(refreshToken)

	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

