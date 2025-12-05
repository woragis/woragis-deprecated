package projects

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers project endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Use the provided router directly (it's already a group with the correct path)
	api.Post("/", handler.CreateProject)
	api.Get("/", handler.ListProjects)
	api.Get("/slug/:slug", handler.GetProjectBySlug)
	api.Get("/slug", handler.SearchProjectsBySlug)
	api.Patch("/:id/status", handler.UpdateStatus)
	api.Patch("/:id/metrics", handler.UpdateMetrics)
	api.Delete("/:id", handler.DeleteProject)

	api.Post("/:id/milestones", handler.AddMilestone)
	api.Get("/:id/milestones", handler.ListMilestones)
	api.Patch("/milestones/:milestoneID", handler.ToggleMilestoneCompletion)
	api.Post("/:id/milestones/bulk", handler.BulkUpdateMilestones)

	api.Get("/:id/kanban", handler.GetKanbanBoard)
	api.Post("/:id/kanban/columns", handler.CreateKanbanColumn)
	api.Patch("/:id/kanban/columns/:columnID", handler.UpdateKanbanColumn)
	api.Patch("/:id/kanban/columns/reorder", handler.ReorderKanbanColumns)
	api.Delete("/:id/kanban/columns/:columnID", handler.DeleteKanbanColumn)

	api.Post("/:id/kanban/cards", handler.CreateKanbanCard)
	api.Patch("/:id/kanban/cards/:cardID", handler.UpdateKanbanCard)
	api.Patch("/:id/kanban/cards/:cardID/move", handler.MoveKanbanCard)
	api.Delete("/:id/kanban/cards/:cardID", handler.DeleteKanbanCard)

	api.Post("/:id/dependencies", handler.CreateDependency)
	api.Get("/:id/dependencies", handler.ListDependencies)
	api.Delete("/:id/dependencies/:dependencyID", handler.DeleteDependency)

	api.Post("/:id/duplicate", handler.DuplicateProject)

	// Documentation routes
	api.Post("/:id/documentation", handler.CreateDocumentation)
	api.Get("/:id/documentation", handler.GetDocumentation)
	api.Patch("/:id/documentation/visibility", handler.UpdateDocumentationVisibility)
	api.Delete("/:id/documentation", handler.DeleteDocumentation)
	api.Get("/slug/:slug/documentation", handler.GetPublicDocumentation) // Public access

	// Documentation Section routes
	api.Post("/:id/documentation/sections", handler.CreateDocumentationSection)
	api.Get("/:id/documentation/sections", handler.ListDocumentationSections)
	api.Patch("/documentation/sections/:sectionID", handler.UpdateDocumentationSection)
	api.Delete("/documentation/sections/:sectionID", handler.DeleteDocumentationSection)
	api.Patch("/:id/documentation/sections/reorder", handler.ReorderDocumentationSections)

	// Technology routes
	api.Post("/:id/technologies", handler.CreateTechnology)
	api.Get("/:id/technologies", handler.ListTechnologies)
	api.Patch("/technologies/:techID", handler.UpdateTechnology)
	api.Delete("/technologies/:techID", handler.DeleteTechnology)
	api.Post("/:id/technologies/bulk", handler.BulkCreateTechnologies)
	api.Patch("/:id/technologies/bulk", handler.BulkUpdateTechnologies)

	// File Structure routes
	api.Post("/:id/file-structures", handler.CreateFileStructure)
	api.Get("/:id/file-structures", handler.ListFileStructures)
	api.Patch("/file-structures/:fileStructureID", handler.UpdateFileStructure)
	api.Delete("/file-structures/:fileStructureID", handler.DeleteFileStructure)
	api.Post("/:id/file-structures/bulk", handler.BulkCreateFileStructures)
	api.Patch("/:id/file-structures/bulk", handler.BulkUpdateFileStructures)

	// Architecture Diagram routes
	api.Post("/:id/architecture-diagrams", handler.CreateArchitectureDiagram)
	api.Get("/:id/architecture-diagrams", handler.ListArchitectureDiagrams)
	api.Get("/architecture-diagrams/:diagramID", handler.GetArchitectureDiagram)
	api.Patch("/architecture-diagrams/:diagramID", handler.UpdateArchitectureDiagram)
	api.Delete("/architecture-diagrams/:diagramID", handler.DeleteArchitectureDiagram)

	// Project Case Study routes (handled by projectcasestudies subdomain)
	// Routes are registered separately in main.go after project handler is created
}
