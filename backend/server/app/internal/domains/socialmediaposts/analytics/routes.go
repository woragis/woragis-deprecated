package analytics

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers analytics endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Post("/", handler.RecordAnalytics)
	api.Get("/posts/:id", handler.GetPostAnalytics)
	api.Get("/summary", handler.GetAnalyticsSummary)
	api.Get("/top-posts", handler.GetTopPosts)
}
