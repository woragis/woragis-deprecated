package scheduler

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes scheduler endpoints.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler constructs a Handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type createSchedulePayload struct {
	ReportType  string          `json:"report_type"`
	AgentAlias  string          `json:"agent_alias"`
	Frequency   string          `json:"frequency"`
	Weekday     string          `json:"weekday"`
	TimeOfDay   string          `json:"time_of_day"`
	Timezone    string          `json:"timezone"`
	RRule       string          `json:"rrule"`
	Priority    int             `json:"priority"`
	Channels    map[string]bool `json:"channels"`
	Email       string          `json:"email"`
	PhoneNumber string          `json:"phone_number"`
}

type updateSchedulePayload struct {
	ReportType  string          `json:"report_type"`
	AgentAlias  string          `json:"agent_alias"`
	Frequency   string          `json:"frequency"`
	Weekday     string          `json:"weekday"`
	TimeOfDay   string          `json:"time_of_day"`
	Timezone    string          `json:"timezone"`
	RRule       string          `json:"rrule"`
	Priority    *int            `json:"priority"`
	Channels    map[string]bool `json:"channels"`
	Email       string          `json:"email"`
	PhoneNumber string          `json:"phone_number"`
	Active      *bool           `json:"active"`
	Paused      *bool           `json:"paused"`
}

type bulkIDsPayload struct {
	ScheduleIDs []string `json:"schedule_ids"`
}

type runListQuery struct {
	ScheduleID string `query:"schedule_id"`
	Status     string `query:"status"`
	Limit      string `query:"limit"`
	Offset     string `query:"offset"`
}

type scheduleResponse struct {
	ID          string          `json:"id"`
	ReportType  string          `json:"report_type"`
	AgentAlias  string          `json:"agent_alias"`
	Frequency   string          `json:"frequency"`
	Weekday     string          `json:"weekday"`
	TimeOfDay   string          `json:"time_of_day"`
	Timezone    string          `json:"timezone"`
	RRule       string          `json:"rrule,omitempty"`
	Priority    int             `json:"priority"`
	Channels    map[string]bool `json:"channels"`
	Email       string          `json:"email"`
	PhoneNumber string          `json:"phone_number"`
	Active      bool            `json:"active"`
	Paused      bool            `json:"paused"`
	NextRun     string          `json:"next_run"`
	LastRun     *string         `json:"last_run,omitempty"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

type runResponse struct {
	ID          string         `json:"id"`
	ScheduleID  string         `json:"schedule_id"`
	Status      string         `json:"status"`
	Output      string         `json:"output,omitempty"`
	Error       string         `json:"error,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	StartedAt   *string        `json:"started_at,omitempty"`
	CompletedAt *string        `json:"completed_at,omitempty"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

// PostSchedule creates a schedule.
func (h *Handler) PostSchedule(c *fiber.Ctx) error {
	var payload createSchedulePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	schedule, err := h.service.Create(c.Context(), CreateRequest{
		UserID:      userID,
		ReportType:  payload.ReportType,
		AgentAlias:  payload.AgentAlias,
		Frequency:   payload.Frequency,
		Weekday:     payload.Weekday,
		TimeOfDay:   payload.TimeOfDay,
		Timezone:    payload.Timezone,
		RRule:       payload.RRule,
		Priority:    payload.Priority,
		Channels:    payload.Channels,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toScheduleResponse(schedule))
}

// PatchSchedule updates schedule metadata.
func (h *Handler) PatchSchedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var payload updateSchedulePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	schedule, err := h.service.Update(c.Context(), UpdateRequest{
		UserID:      userID,
		ScheduleID:  id,
		ReportType:  payload.ReportType,
		AgentAlias:  payload.AgentAlias,
		Frequency:   payload.Frequency,
		Weekday:     payload.Weekday,
		TimeOfDay:   payload.TimeOfDay,
		Timezone:    payload.Timezone,
		RRule:       payload.RRule,
		Priority:    payload.Priority,
		Channels:    payload.Channels,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
		Active:      payload.Active,
		Paused:      payload.Paused,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toScheduleResponse(schedule))
}

// DeleteSchedule removes a schedule.
func (h *Handler) DeleteSchedule(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	if err := h.service.Delete(c.Context(), id, userID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{"message": "Schedule deleted successfully"})
}

// GetSchedules lists schedules for the authenticated user.
func (h *Handler) GetSchedules(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	schedules, err := h.service.List(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]scheduleResponse, 0, len(schedules))
	for i := range schedules {
		resp = append(resp, toScheduleResponse(&schedules[i]))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

// BulkActivate activates schedules.
func (h *Handler) BulkActivate(c *fiber.Ctx) error {
	req, err := h.parseBulkIDs(c)
	if err != nil {
		return err
	}
	if err := h.service.BulkSetActive(c.Context(), req.UserID, req.ScheduleIDs, true); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "activated"})
}

// BulkDeactivate deactivates schedules.
func (h *Handler) BulkDeactivate(c *fiber.Ctx) error {
	req, err := h.parseBulkIDs(c)
	if err != nil {
		return err
	}
	if err := h.service.BulkSetActive(c.Context(), req.UserID, req.ScheduleIDs, false); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "deactivated"})
}

// BulkPause pauses schedules.
func (h *Handler) BulkPause(c *fiber.Ctx) error {
	req, err := h.parseBulkIDs(c)
	if err != nil {
		return err
	}
	if err := h.service.BulkPause(c.Context(), req.UserID, req.ScheduleIDs, true); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "paused"})
}

// BulkResume resumes paused schedules.
func (h *Handler) BulkResume(c *fiber.Ctx) error {
	req, err := h.parseBulkIDs(c)
	if err != nil {
		return err
	}
	if err := h.service.BulkPause(c.Context(), req.UserID, req.ScheduleIDs, false); err != nil {
		return h.handleError(c, err)
	}
	return response.Success(c, fiber.StatusOK, fiber.Map{"status": "resumed"})
}

// ListRuns returns execution history.
func (h *Handler) ListRuns(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	var query runListQuery
	if err := c.QueryParser(&query); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	var scheduleID uuid.UUID
	if strings.TrimSpace(query.ScheduleID) != "" {
		scheduleID, err = uuid.Parse(query.ScheduleID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "invalid schedule_id"})
		}
	}

	limit, offset := parseLimitOffset(query.Limit, query.Offset, 50)

	runs, err := h.service.ListRuns(c.Context(), ListRunsRequest{
		UserID:     userID,
		ScheduleID: scheduleID,
		Status:     query.Status,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]runResponse, 0, len(runs))
	for _, run := range runs {
		resp = append(resp, toRunResponse(run))
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *Handler) parseBulkIDs(c *fiber.Ctx) (struct {
	UserID      uuid.UUID
	ScheduleIDs []uuid.UUID
}, error) {
	var payload bulkIDsPayload
	if err := c.BodyParser(&payload); err != nil {
		return struct {
			UserID      uuid.UUID
			ScheduleIDs []uuid.UUID
		}{}, response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return struct {
			UserID      uuid.UUID
			ScheduleIDs []uuid.UUID
		}{}, response.Error(c, fiber.StatusUnauthorized, ErrCodeInvalidPayload, nil)
	}

	ids := make([]uuid.UUID, 0, len(payload.ScheduleIDs))
	for _, raw := range payload.ScheduleIDs {
		id, err := uuid.Parse(raw)
		if err != nil {
			return struct {
				UserID      uuid.UUID
				ScheduleIDs []uuid.UUID
			}{}, response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{"message": "invalid schedule_ids"})
		}
		ids = append(ids, id)
	}

	return struct {
		UserID      uuid.UUID
		ScheduleIDs []uuid.UUID
	}{
		UserID:      userID,
		ScheduleIDs: ids,
	}, nil
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromError(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, fiber.Map{"message": domainErr.Message})
	}

	h.logError("scheduler: unexpected error", err)
	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, nil)
}

func statusFromError(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidReport, ErrCodeInvalidAgent, ErrCodeInvalidFrequency:
		return fiber.StatusBadRequest
	case ErrCodeNotFound:
		return fiber.StatusNotFound
	default:
		return fiber.StatusInternalServerError
	}
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

func toScheduleResponse(schedule *Schedule) scheduleResponse {
	resp := scheduleResponse{
		ID:          schedule.ID.String(),
		ReportType:  schedule.ReportType,
		AgentAlias:  schedule.AgentAlias,
		Frequency:   schedule.Frequency,
		Weekday:     schedule.Weekday,
		TimeOfDay:   schedule.TimeOfDay,
		Timezone:    schedule.Timezone,
		RRule:       schedule.RRule,
		Priority:    schedule.Priority,
		Channels:    schedule.Channels.Data,
		Email:       schedule.Email,
		PhoneNumber: schedule.PhoneNumber,
		Active:      schedule.Active,
		Paused:      schedule.Paused,
		NextRun:     schedule.NextRun.Format(time.RFC3339),
		CreatedAt:   schedule.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   schedule.UpdatedAt.Format(time.RFC3339),
	}
	if schedule.LastRun != nil {
		str := schedule.LastRun.Format(time.RFC3339)
		resp.LastRun = &str
	}
	return resp
}

func toRunResponse(run ExecutionRun) runResponse {
	resp := runResponse{
		ID:         run.ID.String(),
		ScheduleID: run.ScheduleID.String(),
		Status:     run.Status,
		Output:     run.Output,
		Error:      run.ErrorMessage,
		Metadata:   run.Metadata.Data,
		CreatedAt:  run.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  run.UpdatedAt.Format(time.RFC3339),
	}
	if run.StartedAt != nil {
		str := run.StartedAt.Format(time.RFC3339)
		resp.StartedAt = &str
	}
	if run.CompletedAt != nil {
		str := run.CompletedAt.Format(time.RFC3339)
		resp.CompletedAt = &str
	}
	return resp
}

func parseLimitOffset(limitRaw, offsetRaw string, defaultLimit int) (int, int) {
	limit := defaultLimit
	if l, err := strconv.Atoi(limitRaw); err == nil && l > 0 {
		limit = l
	}
	offset := 0
	if o, err := strconv.Atoi(offsetRaw); err == nil && o >= 0 {
		offset = o
	}
	return limit, offset
}
