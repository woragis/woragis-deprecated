package impactmetrics

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	translationsdomain "github.com/woragis/backend/server/app/internal/domains/translations"
	translationenricher "github.com/woragis/backend/server/app/pkg/translations"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes impact metric endpoints.
type Handler interface {
	CreateImpactMetric(c *fiber.Ctx) error
	UpdateImpactMetric(c *fiber.Ctx) error
	GetImpactMetric(c *fiber.Ctx) error
	ListImpactMetrics(c *fiber.Ctx) error
	ListFeaturedImpactMetrics(c *fiber.Ctx) error
	GetMetricsByEntity(c *fiber.Ctx) error
	DeleteImpactMetric(c *fiber.Ctx) error
	// Dashboard endpoints
	GetDashboardMetrics(c *fiber.Ctx) error
	GetMetricsByType(c *fiber.Ctx) error
	GetTotalValueByType(c *fiber.Ctx) error
}

type handler struct {
	service          Service
	enricher         *translationenricher.Enricher
	translationService translationsdomain.Service
	logger           *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs an impact metric handler.
func NewHandler(service Service, enricher *translationenricher.Enricher, translationService translationsdomain.Service, logger *slog.Logger) Handler {
	return &handler{
		service:           service,
		enricher:          enricher,
		translationService: translationService,
		logger:            logger,
	}
}

// Payloads

type createImpactMetricPayload struct {
	Type        MetricType  `json:"type"`
	Value       float64     `json:"value"`
	Unit        MetricUnit  `json:"unit"`
	Description string      `json:"description,omitempty"`
	EntityType  *EntityType `json:"entityType,omitempty"`
	EntityID    *string     `json:"entityId,omitempty"`
	PeriodStart *string     `json:"periodStart,omitempty"` // ISO 8601 date string
	PeriodEnd   *string     `json:"periodEnd,omitempty"`   // ISO 8601 date string
	Featured    bool        `json:"featured,omitempty"`
	DisplayOrder int        `json:"displayOrder,omitempty"`
}

type updateImpactMetricPayload struct {
	Type        *MetricType  `json:"type,omitempty"`
	Value       *float64     `json:"value,omitempty"`
	Unit        *MetricUnit  `json:"unit,omitempty"`
	Description *string      `json:"description,omitempty"`
	EntityType  *EntityType  `json:"entityType,omitempty"`
	EntityID    *string      `json:"entityId,omitempty"`
	PeriodStart *string      `json:"periodStart,omitempty"`
	PeriodEnd   *string      `json:"periodEnd,omitempty"`
	Featured    *bool        `json:"featured,omitempty"`
	DisplayOrder *int         `json:"displayOrder,omitempty"`
}

// Handlers

func (h *handler) CreateImpactMetric(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	var payload createImpactMetricPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req := CreateImpactMetricRequest{
		Type:        payload.Type,
		Value:       payload.Value,
		Unit:        payload.Unit,
		Description: payload.Description,
		EntityType:  payload.EntityType,
		EntityID:    payload.EntityID,
		Featured:    payload.Featured,
		DisplayOrder: payload.DisplayOrder,
	}

	// Parse dates
	if payload.PeriodStart != nil && *payload.PeriodStart != "" {
		periodStart, err := parseDate(*payload.PeriodStart)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidDate, fiber.Map{
				"message": "invalid period start date format",
			})
		}
		req.PeriodStart = &periodStart
	}

	if payload.PeriodEnd != nil && *payload.PeriodEnd != "" {
		periodEnd, err := parseDate(*payload.PeriodEnd)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidDate, fiber.Map{
				"message": "invalid period end date format",
			})
		}
		req.PeriodEnd = &periodEnd
	}

	metric, err := h.service.CreateImpactMetric(c.Context(), userID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, metric)
}

func (h *handler) UpdateImpactMetric(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	metricID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid metric id",
		})
	}

	var payload updateImpactMetricPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req := UpdateImpactMetricRequest{}

	if payload.Type != nil {
		req.Type = payload.Type
	}
	if payload.Value != nil {
		req.Value = payload.Value
	}
	if payload.Unit != nil {
		req.Unit = payload.Unit
	}
	if payload.Description != nil {
		req.Description = payload.Description
	}
	if payload.EntityType != nil {
		req.EntityType = payload.EntityType
	}
	if payload.EntityID != nil {
		req.EntityID = payload.EntityID
	}
	if payload.PeriodStart != nil {
		if *payload.PeriodStart != "" {
			periodStart, err := parseDate(*payload.PeriodStart)
			if err != nil {
				return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidDate, fiber.Map{
					"message": "invalid period start date format",
				})
			}
			req.PeriodStart = &periodStart
		} else {
			req.PeriodStart = nil
		}
	}
	if payload.PeriodEnd != nil {
		if *payload.PeriodEnd != "" {
			periodEnd, err := parseDate(*payload.PeriodEnd)
			if err != nil {
				return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidDate, fiber.Map{
					"message": "invalid period end date format",
				})
			}
			req.PeriodEnd = &periodEnd
		} else {
			req.PeriodEnd = nil
		}
	}
	if payload.Featured != nil {
		req.Featured = payload.Featured
	}
	if payload.DisplayOrder != nil {
		req.DisplayOrder = payload.DisplayOrder
	}

	metric, err := h.service.UpdateImpactMetric(c.Context(), userID, metricID, req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, metric)
}

func (h *handler) GetImpactMetric(c *fiber.Ctx) error {
	metricID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid metric id",
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	metric, err := h.service.GetImpactMetric(c.Context(), metricID, userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		fieldMap := map[string]*string{
			"description": &metric.Description,
		}
		_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeImpactMetric, metric.ID, language, fieldMap)
	}

	return response.Success(c, fiber.StatusOK, metric)
}

func (h *handler) ListImpactMetrics(c *fiber.Ctx) error {
	var userID *uuid.UUID
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		uid := apiKey.UserID
		userID = &uid
	} else if uid, err := authdomain.UserIDFromContext(c); err == nil {
		userID = &uid
	}

	filters := ListImpactMetricsFilters{
		UserID: userID,
	}

	// Parse query parameters
	if typeStr := c.Query("type"); typeStr != "" {
		metricType := MetricType(typeStr)
		filters.Type = &metricType
	}

	if entityTypeStr := c.Query("entityType"); entityTypeStr != "" {
		entityType := EntityType(entityTypeStr)
		filters.EntityType = &entityType
	}

	if entityIDStr := c.Query("entityId"); entityIDStr != "" {
		entityID, err := uuid.Parse(entityIDStr)
		if err == nil {
			filters.EntityID = &entityID
		}
	}

	if featuredStr := c.Query("featured"); featuredStr != "" {
		featured := featuredStr == "true"
		filters.Featured = &featured
	}

	if periodStartStr := c.Query("periodStart"); periodStartStr != "" {
		periodStart, err := parseDate(periodStartStr)
		if err == nil {
			filters.PeriodStart = &periodStart
		}
	}

	if periodEndStr := c.Query("periodEnd"); periodEndStr != "" {
		periodEnd, err := parseDate(periodEndStr)
		if err == nil {
			filters.PeriodEnd = &periodEnd
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	if orderBy := c.Query("orderBy"); orderBy != "" {
		filters.OrderBy = orderBy
	}

	if order := c.Query("order"); order != "" {
		filters.Order = order
	}

	metrics, err := h.service.ListImpactMetrics(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range metrics {
			fieldMap := map[string]*string{
				"description": &metrics[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeImpactMetric, metrics[i].ID, language, fieldMap)
		}
	}

	return response.Success(c, fiber.StatusOK, metrics)
}

func (h *handler) ListFeaturedImpactMetrics(c *fiber.Ctx) error {
	metrics, err := h.service.ListFeaturedImpactMetrics(c.Context())
	if err != nil {
		return h.handleError(c, err)
	}

	// Apply translations if enricher is available
	if h.enricher != nil {
		language := translationsdomain.LanguageFromContext(c)
		for i := range metrics {
			fieldMap := map[string]*string{
				"description": &metrics[i].Description,
			}
			_ = h.enricher.EnrichEntityFields(c.Context(), translationsdomain.EntityTypeImpactMetric, metrics[i].ID, language, fieldMap)
		}
	}

	return response.Success(c, fiber.StatusOK, metrics)

func (h *handler) GetMetricsByEntity(c *fiber.Ctx) error {
	entityTypeStr := c.Params("entityType")
	entityIDStr := c.Params("entityId")

	entityType := EntityType(entityTypeStr)
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid entity id",
		})
	}

	metrics, err := h.service.GetMetricsByEntity(c.Context(), entityType, entityID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, metrics)
}

func (h *handler) DeleteImpactMetric(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	metricID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid metric id",
		})
	}

	if err := h.service.DeleteImpactMetric(c.Context(), userID, metricID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "metric deleted successfully",
	})
}

// Dashboard handlers

func (h *handler) GetDashboardMetrics(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	dashboard, err := h.service.GetDashboardMetrics(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, dashboard)
}

func (h *handler) GetMetricsByType(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	metricTypeStr := c.Params("type")
	metricType := MetricType(metricTypeStr)

	if !isValidMetricType(metricType) {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidType, fiber.Map{
			"message": "invalid metric type",
		})
	}

	metrics, err := h.service.GetMetricsByType(c.Context(), userID, metricType)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, metrics)
}

func (h *handler) GetTotalValueByType(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, fiber.Map{
			"message": "authentication required",
		})
	}

	metricTypeStr := c.Params("type")
	metricType := MetricType(metricTypeStr)

	if !isValidMetricType(metricType) {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidType, fiber.Map{
			"message": "invalid metric type",
		})
	}

	total, err := h.service.GetTotalValueByType(c.Context(), userID, metricType)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"type":  metricType,
		"total": total,
	})
}

// Helper functions

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	domainErr, ok := AsDomainError(err)
	if !ok {
		h.logger.Error("unexpected error in impact metric handler", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{
			"message": "internal server error",
		})
	}

	statusCode := fiber.StatusInternalServerError
	switch domainErr.Code {
	case ErrCodeInvalidPayload, ErrCodeInvalidType, ErrCodeInvalidUnit, ErrCodeInvalidValue, ErrCodeInvalidEntityType, ErrCodeInvalidDate:
		statusCode = fiber.StatusBadRequest
	case ErrCodeNotFound:
		statusCode = fiber.StatusNotFound
	case ErrCodeUnauthorized:
		statusCode = fiber.StatusUnauthorized
	case ErrCodeConflict:
		statusCode = fiber.StatusConflict
	}

	return response.Error(c, statusCode, domainErr.Code, fiber.Map{
		"message": domainErr.Message,
	})
}

func parseDate(dateStr string) (time.Time, error) {
	// Try ISO 8601 format first (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t, nil
	}
	// Try RFC3339 format
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t, nil
	}
	return time.Time{}, fiber.NewError(fiber.StatusBadRequest, "invalid date format")
}

