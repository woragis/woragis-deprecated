package monitoring

import "github.com/gofiber/fiber/v2"

// SetupRoutes wires monitoring routes (events endpoint only when repo available).
func SetupRoutes(api fiber.Router, handler *Handler) {
	group := api.Group("/monitoring")
	group.Get("/events", handler.ListEvents)
}
