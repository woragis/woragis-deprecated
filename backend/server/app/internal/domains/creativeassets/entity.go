package creativeassets

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// AssetType represents the type of creative asset
type AssetType string

const (
	AssetTypeImage   AssetType = "image"
	AssetTypeDiagram AssetType = "diagram"
	AssetTypeVideo   AssetType = "video"
)

// AssetPurpose represents what the asset is used for
type AssetPurpose string

const (
	PurposeThumbnail     AssetPurpose = "thumbnail"
	PurposeFeaturedImage AssetPurpose = "featured_image"
	PurposeOGImage       AssetPurpose = "og_image"
	PurposeDiagram       AssetPurpose = "diagram"
	PurposeCover         AssetPurpose = "cover"
	PurposeIllustration  AssetPurpose = "illustration"
)

// EntityType represents the entity this asset is associated with
type EntityType string

const (
	EntityTypePost              EntityType = "post"
	EntityTypeCaseStudy         EntityType = "case_study"
	EntityTypeTechnicalWriting  EntityType = "technical_writing"
	EntityTypeProblemSolution   EntityType = "problem_solution"
	EntityTypeSystemDesign      EntityType = "system_design"
)

// CreativeAsset represents a generated creative asset (image, diagram, video)
type CreativeAsset struct {
	ID              uuid.UUID   `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID   `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	EntityType      EntityType  `gorm:"column:entity_type;type:varchar(50);not null;index:idx_asset_entity" json:"entityType"`
	EntityID        uuid.UUID   `gorm:"column:entity_id;type:uuid;not null;index:idx_asset_entity" json:"entityId"`
	AssetType       AssetType   `gorm:"column:asset_type;type:varchar(50);not null;index" json:"assetType"`
	Purpose         AssetPurpose `gorm:"column:purpose;type:varchar(50);not null;index" json:"purpose"`
	// Storage
	B64Data         string      `gorm:"column:b64_data;type:text" json:"-"` // Base64 data, not exposed in JSON by default
	URL             string      `gorm:"column:url;size:512" json:"url,omitempty"`
	// Metadata
	Prompt          string      `gorm:"column:prompt;type:text" json:"prompt,omitempty"` // Original generation prompt
	Provider        string      `gorm:"column:provider;size:50" json:"provider,omitempty"` // e.g., "openai", "mermaid"
	Format          string      `gorm:"column:format;size:20" json:"format,omitempty"` // e.g., "png", "svg", "mp4"
	Width           int         `gorm:"column:width" json:"width,omitempty"`
	Height          int         `gorm:"column:height" json:"height,omitempty"`
	SizeBytes       int64       `gorm:"column:size_bytes" json:"sizeBytes,omitempty"`
	// Additional data for diagrams
	DiagramCode     string      `gorm:"column:diagram_code;type:text" json:"diagramCode,omitempty"` // Mermaid/Graphviz code
	DiagramType     string      `gorm:"column:diagram_type;size:50" json:"diagramType,omitempty"` // "mermaid" or "graphviz"
	// Timestamps
	CreatedAt       time.Time   `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time   `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for CreativeAsset
func (CreativeAsset) TableName() string {
	return "creative_assets"
}

// NewCreativeAsset creates a new creative asset entity
func NewCreativeAsset(
	userID uuid.UUID,
	entityType EntityType,
	entityID uuid.UUID,
	assetType AssetType,
	purpose AssetPurpose,
) *CreativeAsset {
	return &CreativeAsset{
		ID:         uuid.New(),
		UserID:     userID,
		EntityType: entityType,
		EntityID:   entityID,
		AssetType:  assetType,
		Purpose:    purpose,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

// Validate ensures creative asset invariants hold
func (c *CreativeAsset) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilCreativeAsset)
	}

	if c.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyAssetID)
	}

	if c.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if c.EntityID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyEntityID)
	}

	if !isValidEntityType(c.EntityType) {
		return NewDomainError(ErrCodeInvalidType, ErrUnsupportedEntityType)
	}

	if !isValidAssetType(c.AssetType) {
		return NewDomainError(ErrCodeInvalidType, ErrUnsupportedAssetType)
	}

	if !isValidPurpose(c.Purpose) {
		return NewDomainError(ErrCodeInvalidType, ErrUnsupportedPurpose)
	}

	return nil
}

// SetImageData sets image data from base64 or URL
func (c *CreativeAsset) SetImageData(b64Data, url string, width, height int, sizeBytes int64) {
	if b64Data != "" {
		c.B64Data = strings.TrimSpace(b64Data)
	}
	if url != "" {
		c.URL = strings.TrimSpace(url)
	}
	if width > 0 {
		c.Width = width
	}
	if height > 0 {
		c.Height = height
	}
	if sizeBytes > 0 {
		c.SizeBytes = sizeBytes
	}
	c.UpdatedAt = time.Now().UTC()
}

// SetDiagramData sets diagram-specific data
func (c *CreativeAsset) SetDiagramData(code, diagramType string) {
	if code != "" {
		c.DiagramCode = strings.TrimSpace(code)
	}
	if diagramType != "" {
		c.DiagramType = strings.TrimSpace(diagramType)
	}
	c.UpdatedAt = time.Now().UTC()
}

// SetMetadata sets generation metadata
func (c *CreativeAsset) SetMetadata(prompt, provider, format string) {
	if prompt != "" {
		c.Prompt = strings.TrimSpace(prompt)
	}
	if provider != "" {
		c.Provider = strings.TrimSpace(provider)
	}
	if format != "" {
		c.Format = strings.TrimSpace(format)
	}
	c.UpdatedAt = time.Now().UTC()
}

// GetB64Data returns the base64 data (for specific use cases)
func (c *CreativeAsset) GetB64Data() string {
	return c.B64Data
}

// Validation helpers

func isValidEntityType(et EntityType) bool {
	switch et {
	case EntityTypePost, EntityTypeCaseStudy, EntityTypeTechnicalWriting,
		EntityTypeProblemSolution, EntityTypeSystemDesign:
		return true
	}
	return false
}

func isValidAssetType(at AssetType) bool {
	switch at {
	case AssetTypeImage, AssetTypeDiagram, AssetTypeVideo:
		return true
	}
	return false
}

func isValidPurpose(p AssetPurpose) bool {
	switch p {
	case PurposeThumbnail, PurposeFeaturedImage, PurposeOGImage,
		PurposeDiagram, PurposeCover, PurposeIllustration:
		return true
	}
	return false
}

