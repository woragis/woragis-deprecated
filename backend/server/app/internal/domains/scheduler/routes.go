package scheduler

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers scheduler endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	group := api.Group("/scheduler")

	group.Post("/", handler.PostSchedule)
	group.Get("/", handler.GetSchedules)
	group.Patch("/:id", handler.PatchSchedule)
	group.Delete("/:id", handler.DeleteSchedule)
	group.Post("/bulk/activate", handler.BulkActivate)
	group.Post("/bulk/deactivate", handler.BulkDeactivate)
	group.Post("/bulk/pause", handler.BulkPause)
	group.Post("/bulk/resume", handler.BulkResume)
	group.Get("/runs", handler.ListRuns)
}
