package platforms

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers platform configuration endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Get("/", handler.ListConfigs)
	api.Get("/:id", handler.GetConfig)
	api.Get("/by-name/:name", handler.GetConfigByName)
	api.Patch("/:id", handler.UpdateConfig)
	api.Get("/:name/optimal-times", handler.GetOptimalTimes)
}
