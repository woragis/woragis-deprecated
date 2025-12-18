package assets

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes content asset endpoints.
type Handler interface {
	CreateAsset(c *fiber.Ctx) error
	GetAsset(c *fiber.Ctx) error
	GetAssetsByContentPost(c *fiber.Ctx) error
	GetAssetsBySocialPost(c *fiber.Ctx) error
	UpdateAsset(c *fiber.Ctx) error
	DeleteAsset(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs an assets handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type createAssetPayload struct {
	ContentPostID *string `json:"contentPostId,omitempty"`
	SocialPostID   *string `json:"socialPostId,omitempty"`
	AssetType      string  `json:"assetType"`
	FilePath       string  `json:"filePath"`
	FileURL        string  `json:"fileUrl,omitempty"`
	AltText        string  `json:"altText,omitempty"`
}

type updateAssetPayload struct {
	FileURL *string `json:"fileUrl,omitempty"`
	AltText *string `json:"altText,omitempty"`
}

// Handlers

func (h *handler) CreateAsset(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var payload createAssetPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req := CreateAssetRequest{
		AssetType: AssetType(payload.AssetType),
		FilePath:  payload.FilePath,
		FileURL:   payload.FileURL,
		AltText:   payload.AltText,
	}

	if payload.ContentPostID != nil {
		contentPostID, err := uuid.Parse(*payload.ContentPostID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		req.ContentPostID = &contentPostID
	}

	if payload.SocialPostID != nil {
		socialPostID, err := uuid.Parse(*payload.SocialPostID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		req.SocialPostID = &socialPostID
	}

	asset, err := h.service.CreateAsset(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toAssetResponse(asset))
}

func (h *handler) GetAsset(c *fiber.Ctx) error {
	assetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	asset, err := h.service.GetAsset(c.Context(), assetID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toAssetResponse(asset))
}

func (h *handler) GetAssetsByContentPost(c *fiber.Ctx) error {
	contentPostID, err := uuid.Parse(c.Params("contentPostId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	assets, err := h.service.GetAssetsByContentPost(c.Context(), contentPostID)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]assetResponse, len(assets))
	for i := range assets {
		responses[i] = toAssetResponse(&assets[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) GetAssetsBySocialPost(c *fiber.Ctx) error {
	socialPostID, err := uuid.Parse(c.Params("socialPostId"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	assets, err := h.service.GetAssetsBySocialPost(c.Context(), socialPostID)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]assetResponse, len(assets))
	for i := range assets {
		responses[i] = toAssetResponse(&assets[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) UpdateAsset(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	assetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateAssetPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	asset, err := h.service.UpdateAsset(c.Context(), UpdateAssetRequest{
		AssetID: assetID,
		FileURL: payload.FileURL,
		AltText: payload.AltText,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toAssetResponse(asset))
}

func (h *handler) DeleteAsset(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	assetID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.DeleteAsset(c.Context(), assetID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "asset deleted"})
}

// Response helpers

type assetResponse struct {
	ID            string  `json:"id"`
	ContentPostID *string `json:"contentPostId,omitempty"`
	SocialPostID  *string `json:"socialPostId,omitempty"`
	AssetType     string  `json:"assetType"`
	FilePath      string  `json:"filePath"`
	FileURL       string  `json:"fileUrl,omitempty"`
	AltText       string  `json:"altText,omitempty"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

func toAssetResponse(asset *ContentAsset) assetResponse {
	var contentPostID, socialPostID *string
	if asset.ContentPostID != nil {
		idStr := asset.ContentPostID.String()
		contentPostID = &idStr
	}
	if asset.SocialPostID != nil {
		idStr := asset.SocialPostID.String()
		socialPostID = &idStr
	}

	return assetResponse{
		ID:            asset.ID.String(),
		ContentPostID: contentPostID,
		SocialPostID:  socialPostID,
		AssetType:     string(asset.AssetType),
		FilePath:      asset.FilePath,
		FileURL:       asset.FileURL,
		AltText:       asset.AltText,
		CreatedAt:     asset.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     asset.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Error handling

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	domainErr, ok := AsDomainError(err)
	if !ok {
		h.logger.Error("unexpected error", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 0, nil)
	}

	statusCode := fiber.StatusInternalServerError
	switch domainErr.Code {
	case ErrCodeInvalidPayload, ErrCodeInvalidAssetType:
		statusCode = fiber.StatusBadRequest
	case ErrCodeNotFound:
		statusCode = fiber.StatusNotFound
	case ErrCodeConflict:
		statusCode = fiber.StatusConflict
	case ErrCodeRepositoryFailure:
		statusCode = fiber.StatusInternalServerError
	}

	return response.Error(c, statusCode, domainErr.Code, fiber.Map{"message": domainErr.Message})
}

func unauthorizedResponse(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusUnauthorized, 0, fiber.Map{"message": "unauthorized"})
}
