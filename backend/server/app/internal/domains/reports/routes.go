package reports

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers report endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	group := api.Group("/reports")

	group.Post("/summary", handler.PostSummary)

	group.Post("/", handler.CreateDefinition)
	group.Get("/", handler.ListDefinitions)
	group.Get("/:id", handler.GetDefinition)
	group.Put("/:id", handler.UpdateDefinition)
	group.Post("/archive", handler.ArchiveDefinitions)
	group.Post("/restore", handler.RestoreDefinitions)
	group.Post("/delete", handler.DeleteDefinitions)
	group.Post("/favorite", handler.ToggleFavorite)

	group.Post("/:id/schedules", handler.CreateSchedule)
	group.Get("/:id/schedules", handler.ListSchedules)
	group.Put("/schedules/:scheduleID", handler.UpdateSchedule)
	group.Post("/schedules/:scheduleID/toggle", handler.ToggleSchedule)
	group.Delete("/schedules/:scheduleID", handler.DeleteSchedule)

	group.Post("/:id/deliveries", handler.CreateDelivery)
	group.Get("/:id/deliveries", handler.ListDeliveries)
	group.Put("/deliveries/:deliveryID", handler.UpdateDelivery)
	group.Post("/deliveries/:deliveryID/toggle", handler.ToggleDelivery)
	group.Delete("/deliveries/:deliveryID", handler.DeleteDelivery)

	group.Post("/runs/bulk", handler.QueueRuns)
	group.Get("/:id/runs", handler.ListRuns)
}
