package creativeassets

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"woragis-social-media-service/pkg/middleware"
	"woragis-social-media-service/pkg/response"
)

// Handler exposes creative asset endpoints
type Handler interface {
	CreateAsset(c *fiber.Ctx) error
	GetAsset(c *fiber.Ctx) error
	GetAssetsByEntity(c *fiber.Ctx) error
	GetAssetByEntityAndPurpose(c *fiber.Ctx) error
	DeleteAsset(c *fiber.Ctx) error
	GenerateImage(c *fiber.Ctx) error
	GenerateThumbnail(c *fiber.Ctx) error
	GenerateDiagram(c *fiber.Ctx) error
	GetAssetData(c *fiber.Ctx) error // Returns base64 data for serving images
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a creative asset handler
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type createAssetPayload struct {
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

type generateImagePayload struct {
	EntityType EntityType   `json:"entityType"`
	EntityID   uuid.UUID    `json:"entityId"`
	Purpose    AssetPurpose `json:"purpose"`
	Prompt     string       `json:"prompt"`
	Context    string       `json:"context,omitempty"`
}

type generateThumbnailPayload struct {
	EntityType EntityType `json:"entityType"`
	EntityID   uuid.UUID  `json:"entityId"`
	Prompt     string     `json:"prompt"`
	Context    string     `json:"context,omitempty"`
}

type generateDiagramPayload struct {
	EntityType  EntityType `json:"entityType"`
	EntityID    uuid.UUID  `json:"entityId"`
	Description string     `json:"description"`
	DiagramKind string     `json:"diagramKind,omitempty"`
}

// Handlers

func (h *handler) CreateAsset(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
			"message": "unauthorized",
		})
	}

	var payload createAssetPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	req := CreateAssetRequest{
		EntityType: payload.EntityType,
		EntityID:   payload.EntityID,
		AssetType:  payload.AssetType,
		Purpose:    payload.Purpose,
		B64Data:    payload.B64Data,
		URL:        payload.URL,
		Prompt:     payload.Prompt,
		Provider:   payload.Provider,
		Format:     payload.Format,
	}

	asset, err := h.service.CreateAsset(c.UserContext(), userID, req)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			switch domainErr.Code {
			case ErrCodeNotFound:
				return response.Error(c, fiber.StatusNotFound, ErrCodeNotFound, fiber.Map{
					"message": domainErr.Message,
				})
			case ErrCodeUnauthorized:
				return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
					"message": domainErr.Message,
				})
			case ErrCodeInvalidPayload, ErrCodeInvalidType:
				return response.Error(c, fiber.StatusBadRequest, domainErr.Code, fiber.Map{
					"message": domainErr.Message,
				})
			}
		}
		h.logger.Error("failed to create asset", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to create asset",
			"error":   err.Error(),
		})
	}

	return response.Success(c, fiber.StatusCreated, asset)
}

func (h *handler) GetAsset(c *fiber.Ctx) error {
	assetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid asset ID",
			"error":   err.Error(),
		})
	}

	asset, err := h.service.GetAsset(c.UserContext(), assetID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok && domainErr.Code == ErrCodeNotFound {
			return response.Error(c, fiber.StatusNotFound, ErrCodeNotFound, fiber.Map{
				"message": domainErr.Message,
			})
		}
		h.logger.Error("failed to get asset", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to get asset",
			"error":   err.Error(),
		})
	}

	return response.Success(c, fiber.StatusOK, asset)
}

func (h *handler) GetAssetsByEntity(c *fiber.Ctx) error {
	entityType := EntityType(c.Params("entityType"))
	entityID, err := uuid.Parse(c.Params("entityId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid entity ID",
			"error":   err.Error(),
		})
	}

	assets, err := h.service.GetAssetsByEntity(c.UserContext(), entityType, entityID)
	if err != nil {
		h.logger.Error("failed to get assets by entity", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to get assets",
			"error":   err.Error(),
		})
	}

	return response.Success(c, fiber.StatusOK, assets)
}

func (h *handler) GetAssetByEntityAndPurpose(c *fiber.Ctx) error {
	entityType := EntityType(c.Params("entityType"))
	entityID, err := uuid.Parse(c.Params("entityId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid entity ID",
			"error":   err.Error(),
		})
	}
	purpose := AssetPurpose(c.Query("purpose"))
	if purpose == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "purpose query parameter is required",
		})
	}

	asset, err := h.service.GetAssetByEntityAndPurpose(c.UserContext(), entityType, entityID, purpose)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok && domainErr.Code == ErrCodeNotFound {
			return response.Error(c, fiber.StatusNotFound, ErrCodeNotFound, fiber.Map{
				"message": domainErr.Message,
			})
		}
		h.logger.Error("failed to get asset by entity and purpose", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to get asset",
			"error":   err.Error(),
		})
	}

	return response.Success(c, fiber.StatusOK, asset)
}

func (h *handler) DeleteAsset(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
			"message": "unauthorized",
		})
	}

	assetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid asset ID",
			"error":   err.Error(),
		})
	}

	if err := h.service.DeleteAsset(c.UserContext(), userID, assetID); err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			switch domainErr.Code {
			case ErrCodeNotFound:
				return response.Error(c, fiber.StatusNotFound, ErrCodeNotFound, fiber.Map{
					"message": domainErr.Message,
				})
			case ErrCodeUnauthorized:
				return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
					"message": domainErr.Message,
				})
			}
		}
		h.logger.Error("failed to delete asset", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to delete asset",
			"error":   err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *handler) GenerateImage(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
			"message": "unauthorized",
		})
	}

	var payload generateImagePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	asset, err := h.service.GenerateAndStoreImage(
		c.UserContext(),
		userID,
		payload.EntityType,
		payload.EntityID,
		payload.Purpose,
		payload.Prompt,
		payload.Context,
	)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			switch domainErr.Code {
			case ErrCodeInvalidPayload, ErrCodeInvalidType:
				return response.Error(c, fiber.StatusBadRequest, domainErr.Code, fiber.Map{
					"message": domainErr.Message,
				})
			}
		}
		h.logger.Error("failed to generate image", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to generate image",
			"error":   err.Error(),
		})
	}

	return response.Success(c, fiber.StatusCreated, asset)
}

func (h *handler) GenerateThumbnail(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
			"message": "unauthorized",
		})
	}

	var payload generateThumbnailPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	asset, err := h.service.GenerateAndStoreThumbnail(
		c.UserContext(),
		userID,
		payload.EntityType,
		payload.EntityID,
		payload.Prompt,
		payload.Context,
	)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			switch domainErr.Code {
			case ErrCodeInvalidPayload, ErrCodeInvalidType:
				return response.Error(c, fiber.StatusBadRequest, domainErr.Code, fiber.Map{
					"message": domainErr.Message,
				})
			}
		}
		h.logger.Error("failed to generate thumbnail", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to generate thumbnail",
			"error":   err.Error(),
		})
	}

	return response.Success(c, fiber.StatusCreated, asset)
}

func (h *handler) GenerateDiagram(c *fiber.Ctx) error {
	userID, err := middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
			"message": "unauthorized",
		})
	}

	var payload generateDiagramPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request body",
			"error":   err.Error(),
		})
	}

	if payload.DiagramKind == "" {
		payload.DiagramKind = "flowchart"
	}

	asset, err := h.service.GenerateAndStoreDiagram(
		c.UserContext(),
		userID,
		payload.EntityType,
		payload.EntityID,
		payload.Description,
		payload.DiagramKind,
	)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok {
			switch domainErr.Code {
			case ErrCodeInvalidPayload, ErrCodeInvalidType:
				return response.Error(c, fiber.StatusBadRequest, domainErr.Code, fiber.Map{
					"message": domainErr.Message,
				})
			}
		}
		h.logger.Error("failed to generate diagram", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to generate diagram",
			"error":   err.Error(),
		})
	}

	return response.Success(c, fiber.StatusCreated, asset)
}

func (h *handler) GetAssetData(c *fiber.Ctx) error {
	assetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid asset ID",
			"error":   err.Error(),
		})
	}

	asset, err := h.service.GetAsset(c.UserContext(), assetID)
	if err != nil {
		if domainErr, ok := err.(*DomainError); ok && domainErr.Code == ErrCodeNotFound {
			return response.Error(c, fiber.StatusNotFound, ErrCodeNotFound, fiber.Map{
				"message": domainErr.Message,
			})
		}
		h.logger.Error("failed to get asset", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to get asset",
			"error":   err.Error(),
		})
	}

	// Return base64 data for serving as image
	b64Data := asset.GetB64Data()
	if b64Data == "" && asset.URL != "" {
		// If no base64 but URL exists, redirect to URL
		return c.Redirect(asset.URL, fiber.StatusFound)
	}

	if b64Data == "" {
		return response.Error(c, fiber.StatusNotFound, ErrCodeNotFound, fiber.Map{
			"message": "asset data not found",
		})
	}

	// Set appropriate content type
	contentType := "image/png"
	switch asset.Format {
	case "jpeg", "jpg":
		contentType = "image/jpeg"
	case "webp":
		contentType = "image/webp"
	case "svg":
		contentType = "image/svg+xml"
	case "gif":
		contentType = "image/gif"
	}

	c.Set("Content-Type", contentType)
	return c.SendString(b64Data)
}

