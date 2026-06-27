package assets

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AssetType represents the type of content asset.
type AssetType string

const (
	AssetTypeImage    AssetType = "image"
	AssetTypeVideo    AssetType = "video"
	AssetTypeDocument AssetType = "document"
	AssetTypeOther    AssetType = "other"
)

// ContentAsset represents an asset associated with a content post or social media post.
type ContentAsset struct {
	ID            uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ContentPostID *uuid.UUID `gorm:"column:content_post_id;type:uuid;index" json:"contentPostId,omitempty"`
	SocialPostID  *uuid.UUID `gorm:"column:social_post_id;type:uuid;index" json:"socialPostId,omitempty"`
	AssetType     AssetType  `gorm:"column:asset_type;type:varchar(20);not null;index" json:"assetType"`
	FilePath      string     `gorm:"column:file_path;type:varchar(512);not null" json:"filePath"`
	FileURL       string     `gorm:"column:file_url;type:varchar(512)" json:"fileUrl,omitempty"`
	AltText       string     `gorm:"column:alt_text;type:varchar(255)" json:"altText,omitempty"`
	Metadata      datatypes.JSON `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for ContentAsset.
func (ContentAsset) TableName() string {
	return "content_assets"
}

// NewContentAsset creates a new content asset entity.
func NewContentAsset(assetType AssetType, filePath string) (*ContentAsset, error) {
	asset := &ContentAsset{
		ID:        uuid.New(),
		AssetType: assetType,
		FilePath:  strings.TrimSpace(filePath),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return asset, asset.Validate()
}

// Validate ensures content asset invariants hold.
func (a *ContentAsset) Validate() error {
	if a == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilContentAsset)
	}

	if a.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyContentAssetID)
	}

	if a.ContentPostID == nil && a.SocialPostID == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrMissingPostID)
	}

	if !isValidAssetType(a.AssetType) {
		return NewDomainError(ErrCodeInvalidAssetType, ErrUnsupportedAssetType)
	}

	if strings.TrimSpace(a.FilePath) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyFilePath)
	}

	return nil
}

// SetContentPostID links the asset to a content post.
func (a *ContentAsset) SetContentPostID(contentPostID uuid.UUID) {
	a.ContentPostID = &contentPostID
	a.UpdatedAt = time.Now().UTC()
}

// SetSocialPostID links the asset to a social media post.
func (a *ContentAsset) SetSocialPostID(socialPostID uuid.UUID) {
	a.SocialPostID = &socialPostID
	a.UpdatedAt = time.Now().UTC()
}

// SetFileURL sets the public URL for the asset.
func (a *ContentAsset) SetFileURL(url string) {
	a.FileURL = strings.TrimSpace(url)
	a.UpdatedAt = time.Now().UTC()
}

// SetAltText sets the alt text for the asset.
func (a *ContentAsset) SetAltText(altText string) {
	a.AltText = strings.TrimSpace(altText)
	a.UpdatedAt = time.Now().UTC()
}

// Validation helpers

func isValidAssetType(at AssetType) bool {
	switch at {
	case AssetTypeImage, AssetTypeVideo, AssetTypeDocument, AssetTypeOther:
		return true
	}
	return false
}
