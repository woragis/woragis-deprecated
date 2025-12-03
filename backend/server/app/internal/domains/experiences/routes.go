package experiences

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers experience endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Experience CRUD operations
	api.Post("/", handler.CreateExperience)
	api.Get("/", handler.ListExperiences)
	api.Get("/:id", handler.GetExperience)
	api.Patch("/:id", handler.UpdateExperience)
	api.Delete("/:id", handler.DeleteExperience)
}

