package problemsolutions

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

// Handler exposes problem solution endpoints.
type Handler interface {
	CreateProblemSolution(c *fiber.Ctx) error
	UpdateProblemSolution(c *fiber.Ctx) error
	GetProblemSolution(c *fiber.Ctx) error
	GetProblemSolutionPublic(c *fiber.Ctx) error
	ListProblemSolutions(c *fiber.Ctx) error
	ListFeaturedProblemSolutions(c *fiber.Ctx) error
	DeleteProblemSolution(c *fiber.Ctx) error
	GetProblemSolutionMatrix(c *fiber.Ctx) error
}

type handler struct {
	service           Service
	enricher          *translationenricher.Enricher
	translationService translationsdomain.Service
	logger            *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a problem solution handler.
func NewHandler(service Service, enricher *translationenricher.Enricher, translationService translationsdomain.Service, logger *slog.Logger) Handler {
	return &handler{
		service:           service,
		enricher:          enricher,
		translationService: translationService,
		logger:            logger,
	}
}

// Payloads

type createProblemSolutionPayload struct {
	Problem     string       `json:"problem"`
	Context     string       `json:"context"`
	Solution    string       `json:"solution"`
	Technologies []string    `json:"technologies"`
	Impact      string       `json:"impact"`
	Metrics     *MetricsData `json:"metrics,omitempty"`
	Featured    bool         `json:"featured"`
}

type updateProblemSolutionPayload struct {
	Problem     *string      `json:"problem"`
	Context     *string      `json:"context"`
	Solution    *string      `json:"solution"`
	Technologies []string     `json:"technologies"`
	Impact      *string       `json:"impact"`
	Metrics     *MetricsData `json:"metrics,omitempty"`
	Featured    *bool        `json:"featured"`
}

// Responses

type problemSolutionResponse struct {
	ID          string       `json:"id"`
	UserID      string       `json:"userId"`
	Problem     string       `json:"problem"`
	Context     string       `json:"context"`
	Solution    string       `json:"solution"`
	Technologies []string     `json:"technologies"`
	Impact      string       `json:"impact"`
	Metrics     *MetricsData `json:"metrics,omitempty"`
	Featured    bool         `json:"featured"`
	CreatedAt   string       `json:"createdAt"`
	UpdatedAt   string       `json:"updatedAt"`
}

func toProblemSolutionResponse(ps *ProblemSolution) problemSolutionResponse {
	return problemSolutionResponse{
		ID:          ps.ID.String(),
		UserID:      ps.UserID.String(),
		Problem:     ps.Problem,
		Context:     ps.Context,
		Solution:    ps.Solution,
		Technologies: []string(ps.Technologies),
		Impact:      ps.Impact,
		Metrics:     ps.Metrics,
		Featured:    ps.Featured,
		CreatedAt:   ps.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   ps.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// Handlers

func (h *handler) CreateProblemSolution(c *fiber.Ctx) error {
	var payload createProblemSolutionPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	problemSolution, err := h.service.CreateProblemSolution(c.Context(), CreateProblemSolutionRequest{
		UserID:      userID,
		Problem:     payload.Problem,
		Context:     payload.Context,
		Solution:    payload.Solution,
		Technologies: payload.Technologies,
		Impact:      payload.Impact,
		Metrics:     payload.Metrics,
		Featured:    payload.Featured,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	// Automatically trigger translations for all supported languages
	if h.translationService != nil {
		sourceText := make(map[string]string)
		fields := []string{}

		if problemSolution.Problem != "" {
			sourceText["problem"] = problemSolution.Problem
			fields = append(fields, "problem")
		}
		if problemSolution.Context != "" {
			sourceText["context"] = problemSolution.Context
			fields = append(fields, "context")
		}
		if problemSolution.Solution != "" {
			sourceText["solution"] = problemSolution.Solution
			fields = append(fields, "solution")
		}
		if problemSolution.Impact != "" {
			sourceText["impact"] = problemSolution.Impact
			fields = append(fields, "impact")
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
					translationsdomain.EntityTypeProblemSolution,
					problemSolution.ID,
					lang,
					fields,
					sourceText,
				); err != nil {
					h.logger.Warn("Failed to queue translation",
						slog.String("problemSolutionId", problemSolution.ID.String()),
						slog.String("language", string(lang)),
						slog.Any("error", err),
					)
				}
			}
		}()
	}

	return response.Success(c, fiber.StatusCreated, toProblemSolutionResponse(problemSolution))
}

func (h *handler) UpdateProblemSolution(c *fiber.Ctx) error {
	problemSolutionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateProblemSolutionPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	problemSolution, err := h.service.UpdateProblemSolution(c.Context(), UpdateProblemSolutionRequest{
		ProblemSolutionID: problemSolutionID,
		UserID:             userID,
		Problem:             payload.Problem,
		Context:            payload.Context,
		Solution:           payload.Solution,
		Technologies:       payload.Technologies,
		Impact:             payload.Impact,
		Metrics:            payload.Metrics,
		Featured:           payload.Featured,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	// Trigger translations for updated fields
	if h.translationService != nil && (payload.Problem != nil || payload.Context != nil || payload.Solution != nil || payload.Impact != nil) {
		sourceText := make(map[string]string)
		fields := []string{}

		if payload.Problem != nil && *payload.Problem != "" {
			sourceText["problem"] = *payload.Problem
			fields = append(fields, "problem")
		}
		if payload.Context != nil && *payload.Context != "" {
			sourceText["context"] = *payload.Context
			fields = append(fields, "context")
		}
		if payload.Solution != nil && *payload.Solution != "" {
			sourceText["solution"] = *payload.Solution
			fields = append(fields, "solution")
		}
		if payload.Impact != nil && *payload.Impact != "" {
			sourceText["impact"] = *payload.Impact
			fields = append(fields, "impact")
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
						translationsdomain.EntityTypeProblemSolution,
						problemSolution.ID,
						lang,
						fields,
						sourceText,
					); err != nil {
						h.logger.Warn("Failed to queue translation",
							slog.String("problemSolutionId", problemSolution.ID.String()),
							slog.String("language", string(lang)),
							slog.Any("error", err),
						)
					}
				}
			}()
		}
	}

	return response.Success(c, fiber.StatusOK, toProblemSolutionResponse(problemSolution))
}

func (h *handler) GetProblemSolution(c *fiber.Ctx) error {
	problemSolutionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	problemSolution, err := h.service.GetProblemSolution(c.Context(), problemSolutionID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"problem":  &problemSolution.Problem,
			"context":  &problemSolution.Context,
			"solution": &problemSolution.Solution,
			"impact":   &problemSolution.Impact,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeProblemSolution, problemSolution.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, toProblemSolutionResponse(problemSolution))
}

func (h *handler) GetProblemSolutionPublic(c *fiber.Ctx) error {
	problemSolutionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	problemSolution, err := h.service.GetProblemSolutionPublic(c.Context(), problemSolutionID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"problem":  &problemSolution.Problem,
			"context":  &problemSolution.Context,
			"solution": &problemSolution.Solution,
			"impact":   &problemSolution.Impact,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeProblemSolution, problemSolution.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, toProblemSolutionResponse(problemSolution))
}

func (h *handler) ListProblemSolutions(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	problemSolutions, err := h.service.ListProblemSolutions(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range problemSolutions {
			fieldMap := map[string]*string{
				"problem":  &problemSolutions[i].Problem,
				"context":  &problemSolutions[i].Context,
				"solution": &problemSolutions[i].Solution,
				"impact":   &problemSolutions[i].Impact,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeProblemSolution, problemSolutions[i].ID, language, fieldMap)
		}
	}

	resp := make([]problemSolutionResponse, 0, len(problemSolutions))
	for _, ps := range problemSolutions {
		resp = append(resp, toProblemSolutionResponse(&ps))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) ListFeaturedProblemSolutions(c *fiber.Ctx) error {
	problemSolutions, err := h.service.ListFeaturedProblemSolutions(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range problemSolutions {
			fieldMap := map[string]*string{
				"problem":  &problemSolutions[i].Problem,
				"context":  &problemSolutions[i].Context,
				"solution": &problemSolutions[i].Solution,
				"impact":   &problemSolutions[i].Impact,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeProblemSolution, problemSolutions[i].ID, language, fieldMap)
		}
	}

	resp := make([]problemSolutionResponse, 0, len(problemSolutions))
	for _, ps := range problemSolutions {
		resp = append(resp, toProblemSolutionResponse(&ps))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) DeleteProblemSolution(c *fiber.Ctx) error {
	problemSolutionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	if err := h.service.DeleteProblemSolution(c.Context(), DeleteProblemSolutionRequest{
		ProblemSolutionID: problemSolutionID,
		UserID:             userID,
	}); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"id": problemSolutionID.String()})
}

func (h *handler) GetProblemSolutionMatrix(c *fiber.Ctx) error {
	matrix, err := h.service.GetProblemSolutionMatrix(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, matrix)
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

