package scheduling

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers scheduling endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Post("/", handler.SchedulePost)
	api.Get("/", handler.GetScheduleForDateRange)
	api.Get("/upcoming", handler.GetUpcomingPosts)
	api.Get("/check-conflicts", handler.CheckConflicts)
	api.Get("/:id", handler.GetSchedule)
	api.Patch("/:id", handler.UpdateSchedule)
	api.Delete("/:id", handler.CancelSchedule)
	api.Post("/:id/auto", handler.AutoSchedule)
}
