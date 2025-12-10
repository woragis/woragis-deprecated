package chats

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes registers chat endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	group := api.Group("/chats")

	group.Post("/conversations", handler.CreateConversation)
	group.Get("/conversations", handler.ListConversations)
	group.Get("/conversations/search", handler.SearchConversations)
	group.Get("/conversations/:id", handler.GetConversation)

	group.Post("/conversations/:id/messages", handler.AppendMessage)
	group.Get("/conversations/:id/messages", handler.ListMessages)

	group.Post("/conversations/archive", handler.ArchiveConversations)
	group.Post("/conversations/delete", handler.DeleteConversations)
	group.Post("/conversations/restore", handler.RestoreConversations)

	group.Post("/conversations/:id/transcripts", handler.ShareTranscript)
	group.Get("/conversations/:id/transcripts", handler.ListTranscripts)
	group.Get("/transcripts/:code", handler.GetTranscript)

	group.Post("/conversations/:id/assign", handler.AssignConversation)
	group.Post("/conversations/:id/unassign", handler.UnassignConversation)
	group.Get("/conversations/:id/assignments", handler.ListAssignments)

	group.Get("/conversations/:id/context", handler.GetContextPreview)

	// Note: WebSocket stream route is registered separately in main.go
	// with auth middleware before upgrade check
	// group.Get("/conversations/:id/stream", websocket.New(handler.StreamConversation))
}
