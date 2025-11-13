package finances

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers finance endpoints under /finance.
func SetupRoutes(api fiber.Router, handler Handler) {
	finance := api.Group("/finance")

	finance.Post("/transactions", handler.RecordTransaction)
	finance.Post("/transactions/bulk", handler.BulkRecord)
	finance.Post("/transactions/import", handler.ImportTransactions)
	finance.Patch("/transactions/:id", handler.UpdateTransaction)
	finance.Patch("/transactions/:id/archive", handler.ToggleArchived)
	finance.Patch("/transactions/:id/recurring", handler.ToggleRecurring)
	finance.Patch("/transactions/:id/essential", handler.ToggleEssential)
	finance.Patch("/transactions/bulk/category", handler.BulkUpdateCategory)
	finance.Patch("/transactions/bulk/type", handler.BulkUpdateType)
	finance.Delete("/transactions/bulk", handler.BulkDelete)
	finance.Get("/transactions", handler.ListTransactions)
	finance.Get("/summary", handler.Summary)

	finance.Post("/templates", handler.CreateTemplate)
	finance.Get("/templates", handler.ListTemplates)
	finance.Put("/templates/:id", handler.UpdateTemplate)
	finance.Delete("/templates/:id", handler.DeleteTemplate)

	finance.Get("/reports/cashflow", handler.Cashflow)
}
