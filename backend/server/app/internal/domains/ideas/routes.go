package ideas

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers ideas endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	group := api.Group("/ideas")

	group.Post("/", handler.PostIdea)
	group.Get("/", handler.ListIdeas)
	group.Get("/slug/:slug", handler.GetIdeaBySlug)
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

	// IdeaNode routes
	group.Post("/nodes", handler.PostIdeaNode)
	group.Get("/:id/nodes", handler.GetIdeaNodes)
	group.Patch("/nodes/:id", handler.PatchIdeaNode)
	group.Patch("/nodes/:id/position", handler.PatchIdeaNodePosition)
	group.Patch("/nodes/:id/resize", handler.PatchIdeaNodeResize)
	group.Delete("/nodes/:id", handler.DeleteIdeaNode)

	// IdeaNodeConnection routes
	group.Post("/node-connections", handler.PostIdeaNodeConnection)
	group.Get("/:id/node-connections", handler.GetIdeaNodeConnections)
	group.Delete("/node-connections/:id", handler.DeleteIdeaNodeConnection)

	// Document routes
	group.Post("/documents", handler.PostDocument)
	group.Get("/:id/documents", handler.GetDocuments)
	group.Patch("/documents/:id", handler.PatchDocument)
	group.Delete("/documents/:id", handler.DeleteDocument)
}
