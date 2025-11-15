package clients

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers client endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	clients := api.Group("/clients")

	clients.Post("/", handler.CreateClient)
	clients.Get("/", handler.ListClients)
	clients.Get("/:id", handler.GetClient)
	clients.Patch("/:id", handler.UpdateClient)
	clients.Patch("/:id/archive", handler.ToggleArchived)
	clients.Delete("/:id", handler.DeleteClient)
}

