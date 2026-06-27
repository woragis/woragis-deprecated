package interests

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"woragis-skills-service/pkg/middleware"
	"woragis-skills-service/pkg/response"
)

// Handler exposes interest endpoints.
type Handler interface {
	CreateInterest(c *fiber.Ctx) error
	UpdateInterest(c *fiber.Ctx) error
	DeleteInterest(c *fiber.Ctx) error
	GetInterest(c *fiber.Ctx) error
	GetInterestBySlug(c *fiber.Ctx) error
	ListInterests(c *fiber.Ctx) error
	ListFeaturedInterests(c *fiber.Ctx) error
	SearchInterests(c *fiber.Ctx) error
}

type handler struct {
	service           Service
	enricher          interface{} // *translationenricher.Enricher - TODO: Re-enable when translation service is implemented
	translationService interface{} // translationsdomain.Service - TODO: Re-enable when translation service is implemented
	logger            *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs an interest handler.
func NewHandler(service Service, enricher interface{}, translationService interface{}, logger *slog.Logger) Handler {
	return &handler{
		service:           service,
		enricher:          enricher,
		translationService: translationService,
		logger:            logger,
	}
}

// Payloads

type createInterestPayload struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	Icon             string `json:"icon,omitempty"`
	Color            string `json:"color,omitempty"`
	BgGradient       string `json:"bgGradient,omitempty"`
	BorderColor      string `json:"borderColor,omitempty"`
	HoverBorderColor string `json:"hoverBorderColor,omitempty"`
	ShadowColor      string `json:"shadowColor,omitempty"`
	FullWidth        bool   `json:"fullWidth"`
	Featured         bool   `json:"featured"`
}

type updateInterestPayload struct {
	Title            string `json:"title,omitempty"`
	Description      string `json:"description,omitempty"`
	Icon             string `json:"icon,omitempty"`
	Color            string `json:"color,omitempty"`
	BgGradient       string `json:"bgGradient,omitempty"`
	BorderColor      string `json:"borderColor,omitempty"`
	HoverBorderColor string `json:"hoverBorderColor,omitempty"`
	ShadowColor      string `json:"shadowColor,omitempty"`
	FullWidth        *bool  `json:"fullWidth,omitempty"`
	Featured         *bool  `json:"featured,omitempty"`
}

// Handlers

func (h *handler) CreateInterest(c *fiber.Ctx) error {
	_, err := middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var payload createInterestPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	interest, err := h.service.CreateInterest(c.Context(), CreateInterestRequest{
		Title:            payload.Title,
		Description:      payload.Description,
		Icon:             payload.Icon,
		Color:            payload.Color,
		BgGradient:       payload.BgGradient,
		BorderColor:      payload.BorderColor,
		HoverBorderColor: payload.HoverBorderColor,
		ShadowColor:      payload.ShadowColor,
		FullWidth:        payload.FullWidth,
		Featured:         payload.Featured,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toInterestResponse(interest))
}

func (h *handler) UpdateInterest(c *fiber.Ctx) error {
	_, err := middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	interestID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateInterestPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	interest, err := h.service.UpdateInterest(c.Context(), UpdateInterestRequest{
		InterestID:       interestID,
		Title:            payload.Title,
		Description:      payload.Description,
		Icon:             payload.Icon,
		Color:            payload.Color,
		BgGradient:       payload.BgGradient,
		BorderColor:      payload.BorderColor,
		HoverBorderColor: payload.HoverBorderColor,
		ShadowColor:      payload.ShadowColor,
		FullWidth:        payload.FullWidth,
		Featured:         payload.Featured,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toInterestResponse(interest))
}

func (h *handler) DeleteInterest(c *fiber.Ctx) error {
	interestIDStr := c.Params("id")
	if interestIDStr == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	interestID, err := uuid.Parse(interestIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	_, err = middleware.GetUserIDFromFiberContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 0, nil)
	}

	if err := h.service.DeleteInterest(c.Context(), interestID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, map[string]string{"message": "Interest deleted successfully"})
}

func (h *handler) GetInterest(c *fiber.Ctx) error {
	interestID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	interest, err := h.service.GetInterest(c.Context(), interestID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment
	// TODO: Re-enable when translation service is implemented
	// language := translationsdomain.LanguageFromContext(c)
	_ = c // Avoid unused variable
	// TODO: Re-enable when translation service is implemented
	// fieldMap := map[string]*string{
	// 	"title":       &interest.Title,
	// 	"description": &interest.Description,
	// }
	// if err := h.enricher.EnrichEntityFields(
	// 	c.Context(),
	// 	translationsdomain.EntityTypeInterest,
	// 	interest.ID,
	// 	language,
	// 	fieldMap,
	// 		// ); err != nil {
		// 	h.logger.Warn("Failed to enrich interest with translations", slog.Any("error", err))
		// }

	return response.Success(c, fiber.StatusOK, toInterestResponse(interest))
}

func (h *handler) GetInterestBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	interest, err := h.service.GetInterestBySlug(c.Context(), slug)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment
	// TODO: Re-enable when translation service is implemented
	// language := translationsdomain.LanguageFromContext(c)
	_ = c // Avoid unused variable
	// TODO: Re-enable when translation service is implemented
	// fieldMap := map[string]*string{
	// 	"title":       &interest.Title,
	// 	"description": &interest.Description,
	// }
	// if err := h.enricher.EnrichEntityFields(
	// 	c.Context(),
	// 	translationsdomain.EntityTypeInterest,
	// 	interest.ID,
	// 	language,
	// 	fieldMap,
	// 		// ); err != nil {
		// 	h.logger.Warn("Failed to enrich interest with translations", slog.Any("error", err))
		// }

	return response.Success(c, fiber.StatusOK, toInterestResponse(interest))
}

func (h *handler) ListInterests(c *fiber.Ctx) error {
	interests, err := h.service.ListInterests(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment to each interest
	// TODO: Re-enable when translation service is implemented
	// language := translationsdomain.LanguageFromContext(c)
	_ = c // Avoid unused variable
	for range interests {
		// TODO: Re-enable when translation service is implemented
		// fieldMap := map[string]*string{
		// 	"title":       &interests[i].Title,
		// 	"description": &interests[i].Description,
		// }
		// TODO: Re-enable when translation service is implemented
		// TODO: Re-enable when translation service is implemented
		// if err := h.enricher.EnrichEntityFields(
		// 	c.Context(),
		// 	translationsdomain.EntityTypeInterest,
		// 	interests[i].ID,
		// 	language,
		// 	fieldMap,
		// ); err != nil {
		// 	h.logger.Warn("Failed to enrich interest with translations",
		// 		slog.String("interestId", interests[i].ID.String()),
		// 		slog.Any("error", err),
		// 	)
		// }
	}

	responses := make([]interestResponse, len(interests))
	for i := range interests {
		responses[i] = toInterestResponse(&interests[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) ListFeaturedInterests(c *fiber.Ctx) error {
	interests, err := h.service.ListFeaturedInterests(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment to each interest
	// TODO: Re-enable when translation service is implemented
	// language := translationsdomain.LanguageFromContext(c)
	_ = c // Avoid unused variable
	for range interests {
		// TODO: Re-enable when translation service is implemented
		// fieldMap := map[string]*string{
		// 	"title":       &interests[i].Title,
		// 	"description": &interests[i].Description,
		// }
		// TODO: Re-enable when translation service is implemented
		// TODO: Re-enable when translation service is implemented
		// if err := h.enricher.EnrichEntityFields(
		// 	c.Context(),
		// 	translationsdomain.EntityTypeInterest,
		// 	interests[i].ID,
		// 	language,
		// 	fieldMap,
		// ); err != nil {
		// 	h.logger.Warn("Failed to enrich interest with translations",
		// 		slog.String("interestId", interests[i].ID.String()),
		// 		slog.Any("error", err),
		// 	)
		// }
	}

	responses := make([]interestResponse, len(interests))
	for i := range interests {
		responses[i] = toInterestResponse(&interests[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) SearchInterests(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	interests, err := h.service.SearchInterests(c.Context(), query)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translation enrichment to each interest
	// TODO: Re-enable when translation service is implemented
	// language := translationsdomain.LanguageFromContext(c)
	_ = c // Avoid unused variable
	for range interests {
		// TODO: Re-enable when translation service is implemented
		// fieldMap := map[string]*string{
		// 	"title":       &interests[i].Title,
		// 	"description": &interests[i].Description,
		// }
		// TODO: Re-enable when translation service is implemented
		// TODO: Re-enable when translation service is implemented
		// if err := h.enricher.EnrichEntityFields(
		// 	c.Context(),
		// 	translationsdomain.EntityTypeInterest,
		// 	interests[i].ID,
		// 	language,
		// 	fieldMap,
		// ); err != nil {
		// 	h.logger.Warn("Failed to enrich interest with translations",
		// 		slog.String("interestId", interests[i].ID.String()),
		// 		slog.Any("error", err),
		// 	)
		// }
	}

	responses := make([]interestResponse, len(interests))
	for i := range interests {
		responses[i] = toInterestResponse(&interests[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

// Response helpers

type interestResponse struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	Description     string `json:"description"`
	Icon            string `json:"icon,omitempty"`
	Color           string `json:"color,omitempty"`
	BgGradient      string `json:"bgGradient,omitempty"`
	BorderColor     string `json:"borderColor,omitempty"`
	HoverBorderColor string `json:"hoverBorderColor,omitempty"`
	ShadowColor     string `json:"shadowColor,omitempty"`
	FullWidth       bool   `json:"fullWidth"`
	Featured        bool   `json:"featured"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

func toInterestResponse(interest *Interest) interestResponse {
	return interestResponse{
		ID:              interest.ID.String(),
		Title:           interest.Title,
		Slug:            interest.Slug,
		Description:     interest.Description,
		Icon:            interest.Icon,
		Color:           interest.Color,
		BgGradient:      interest.BgGradient,
		BorderColor:     interest.BorderColor,
		HoverBorderColor: interest.HoverBorderColor,
		ShadowColor:     interest.ShadowColor,
		FullWidth:       interest.FullWidth,
		Featured:        interest.Featured,
		CreatedAt:       interest.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       interest.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
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
	case ErrCodeInvalidPayload, ErrCodeInvalidTitle:
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

