package whatsapp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// Service provides high-level WhatsApp messaging operations.
type Service interface {
	SendToUser(ctx context.Context, userID uuid.UUID, message string) error
	SendToClient(ctx context.Context, userID uuid.UUID, clientID uuid.UUID, message string) error
	SendToPhoneNumber(ctx context.Context, phoneNumber string, message string) error
}

// service implements WhatsApp service operations.
type service struct {
	notifier     Notifier
	userRepo     UserRepository
	clientRepo   ClientRepository
	logger       *slog.Logger
}

// UserRepository defines operations to fetch user phone numbers.
type UserRepository interface {
	GetUserPhoneNumber(ctx context.Context, userID uuid.UUID) (string, error)
}

// ClientRepository defines operations to fetch client phone numbers.
type ClientRepository interface {
	GetClientPhoneNumber(ctx context.Context, userID, clientID uuid.UUID) (string, error)
}

// NewService creates a new WhatsApp service.
func NewService(notifier Notifier, userRepo UserRepository, clientRepo ClientRepository, logger *slog.Logger) Service {
	return &service{
		notifier:   notifier,
		userRepo:   userRepo,
		clientRepo: clientRepo,
		logger:     logger,
	}
}

// SendToUser sends a WhatsApp message to a user by their user ID.
func (s *service) SendToUser(ctx context.Context, userID uuid.UUID, message string) error {
	// Validate inputs
	if err := ValidateMessage(message); err != nil {
		return NewDomainError(ErrCodeInvalidPayload, fmt.Sprintf("message validation failed: %v", err))
	}

	if s.notifier == nil {
		return NewDomainError(ErrCodeServiceNotConfigured, ErrServiceNotConfigured)
	}

	if s.userRepo == nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrRepositoryFailure)
	}

	phoneNumber, err := s.userRepo.GetUserPhoneNumber(ctx, userID)
	if err != nil {
		// Check if it's a domain error from the repository adapter
		if _, ok := AsDomainError(err); ok {
			return NewDomainError(ErrCodeUserNotFound, ErrUserNotFound)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrRepositoryFailure)
	}

	if phoneNumber == "" {
		return NewDomainError(ErrCodeNoPhoneNumber, ErrUserNoPhoneNumber)
	}

	// Validate phone number
	if err := ValidatePhoneNumber(phoneNumber); err != nil {
		return NewDomainError(ErrCodeInvalidPayload, fmt.Sprintf("phone number validation failed: %v", err))
	}

	if err := s.notifier.Send(ctx, phoneNumber, message); err != nil {
		// Check if it's already a domain error
		if domainErr, ok := AsDomainError(err); ok {
			return domainErr
		}
		return NewDomainError(ErrCodeSendFailure, ErrSendFailure)
	}

	return nil
}

// SendToClient sends a WhatsApp message to a client by their client ID.
func (s *service) SendToClient(ctx context.Context, userID uuid.UUID, clientID uuid.UUID, message string) error {
	// Validate inputs
	if err := ValidateMessage(message); err != nil {
		return NewDomainError(ErrCodeInvalidPayload, fmt.Sprintf("message validation failed: %v", err))
	}

	if s.notifier == nil {
		return NewDomainError(ErrCodeServiceNotConfigured, ErrServiceNotConfigured)
	}

	if s.clientRepo == nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrRepositoryFailure)
	}

	phoneNumber, err := s.clientRepo.GetClientPhoneNumber(ctx, userID, clientID)
	if err != nil {
		// Check if it's a domain error from the repository adapter
		if _, ok := AsDomainError(err); ok {
			return NewDomainError(ErrCodeClientNotFound, ErrClientNotFound)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrRepositoryFailure)
	}

	if phoneNumber == "" {
		return NewDomainError(ErrCodeNoPhoneNumber, ErrClientNoPhoneNumber)
	}

	// Validate phone number
	if err := ValidatePhoneNumber(phoneNumber); err != nil {
		return NewDomainError(ErrCodeInvalidPayload, fmt.Sprintf("phone number validation failed: %v", err))
	}

	if err := s.notifier.Send(ctx, phoneNumber, message); err != nil {
		// Check if it's already a domain error
		if domainErr, ok := AsDomainError(err); ok {
			return domainErr
		}
		return NewDomainError(ErrCodeSendFailure, ErrSendFailure)
	}

	return nil
}

// SendToPhoneNumber sends a WhatsApp message directly to a phone number.
func (s *service) SendToPhoneNumber(ctx context.Context, phoneNumber string, message string) error {
	if s.notifier == nil {
		return NewDomainError(ErrCodeServiceNotConfigured, ErrServiceNotConfigured)
	}

	if phoneNumber == "" {
		return NewDomainError(ErrCodeInvalidPhoneNumber, ErrInvalidPhoneNumber)
	}

	if err := s.notifier.Send(ctx, phoneNumber, message); err != nil {
		// Check if it's already a domain error
		if domainErr, ok := AsDomainError(err); ok {
			return domainErr
		}
		return NewDomainError(ErrCodeSendFailure, ErrSendFailure)
	}

	return nil
}

