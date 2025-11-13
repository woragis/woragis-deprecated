package projects

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers project endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	group := api.Group("/projects")

	group.Post("/", handler.CreateProject)
	group.Get("/", handler.ListProjects)
	group.Patch("/:id/status", handler.UpdateStatus)
	group.Patch("/:id/metrics", handler.UpdateMetrics)
	group.Post("/:id/milestones", handler.AddMilestone)
	group.Get("/:id/milestones", handler.ListMilestones)
	group.Patch("/milestones/:milestoneID", handler.ToggleMilestoneCompletion)
}
