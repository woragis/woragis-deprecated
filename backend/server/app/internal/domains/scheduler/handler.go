package scheduler

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

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
	UserID      string `json:"user_id"`
	ReportType  string `json:"report_type"`
	AgentAlias  string `json:"agent_alias"`
	Frequency   string `json:"frequency"`
	Weekday     string `json:"weekday"`
	TimeOfDay   string `json:"time_of_day"`
	Timezone    string `json:"timezone"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type updateSchedulePayload struct {
	UserID      string `json:"user_id"`
	ReportType  string `json:"report_type"`
	AgentAlias  string `json:"agent_alias"`
	Frequency   string `json:"frequency"`
	Weekday     string `json:"weekday"`
	TimeOfDay   string `json:"time_of_day"`
	Timezone    string `json:"timezone"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Active      *bool  `json:"active"`
}

// PostSchedule creates a schedule.
func (h *Handler) PostSchedule(c *fiber.Ctx) error {
	var payload createSchedulePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	schedule, err := h.service.Create(c.Context(), CreateRequest{
		UserID:      userID,
		ReportType:  payload.ReportType,
		AgentAlias:  payload.AgentAlias,
		Frequency:   payload.Frequency,
		Weekday:     payload.Weekday,
		TimeOfDay:   payload.TimeOfDay,
		Timezone:    payload.Timezone,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, schedule)
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

	userID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
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
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
		Active:      payload.Active,
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, schedule)
}

// GetSchedules lists schedules for a user.
func (h *Handler) GetSchedules(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Query("user_id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, nil)
	}

	schedules, err := h.service.List(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, schedules)
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		status := statusFromError(domainErr.Code)
		h.logWarn(domainErr.Message)
		return response.Error(c, status, domainErr.Code, nil)
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
