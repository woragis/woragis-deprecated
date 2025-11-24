package aimlintegrations

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers AI/ML integration endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// AI/ML Integration routes
	api.Post("/", handler.CreateAIMLIntegration)
	api.Get("/", handler.ListAIMLIntegrations)
	api.Get("/featured", handler.ListFeaturedAIMLIntegrations) // Public access - Showcase
	api.Get("/type/:type", handler.GetIntegrationsByType)      // Public access - Filter by type
	api.Get("/framework/:framework", handler.GetIntegrationsByFramework) // Public access - Filter by framework
	api.Get("/project/:projectId", handler.GetIntegrationsByProject)     // Public access - Get integrations by project
	api.Get("/:id", handler.GetAIMLIntegration)
	api.Get("/:id/public", handler.GetAIMLIntegrationPublic) // Public access
	api.Patch("/:id", handler.UpdateAIMLIntegration)
	api.Delete("/:id", handler.DeleteAIMLIntegration)
}

