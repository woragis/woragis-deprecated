package languages

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler wires services to HTTP endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler builds a new Handler instance.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type studySessionPayload struct {
	UserID       string `json:"user_id"`
	LanguageCode string `json:"language_code"`
	SkillFocus   string `json:"skill_focus"`
	DurationMin  int    `json:"duration_min"`
	Notes        string `json:"notes"`
	CompletedAt  string `json:"completed_at"`
}

type vocabularyPayload struct {
	UserID       string `json:"user_id"`
	LanguageCode string `json:"language_code"`
	Term         string `json:"term"`
	Translation  string `json:"translation"`
	Context      string `json:"context"`
	ReviewAt     string `json:"review_at"`
}

type sessionsQuery struct {
	UserID   string `query:"user_id"`
	Language string `query:"language"`
	From     string `query:"from"`
	To       string `query:"to"`
}

type vocabularyQuery struct {
	UserID   string `query:"user_id"`
	Language string `query:"language"`
	DueOnly  string `query:"due_only"`
}

type summaryQuery struct {
	UserID string `query:"user_id"`
}

type studySessionResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	LanguageCode string    `json:"language_code"`
	SkillFocus   string    `json:"skill_focus"`
	DurationMin  int       `json:"duration_min"`
	Notes        string    `json:"notes,omitempty"`
	CompletedAt  time.Time `json:"completed_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type vocabularyResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	LanguageCode string    `json:"language_code"`
	Term         string    `json:"term"`
	Translation  string    `json:"translation"`
	Context      string    `json:"context,omitempty"`
	AddedAt      time.Time `json:"added_at"`
	ReviewAt     time.Time `json:"review_at"`
}

type summaryResponse struct {
	LanguageCode    string `json:"language_code"`
	TotalMinutes    int64  `json:"total_minutes"`
	SessionCount    int64  `json:"session_count"`
	VocabularyCount int64  `json:"vocabulary_count"`
}

// PostStudySession handles session logging.
func (h *Handler) PostStudySession(c *fiber.Ctx) error {
	var payload studySessionPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var completedAt time.Time
	if payload.CompletedAt != "" {
		if completedAt, err = time.Parse(time.RFC3339, payload.CompletedAt); err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	session, err := h.service.LogStudySession(c.Context(), StudySessionRequest{
		UserID:       userID,
		LanguageCode: payload.LanguageCode,
		SkillFocus:   payload.SkillFocus,
		DurationMin:  payload.DurationMin,
		Notes:        payload.Notes,
		CompletedAt:  completedAt,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, studySessionResponse{
		ID:           session.ID.String(),
		UserID:       session.UserID.String(),
		LanguageCode: session.LanguageCode,
		SkillFocus:   session.SkillFocus,
		DurationMin:  session.DurationMin,
		Notes:        session.Notes,
		CompletedAt:  session.CompletedAt,
		CreatedAt:    session.CreatedAt,
	})
}

// GetStudySessions returns sessions for a user.
func (h *Handler) GetStudySessions(c *fiber.Ctx) error {
	var query sessionsQuery
	if err := c.QueryParser(&query); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	from, to, err := parseRange(query.From, query.To)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	sessions, err := h.service.GetStudySessions(c.Context(), userID, query.Language, from, to)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]studySessionResponse, 0, len(sessions))
	for _, session := range sessions {
		resp = append(resp, studySessionResponse{
			ID:           session.ID.String(),
			UserID:       session.UserID.String(),
			LanguageCode: session.LanguageCode,
			SkillFocus:   session.SkillFocus,
			DurationMin:  session.DurationMin,
			Notes:        session.Notes,
			CompletedAt:  session.CompletedAt,
			CreatedAt:    session.CreatedAt,
		})
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// PostVocabulary handles vocabulary registration.
func (h *Handler) PostVocabulary(c *fiber.Ctx) error {
	var payload vocabularyPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var reviewAt time.Time
	if payload.ReviewAt != "" {
		if reviewAt, err = time.Parse(time.RFC3339, payload.ReviewAt); err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	entry, err := h.service.AddVocabularyEntry(c.Context(), VocabularyRequest{
		UserID:       userID,
		LanguageCode: payload.LanguageCode,
		Term:         payload.Term,
		Translation:  payload.Translation,
		Context:      payload.Context,
		ReviewAt:     reviewAt,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, vocabularyResponse{
		ID:           entry.ID.String(),
		UserID:       entry.UserID.String(),
		LanguageCode: entry.LanguageCode,
		Term:         entry.Term,
		Translation:  entry.Translation,
		Context:      entry.Context,
		AddedAt:      entry.AddedAt,
		ReviewAt:     entry.ReviewAt,
	})
}

// GetVocabulary lists vocabulary entries.
func (h *Handler) GetVocabulary(c *fiber.Ctx) error {
	var query vocabularyQuery
	if err := c.QueryParser(&query); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	dueOnly := false
	if query.DueOnly != "" {
		dueOnly, err = strconv.ParseBool(query.DueOnly)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
	}

	entries, err := h.service.GetVocabulary(c.Context(), userID, query.Language, dueOnly)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]vocabularyResponse, 0, len(entries))
	for _, entry := range entries {
		resp = append(resp, vocabularyResponse{
			ID:           entry.ID.String(),
			UserID:       entry.UserID.String(),
			LanguageCode: entry.LanguageCode,
			Term:         entry.Term,
			Translation:  entry.Translation,
			Context:      entry.Context,
			AddedAt:      entry.AddedAt,
			ReviewAt:     entry.ReviewAt,
		})
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// GetSummary aggregates metrics per language.
func (h *Handler) GetSummary(c *fiber.Ctx) error {
	var query summaryQuery
	if err := c.QueryParser(&query); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	summary, err := h.service.GetSummary(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]summaryResponse, 0, len(summary))
	for _, s := range summary {
		resp = append(resp, summaryResponse{
			LanguageCode:    s.LanguageCode,
			TotalMinutes:    s.TotalMinutes,
			SessionCount:    s.SessionCount,
			VocabularyCount: s.VocabularyCount,
		})
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromError(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, nil)
	}

	h.logError("languages: unexpected error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func statusFromError(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidLanguage, ErrCodeInvalidDuration, ErrCodeInvalidCompletedAt, ErrCodeInvalidVocabulary, ErrCodeInvalidReviewAt:
		return fiber.StatusBadRequest
	case ErrCodeRepositoryFailure:
		return fiber.StatusInternalServerError
	case ErrCodeSummaryFailure:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusInternalServerError
	}
}

func parseRange(fromRaw, toRaw string) (time.Time, time.Time, error) {
	var (
		from time.Time
		to   time.Time
		err  error
	)

	if fromRaw != "" {
		if from, err = time.Parse(time.RFC3339, fromRaw); err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if toRaw != "" {
		if to, err = time.Parse(time.RFC3339, toRaw); err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	return from, to, nil
}

func (h *Handler) logWarn(message string) {
	if h.logger != nil {
		h.logger.Warn(message)
	}
}

func (h *Handler) logError(message string, err error) {
	if h.logger != nil {
		h.logger.Error(message, slog.Any("error", err))
	}
}
