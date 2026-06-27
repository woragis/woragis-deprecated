package skills

import (
	"errors"
	"fmt"
	"time"

	"woragis-skills-service/pkg/auth"
	"woragis-skills-service/pkg/crypto"
	appmetrics "woragis-skills-service/pkg/metrics"

	"github.com/google/uuid"
)

var (
	ErrInvalidPassword     = errors.New("invalid password")
	ErrUserNotVerified     = errors.New("user email not verified")
	ErrUserInactive        = errors.New("user account is inactive")
	ErrSessionExpired      = errors.New("session has expired")
	ErrTokenExpired        = errors.New("verification token has expired")
	ErrTokenAlreadyUsed    = errors.New("verification token already used")
	ErrPasswordTooWeak     = errors.New("password does not meet strength requirements")
)

// Service defines the interface for auth service operations
type Service interface {
	register(req *RegisterRequest) (*AuthResponse, error)
	login(req *LoginRequest, userAgent, ipAddress string) (*AuthResponse, error)
	refreshToken(refreshToken string) (*AuthResponse, error)
	logout(refreshToken string) error
	logoutAll(userID uuid.UUID) error
	getUserProfile(userID uuid.UUID) (*Profile, error)
	updateUserProfile(userID uuid.UUID, req *ProfileUpdateRequest) (*Profile, error)
	changePassword(userID uuid.UUID, req *PasswordChangeRequest) error
	createVerificationToken(userID uuid.UUID, tokenType string) (*VerificationToken, error)
	verifyEmail(token string) error
	getUserByID(userID uuid.UUID) (*User, error)
	listUsers(page, limit int, search string) ([]User, int64, error)
	cleanupExpiredSessions() error
}

// serviceImpl implements the Service interface
type serviceImpl struct {
	repo       Repository
	jwtManager *auth.JWTManager
	bcryptCost int
}

// NewService creates a new auth service
func NewService(repo Repository, jwtManager *auth.JWTManager, bcryptCost int) Service {
	return &serviceImpl{
		repo:       repo,
		jwtManager: jwtManager,
		bcryptCost: bcryptCost,
	}
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=30"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=2,max=50"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Username string `json:"username" validate:"omitempty,min=3,max=30"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

// ProfileUpdateRequest represents profile update request
type ProfileUpdateRequest struct {
	Avatar      string     `json:"avatar"`
	Bio         string     `json:"bio"`
	DateOfBirth *time.Time `json:"date_of_birth"`
	Gender      string     `json:"gender"`
	Phone       string     `json:"phone"`
	Location    string     `json:"location"`
	Website     string     `json:"website"`
	SocialLinks string     `json:"social_links"`
	Preferences string     `json:"preferences"`
}

// PasswordChangeRequest represents password change request
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// register registers a new user
func (s *serviceImpl) register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.repo.getUserByEmail(req.Email)
	if err != nil && err != ErrUserNotFound {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Validate password strength
	if err := auth.CheckPasswordStrength(req.Password); err != nil {
		return nil, ErrPasswordTooWeak
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password, s.bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &User{
		ID:        uuid.New(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
		IsActive:  true,
		IsVerified: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.createUser(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Record user registration metric
	appmetrics.RecordUserRegistration()

	// Create default profile
	profile := &Profile{
		ID:           uuid.New(),
		UserID:       user.ID,
		SocialLinks:  "{}", // Empty JSON object
		Preferences:  "{}", // Empty JSON object
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.createProfile(profile); err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	// Generate tokens
	accessToken, refreshToken, err := s.jwtManager.Generate(user.ID, user.Email, user.Role, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session
	session := &Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.createSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	if err := s.repo.updateLastLogin(user.ID); err != nil {
		return nil, fmt.Errorf("failed to update last login: %w", err)
	}

	// Load user with profile
	user, err = s.repo.getUserByID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(), // 24 hours
	}, nil
}

// login authenticates a user
func (s *serviceImpl) login(req *LoginRequest, userAgent, ipAddress string) (*AuthResponse, error) {
	var user *User
	var err error

	// Validate that either email or username is provided
	if req.Email == "" && req.Username == "" {
		return nil, ErrInvalidCredentials
	}

	// Get user by email or username
	if req.Email != "" {
		user, err = s.repo.getUserByEmail(req.Email)
	} else {
		user, err = s.repo.getUserByUsername(req.Username)
	}

	if err != nil {
		if err == ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	if err := auth.VerifyPassword(req.Password, user.Password); err != nil {
		appmetrics.RecordUserLogin(false)
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, refreshToken, err := s.jwtManager.Generate(user.ID, user.Email, user.Role, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session
	session := &Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    ipAddress,
		IsActive:     true,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.createSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	if err := s.repo.updateLastLogin(user.ID); err != nil {
		return nil, fmt.Errorf("failed to update last login: %w", err)
	}

	// Record successful login
	appmetrics.RecordUserLogin(true)

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(), // 24 hours
	}, nil
}

// refreshToken refreshes an access token
func (s *serviceImpl) refreshToken(refreshToken string) (*AuthResponse, error) {
	// Get session by refresh token
	session, err := s.repo.getSessionByRefreshToken(refreshToken)
	if err != nil {
		if err == ErrSessionNotFound {
			return nil, ErrSessionExpired
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Check if session is valid
	if !session.IsSessionValid() {
		return nil, ErrSessionExpired
	}

	// Get user
	user, err := s.repo.getUserByID(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Generate new access token
	newAccessToken, err := s.jwtManager.Refresh(refreshToken)
	if err != nil {
		appmetrics.RecordTokenRefresh(false)
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	// Record successful token refresh
	appmetrics.RecordTokenRefresh(true)

	return &AuthResponse{
		User:         user,
		AccessToken:  newAccessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour).Unix(), // 24 hours
	}, nil
}

// logout logs out a user by deactivating their session
func (s *serviceImpl) logout(refreshToken string) error {
	session, err := s.repo.getSessionByRefreshToken(refreshToken)
	if err != nil {
		if err == ErrSessionNotFound {
			return nil // Already logged out
		}
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Revoke refresh token (add to blacklist)
	// Use refresh token expiry duration (typically 7 days)
	if err := s.jwtManager.RevokeToken(refreshToken, 7*24*time.Hour); err != nil {
		// Log error but continue with session deactivation
		// Token revocation is best-effort
	} else {
		// Record token revocation
		appmetrics.RecordTokenRevocation()
	}

	return s.repo.deactivateSession(session.ID)
}

// logoutAll logs out a user from all devices
func (s *serviceImpl) logoutAll(userID uuid.UUID) error {
	// Revoke all tokens for the user (add to blacklist)
	// Use refresh token expiry duration (typically 7 days)
	if err := s.jwtManager.RevokeUserTokens(userID, 7*24*time.Hour); err != nil {
		// Log error but continue with session deactivation
		// Token revocation is best-effort
	} else {
		// Record token revocation
		appmetrics.RecordTokenRevocation()
	}

	return s.repo.deactivateAllUserSessions(userID)
}

// getUserProfile retrieves user profile
func (s *serviceImpl) getUserProfile(userID uuid.UUID) (*Profile, error) {
	return s.repo.getProfileByUserID(userID)
}

// updateUserProfile updates user profile
func (s *serviceImpl) updateUserProfile(userID uuid.UUID, req *ProfileUpdateRequest) (*Profile, error) {
	profile, err := s.repo.getProfileByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	// Update profile fields
	if req.Avatar != "" {
		profile.Avatar = req.Avatar
	}
	if req.Bio != "" {
		profile.Bio = req.Bio
	}
	if req.DateOfBirth != nil {
		profile.DateOfBirth = req.DateOfBirth
	}
	if req.Gender != "" {
		profile.Gender = req.Gender
	}
	if req.Phone != "" {
		profile.Phone = req.Phone
	}
	if req.Location != "" {
		profile.Location = req.Location
	}
	if req.Website != "" {
		profile.Website = req.Website
	}
	if req.SocialLinks != "" {
		profile.SocialLinks = req.SocialLinks
	}
	if req.Preferences != "" {
		profile.Preferences = req.Preferences
	}

	profile.UpdatedAt = time.Now()

	if err := s.repo.updateProfile(profile); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return profile, nil
}

// changePassword changes user password
func (s *serviceImpl) changePassword(userID uuid.UUID, req *PasswordChangeRequest) error {
	// Get user
	user, err := s.repo.getUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Verify current password
	if err := auth.VerifyPassword(req.CurrentPassword, user.Password); err != nil {
		return ErrInvalidPassword
	}

	// Validate new password strength
	if err := auth.CheckPasswordStrength(req.NewPassword); err != nil {
		return ErrPasswordTooWeak
	}

	// Hash new password
	hashedPassword, err := auth.HashPassword(req.NewPassword, s.bcryptCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	user.Password = hashedPassword
	user.UpdatedAt = time.Now()

	if err := s.repo.updateUser(user); err != nil {
		appmetrics.RecordPasswordChange(false)
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Record successful password change
	appmetrics.RecordPasswordChange(true)

	// Logout from all devices for security
	return s.logoutAll(userID)
}

// createVerificationToken creates a verification token for email verification
func (s *serviceImpl) createVerificationToken(userID uuid.UUID, tokenType string) (*VerificationToken, error) {
	// Generate random token
	token, err := crypto.GenerateRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	verificationToken := &VerificationToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		Type:      tokenType,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours
		IsUsed:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.createVerificationToken(verificationToken); err != nil {
		return nil, fmt.Errorf("failed to create verification token: %w", err)
	}

	return verificationToken, nil
}

// verifyEmail verifies user email using verification token
func (s *serviceImpl) verifyEmail(token string) error {
	verificationToken, err := s.repo.getVerificationToken(token)
	if err != nil {
		if err == ErrTokenNotFound {
			return ErrTokenExpired
		}
		return fmt.Errorf("failed to get verification token: %w", err)
	}

	// Check if token is valid
	if !verificationToken.IsTokenValid() {
		if verificationToken.IsUsed {
			return ErrTokenAlreadyUsed
		}
		return ErrTokenExpired
	}

	// Mark token as used
	if err := s.repo.markTokenAsUsed(verificationToken.ID); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	// Verify user email
	if err := s.repo.verifyUserEmail(verificationToken.UserID); err != nil {
		appmetrics.RecordEmailVerification(false)
		return fmt.Errorf("failed to verify user email: %w", err)
	}

	// Record successful email verification
	appmetrics.RecordEmailVerification(true)

	return nil
}

// getUserByID retrieves a user by ID
func (s *serviceImpl) getUserByID(userID uuid.UUID) (*User, error) {
	return s.repo.getUserByID(userID)
}

// listUsers retrieves users with pagination
func (s *serviceImpl) listUsers(page, limit int, search string) ([]User, int64, error) {
	offset := (page - 1) * limit
	return s.repo.listUsers(offset, limit, search)
}

// cleanupExpiredSessions removes expired sessions and tokens
func (s *serviceImpl) cleanupExpiredSessions() error {
	if err := s.repo.deleteExpiredSessions(); err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	if err := s.repo.deleteExpiredTokens(); err != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	if err := s.repo.deleteUsedTokens(); err != nil {
		return fmt.Errorf("failed to delete used tokens: %w", err)
	}

	return nil
}
