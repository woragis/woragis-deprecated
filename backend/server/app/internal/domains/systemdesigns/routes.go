package systemdesigns

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers system design endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// System design routes
	api.Post("/", handler.CreateSystemDesign)
	api.Get("/", handler.ListSystemDesigns)
	api.Get("/featured", handler.ListFeaturedSystemDesigns) // Public access
	api.Get("/:id", handler.GetSystemDesign)
	api.Get("/:id/public", handler.GetSystemDesignPublic) // Public access
	api.Patch("/:id", handler.UpdateSystemDesign)
	api.Delete("/:id", handler.DeleteSystemDesign)
}

