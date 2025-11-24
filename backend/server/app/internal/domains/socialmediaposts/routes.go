package socialmediaposts

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers social media post endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Post CRUD operations
	api.Post("/", handler.CreatePost)
	api.Get("/", handler.ListPosts)
	api.Get("/by-url", handler.GetPostByURL)
	api.Get("/:id", handler.GetPost)
	api.Patch("/:id", handler.UpdatePost)
	api.Patch("/:id/engagement", handler.UpdatePostEngagement)
	api.Delete("/:id", handler.DeletePost)

	// Link operations
	api.Post("/links", handler.CreateLink)
	api.Patch("/links/:id", handler.UpdateLink)
	api.Delete("/links/:id", handler.DeleteLink)
	api.Get("/posts/:postId/links", handler.GetLinksByPost)
	api.Get("/posts/:postId/entities", handler.GetEntitiesByPost)
	api.Get("/links/by-entity", handler.GetLinksByEntity)
	api.Get("/posts/by-entity", handler.GetPostsByEntity)
}

