package userprofiles

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers user profile endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Get("/profile", handler.GetProfile)
	api.Put("/profile", handler.UpsertProfile)
}

