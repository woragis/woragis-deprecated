package auth

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	emailservice "github.com/woragis/backend/server/app/internal/services/email"
	"github.com/woragis/backend/server/app/internal/monitoring"
)

// Service exposes use-cases for the auth domain.
type Service struct {
	repo        Repository
	emailSender emailservice.Sender
	tokenStore  TokenStore
	monitor     monitoring.Tracker
	logger      *slog.Logger
	publicURL   string
	tokenTTL    time.Duration
}

// NewService wires a Service with its collaborators.
func NewService(repo Repository, emailSender emailservice.Sender, tokenStore TokenStore, monitor monitoring.Tracker, publicURL string, logger *slog.Logger) *Service {
	if publicURL == "" {
		publicURL = "http://localhost:8080"
	}
	return &Service{
		repo:        repo,
		emailSender: emailSender,
		tokenStore:  tokenStore,
		monitor:    monitor,
		logger:     logger,
		publicURL:  publicURL,
		tokenTTL:    30 * time.Minute,
	}
}

// RegisterRequest transports user registration data.
type RegisterRequest struct {
	Email    string
	Password string
}

// PasswordResetRequest is used to initiate a reset flow.
type PasswordResetRequest struct {
	Email string
}

// PasswordResetConfirmRequest finalises password change.
type PasswordResetConfirmRequest struct {
	Token    string
	Password string
}

// Register registers a new user account and sends a welcome e-mail.
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" {
		return nil, NewDomainError(ErrCodeInvalidEmail, ErrEmptyEmail)
	}

	if len(req.Password) < 8 {
		return nil, NewDomainError(ErrCodeInvalidPassword, ErrPasswordTooShort)
	}

	if len(req.Password) > 72 {
		return nil, NewDomainError(ErrCodeInvalidPassword, ErrPasswordTooLong)
	}

	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		return nil, NewDomainError(ErrCodeEmailAlreadyExists, ErrUserAlreadyExists)
	} else if domainErr, ok := AsDomainError(err); ok {
		if domainErr.Code != ErrCodeUserNotFound {
			return nil, domainErr
		}
	} else if err != nil {
		return nil, err
	}

	passwordHash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := NewUser(email, passwordHash)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	s.sendEmail(ctx, user.Email, "Welcome to Woragis", "Your account has been created.")

	if s.monitor != nil {
		s.monitor.RecordUserRegistration(ctx, user.ID)
	}

	return user, nil
}

// Authenticate validates credentials and returns the user aggregate.
func (s *Service) Authenticate(ctx context.Context, email, password string) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return nil, NewDomainError(ErrCodeInvalidCredentials, ErrCredentialsMismatch)
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, NewDomainError(ErrCodeInvalidCredentials, ErrCredentialsMismatch)
	}

	return user, nil
}

// RequestPasswordReset sends a reset link to the supplied email.
func (s *Service) RequestPasswordReset(ctx context.Context, req PasswordResetRequest) error {
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" {
		return NewDomainError(ErrCodeInvalidEmail, ErrEmptyEmail)
	}

	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		// Do not leak user existence; log and return nil
		if s.logger != nil {
			s.logger.Info("password reset requested for unknown email", slog.String("email", email))
		}
		return nil
	}

	if s.tokenStore == nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrInvalidResetToken)
	}

	token, err := s.tokenStore.CreateToken(ctx, user.ID, s.tokenTTL)
	if err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrInvalidResetToken)
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", strings.TrimRight(s.publicURL, "/"), token)
	body := fmt.Sprintf("Hello,\n\nWe received a request to reset your Woragis password. Use the link below to choose a new password (valid for %d minutes):\n%s\n\nIf you didn't request this, you can safely ignore this email.\n", int(s.tokenTTL.Minutes()), resetURL)

	s.sendEmail(ctx, user.Email, "Reset your Woragis password", body)

	return nil
}

// ResetPassword validates the token and updates the password hash.
func (s *Service) ResetPassword(ctx context.Context, req PasswordResetConfirmRequest) error {
	if len(req.Password) < 8 {
		return NewDomainError(ErrCodeInvalidPassword, ErrPasswordTooShort)
	}
	if len(req.Password) > 72 {
		return NewDomainError(ErrCodeInvalidPassword, ErrPasswordTooLong)
	}

	if s.tokenStore == nil {
		return NewDomainError(ErrCodeInvalidToken, ErrInvalidResetToken)
	}

	userID, err := s.tokenStore.ValidateToken(ctx, req.Token)
	if err != nil {
		return NewDomainError(ErrCodeInvalidToken, ErrInvalidResetToken)
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePasswordHash(ctx, userID, hash); err != nil {
		return err
	}

	_ = s.tokenStore.InvalidateToken(ctx, req.Token)

	user, err := s.repo.FindByID(ctx, userID)
	if err == nil {
		s.sendEmail(ctx, user.Email, "Your Woragis password was changed", "If you didn't perform this change please contact support immediately.")
	}

	return nil
}

func (s *Service) sendEmail(ctx context.Context, to, subject, body string) {
	if s.emailSender == nil {
		return
	}

	if err := s.emailSender.Send(ctx, to, subject, body); err != nil && s.logger != nil {
		s.logger.Warn("auth: email send failed", slog.Any("error", err), slog.String("email", to))
	}
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", NewDomainError(ErrCodePasswordHashFailure, ErrHashPassword)
	}

	return string(hash), nil
}
