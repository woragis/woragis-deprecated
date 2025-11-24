package impactmetrics

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers impact metric endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Impact metric routes
	api.Post("/", handler.CreateImpactMetric)
	api.Get("/", handler.ListImpactMetrics)
	api.Get("/featured", handler.ListFeaturedImpactMetrics) // Public access
	api.Get("/dashboard", handler.GetDashboardMetrics)      // Dashboard aggregation
	api.Get("/type/:type", handler.GetMetricsByType)        // Get metrics by type
	api.Get("/type/:type/total", handler.GetTotalValueByType) // Get total value by type
	api.Get("/entity/:entityType/:entityId", handler.GetMetricsByEntity) // Get metrics by entity (public)
	api.Get("/:id", handler.GetImpactMetric)
	api.Patch("/:id", handler.UpdateImpactMetric)
	api.Delete("/:id", handler.DeleteImpactMetric)
}

