package ideas

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers ideas endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	group := api.Group("/ideas")

	group.Post("/", handler.PostIdea)
	group.Get("/", handler.ListIdeas)
	group.Patch("/:id", handler.PatchIdea)
	group.Patch("/:id/position", handler.PatchIdeaPosition)
	group.Post("/links", handler.PostLink)
	group.Get("/links", handler.ListLinks)
}
