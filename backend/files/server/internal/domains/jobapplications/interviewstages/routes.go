package interviewstages

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers interview stage endpoints.
// The routes are nested under /job-applications/:applicationId/interview-stages
// so applicationId is available in the route params.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Post("/", handler.CreateStage)
	api.Get("/", handler.ListStages)
	api.Get("/:id", handler.GetStage)
	api.Patch("/:id", handler.UpdateStage)
	api.Delete("/:id", handler.DeleteStage)
	api.Post("/:id/schedule", handler.ScheduleStage)
	api.Post("/:id/complete", handler.CompleteStage)
}

