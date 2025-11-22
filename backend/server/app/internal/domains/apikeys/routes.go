package apikeys

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers API key endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	group := api.Group("/api-keys")

	group.Post("/", handler.CreateAPIKey)
	group.Get("/", handler.ListAPIKeys)
	group.Get("/:id", handler.GetAPIKey)
	group.Patch("/:id", handler.UpdateAPIKey)
	group.Delete("/:id", handler.DeleteAPIKey)
}

