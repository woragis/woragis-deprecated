package projects

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers project endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	group := api.Group("/projects")

	group.Post("/", handler.CreateProject)
	group.Get("/", handler.ListProjects)
	group.Get("/slug/:slug", handler.GetProjectBySlug)
	group.Get("/slug", handler.SearchProjectsBySlug)
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

	// Documentation routes
	group.Post("/:id/documentation", handler.CreateDocumentation)
	group.Get("/:id/documentation", handler.GetDocumentation)
	group.Patch("/:id/documentation/visibility", handler.UpdateDocumentationVisibility)
	group.Delete("/:id/documentation", handler.DeleteDocumentation)
	group.Get("/slug/:slug/documentation", handler.GetPublicDocumentation) // Public access

	// Documentation Section routes
	group.Post("/:id/documentation/sections", handler.CreateDocumentationSection)
	group.Get("/:id/documentation/sections", handler.ListDocumentationSections)
	group.Patch("/documentation/sections/:sectionID", handler.UpdateDocumentationSection)
	group.Delete("/documentation/sections/:sectionID", handler.DeleteDocumentationSection)
	group.Patch("/:id/documentation/sections/reorder", handler.ReorderDocumentationSections)

	// Technology routes
	group.Post("/:id/technologies", handler.CreateTechnology)
	group.Get("/:id/technologies", handler.ListTechnologies)
	group.Patch("/technologies/:techID", handler.UpdateTechnology)
	group.Delete("/technologies/:techID", handler.DeleteTechnology)
	group.Post("/:id/technologies/bulk", handler.BulkCreateTechnologies)
	group.Patch("/:id/technologies/bulk", handler.BulkUpdateTechnologies)

	// File Structure routes
	group.Post("/:id/file-structures", handler.CreateFileStructure)
	group.Get("/:id/file-structures", handler.ListFileStructures)
	group.Patch("/file-structures/:fileStructureID", handler.UpdateFileStructure)
	group.Delete("/file-structures/:fileStructureID", handler.DeleteFileStructure)
	group.Post("/:id/file-structures/bulk", handler.BulkCreateFileStructures)
	group.Patch("/:id/file-structures/bulk", handler.BulkUpdateFileStructures)

	// Architecture Diagram routes
	group.Post("/:id/architecture-diagrams", handler.CreateArchitectureDiagram)
	group.Get("/:id/architecture-diagrams", handler.ListArchitectureDiagrams)
	group.Get("/architecture-diagrams/:diagramID", handler.GetArchitectureDiagram)
	group.Patch("/architecture-diagrams/:diagramID", handler.UpdateArchitectureDiagram)
	group.Delete("/architecture-diagrams/:diagramID", handler.DeleteArchitectureDiagram)
}
