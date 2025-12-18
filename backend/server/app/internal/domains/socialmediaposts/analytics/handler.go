package analytics

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes analytics endpoints.
type Handler interface {
	RecordAnalytics(c *fiber.Ctx) error
	GetPostAnalytics(c *fiber.Ctx) error
	GetAnalyticsSummary(c *fiber.Ctx) error
	GetTopPosts(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs an analytics handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type recordAnalyticsPayload struct {
	SocialPostID string  `json:"socialPostId"`
	MetricDate   string  `json:"metricDate"`
	Likes        *int64  `json:"likes,omitempty"`
	Comments     *int64  `json:"comments,omitempty"`
	Shares       *int64  `json:"shares,omitempty"`
	Views        *int64  `json:"views,omitempty"`
	Clicks       *int64  `json:"clicks,omitempty"`
	Reach        *int64  `json:"reach,omitempty"`
	Impressions  *int64  `json:"impressions,omitempty"`
}

// Handlers

func (h *handler) RecordAnalytics(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var payload recordAnalyticsPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	socialPostID, err := uuid.Parse(payload.SocialPostID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	metricDate, err := time.Parse("2006-01-02", payload.MetricDate)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	analytics, err := h.service.RecordAnalytics(c.Context(), RecordAnalyticsRequest{
		SocialPostID: socialPostID,
		MetricDate:   metricDate,
		Likes:        payload.Likes,
		Comments:     payload.Comments,
		Shares:       payload.Shares,
		Views:        payload.Views,
		Clicks:       payload.Clicks,
		Reach:        payload.Reach,
		Impressions:  payload.Impressions,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toAnalyticsResponse(analytics))
}

func (h *handler) GetPostAnalytics(c *fiber.Ctx) error {
	socialPostID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var startDate, endDate *time.Time
	if startDateStr := c.Query("startDate"); startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			startDate = &parsed
		}
	}
	if endDateStr := c.Query("endDate"); endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			endDate = &parsed
		}
	}

	analytics, err := h.service.GetPostAnalytics(c.Context(), socialPostID, startDate, endDate)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]analyticsResponse, len(analytics))
	for i := range analytics {
		responses[i] = toAnalyticsResponse(&analytics[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) GetAnalyticsSummary(c *fiber.Ctx) error {
	filters := AnalyticsFilters{}

	if postIDStr := c.Query("socialPostId"); postIDStr != "" {
		postID, err := uuid.Parse(postIDStr)
		if err == nil {
			filters.SocialPostID = &postID
		}
	}

	if startDateStr := c.Query("startDate"); startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			filters.StartDate = &parsed
		}
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			filters.EndDate = &parsed
		}
	}

	summary, err := h.service.GetAnalyticsSummary(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, summary)
}

func (h *handler) GetTopPosts(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	metric := c.Query("metric", "likes")

	topPosts, err := h.service.GetTopPosts(c.Context(), limit, metric)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, topPosts)
}

// Response helpers

type analyticsResponse struct {
	ID             string   `json:"id"`
	SocialPostID   string   `json:"socialPostId"`
	MetricDate     string   `json:"metricDate"`
	Likes          *int64   `json:"likes,omitempty"`
	Comments       *int64   `json:"comments,omitempty"`
	Shares         *int64   `json:"shares,omitempty"`
	Views          *int64   `json:"views,omitempty"`
	Clicks         *int64   `json:"clicks,omitempty"`
	EngagementRate *float64 `json:"engagementRate,omitempty"`
	Reach          *int64   `json:"reach,omitempty"`
	Impressions    *int64   `json:"impressions,omitempty"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

func toAnalyticsResponse(analytics *PostAnalytics) analyticsResponse {
	return analyticsResponse{
		ID:             analytics.ID.String(),
		SocialPostID:   analytics.SocialPostID.String(),
		MetricDate:     analytics.MetricDate.Format("2006-01-02"),
		Likes:          analytics.Likes,
		Comments:       analytics.Comments,
		Shares:         analytics.Shares,
		Views:          analytics.Views,
		Clicks:         analytics.Clicks,
		EngagementRate: analytics.EngagementRate,
		Reach:          analytics.Reach,
		Impressions:    analytics.Impressions,
		CreatedAt:      analytics.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      analytics.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
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
	case ErrCodeInvalidPayload:
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
