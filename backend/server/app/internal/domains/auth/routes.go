package auth

import "github.com/gofiber/fiber/v2"

// SetupRoutes wires the public auth handlers to the supplied Fiber router.
func SetupRoutes(api fiber.Router, handler *Handler) {
	authGroup := api.Group("/auth")

	authGroup.Post("/register", handler.Register)
	authGroup.Post("/confirm", handler.ConfirmEmail)
	authGroup.Post("/confirm/resend", handler.ResendConfirmation)
	authGroup.Post("/login", handler.Login)
	authGroup.Post("/refresh", handler.RefreshSession)
	authGroup.Post("/password/reset/request", handler.RequestPasswordReset)
	authGroup.Post("/password/reset/confirm", handler.ConfirmPasswordReset)
}

// SetupProtectedRoutes wires protected auth handlers requiring JWT middleware.
func SetupProtectedRoutes(api fiber.Router, handler *Handler) {
	authGroup := api.Group("/auth")

	authGroup.Post("/logout", handler.Logout)
	authGroup.Get("/sessions", handler.ListSessions)
	authGroup.Post("/sessions/revoke", handler.RevokeOtherSessions)
	authGroup.Post("/mfa/enable", handler.EnableMFA)
	authGroup.Post("/mfa/verify", handler.VerifyMFA)
	authGroup.Post("/mfa/disable", handler.DisableMFA)
}
