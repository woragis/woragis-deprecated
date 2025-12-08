package interviewstages

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes interview stage endpoints.
type Handler interface {
	CreateStage(c *fiber.Ctx) error
	GetStage(c *fiber.Ctx) error
	ListStages(c *fiber.Ctx) error
	UpdateStage(c *fiber.Ctx) error
	DeleteStage(c *fiber.Ctx) error
	ScheduleStage(c *fiber.Ctx) error
	CompleteStage(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

// NewHandler constructs an interview stage handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

type createStagePayload struct {
	JobApplicationID string    `json:"jobApplicationId"`
	StageType        StageType `json:"stageType"`
	ScheduledDate    string    `json:"scheduledDate,omitempty"` // ISO 8601 format
	InterviewerName  string    `json:"interviewerName,omitempty"`
	InterviewerEmail string    `json:"interviewerEmail,omitempty"`
	Location         string    `json:"location,omitempty"`
	Notes            string    `json:"notes,omitempty"`
}

type updateStagePayload struct {
	ScheduledDate    *string `json:"scheduledDate,omitempty"` // ISO 8601 format
	InterviewerName  *string `json:"interviewerName,omitempty"`
	InterviewerEmail *string `json:"interviewerEmail,omitempty"`
	Location         *string `json:"location,omitempty"`
	Notes            *string `json:"notes,omitempty"`
	Feedback         *string `json:"feedback,omitempty"`
}

type scheduleStagePayload struct {
	ScheduledDate string `json:"scheduledDate"` // ISO 8601 format
}

type completeStagePayload struct {
	CompletedDate string       `json:"completedDate"` // ISO 8601 format
	Outcome       StageOutcome `json:"outcome"`
}

func (h *handler) CreateStage(c *fiber.Ctx) error {
	var payload createStagePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request payload",
		})
	}

	// Get applicationId from route params (when nested) or payload
	var applicationID uuid.UUID
	var err error
	if applicationIDStr := c.Params("applicationId"); applicationIDStr != "" {
		applicationID, err = uuid.Parse(applicationIDStr)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
				"message": "invalid job application id in route",
			})
		}
	} else {
		applicationID, err = uuid.Parse(payload.JobApplicationID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
				"message": "invalid job application id",
			})
		}
	}

	stage, err := h.service.CreateStage(c.Context(), applicationID, payload.StageType)
	if err != nil {
		return h.handleError(c, err)
	}

	// Update additional fields if provided
	if payload.ScheduledDate != "" || payload.InterviewerName != "" || payload.InterviewerEmail != "" || payload.Location != "" || payload.Notes != "" {
		updates := UpdateStageRequest{}
		if payload.ScheduledDate != "" {
			parsedDate, err := time.Parse(time.RFC3339, payload.ScheduledDate)
			if err != nil {
				return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
					"message": "invalid scheduled date format, use ISO 8601",
				})
			}
			updates.ScheduledDate = &parsedDate
		}
		if payload.InterviewerName != "" {
			updates.InterviewerName = &payload.InterviewerName
		}
		if payload.InterviewerEmail != "" {
			updates.InterviewerEmail = &payload.InterviewerEmail
		}
		if payload.Location != "" {
			updates.Location = &payload.Location
		}
		if payload.Notes != "" {
			updates.Notes = &payload.Notes
		}
		stage, err = h.service.UpdateStage(c.Context(), stage.ID, updates)
		if err != nil {
			return h.handleError(c, err)
		}
	}

	return response.Success(c, fiber.StatusCreated, stage)
}

func (h *handler) GetStage(c *fiber.Ctx) error {
	stageID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid stage id",
		})
	}

	stage, err := h.service.GetStage(c.Context(), stageID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, stage)
}

func (h *handler) ListStages(c *fiber.Ctx) error {
	filters := StageFilters{}

	// Get applicationId from route params (when nested) or query parameter
	if applicationIDStr := c.Params("applicationId"); applicationIDStr != "" {
		applicationID, err := uuid.Parse(applicationIDStr)
		if err == nil {
			filters.JobApplicationID = &applicationID
		}
	} else if applicationIDStr := c.Query("jobApplicationId"); applicationIDStr != "" {
		applicationID, err := uuid.Parse(applicationIDStr)
		if err == nil {
			filters.JobApplicationID = &applicationID
		}
	}
	if stageTypeStr := c.Query("stageType"); stageTypeStr != "" {
		stageType := StageType(stageTypeStr)
		filters.StageType = &stageType
	}
	if outcomeStr := c.Query("outcome"); outcomeStr != "" {
		outcome := StageOutcome(outcomeStr)
		filters.Outcome = &outcome
	}

	// Pagination
	if limit := c.QueryInt("limit", 50); limit > 0 {
		filters.Limit = limit
	}
	if offset := c.QueryInt("offset", 0); offset > 0 {
		filters.Offset = offset
	}

	stages, err := h.service.ListStages(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"stages": stages,
		"count":  len(stages),
	})
}

func (h *handler) UpdateStage(c *fiber.Ctx) error {
	stageID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid stage id",
		})
	}

	var payload updateStagePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request payload",
		})
	}

	updates := UpdateStageRequest{}
	if payload.ScheduledDate != nil {
		parsedDate, err := time.Parse(time.RFC3339, *payload.ScheduledDate)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
				"message": "invalid scheduled date format, use ISO 8601",
			})
		}
		updates.ScheduledDate = &parsedDate
	}
	updates.InterviewerName = payload.InterviewerName
	updates.InterviewerEmail = payload.InterviewerEmail
	updates.Location = payload.Location
	updates.Notes = payload.Notes
	updates.Feedback = payload.Feedback

	stage, err := h.service.UpdateStage(c.Context(), stageID, updates)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, stage)
}

func (h *handler) DeleteStage(c *fiber.Ctx) error {
	stageID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid stage id",
		})
	}

	if err := h.service.DeleteStage(c.Context(), stageID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "stage deleted successfully",
	})
}


func (h *handler) ScheduleStage(c *fiber.Ctx) error {
	stageID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid stage id",
		})
	}

	var payload scheduleStagePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request payload",
		})
	}

	scheduledDate, err := time.Parse(time.RFC3339, payload.ScheduledDate)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid scheduled date format, use ISO 8601",
		})
	}

	stage, err := h.service.ScheduleStage(c.Context(), stageID, scheduledDate)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, stage)
}

func (h *handler) CompleteStage(c *fiber.Ctx) error {
	stageID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid stage id",
		})
	}

	var payload completeStagePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request payload",
		})
	}

	completedDate, err := time.Parse(time.RFC3339, payload.CompletedDate)
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid completed date format, use ISO 8601",
		})
	}

	stage, err := h.service.CompleteStage(c.Context(), stageID, completedDate, payload.Outcome)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, stage)
}

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		statusCode := fiber.StatusInternalServerError
		switch domainErr.Code {
		case ErrCodeNotFound:
			statusCode = fiber.StatusNotFound
		case ErrCodeInvalidPayload, ErrCodeInvalidStageType, ErrCodeInvalidOutcome:
			statusCode = fiber.StatusBadRequest
		}

		return response.Error(c, statusCode, domainErr.Code, fiber.Map{
			"message": domainErr.Message,
		})
	}

	h.logger.Error("unhandled error", slog.Any("error", err))
	return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
		"message": "internal server error",
	})
}

