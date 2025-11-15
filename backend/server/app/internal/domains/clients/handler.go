package clients

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler exposes client HTTP endpoints.
type Handler struct {
	service Service
	logger  *slog.Logger
}

// NewHandler constructs a Handler.
func NewHandler(service Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type createClientPayload struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Company     string `json:"company"`
	Notes       string `json:"notes"`
}

type updateClientPayload struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Company     *string `json:"company"`
	Notes       *string `json:"notes"`
}

type clientResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number"`
	Company     string `json:"company,omitempty"`
	Notes       string `json:"notes,omitempty"`
	IsArchived  bool   `json:"is_archived"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func toClientResponse(c *Client) clientResponse {
	return clientResponse{
		ID:          c.ID.String(),
		UserID:      c.UserID.String(),
		Name:        c.Name,
		Email:       c.Email,
		PhoneNumber: c.PhoneNumber,
		Company:     c.Company,
		Notes:       c.Notes,
		IsArchived:  c.IsArchived,
		CreatedAt:   c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   c.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreateClient handles POST /clients
func (h *Handler) CreateClient(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	var payload createClientPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "Invalid request payload",
		})
	}

	client, err := h.service.CreateClient(c.Context(), CreateClientRequest{
		UserID:      userID,
		Name:        payload.Name,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
		Company:     payload.Company,
		Notes:       payload.Notes,
	})

	if err != nil {
		return handleError(c, err)
	}

	return response.Success(c, fiber.StatusCreated, toClientResponse(client))
}

// UpdateClient handles PATCH /clients/:id
func (h *Handler) UpdateClient(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "Invalid client ID",
		})
	}

	var payload updateClientPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "Invalid request payload",
		})
	}

	client, err := h.service.UpdateClient(c.Context(), UpdateClientRequest{
		UserID:      userID,
		ClientID:    id,
		Name:        payload.Name,
		Email:       payload.Email,
		PhoneNumber: payload.PhoneNumber,
		Company:     payload.Company,
		Notes:       payload.Notes,
	})

	if err != nil {
		return handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toClientResponse(client))
}

// GetClient handles GET /clients/:id
func (h *Handler) GetClient(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "Invalid client ID",
		})
	}

	client, err := h.service.GetClient(c.Context(), userID, id)
	if err != nil {
		return handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, toClientResponse(client))
}

// ListClients handles GET /clients
func (h *Handler) ListClients(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	includeArchived := c.Query("include_archived") == "true"

	clients, err := h.service.ListClients(c.Context(), userID, includeArchived)
	if err != nil {
		return handleError(c, err)
	}

	responses := make([]clientResponse, len(clients))
	for i := range clients {
		responses[i] = toClientResponse(&clients[i])
	}

	return response.Success(c, fiber.StatusOK, responses)
}

// DeleteClient handles DELETE /clients/:id
func (h *Handler) DeleteClient(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "Invalid client ID",
		})
	}

	if err := h.service.DeleteClient(c.Context(), userID, id); err != nil {
		return handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "Client deleted successfully",
	})
}

// ToggleArchived handles PATCH /clients/:id/archive
func (h *Handler) ToggleArchived(c *fiber.Ctx) error {
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return unauthorizedResponse(c)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "Invalid client ID",
		})
	}

	var payload struct {
		Archived bool `json:"archived"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "Invalid request payload",
		})
	}

	if err := h.service.ToggleArchived(c.Context(), userID, id, payload.Archived); err != nil {
		return handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "Client archive status updated",
	})
}

func unauthorizedResponse(c *fiber.Ctx) error {
	return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
		"message": "authentication required",
	})
}

func handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := err.(*DomainError); ok {
		status := statusFromErrorCode(domainErr.Code)
		return response.Error(c, status, domainErr.Code, fiber.Map{
			"message": domainErr.Message,
		})
	}

	return response.Error(c, fiber.StatusInternalServerError, ErrCodeRepositoryFailure, fiber.Map{
		"message": "An unexpected error occurred",
	})
}

func statusFromErrorCode(code int) int {
	switch code {
	case ErrCodeInvalidPayload, ErrCodeInvalidName, ErrCodeInvalidPhoneNumber, ErrCodeInvalidEmail:
		return fiber.StatusBadRequest
	case ErrCodeNotFound:
		return fiber.StatusNotFound
	case ErrCodeRepositoryFailure:
		return fiber.StatusInternalServerError
	default:
		return fiber.StatusInternalServerError
	}
}

