package interests

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers interest endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Interest CRUD operations
	api.Post("/", handler.CreateInterest)
	api.Get("/", handler.ListInterests)
	api.Get("/featured", handler.ListFeaturedInterests)
	api.Get("/search", handler.SearchInterests)
	api.Get("/:id", handler.GetInterest)
	api.Get("/slug/:slug", handler.GetInterestBySlug)
	api.Patch("/:id", handler.UpdateInterest)
	api.Delete("/:id", handler.DeleteInterest)
}

