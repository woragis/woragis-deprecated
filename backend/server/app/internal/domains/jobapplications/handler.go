package jobapplications

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes job application endpoints.
type Handler interface {
	CreateJobApplication(c *fiber.Ctx) error
	GetJobApplication(c *fiber.Ctx) error
	ListJobApplications(c *fiber.Ctx) error
	UpdateJobApplicationStatus(c *fiber.Ctx) error
	DeleteJobApplication(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

// NewHandler constructs a job application handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

type createJobApplicationPayload struct {
	CompanyName string `json:"companyName"`
	Location    string `json:"location"`
	JobTitle    string `json:"jobTitle"`
	JobURL      string `json:"jobUrl"`
	Website     string `json:"website"`
}

type updateStatusPayload struct {
	Status ApplicationStatus `json:"status"`
}

func (h *handler) CreateJobApplication(c *fiber.Ctx) error {
	// Get user ID from context
	var userID uuid.UUID
	var err error
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		userID = apiKey.UserID
	} else {
		userID, err = authdomain.UserIDFromContext(c)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
				"message": "authentication required",
			})
		}
	}

	var payload createJobApplicationPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request payload",
		})
	}

	application, err := h.service.RequestJobApplication(
		c.Context(),
		userID,
		payload.CompanyName,
		payload.Location,
		payload.JobTitle,
		payload.JobURL,
		payload.Website,
	)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, application)
}

func (h *handler) GetJobApplication(c *fiber.Ctx) error {
	applicationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid application id",
		})
	}

	application, err := h.service.GetJobApplication(c.Context(), applicationID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, application)
}

func (h *handler) ListJobApplications(c *fiber.Ctx) error {
	// Get user ID from context
	var userID *uuid.UUID
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		id := apiKey.UserID
		userID = &id
	} else if uid, err := authdomain.UserIDFromContext(c); err == nil {
		userID = &uid
	}

	filters := JobApplicationFilters{
		UserID: userID,
	}

	// Optional query parameters
	if website := c.Query("website"); website != "" {
		filters.Website = &website
	}
	if status := c.Query("status"); status != "" {
		appStatus := ApplicationStatus(status)
		filters.Status = &appStatus
	}

	// Pagination
	if limit := c.QueryInt("limit", 50); limit > 0 {
		filters.Limit = limit
	}
	if offset := c.QueryInt("offset", 0); offset > 0 {
		filters.Offset = offset
	}

	applications, err := h.service.ListJobApplications(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"applications": applications,
		"count":        len(applications),
	})
}

func (h *handler) UpdateJobApplicationStatus(c *fiber.Ctx) error {
	applicationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid application id",
		})
	}

	var payload updateStatusPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request payload",
		})
	}

	if err := h.service.UpdateJobApplicationStatus(c.Context(), applicationID, payload.Status); err != nil {
		return h.handleError(c, err)
	}

	application, err := h.service.GetJobApplication(c.Context(), applicationID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, application)
}

func (h *handler) DeleteJobApplication(c *fiber.Ctx) error {
	_, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid application id",
		})
	}

	// Note: We don't have DeleteJobApplication in service, but we can add it if needed
	// For now, just return not implemented
	return response.Error(c, fiber.StatusNotImplemented, 501, fiber.Map{
		"message": "delete not implemented yet",
	})
}

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		statusCode := fiber.StatusInternalServerError
		switch domainErr.Code {
		case ErrCodeNotFound:
			statusCode = fiber.StatusNotFound
		case ErrCodeInvalidPayload, ErrCodeInvalidStatus:
			statusCode = fiber.StatusBadRequest
		case ErrCodeJobQueueFailure, ErrCodeAIServiceFailure, ErrCodePlaywrightFailure:
			statusCode = fiber.StatusServiceUnavailable
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

