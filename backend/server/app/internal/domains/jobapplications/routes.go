package jobapplications

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers job application endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Post("/", handler.CreateJobApplication)
	api.Get("/", handler.ListJobApplications)
	api.Get("/:id", handler.GetJobApplication)
	api.Patch("/:id/status", handler.UpdateJobApplicationStatus)
	api.Delete("/:id", handler.DeleteJobApplication)
}

