package assets

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates content asset workflows.
type Service interface {
	CreateAsset(ctx context.Context, req CreateAssetRequest) (*ContentAsset, error)
	UpdateAsset(ctx context.Context, req UpdateAssetRequest) (*ContentAsset, error)
	GetAsset(ctx context.Context, assetID uuid.UUID) (*ContentAsset, error)
	GetAssetsByContentPost(ctx context.Context, contentPostID uuid.UUID) ([]ContentAsset, error)
	GetAssetsBySocialPost(ctx context.Context, socialPostID uuid.UUID) ([]ContentAsset, error)
	DeleteAsset(ctx context.Context, assetID uuid.UUID) error
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

type CreateAssetRequest struct {
	ContentPostID *uuid.UUID `json:"contentPostId,omitempty"`
	SocialPostID  *uuid.UUID `json:"socialPostId,omitempty"`
	AssetType     AssetType  `json:"assetType"`
	FilePath      string     `json:"filePath"`
	FileURL       string     `json:"fileUrl,omitempty"`
	AltText       string     `json:"altText,omitempty"`
}

type UpdateAssetRequest struct {
	AssetID  uuid.UUID `json:"-"`
	FileURL  *string   `json:"fileUrl,omitempty"`
	AltText  *string   `json:"altText,omitempty"`
}

// CreateAsset creates a new content asset.
func (s *service) CreateAsset(ctx context.Context, req CreateAssetRequest) (*ContentAsset, error) {
	asset, err := NewContentAsset(req.AssetType, req.FilePath)
	if err != nil {
		return nil, err
	}

	if req.ContentPostID != nil {
		asset.SetContentPostID(*req.ContentPostID)
	}

	if req.SocialPostID != nil {
		asset.SetSocialPostID(*req.SocialPostID)
	}

	if req.FileURL != "" {
		asset.SetFileURL(req.FileURL)
	}

	if req.AltText != "" {
		asset.SetAltText(req.AltText)
	}

	if err := s.repo.CreateAsset(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// UpdateAsset updates an existing content asset.
func (s *service) UpdateAsset(ctx context.Context, req UpdateAssetRequest) (*ContentAsset, error) {
	asset, err := s.repo.GetAsset(ctx, req.AssetID)
	if err != nil {
		return nil, err
	}

	if req.FileURL != nil {
		asset.SetFileURL(*req.FileURL)
	}

	if req.AltText != nil {
		asset.SetAltText(*req.AltText)
	}

	asset.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateAsset(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// GetAsset retrieves a content asset.
func (s *service) GetAsset(ctx context.Context, assetID uuid.UUID) (*ContentAsset, error) {
	return s.repo.GetAsset(ctx, assetID)
}

// GetAssetsByContentPost retrieves assets for a content post.
func (s *service) GetAssetsByContentPost(ctx context.Context, contentPostID uuid.UUID) ([]ContentAsset, error) {
	return s.repo.GetAssetsByContentPost(ctx, contentPostID)
}

// GetAssetsBySocialPost retrieves assets for a social media post.
func (s *service) GetAssetsBySocialPost(ctx context.Context, socialPostID uuid.UUID) ([]ContentAsset, error) {
	return s.repo.GetAssetsBySocialPost(ctx, socialPostID)
}

// DeleteAsset deletes a content asset.
func (s *service) DeleteAsset(ctx context.Context, assetID uuid.UUID) error {
	return s.repo.DeleteAsset(ctx, assetID)
}
