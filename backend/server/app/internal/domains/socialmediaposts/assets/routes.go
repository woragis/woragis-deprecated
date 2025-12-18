package assets

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers content asset endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Post("/", handler.CreateAsset)
	api.Get("/:id", handler.GetAsset)
	api.Get("/content-posts/:contentPostId", handler.GetAssetsByContentPost)
	api.Get("/social-posts/:socialPostId", handler.GetAssetsBySocialPost)
	api.Patch("/:id", handler.UpdateAsset)
	api.Delete("/:id", handler.DeleteAsset)
}
