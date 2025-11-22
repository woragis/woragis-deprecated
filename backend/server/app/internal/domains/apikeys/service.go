package apikeys

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates API key workflows.
type Service interface {
	CreateAPIKey(ctx context.Context, userID uuid.UUID, name string, expiresAt *time.Time) (*APIKeyWithToken, error)
	UpdateAPIKey(ctx context.Context, userID uuid.UUID, apiKeyID uuid.UUID, name string) (*APIKey, error)
	DeleteAPIKey(ctx context.Context, userID uuid.UUID, apiKeyID uuid.UUID) error
	GetAPIKey(ctx context.Context, userID uuid.UUID, apiKeyID uuid.UUID) (*APIKey, error)
	ListAPIKeys(ctx context.Context, userID uuid.UUID) ([]APIKey, error)
	ValidateAPIKey(ctx context.Context, token string) (*APIKey, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

var _ Service = (*service)(nil)

// NewService constructs a Service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) CreateAPIKey(ctx context.Context, userID uuid.UUID, name string, expiresAt *time.Time) (*APIKeyWithToken, error) {
	apiKey, token, err := NewAPIKey(userID, name, expiresAt)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateAPIKey(ctx, apiKey); err != nil {
		return nil, err
	}

	return &APIKeyWithToken{
		APIKey: *apiKey,
		Token:  token,
	}, nil
}

func (s *service) UpdateAPIKey(ctx context.Context, userID uuid.UUID, apiKeyID uuid.UUID, name string) (*APIKey, error) {
	apiKey, err := s.repo.GetAPIKey(ctx, apiKeyID, userID)
	if err != nil {
		return nil, err
	}

	apiKey.Name = name
	apiKey.UpdatedAt = time.Now().UTC()

	if err := apiKey.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateAPIKey(ctx, apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

func (s *service) DeleteAPIKey(ctx context.Context, userID uuid.UUID, apiKeyID uuid.UUID) error {
	return s.repo.DeleteAPIKey(ctx, apiKeyID, userID)
}

func (s *service) GetAPIKey(ctx context.Context, userID uuid.UUID, apiKeyID uuid.UUID) (*APIKey, error) {
	return s.repo.GetAPIKey(ctx, apiKeyID, userID)
}

func (s *service) ListAPIKeys(ctx context.Context, userID uuid.UUID) ([]APIKey, error) {
	return s.repo.ListAPIKeys(ctx, userID)
}

func (s *service) ValidateAPIKey(ctx context.Context, token string) (*APIKey, error) {
	// Hash the provided token to look it up
	keyHash := hashToken(token)
	
	apiKey, err := s.repo.GetAPIKeyByHash(ctx, keyHash)
	if err != nil {
		return nil, NewDomainError(ErrCodeInvalidToken, ErrInvalidAPIKey)
	}

	// Check if expired
	if apiKey.IsExpired() {
		return nil, NewDomainError(ErrCodeExpiredToken, ErrExpiredAPIKey)
	}

	// Update last used timestamp
	apiKey.UpdateLastUsed()
	if err := s.repo.UpdateAPIKey(ctx, apiKey); err != nil {
		// Log but don't fail the request
		if s.logger != nil {
			s.logger.Warn("failed to update API key last used", "error", err)
		}
	}

	return apiKey, nil
}

