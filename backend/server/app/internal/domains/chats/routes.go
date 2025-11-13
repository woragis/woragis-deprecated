package chats

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers chat endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	group := api.Group("/chats")

	group.Post("/conversations", handler.CreateConversation)
	group.Get("/conversations", handler.ListConversations)
	group.Get("/conversations/:id/messages", handler.ListMessages)
	group.Post("/conversations/:id/messages", handler.AppendMessage)
}
