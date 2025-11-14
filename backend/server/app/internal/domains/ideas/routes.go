package ideas

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers ideas endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	group := api.Group("/ideas")

	group.Post("/", handler.PostIdea)
	group.Get("/", handler.ListIdeas)
	group.Get("/:id/versions", handler.GetIdeaVersions)
	group.Patch("/:id", handler.PatchIdea)
	group.Patch("/:id/position", handler.PatchIdeaPosition)
	group.Post("/bulk/move", handler.PostBulkMove)
	group.Post("/bulk/update", handler.PostBulkUpdate)
	group.Post("/bulk/delete", handler.PostBulkDelete)
	group.Post("/bulk/restore", handler.PostBulkRestore)
	group.Post("/links", handler.PostLink)
	group.Get("/links", handler.ListLinks)
	group.Get("/collaborators", handler.ListCollaborators)
	group.Post("/collaborators", handler.PostCollaborator)
	group.Delete("/collaborators/:collaborator_id", handler.DeleteCollaborator)
}
