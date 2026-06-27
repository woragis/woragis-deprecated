package responses

import (
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"woragis-files-service/pkg/response"
)

// Handler exposes response endpoints.
type Handler interface {
	CreateResponse(c *fiber.Ctx) error
	GetResponse(c *fiber.Ctx) error
	ListResponses(c *fiber.Ctx) error
	UpdateResponse(c *fiber.Ctx) error
	DeleteResponse(c *fiber.Ctx) error
}

type handler struct {
	service Service
	logger  *slog.Logger
}

// NewHandler constructs a response handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

type createResponsePayload struct {
	JobApplicationID string       `json:"jobApplicationId"`
	ResponseType     ResponseType `json:"responseType"`
	ResponseDate     string       `json:"responseDate"` // ISO 8601 format
	Message          string       `json:"message,omitempty"`
	ContactPerson    string       `json:"contactPerson,omitempty"`
	ContactEmail     string       `json:"contactEmail,omitempty"`
	ContactPhone     string       `json:"contactPhone,omitempty"`
	ResponseChannel  string       `json:"responseChannel,omitempty"`
}

type updateResponsePayload struct {
	Message         *string `json:"message,omitempty"`
	ContactPerson   *string `json:"contactPerson,omitempty"`
	ContactEmail    *string `json:"contactEmail,omitempty"`
	ContactPhone    *string `json:"contactPhone,omitempty"`
	ResponseChannel *string `json:"responseChannel,omitempty"`
}

func (h *handler) CreateResponse(c *fiber.Ctx) error {
	var payload createResponsePayload
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

	responseDate := time.Now().UTC()
	if payload.ResponseDate != "" {
		parsedDate, err := time.Parse(time.RFC3339, payload.ResponseDate)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
				"message": "invalid response date format, use ISO 8601",
			})
		}
		responseDate = parsedDate
	}

	resp, err := h.service.CreateResponse(c.Context(), applicationID, payload.ResponseType, responseDate)
	if err != nil {
		return h.handleError(c, err)
	}

	// Update additional fields if provided
	if payload.Message != "" || payload.ContactPerson != "" || payload.ContactEmail != "" || payload.ContactPhone != "" || payload.ResponseChannel != "" {
		updates := UpdateResponseRequest{}
		if payload.Message != "" {
			updates.Message = &payload.Message
		}
		if payload.ContactPerson != "" {
			updates.ContactPerson = &payload.ContactPerson
		}
		if payload.ContactEmail != "" {
			updates.ContactEmail = &payload.ContactEmail
		}
		if payload.ContactPhone != "" {
			updates.ContactPhone = &payload.ContactPhone
		}
		if payload.ResponseChannel != "" {
			updates.ResponseChannel = &payload.ResponseChannel
		}
		resp, err = h.service.UpdateResponse(c.Context(), resp.ID, updates)
		if err != nil {
			return h.handleError(c, err)
		}
	}

	return response.Success(c, fiber.StatusCreated, resp)
}

func (h *handler) GetResponse(c *fiber.Ctx) error {
	responseID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid response id",
		})
	}

	resp, err := h.service.GetResponse(c.Context(), responseID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) ListResponses(c *fiber.Ctx) error {
	filters := ResponseFilters{}

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
	if responseTypeStr := c.Query("responseType"); responseTypeStr != "" {
		responseType := ResponseType(responseTypeStr)
		filters.ResponseType = &responseType
	}

	// Pagination
	if limit := c.QueryInt("limit", 50); limit > 0 {
		filters.Limit = limit
	}
	if offset := c.QueryInt("offset", 0); offset > 0 {
		filters.Offset = offset
	}

	responses, err := h.service.ListResponses(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"responses": responses,
		"count":     len(responses),
	})
}

func (h *handler) UpdateResponse(c *fiber.Ctx) error {
	responseID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid response id",
		})
	}

	var payload updateResponsePayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request payload",
		})
	}

	updates := UpdateResponseRequest{
		Message:         payload.Message,
		ContactPerson:   payload.ContactPerson,
		ContactEmail:    payload.ContactEmail,
		ContactPhone:    payload.ContactPhone,
		ResponseChannel: payload.ResponseChannel,
	}

	resp, err := h.service.UpdateResponse(c.Context(), responseID, updates)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, resp)
}

func (h *handler) DeleteResponse(c *fiber.Ctx) error {
	responseID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid response id",
		})
	}

	if err := h.service.DeleteResponse(c.Context(), responseID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "response deleted successfully",
	})
}


func (h *handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		statusCode := fiber.StatusInternalServerError
		switch domainErr.Code {
		case ErrCodeNotFound:
			statusCode = fiber.StatusNotFound
		case ErrCodeInvalidPayload, ErrCodeInvalidResponseType:
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

