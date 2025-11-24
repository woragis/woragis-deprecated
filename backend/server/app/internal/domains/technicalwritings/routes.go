package technicalwritings

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers technical writing endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Technical Writing routes
	api.Post("/", handler.CreateTechnicalWriting)
	api.Get("/", handler.ListTechnicalWritings)
	api.Get("/search", handler.SearchTechnicalWritings) // Public access - Search writings
	api.Get("/featured", handler.ListFeaturedTechnicalWritings) // Public access - Portfolio showcase
	api.Get("/type/:type", handler.GetWritingsByType) // Public access - Filter by type (article, tutorial, etc.)
	api.Get("/platform/:platform", handler.GetWritingsByPlatform) // Public access - Filter by platform (medium, dev.to, etc.)
	api.Get("/project/:projectId", handler.GetWritingsByProject) // Public access - Get writings by project
	api.Get("/:id", handler.GetTechnicalWriting)
	api.Get("/:id/public", handler.GetTechnicalWritingPublic) // Public access
	api.Patch("/:id", handler.UpdateTechnicalWriting)
	api.Delete("/:id", handler.DeleteTechnicalWriting)
}

