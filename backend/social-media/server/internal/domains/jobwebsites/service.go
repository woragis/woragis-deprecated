package jobwebsites

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service orchestrates job website operations.
type Service interface {
	CreateJobWebsite(ctx context.Context, name, displayName, baseURL, loginURL string, dailyLimit int) (*JobWebsite, error)
	GetJobWebsite(ctx context.Context, websiteID uuid.UUID) (*JobWebsite, error)
	GetJobWebsiteByName(ctx context.Context, name string) (*JobWebsite, error)
	ListJobWebsites(ctx context.Context, enabledOnly bool) ([]JobWebsite, error)
	UpdateJobWebsite(ctx context.Context, websiteID uuid.UUID, updates JobWebsiteUpdates) (*JobWebsite, error)
	IncrementCount(ctx context.Context, websiteName string) error
	ResetCount(ctx context.Context, websiteID uuid.UUID) error
	DeleteJobWebsite(ctx context.Context, websiteID uuid.UUID) error
}

type JobWebsiteUpdates struct {
	DailyLimit  *int
	Enabled     *bool
	BaseURL     *string
	LoginURL    *string
	DisplayName *string
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

func (s *service) CreateJobWebsite(ctx context.Context, name, displayName, baseURL, loginURL string, dailyLimit int) (*JobWebsite, error) {
	website, err := NewJobWebsite(name, displayName, baseURL, loginURL, dailyLimit)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateJobWebsite(ctx, website); err != nil {
		return nil, err
	}

	return website, nil
}

func (s *service) GetJobWebsite(ctx context.Context, websiteID uuid.UUID) (*JobWebsite, error) {
	return s.repo.GetJobWebsite(ctx, websiteID)
}

func (s *service) GetJobWebsiteByName(ctx context.Context, name string) (*JobWebsite, error) {
	return s.repo.GetJobWebsiteByName(ctx, name)
}

func (s *service) ListJobWebsites(ctx context.Context, enabledOnly bool) ([]JobWebsite, error) {
	return s.repo.ListJobWebsites(ctx, enabledOnly)
}

func (s *service) UpdateJobWebsite(ctx context.Context, websiteID uuid.UUID, updates JobWebsiteUpdates) (*JobWebsite, error) {
	website, err := s.repo.GetJobWebsite(ctx, websiteID)
	if err != nil {
		return nil, err
	}

	if updates.DailyLimit != nil {
		if err := website.UpdateDailyLimit(*updates.DailyLimit); err != nil {
			return nil, err
		}
	}
	if updates.Enabled != nil {
		website.SetEnabled(*updates.Enabled)
	}
	if updates.BaseURL != nil {
		website.BaseURL = *updates.BaseURL
		website.UpdatedAt = website.UpdatedAt
	}
	if updates.LoginURL != nil {
		website.LoginURL = *updates.LoginURL
		website.UpdatedAt = website.UpdatedAt
	}
	if updates.DisplayName != nil {
		website.DisplayName = *updates.DisplayName
		website.UpdatedAt = website.UpdatedAt
	}

	if err := s.repo.UpdateJobWebsite(ctx, website); err != nil {
		return nil, err
	}

	return website, nil
}

func (s *service) IncrementCount(ctx context.Context, websiteName string) error {
	website, err := s.repo.GetJobWebsiteByName(ctx, websiteName)
	if err != nil {
		return err
	}

	// Check if should reset (new day)
	if website.ShouldReset() {
		website.ResetCount()
	}

	website.IncrementCount()
	return s.repo.UpdateJobWebsite(ctx, website)
}

func (s *service) ResetCount(ctx context.Context, websiteID uuid.UUID) error {
	website, err := s.repo.GetJobWebsite(ctx, websiteID)
	if err != nil {
		return err
	}

	website.ResetCount()
	return s.repo.UpdateJobWebsite(ctx, website)
}

func (s *service) DeleteJobWebsite(ctx context.Context, websiteID uuid.UUID) error {
	return s.repo.DeleteJobWebsite(ctx, websiteID)
}

