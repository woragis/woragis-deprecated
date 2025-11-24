package systemdesigns

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	translationenricher "github.com/woragis/backend/server/app/pkg/translations"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes system design endpoints.
type Handler interface {
	CreateSystemDesign(c *fiber.Ctx) error
	UpdateSystemDesign(c *fiber.Ctx) error
	GetSystemDesign(c *fiber.Ctx) error
	GetSystemDesignPublic(c *fiber.Ctx) error
	ListSystemDesigns(c *fiber.Ctx) error
	ListFeaturedSystemDesigns(c *fiber.Ctx) error
	DeleteSystemDesign(c *fiber.Ctx) error
}

type handler struct {
	service           Service
	enricher          *translationenricher.Enricher
	translationService translationsdomain.Service
	logger            *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a system design handler.
func NewHandler(service Service, enricher *translationenricher.Enricher, translationService translationsdomain.Service, logger *slog.Logger) Handler {
	return &handler{
		service:           service,
		enricher:          enricher,
		translationService: translationService,
		logger:            logger,
	}
}

// Payloads

type createSystemDesignPayload struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Components  *ComponentsData `json:"components,omitempty"`
	DataFlow    string          `json:"dataFlow,omitempty"`
	Scalability string          `json:"scalability,omitempty"`
	Reliability string          `json:"reliability,omitempty"`
	Diagram     string          `json:"diagram,omitempty"`
	Featured    bool            `json:"featured"`
}

type updateSystemDesignPayload struct {
	Title       *string         `json:"title"`
	Description *string         `json:"description"`
	Components  *ComponentsData `json:"components,omitempty"`
	DataFlow    *string         `json:"dataFlow,omitempty"`
	Scalability *string         `json:"scalability,omitempty"`
	Reliability *string         `json:"reliability,omitempty"`
	Diagram     *string         `json:"diagram,omitempty"`
	Featured    *bool           `json:"featured"`
}

// Responses

type systemDesignResponse struct {
	ID          string          `json:"id"`
	UserID      string          `json:"userId"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Components  *ComponentsData `json:"components,omitempty"`
	DataFlow    string          `json:"dataFlow,omitempty"`
	Scalability string          `json:"scalability,omitempty"`
	Reliability string          `json:"reliability,omitempty"`
	Diagram     string          `json:"diagram,omitempty"`
	Featured    bool            `json:"featured"`
	CreatedAt   string          `json:"createdAt"`
	UpdatedAt   string          `json:"updatedAt"`
}

func toSystemDesignResponse(sd *SystemDesign) systemDesignResponse {
	return systemDesignResponse{
		ID:          sd.ID.String(),
		UserID:      sd.UserID.String(),
		Title:       sd.Title,
		Description: sd.Description,
		Components:  sd.Components,
		DataFlow:    sd.DataFlow,
		Scalability: sd.Scalability,
		Reliability: sd.Reliability,
		Diagram:     sd.Diagram,
		Featured:    sd.Featured,
		CreatedAt:   sd.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   sd.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Handlers

func (h *handler) CreateSystemDesign(c *fiber.Ctx) error {
	var payload createSystemDesignPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	systemDesign, err := h.service.CreateSystemDesign(c.Context(), CreateSystemDesignRequest{
		UserID:      userID,
		Title:       payload.Title,
		Description: payload.Description,
		Components:  payload.Components,
		DataFlow:    payload.DataFlow,
		Scalability: payload.Scalability,
		Reliability: payload.Reliability,
		Diagram:     payload.Diagram,
		Featured:    payload.Featured,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	// Automatically trigger translations for all supported languages
	if h.translationService != nil {
		sourceText := make(map[string]string)
		fields := []string{}

		if systemDesign.Title != "" {
			sourceText["title"] = systemDesign.Title
			fields = append(fields, "title")
		}
		if systemDesign.Description != "" {
			sourceText["description"] = systemDesign.Description
			fields = append(fields, "description")
		}
		if systemDesign.DataFlow != "" {
			sourceText["dataFlow"] = systemDesign.DataFlow
			fields = append(fields, "dataFlow")
		}
		if systemDesign.Scalability != "" {
			sourceText["scalability"] = systemDesign.Scalability
			fields = append(fields, "scalability")
		}
		if systemDesign.Reliability != "" {
			sourceText["reliability"] = systemDesign.Reliability
			fields = append(fields, "reliability")
		}

		supportedLanguages := []translationsdomain.Language{
			translationsdomain.LanguagePTBR,
			translationsdomain.LanguageFR,
			translationsdomain.LanguageES,
			translationsdomain.LanguageDE,
			translationsdomain.LanguageRU,
			translationsdomain.LanguageJA,
			translationsdomain.LanguageKO,
			translationsdomain.LanguageZHCN,
			translationsdomain.LanguageEL,
			translationsdomain.LanguageLA,
		}

		go func() {
			ctx := context.Background()
			for _, lang := range supportedLanguages {
				if err := h.translationService.RequestTranslation(
					ctx,
					translationsdomain.EntityTypeSystemDesign,
					systemDesign.ID,
					lang,
					fields,
					sourceText,
				); err != nil {
					h.logger.Warn("Failed to queue translation",
						slog.String("systemDesignId", systemDesign.ID.String()),
						slog.String("language", string(lang)),
						slog.Any("error", err),
					)
				}
			}
		}()
	}

	return response.Success(c, fiber.StatusCreated, toSystemDesignResponse(systemDesign))
}

func (h *handler) UpdateSystemDesign(c *fiber.Ctx) error {
	systemDesignID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateSystemDesignPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	systemDesign, err := h.service.UpdateSystemDesign(c.Context(), UpdateSystemDesignRequest{
		SystemDesignID: systemDesignID,
		UserID:          userID,
		Title:           payload.Title,
		Description:     payload.Description,
		Components:      payload.Components,
		DataFlow:        payload.DataFlow,
		Scalability:     payload.Scalability,
		Reliability:     payload.Reliability,
		Diagram:         payload.Diagram,
		Featured:        payload.Featured,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	// Trigger translations for updated fields
	if h.translationService != nil && (payload.Title != nil || payload.Description != nil || payload.DataFlow != nil || payload.Scalability != nil || payload.Reliability != nil) {
		sourceText := make(map[string]string)
		fields := []string{}

		if payload.Title != nil && *payload.Title != "" {
			sourceText["title"] = *payload.Title
			fields = append(fields, "title")
		}
		if payload.Description != nil && *payload.Description != "" {
			sourceText["description"] = *payload.Description
			fields = append(fields, "description")
		}
		if payload.DataFlow != nil && *payload.DataFlow != "" {
			sourceText["dataFlow"] = *payload.DataFlow
			fields = append(fields, "dataFlow")
		}
		if payload.Scalability != nil && *payload.Scalability != "" {
			sourceText["scalability"] = *payload.Scalability
			fields = append(fields, "scalability")
		}
		if payload.Reliability != nil && *payload.Reliability != "" {
			sourceText["reliability"] = *payload.Reliability
			fields = append(fields, "reliability")
		}

		if len(fields) > 0 {
			supportedLanguages := []translationsdomain.Language{
				translationsdomain.LanguagePTBR,
				translationsdomain.LanguageFR,
				translationsdomain.LanguageES,
				translationsdomain.LanguageDE,
				translationsdomain.LanguageRU,
				translationsdomain.LanguageJA,
				translationsdomain.LanguageKO,
				translationsdomain.LanguageZHCN,
				translationsdomain.LanguageEL,
				translationsdomain.LanguageLA,
			}

			go func() {
				ctx := context.Background()
				for _, lang := range supportedLanguages {
					if err := h.translationService.RequestTranslation(
						ctx,
						translationsdomain.EntityTypeSystemDesign,
						systemDesign.ID,
						lang,
						fields,
						sourceText,
					); err != nil {
						h.logger.Warn("Failed to queue translation",
							slog.String("systemDesignId", systemDesign.ID.String()),
							slog.String("language", string(lang)),
							slog.Any("error", err),
						)
					}
				}
			}()
		}
	}

	return response.Success(c, fiber.StatusOK, toSystemDesignResponse(systemDesign))
}

func (h *handler) GetSystemDesign(c *fiber.Ctx) error {
	systemDesignID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	systemDesign, err := h.service.GetSystemDesign(c.Context(), systemDesignID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"title":       &systemDesign.Title,
			"description": &systemDesign.Description,
			"dataFlow":    &systemDesign.DataFlow,
			"scalability": &systemDesign.Scalability,
			"reliability": &systemDesign.Reliability,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSystemDesign, systemDesign.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, toSystemDesignResponse(systemDesign))
}

func (h *handler) GetSystemDesignPublic(c *fiber.Ctx) error {
	systemDesignID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	systemDesign, err := h.service.GetSystemDesignPublic(c.Context(), systemDesignID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"title":       &systemDesign.Title,
			"description": &systemDesign.Description,
			"dataFlow":    &systemDesign.DataFlow,
			"scalability": &systemDesign.Scalability,
			"reliability": &systemDesign.Reliability,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSystemDesign, systemDesign.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, toSystemDesignResponse(systemDesign))
}

func (h *handler) ListSystemDesigns(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	systemDesigns, err := h.service.ListSystemDesigns(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range systemDesigns {
			fieldMap := map[string]*string{
				"title":       &systemDesigns[i].Title,
				"description": &systemDesigns[i].Description,
				"dataFlow":    &systemDesigns[i].DataFlow,
				"scalability": &systemDesigns[i].Scalability,
				"reliability": &systemDesigns[i].Reliability,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSystemDesign, systemDesigns[i].ID, language, fieldMap)
		}
	}

	resp := make([]systemDesignResponse, 0, len(systemDesigns))
	for _, sd := range systemDesigns {
		resp = append(resp, toSystemDesignResponse(&sd))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) ListFeaturedSystemDesigns(c *fiber.Ctx) error {
	systemDesigns, err := h.service.ListFeaturedSystemDesigns(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range systemDesigns {
			fieldMap := map[string]*string{
				"title":       &systemDesigns[i].Title,
				"description": &systemDesigns[i].Description,
				"dataFlow":    &systemDesigns[i].DataFlow,
				"scalability": &systemDesigns[i].Scalability,
				"reliability": &systemDesigns[i].Reliability,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeSystemDesign, systemDesigns[i].ID, language, fieldMap)
		}
	}

	resp := make([]systemDesignResponse, 0, len(systemDesigns))
	for _, sd := range systemDesigns {
		resp = append(resp, toSystemDesignResponse(&sd))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) DeleteSystemDesign(c *fiber.Ctx) error {
	systemDesignID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	if err := h.service.DeleteSystemDesign(c.Context(), DeleteSystemDesignRequest{
		SystemDesignID: systemDesignID,
		UserID:          userID,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"id": systemDesignID.String()})
}

// Helper functions

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		switch domainErr.Code {
		case ErrCodeNotFound:
			return response.Error(c, fiber.StatusNotFound, domainErr.Code, nil)
		case ErrCodeUnauthorized:
			return response.Error(c, fiber.StatusUnauthorized, domainErr.Code, nil)
		case ErrCodeConflict:
			return response.Error(c, fiber.StatusConflict, domainErr.Code, nil)
		default:
			return response.Error(c, fiber.StatusBadRequest, domainErr.Code, nil)
		}
	}
	h.logger.Error("Unexpected error", slog.Any("error", err))
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func unauthorizedResponse(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, nil)
}

