package scheduler

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers scheduler endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	group := api.Group("/scheduler")

	group.Post("/", handler.PostSchedule)
	group.Get("/", handler.GetSchedules)
	group.Patch("/:id", handler.PatchSchedule)
}
