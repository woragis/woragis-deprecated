package scheduling

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/internal/domains/socialmediaposts"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes scheduling endpoints.
type Handler interface {
	SchedulePost(c *fiber.Ctx) error
	GetSchedule(c *fiber.Ctx) error
	GetScheduleForDateRange(c *fiber.Ctx) error
	GetUpcomingPosts(c *fiber.Ctx) error
	UpdateSchedule(c *fiber.Ctx) error
	CancelSchedule(c *fiber.Ctx) error
	AutoSchedule(c *fiber.Ctx) error
	CheckConflicts(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

var _ Handler = (*handler)(nil)

// NewHandler constructs a scheduling handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// Payloads

type schedulePostPayload struct {
	SocialPostID string     `json:"socialPostId"`
	ScheduledAt  time.Time  `json:"scheduledAt"`
	PlatformID   *string    `json:"platformId,omitempty"`
}

type updateSchedulePayload struct {
	ScheduledAt *time.Time `json:"scheduledAt,omitempty"`
	Status      *string    `json:"status,omitempty"`
}

type autoSchedulePayload struct {
	Platform string `json:"platform"`
}

// Handlers

func (h *handler) SchedulePost(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var payload schedulePostPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	socialPostID, err := uuid.Parse(payload.SocialPostID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req := SchedulePostRequest{
		SocialPostID: socialPostID,
		ScheduledAt:  payload.ScheduledAt.UTC(),
	}

	if payload.PlatformID != nil {
		platformID, err := uuid.Parse(*payload.PlatformID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
		}
		req.PlatformID = &platformID
	}

	schedule, err := h.service.SchedulePost(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toScheduleResponse(schedule))
}

func (h *handler) GetSchedule(c *fiber.Ctx) error {
	scheduleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	schedule, err := h.service.GetSchedule(c.Context(), scheduleID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toScheduleResponse(schedule))
}

func (h *handler) GetScheduleForDateRange(c *fiber.Ctx) error {
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if startDateStr == "" || endDateStr == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	schedules, err := h.service.GetScheduleForDateRange(c.Context(), startDate, endDate)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]scheduleResponse, len(schedules))
	for i := range schedules {
		responses[i] = toScheduleResponse(&schedules[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) GetUpcomingPosts(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	schedules, err := h.service.GetUpcomingPosts(c.Context(), limit)
	if err != nil {
		return h.handleError(c, err)
	}

	responses := make([]scheduleResponse, len(schedules))
	for i := range schedules {
		responses[i] = toScheduleResponse(&schedules[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

func (h *handler) UpdateSchedule(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	scheduleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateSchedulePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	req := UpdateScheduleRequest{
		ScheduleID: scheduleID,
	}

	if payload.ScheduledAt != nil {
		at := payload.ScheduledAt.UTC()
		req.ScheduledAt = &at
	}

	if payload.Status != nil {
		status := ScheduledPostStatus(*payload.Status)
		req.Status = &status
	}

	schedule, err := h.service.UpdateSchedule(c.Context(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toScheduleResponse(schedule))
}

func (h *handler) CancelSchedule(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	scheduleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.CancelSchedule(c.Context(), scheduleID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "schedule cancelled"})
}

func (h *handler) AutoSchedule(c *fiber.Ctx) error {
	_, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	socialPostID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload autoSchedulePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	platform := socialmediaposts.Platform(payload.Platform)
	schedule, err := h.service.AutoSchedule(c.Context(), socialPostID, platform)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toScheduleResponse(schedule))
}

func (h *handler) CheckConflicts(c *fiber.Ctx) error {
	scheduledAtStr := c.Query("scheduledAt")
	if scheduledAtStr == "" {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	scheduledAt, err := time.Parse(time.RFC3339, scheduledAtStr)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var excludeScheduleID *uuid.UUID
	if excludeStr := c.Query("excludeScheduleId"); excludeStr != "" {
		excludeID, err := uuid.Parse(excludeStr)
		if err == nil {
			excludeScheduleID = &excludeID
		}
	}

	hasConflict, err := h.service.CheckConflicts(c.Context(), scheduledAt.UTC(), excludeScheduleID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"hasConflict": hasConflict})
}

// Response helpers

type scheduleResponse struct {
	ID            string              `json:"id"`
	SocialPostID  string              `json:"socialPostId"`
	PlatformID    *string             `json:"platformId,omitempty"`
	ScheduledDate string              `json:"scheduledDate"`
	ScheduledTime string              `json:"scheduledTime"`
	ScheduledAt   string              `json:"scheduledAt"`
	Status        ScheduledPostStatus `json:"status"`
	CreatedAt     string              `json:"createdAt"`
	UpdatedAt     string              `json:"updatedAt"`
}

func toScheduleResponse(schedule *ScheduledPost) scheduleResponse {
	var platformID *string
	if schedule.PlatformID != nil {
		idStr := schedule.PlatformID.String()
		platformID = &idStr
	}

	return scheduleResponse{
		ID:            schedule.ID.String(),
		SocialPostID:  schedule.SocialPostID.String(),
		PlatformID:    platformID,
		ScheduledDate: schedule.ScheduledDate.Format("2006-01-02"),
		ScheduledTime: schedule.ScheduledTime.Format("15:04:05"),
		ScheduledAt:   schedule.ScheduledAt.Format("2006-01-02T15:04:05Z07:00"),
		Status:        schedule.Status,
		CreatedAt:     schedule.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:     schedule.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
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
	case ErrCodeInvalidPayload, ErrCodeInvalidStatus, ErrCodeInvalidStatusTransition:
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
