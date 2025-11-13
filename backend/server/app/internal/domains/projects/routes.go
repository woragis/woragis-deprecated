package projects

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers project endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	group := api.Group("/projects")

	group.Post("/", handler.CreateProject)
	group.Get("/", handler.ListProjects)
	group.Patch("/:id/status", handler.UpdateStatus)
	group.Patch("/:id/metrics", handler.UpdateMetrics)

	group.Post("/:id/milestones", handler.AddMilestone)
	group.Get("/:id/milestones", handler.ListMilestones)
	group.Patch("/milestones/:milestoneID", handler.ToggleMilestoneCompletion)
	group.Post("/:id/milestones/bulk", handler.BulkUpdateMilestones)

	group.Get("/:id/kanban", handler.GetKanbanBoard)
	group.Post("/:id/kanban/columns", handler.CreateKanbanColumn)
	group.Patch("/:id/kanban/columns/:columnID", handler.UpdateKanbanColumn)
	group.Patch("/:id/kanban/columns/reorder", handler.ReorderKanbanColumns)
	group.Delete("/:id/kanban/columns/:columnID", handler.DeleteKanbanColumn)

	group.Post("/:id/kanban/cards", handler.CreateKanbanCard)
	group.Patch("/:id/kanban/cards/:cardID", handler.UpdateKanbanCard)
	group.Patch("/:id/kanban/cards/:cardID/move", handler.MoveKanbanCard)
	group.Delete("/:id/kanban/cards/:cardID", handler.DeleteKanbanCard)

	group.Post("/:id/dependencies", handler.CreateDependency)
	group.Get("/:id/dependencies", handler.ListDependencies)
	group.Delete("/:id/dependencies/:dependencyID", handler.DeleteDependency)

	group.Post("/:id/duplicate", handler.DuplicateProject)
}
