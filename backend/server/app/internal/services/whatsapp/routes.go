package whatsapp

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers WhatsApp endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	whatsapp := api.Group("/whatsapp")

	whatsapp.Get("/qr", handler.GetQRCode)
	whatsapp.Get("/status", handler.GetStatus)
	whatsapp.Post("/send", handler.SendMessage)
}
