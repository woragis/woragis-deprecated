package content

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers content post endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Post("/posts", handler.CreateContentPost)
	api.Get("/posts", handler.ListContentPosts)
	api.Get("/posts/:id", handler.GetContentPost)
	api.Patch("/posts/:id/priority", handler.UpdateContentPostPriority)
	api.Post("/posts/:id/repurpose", handler.RepurposeToPlatforms)
	api.Get("/posts/:id/repurposing-history", handler.GetRepurposingHistory)
	api.Get("/backlog", handler.GetContentBacklog)
}
