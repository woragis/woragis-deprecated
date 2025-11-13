package auth

import "github.com/gofiber/fiber/v2"

// SetupRoutes wires the auth handlers to the supplied Fiber router.
func SetupRoutes(api fiber.Router, handler *Handler) {
	authGroup := api.Group("/auth")

	authGroup.Post("/register", handler.Register)
	authGroup.Post("/login", handler.Login)
	authGroup.Post("/password/reset/request", handler.RequestPasswordReset)
	authGroup.Post("/password/reset/confirm", handler.ConfirmPasswordReset)
}
