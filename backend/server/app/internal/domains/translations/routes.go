package translations

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers translation endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Post("/request", handler.RequestTranslation)
	api.Post("/translate-entity", handler.TranslateEntity)
	api.Get("/", handler.ListTranslations)
	api.Get("/get", handler.GetTranslation)
}

