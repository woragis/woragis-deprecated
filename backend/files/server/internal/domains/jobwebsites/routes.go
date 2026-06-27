package jobwebsites

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers job website endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Post("/", handler.CreateJobWebsite)
	api.Get("/", handler.ListJobWebsites)
	api.Get("/:id", handler.GetJobWebsite)
	api.Patch("/:id", handler.UpdateJobWebsite)
	api.Post("/:id/reset-counter", handler.ResetCounter)
	api.Delete("/:id", handler.DeleteJobWebsite)
}

