package finances

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers finance endpoints under /finance.
func SetupRoutes(api fiber.Router, handler *Handler) {
	finance := api.Group("/finance")

	finance.Post("/transactions", handler.RecordTransaction)
	finance.Get("/transactions", handler.ListTransactions)
	finance.Get("/summary", handler.Summary)
}
