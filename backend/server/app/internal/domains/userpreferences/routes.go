package userpreferences

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers user preferences endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	group := api.Group("/user-preferences")

	group.Get("", handler.GetPreferences)
	group.Patch("", handler.UpdatePreferences)
}

