package creativeassets

import (
	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes registers creative asset routes
func RegisterRoutes(app *fiber.App, handler Handler, authMiddleware fiber.Handler) {
	api := app.Group("/api/v1/creative-assets")

	// Public routes (for serving images)
	api.Get("/:id/data", handler.GetAssetData)

	// Protected routes
	protected := api.Group("", authMiddleware)
	{
		// CRUD operations
		protected.Post("/", handler.CreateAsset)
		protected.Get("/:id", handler.GetAsset)
		protected.Delete("/:id", handler.DeleteAsset)

		// Entity-based queries
		protected.Get("/entity/:entityType/:entityId", handler.GetAssetsByEntity)
		protected.Get("/entity/:entityType/:entityId/purpose", handler.GetAssetByEntityAndPurpose)

		// Generation endpoints
		protected.Post("/generate/image", handler.GenerateImage)
		protected.Post("/generate/thumbnail", handler.GenerateThumbnail)
		protected.Post("/generate/diagram", handler.GenerateDiagram)
	}
}

// Helper function to register routes (can be called from main setup)
func SetupRoutes(app *fiber.App, handler Handler, authMiddleware fiber.Handler) {
	RegisterRoutes(app, handler, authMiddleware)
}

