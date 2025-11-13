package monitoring

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/woragis/backend/server/app/pkg/response"
)

// Handler offers HTTP handlers for monitoring related endpoints.
type Handler struct {
	service *Service
}

// NewHandler constructs Handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ListEvents returns most recent monitoring events (production only).
func (h *Handler) ListEvents(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "20"))
	if err != nil || limit <= 0 {
		limit = 20
	}

	events, err := h.service.ListRecentEvents(c.Context(), limit)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, 2001, fiber.Map{"error": err.Error()})
	}

	return response.Success(c, fiber.StatusOK, events)
}
