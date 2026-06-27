package creativeassets

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	// creativeservice "woragis-social-media-service/internal/services/creative" // TODO: Implement creative service integration
)

// Service orchestrates creative asset workflows
type Service interface {
	// Asset operations
	CreateAsset(ctx context.Context, userID uuid.UUID, req CreateAssetRequest) (*CreativeAsset, error)
	GetAsset(ctx context.Context, assetID uuid.UUID) (*CreativeAsset, error)
	GetAssetsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]CreativeAsset, error)
	GetAssetByEntityAndPurpose(ctx context.Context, entityType EntityType, entityID uuid.UUID, purpose AssetPurpose) (*CreativeAsset, error)
	DeleteAsset(ctx context.Context, userID uuid.UUID, assetID uuid.UUID) error

	// Generation and storage
	GenerateAndStoreImage(ctx context.Context, userID uuid.UUID, entityType EntityType, entityID uuid.UUID, purpose AssetPurpose, prompt, context string) (*CreativeAsset, error)
	GenerateAndStoreThumbnail(ctx context.Context, userID uuid.UUID, entityType EntityType, entityID uuid.UUID, prompt, context string) (*CreativeAsset, error)
	GenerateAndStoreDiagram(ctx context.Context, userID uuid.UUID, entityType EntityType, entityID uuid.UUID, description, diagramKind string) (*CreativeAsset, error)
	StoreImageData(ctx context.Context, userID uuid.UUID, entityType EntityType, entityID uuid.UUID, purpose AssetPurpose, b64Data, url string) (*CreativeAsset, error)
}

// CreateAssetRequest represents a request to create an asset
type CreateAssetRequest struct {
	EntityType EntityType   `json:"entityType"`
	EntityID   uuid.UUID    `json:"entityId"`
	AssetType  AssetType    `json:"assetType"`
	Purpose    AssetPurpose `json:"purpose"`
	B64Data    string       `json:"b64Data,omitempty"`
	URL        string       `json:"url,omitempty"`
	Prompt     string       `json:"prompt,omitempty"`
	Provider   string       `json:"provider,omitempty"`
	Format     string       `json:"format,omitempty"`
}

type service struct {
	repo             Repository
	creativeClient   interface{} // Placeholder for creative service client
}

// NewService creates a new service
func NewService(repo Repository, creativeClient interface{}) Service {
	return &service{
		repo:           repo,
		creativeClient: creativeClient,
	}
}

func (s *service) CreateAsset(ctx context.Context, userID uuid.UUID, req CreateAssetRequest) (*CreativeAsset, error) {
	asset := NewCreativeAsset(userID, req.EntityType, req.EntityID, req.AssetType, req.Purpose)
	
	if req.B64Data != "" || req.URL != "" {
		// Calculate size if we have base64 data
		var sizeBytes int64
		if req.B64Data != "" {
			decoded, err := base64.StdEncoding.DecodeString(req.B64Data)
			if err == nil {
				sizeBytes = int64(len(decoded))
			}
		}
		asset.SetImageData(req.B64Data, req.URL, 0, 0, sizeBytes)
	}
	
	if req.Prompt != "" || req.Provider != "" || req.Format != "" {
		asset.SetMetadata(req.Prompt, req.Provider, req.Format)
	}

	if err := asset.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateAsset(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *service) GetAsset(ctx context.Context, assetID uuid.UUID) (*CreativeAsset, error) {
	return s.repo.GetAsset(ctx, assetID)
}

func (s *service) GetAssetsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]CreativeAsset, error) {
	return s.repo.GetAssetsByEntity(ctx, entityType, entityID)
}

func (s *service) GetAssetByEntityAndPurpose(ctx context.Context, entityType EntityType, entityID uuid.UUID, purpose AssetPurpose) (*CreativeAsset, error) {
	return s.repo.GetAssetByEntityAndPurpose(ctx, entityType, entityID, purpose)
}

func (s *service) DeleteAsset(ctx context.Context, userID uuid.UUID, assetID uuid.UUID) error {
	return s.repo.DeleteAsset(ctx, assetID, userID)
}

func (s *service) GenerateAndStoreImage(ctx context.Context, userID uuid.UUID, entityType EntityType, entityID uuid.UUID, purpose AssetPurpose, prompt, context string) (*CreativeAsset, error) {
	// TODO: Implement creative service integration
	// Generate image using creative service
	// req := creativeservice.ImageGenerationRequest{
	// 	Provider: creativeservice.ProviderOpenAI,
	// 	Prompt:   prompt,
	// 	Context:  context,
	// 	Style:    creativeservice.StyleTechnical,
	// 	N:        1,
	// }
	//
	// resp, err := s.creativeClient.GenerateImage(ctx, req)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// if len(resp.Data) == 0 {
	// 	return nil, NewDomainError(ErrCodeInvalidPayload, "no image data returned")
	// }
	//
	// // Store the asset
	// asset := NewCreativeAsset(userID, entityType, entityID, AssetTypeImage, purpose)
	// asset.SetImageData(resp.Data[0].B64JSON, resp.Data[0].URL, 0, 0, 0)
	// asset.SetMetadata(prompt, string(resp.Provider), "png")
	//
	// if err := asset.Validate(); err != nil {
	// 	return nil, err
	// }
	//
	// if err := s.repo.CreateAsset(ctx, asset); err != nil {
	// 	return nil, err
	// }
	//
	// return asset, nil
	
	return nil, NewDomainError(ErrCodeNotImplemented, "creative service integration not yet implemented")
}

func (s *service) GenerateAndStoreThumbnail(ctx context.Context, userID uuid.UUID, entityType EntityType, entityID uuid.UUID, prompt, context string) (*CreativeAsset, error) {
	// TODO: Implement creative service integration
	return nil, NewDomainError(ErrCodeNotImplemented, "creative service integration not yet implemented")
}

func (s *service) GenerateAndStoreDiagram(ctx context.Context, userID uuid.UUID, entityType EntityType, entityID uuid.UUID, description, diagramKind string) (*CreativeAsset, error) {
	// TODO: Implement creative service integration
	return nil, NewDomainError(ErrCodeNotImplemented, "creative service integration not yet implemented")
}

func (s *service) StoreImageData(ctx context.Context, userID uuid.UUID, entityType EntityType, entityID uuid.UUID, purpose AssetPurpose, b64Data, url string) (*CreativeAsset, error) {
	asset := NewCreativeAsset(userID, entityType, entityID, AssetTypeImage, purpose)
	
	var sizeBytes int64
	if b64Data != "" {
		decoded, err := base64.StdEncoding.DecodeString(b64Data)
		if err == nil {
			sizeBytes = int64(len(decoded))
		}
	}
	
	asset.SetImageData(b64Data, url, 0, 0, sizeBytes)

	if err := asset.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateAsset(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

