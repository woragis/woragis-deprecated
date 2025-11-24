package projectcasestudies

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	translationenricher "github.com/woragis/backend/server/app/pkg/translations"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes project case study endpoints.
type Handler interface {
	CreateCaseStudy(c *fiber.Ctx) error
	UpdateCaseStudy(c *fiber.Ctx) error
	GetCaseStudy(c *fiber.Ctx) error
	GetCaseStudyByProjectID(c *fiber.Ctx) error
	GetPublicCaseStudyByProjectID(c *fiber.Ctx) error
	ListCaseStudies(c *fiber.Ctx) error
	DeleteCaseStudy(c *fiber.Ctx) error
}

type handler struct {
	service           Service
	enricher          *translationenricher.Enricher
	translationService translationsdomain.Service
	logger            *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a project case study handler.
func NewHandler(service Service, enricher *translationenricher.Enricher, translationService translationsdomain.Service, logger *slog.Logger) Handler {
	return &handler{
		service:           service,
		enricher:          enricher,
		translationService: translationService,
		logger:            logger,
	}
}

// Payloads

type createCaseStudyPayload struct {
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Challenge     string         `json:"challenge"`
	Solution      string         `json:"solution"`
	Technologies  []string       `json:"technologies"`
	Architecture  string         `json:"architecture"`
	Metrics       *MetricsData    `json:"metrics,omitempty"`
	Tradeoffs     *TradeoffsData `json:"tradeoffs,omitempty"`
	LessonsLearned []string      `json:"lessonsLearned"`
}

type updateCaseStudyPayload struct {
	Title         *string        `json:"title"`
	Description   *string        `json:"description"`
	Challenge     *string        `json:"challenge"`
	Solution      *string        `json:"solution"`
	Technologies  []string       `json:"technologies"`
	Architecture  *string        `json:"architecture"`
	Metrics       *MetricsData    `json:"metrics,omitempty"`
	Tradeoffs     *TradeoffsData `json:"tradeoffs,omitempty"`
	LessonsLearned []string      `json:"lessonsLearned"`
}

// Responses

type caseStudyResponse struct {
	ID            string         `json:"id"`
	ProjectID     string         `json:"projectId"`
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Challenge     string         `json:"challenge"`
	Solution      string         `json:"solution"`
	Technologies  []string        `json:"technologies"`
	Architecture  string          `json:"architecture"`
	Metrics       *MetricsData    `json:"metrics,omitempty"`
	Tradeoffs     *TradeoffsData `json:"tradeoffs,omitempty"`
	LessonsLearned []string      `json:"lessonsLearned"`
	CreatedAt     string         `json:"createdAt"`
	UpdatedAt     string         `json:"updatedAt"`
}

func toCaseStudyResponse(cs *ProjectCaseStudy) caseStudyResponse {
	return caseStudyResponse{
		ID:            cs.ID.String(),
		ProjectID:     cs.ProjectID.String(),
		Title:         cs.Title,
		Description:   cs.Description,
		Challenge:     cs.Challenge,
		Solution:      cs.Solution,
		Technologies:  []string(cs.Technologies),
		Architecture:  cs.Architecture,
		Metrics:       cs.Metrics,
		Tradeoffs:     cs.Tradeoffs,
		LessonsLearned: []string(cs.LessonsLearned),
		CreatedAt:     cs.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     cs.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Handlers

func (h *handler) CreateCaseStudy(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload createCaseStudyPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	caseStudy, err := h.service.CreateCaseStudy(c.Context(), CreateCaseStudyRequest{
		ProjectID:     projectID,
		UserID:        userID,
		Title:         payload.Title,
		Description:   payload.Description,
		Challenge:     payload.Challenge,
		Solution:      payload.Solution,
		Technologies:  payload.Technologies,
		Architecture:  payload.Architecture,
		Metrics:       payload.Metrics,
		Tradeoffs:     payload.Tradeoffs,
		LessonsLearned: payload.LessonsLearned,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	// Automatically trigger translations for all supported languages
	if h.translationService != nil {
		sourceText := make(map[string]string)
		fields := []string{}

		if caseStudy.Title != "" {
			sourceText["title"] = caseStudy.Title
			fields = append(fields, "title")
		}
		if caseStudy.Description != "" {
			sourceText["description"] = caseStudy.Description
			fields = append(fields, "description")
		}
		if caseStudy.Challenge != "" {
			sourceText["challenge"] = caseStudy.Challenge
			fields = append(fields, "challenge")
		}
		if caseStudy.Solution != "" {
			sourceText["solution"] = caseStudy.Solution
			fields = append(fields, "solution")
		}
		if caseStudy.Architecture != "" {
			sourceText["architecture"] = caseStudy.Architecture
			fields = append(fields, "architecture")
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
					translationsdomain.EntityTypeProjectCaseStudy,
					caseStudy.ID,
					lang,
					fields,
					sourceText,
				); err != nil {
					h.logger.Warn("Failed to queue translation",
						slog.String("caseStudyId", caseStudy.ID.String()),
						slog.String("language", string(lang)),
						slog.Any("error", err),
					)
				}
			}
		}()
	}

	return response.Success(c, fiber.StatusCreated, toCaseStudyResponse(caseStudy))
}

func (h *handler) UpdateCaseStudy(c *fiber.Ctx) error {
	caseStudyID, err := uuid.Parse(c.Params("caseStudyID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateCaseStudyPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	caseStudy, err := h.service.UpdateCaseStudy(c.Context(), UpdateCaseStudyRequest{
		CaseStudyID:   caseStudyID,
		UserID:        userID,
		Title:         payload.Title,
		Description:   payload.Description,
		Challenge:     payload.Challenge,
		Solution:      payload.Solution,
		Technologies:  payload.Technologies,
		Architecture:  payload.Architecture,
		Metrics:       payload.Metrics,
		Tradeoffs:     payload.Tradeoffs,
		LessonsLearned: payload.LessonsLearned,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	// Trigger translations for updated fields
	if h.translationService != nil && (payload.Title != nil || payload.Description != nil || payload.Challenge != nil || payload.Solution != nil || payload.Architecture != nil) {
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
		if payload.Challenge != nil && *payload.Challenge != "" {
			sourceText["challenge"] = *payload.Challenge
			fields = append(fields, "challenge")
		}
		if payload.Solution != nil && *payload.Solution != "" {
			sourceText["solution"] = *payload.Solution
			fields = append(fields, "solution")
		}
		if payload.Architecture != nil && *payload.Architecture != "" {
			sourceText["architecture"] = *payload.Architecture
			fields = append(fields, "architecture")
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
						translationsdomain.EntityTypeProjectCaseStudy,
						caseStudy.ID,
						lang,
						fields,
						sourceText,
					); err != nil {
						h.logger.Warn("Failed to queue translation",
							slog.String("caseStudyId", caseStudy.ID.String()),
							slog.String("language", string(lang)),
							slog.Any("error", err),
						)
					}
				}
			}()
		}
	}

	return response.Success(c, fiber.StatusOK, toCaseStudyResponse(caseStudy))
}

func (h *handler) GetCaseStudy(c *fiber.Ctx) error {
	caseStudyID, err := uuid.Parse(c.Params("caseStudyID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	caseStudy, err := h.service.GetCaseStudy(c.Context(), caseStudyID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"title":       &caseStudy.Title,
			"description": &caseStudy.Description,
			"challenge":   &caseStudy.Challenge,
			"solution":   &caseStudy.Solution,
			"architecture": &caseStudy.Architecture,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeProjectCaseStudy, caseStudy.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, toCaseStudyResponse(caseStudy))
}

func (h *handler) GetCaseStudyByProjectID(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var userID uuid.UUID
	var err2 error

	// Check if request is authenticated via API key
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		userID = apiKey.UserID
	} else {
		userID, err2 = authdomain.UserIDFromContext(c)
		if err2 != nil {
			return unauthorizedResponse(c)
		}
	}

	caseStudy, err := h.service.GetCaseStudyByProjectID(c.Context(), projectID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"title":       &caseStudy.Title,
			"description": &caseStudy.Description,
			"challenge":   &caseStudy.Challenge,
			"solution":   &caseStudy.Solution,
			"architecture": &caseStudy.Architecture,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeProjectCaseStudy, caseStudy.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, toCaseStudyResponse(caseStudy))
}

func (h *handler) GetPublicCaseStudyByProjectID(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	caseStudy, err := h.service.GetPublicCaseStudyByProjectID(c.Context(), projectID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"title":       &caseStudy.Title,
			"description": &caseStudy.Description,
			"challenge":   &caseStudy.Challenge,
			"solution":   &caseStudy.Solution,
			"architecture": &caseStudy.Architecture,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeProjectCaseStudy, caseStudy.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, toCaseStudyResponse(caseStudy))
}

func (h *handler) ListCaseStudies(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	caseStudies, err := h.service.ListCaseStudies(c.Context(), projectID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range caseStudies {
			fieldMap := map[string]*string{
				"title":       &caseStudies[i].Title,
				"description": &caseStudies[i].Description,
				"challenge":   &caseStudies[i].Challenge,
				"solution":   &caseStudies[i].Solution,
				"architecture": &caseStudies[i].Architecture,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeProjectCaseStudy, caseStudies[i].ID, language, fieldMap)
		}
	}

	resp := make([]caseStudyResponse, 0, len(caseStudies))
	for _, cs := range caseStudies {
		resp = append(resp, toCaseStudyResponse(&cs))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) DeleteCaseStudy(c *fiber.Ctx) error {
	caseStudyID, err := uuid.Parse(c.Params("caseStudyID"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	if err := h.service.DeleteCaseStudy(c.Context(), DeleteCaseStudyRequest{
		CaseStudyID: caseStudyID,
		UserID:      userID,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"id": caseStudyID.String()})
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
