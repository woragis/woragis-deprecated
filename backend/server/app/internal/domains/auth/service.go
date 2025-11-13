package auth

import (
	"context"
	"log/slog"
	"strings"

	"golang.org/x/crypto/bcrypt"

	emailservice "github.com/woragis/backend/server/app/internal/services/email"
)

// Service exposes use-cases for the auth domain.
type Service struct {
	repo        Repository
	emailSender emailservice.Sender
	logger      *slog.Logger
}

// NewService wires a Service with its collaborators.
func NewService(repo Repository, emailSender emailservice.Sender, logger *slog.Logger) *Service {
	return &Service{
		repo:        repo,
		emailSender: emailSender,
		logger:      logger,
	}
}

// RegisterRequest transports user registration data.
type RegisterRequest struct {
	Email    string
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

	if err := s.emailSender.Send(ctx, user.Email, "Welcome to Woragis", "Your account has been created."); err != nil {
		if s.logger != nil {
			s.logger.Warn(ErrUnableToSendEmail, slog.String("email", user.Email), slog.Any("error", err))
		}
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

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", NewDomainError(ErrCodePasswordHashFailure, ErrHashPassword)
	}

	return string(hash), nil
}
