package userprofiles

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service orchestrates user profile operations.
type Service interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
	UpsertProfile(ctx context.Context, userID uuid.UUID, aboutMe string) (*UserProfile, error)
}

// service implements Service.
type service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new user profile service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// GetProfile retrieves a user profile.
func (s *service) GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error) {
	return s.repo.GetProfile(ctx, userID)
}

// UpsertProfile creates or updates a user profile.
func (s *service) UpsertProfile(ctx context.Context, userID uuid.UUID, aboutMe string) (*UserProfile, error) {
	profile, err := NewUserProfile(userID, aboutMe)
	if err != nil {
		return nil, err
	}

	if err := s.repo.UpsertProfile(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}

