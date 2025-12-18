package platforms

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/woragis/backend/server/app/internal/domains/socialmediaposts"
)

// Service orchestrates platform configuration workflows.
type Service interface {
	InitializeDefaultPlatforms(ctx context.Context) error
	GetConfig(ctx context.Context, configID uuid.UUID) (*PlatformConfig, error)
	GetConfigByName(ctx context.Context, name socialmediaposts.Platform) (*PlatformConfig, error)
	ListConfigs(ctx context.Context, activeOnly bool) ([]PlatformConfig, error)
	UpdateConfig(ctx context.Context, req UpdateConfigRequest) (*PlatformConfig, error)
	GetOptimalTimes(ctx context.Context, name socialmediaposts.Platform) (*OptimalTimesResponse, error)
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

// Request payloads

type UpdateConfigRequest struct {
	ConfigID         uuid.UUID                      `json:"-"`
	DisplayName      *string                        `json:"displayName,omitempty"`
	PostingFrequency *int                           `json:"postingFrequency,omitempty"`
	BestDays         []string                       `json:"bestDays,omitempty"`
	BestTimes        []string                       `json:"bestTimes,omitempty"`
	SupportedFormats []socialmediaposts.ContentFormat `json:"supportedFormats,omitempty"`
	IsActive         *bool                          `json:"isActive,omitempty"`
}

type OptimalTimesResponse struct {
	Platform        socialmediaposts.Platform `json:"platform"`
	BestDays        []string                  `json:"bestDays,omitempty"`
	BestTimes       []string                  `json:"bestTimes,omitempty"`
	PostingFrequency *int                     `json:"postingFrequency,omitempty"`
}

// InitializeDefaultPlatforms creates default platform configurations.
func (s *service) InitializeDefaultPlatforms(ctx context.Context) error {
	defaultPlatforms := []struct {
		name             socialmediaposts.Platform
		displayName      string
		supportedFormats []socialmediaposts.ContentFormat
	}{
		{
			socialmediaposts.PlatformLinkedIn,
			"LinkedIn",
			[]socialmediaposts.ContentFormat{
				socialmediaposts.FormatLongForm,
				socialmediaposts.FormatArticle,
				socialmediaposts.FormatPost,
			},
		},
		{
			socialmediaposts.PlatformTwitter,
			"Twitter",
			[]socialmediaposts.ContentFormat{
				socialmediaposts.FormatThread,
				socialmediaposts.FormatPost,
			},
		},
		{
			socialmediaposts.PlatformInstagram,
			"Instagram",
			[]socialmediaposts.ContentFormat{
				socialmediaposts.FormatCarousel,
				socialmediaposts.FormatPost,
			},
		},
		{
			socialmediaposts.PlatformMedium,
			"Medium",
			[]socialmediaposts.ContentFormat{
				socialmediaposts.FormatArticle,
				socialmediaposts.FormatLongForm,
			},
		},
		{
			socialmediaposts.PlatformSubstack,
			"Substack",
			[]socialmediaposts.ContentFormat{
				socialmediaposts.FormatNewsletter,
				socialmediaposts.FormatArticle,
			},
		},
		{
			socialmediaposts.PlatformValete,
			"Valete+",
			[]socialmediaposts.ContentFormat{
				socialmediaposts.FormatArticle,
				socialmediaposts.FormatPost,
			},
		},
		{
			socialmediaposts.PlatformWebsite,
			"Website",
			[]socialmediaposts.ContentFormat{
				socialmediaposts.FormatArticle,
				socialmediaposts.FormatLongForm,
			},
		},
	}

	for _, p := range defaultPlatforms {
		// Check if platform already exists
		existing, err := s.repo.GetConfigByName(ctx, p.name)
		if err == nil && existing != nil {
			s.logger.Debug("Platform already exists, skipping", "platform", p.name)
			continue
		}

		config, err := NewPlatformConfig(p.name, p.displayName, p.supportedFormats)
		if err != nil {
			s.logger.Error("Failed to create platform config", "platform", p.name, "error", err)
			continue
		}

		if err := s.repo.CreateConfig(ctx, config); err != nil {
			s.logger.Error("Failed to persist platform config", "platform", p.name, "error", err)
			continue
		}

		s.logger.Info("Initialized platform config", "platform", p.name)
	}

	return nil
}

func (s *service) GetConfig(ctx context.Context, configID uuid.UUID) (*PlatformConfig, error) {
	return s.repo.GetConfig(ctx, configID)
}

func (s *service) GetConfigByName(ctx context.Context, name socialmediaposts.Platform) (*PlatformConfig, error) {
	return s.repo.GetConfigByName(ctx, name)
}

func (s *service) ListConfigs(ctx context.Context, activeOnly bool) ([]PlatformConfig, error) {
	return s.repo.ListConfigs(ctx, activeOnly)
}

func (s *service) UpdateConfig(ctx context.Context, req UpdateConfigRequest) (*PlatformConfig, error) {
	config, err := s.repo.GetConfig(ctx, req.ConfigID)
	if err != nil {
		return nil, err
	}

	if req.DisplayName != nil {
		config.DisplayName = *req.DisplayName
		config.UpdatedAt = time.Now().UTC()
	}

	if req.PostingFrequency != nil {
		config.UpdatePostingFrequency(*req.PostingFrequency)
	}

	if req.BestDays != nil {
		config.SetBestDays(req.BestDays)
	}

	if req.BestTimes != nil {
		config.SetBestTimes(req.BestTimes)
	}

	if req.SupportedFormats != nil {
		config.SetSupportedFormats(req.SupportedFormats)
	}

	if req.IsActive != nil {
		config.SetActive(*req.IsActive)
	}

	if err := s.repo.UpdateConfig(ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *service) GetOptimalTimes(ctx context.Context, name socialmediaposts.Platform) (*OptimalTimesResponse, error) {
	config, err := s.repo.GetConfigByName(ctx, name)
	if err != nil {
		return nil, err
	}

	response := &OptimalTimesResponse{
		Platform:         config.Name,
		PostingFrequency: config.PostingFrequency,
	}

	// Parse best days
	if len(config.BestDays) > 0 {
		var days []string
		if err := json.Unmarshal(config.BestDays, &days); err == nil {
			response.BestDays = days
		}
	}

	// Parse best times
	if len(config.BestTimes) > 0 {
		var times []string
		if err := json.Unmarshal(config.BestTimes, &times); err == nil {
			response.BestTimes = times
		}
	}

	return response, nil
}
