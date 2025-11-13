package reports

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers report endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	group := api.Group("/reports")

	group.Post("/summary", handler.PostSummary)
}
