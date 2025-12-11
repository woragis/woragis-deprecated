package userpreferences

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service orchestrates user preferences operations.
type Service interface {
	GetPreferences(ctx context.Context, userID uuid.UUID) (*UserPreferences, error)
	GetOrCreatePreferences(ctx context.Context, userID uuid.UUID) (*UserPreferences, error)
	UpdatePreferences(ctx context.Context, userID uuid.UUID, language, currency string) (*UserPreferences, error)
	GetDefaultLanguage(ctx context.Context, userID uuid.UUID) (string, error)
	GetDefaultCurrency(ctx context.Context, userID uuid.UUID) (string, error)
	GetDefaultWebsite(ctx context.Context, userID uuid.UUID) (string, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService constructs a Service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) GetPreferences(ctx context.Context, userID uuid.UUID) (*UserPreferences, error) {
	return s.repo.GetPreferences(ctx, userID)
}

func (s *service) GetOrCreatePreferences(ctx context.Context, userID uuid.UUID) (*UserPreferences, error) {
	prefs, err := s.repo.GetPreferences(ctx, userID)
	if err != nil {
		if domainErr, ok := AsDomainError(err); ok && domainErr.Code == ErrCodeNotFound {
			// Create with defaults
			newPrefs, createErr := NewUserPreferences(userID)
			if createErr != nil {
				return nil, createErr
			}
			if err := s.repo.CreatePreferences(ctx, newPrefs); err != nil {
				return nil, err
			}
			return newPrefs, nil
		}
		return nil, err
	}
	return prefs, nil
}

func (s *service) UpdatePreferences(ctx context.Context, userID uuid.UUID, language, currency string) (*UserPreferences, error) {
	prefs, err := s.GetOrCreatePreferences(ctx, userID)
	if err != nil {
		return nil, err
	}

	if language != "" {
		if err := prefs.UpdateLanguage(language); err != nil {
			return nil, err
		}
	}

	if currency != "" {
		if err := prefs.UpdateCurrency(currency); err != nil {
			return nil, err
		}
	}

	if err := s.repo.UpdatePreferences(ctx, prefs); err != nil {
		return nil, err
	}

	return prefs, nil
}

func (s *service) GetDefaultLanguage(ctx context.Context, userID uuid.UUID) (string, error) {
	prefs, err := s.GetOrCreatePreferences(ctx, userID)
	if err != nil {
		return "en", err // Return default on error
	}
	if prefs.DefaultLanguage == "" {
		return "en", nil
	}
	return prefs.DefaultLanguage, nil
}

func (s *service) GetDefaultCurrency(ctx context.Context, userID uuid.UUID) (string, error) {
	prefs, err := s.GetOrCreatePreferences(ctx, userID)
	if err != nil {
		return "USD", err // Return default on error
	}
	if prefs.DefaultCurrency == "" {
		return "USD", nil
	}
	return prefs.DefaultCurrency, nil
}

func (s *service) GetDefaultWebsite(ctx context.Context, userID uuid.UUID) (string, error) {
	prefs, err := s.GetOrCreatePreferences(ctx, userID)
	if err != nil {
		return "", err // Return empty on error
	}
	return prefs.DefaultWebsite, nil
}

