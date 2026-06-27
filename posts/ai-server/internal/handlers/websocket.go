package handlers

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/woragis/posts-ai-service/internal/services"
)

func ChatWebSocket(chatSvc *services.ChatService) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		chatID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			c.WriteMessage(1, []byte(`{"error":"Invalid chat_id"}`))
			return
		}

		userID := c.Query("user_id")
		if userID == "" {
			c.WriteMessage(1, []byte(`{"error":"Missing user_id"}`))
			return
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			c.WriteMessage(1, []byte(`{"error":"Invalid user_id"}`))
			return
		}

		// Verify user owns this chat (optional - could add auth check)
		// For now, trust the user_id query parameter

		for {
			var msg struct {
				Prompt string `json:"prompt"`
				Agent  string `json:"agent"`
			}

			err := c.ReadJSON(&msg)
			if err != nil {
				break
			}

			// Use background context for WebSocket operations
			ctx := context.Background()
			response, err := chatSvc.ImproveContent(ctx, userUUID, chatID, msg.Prompt, msg.Agent)
			if err != nil {
				c.WriteJSON(map[string]string{"error": err.Error()})
				continue
			}

			// Send response in streaming format
			c.WriteJSON(map[string]string{
				"response": response,
				"done":     "true",
			})
		}
	})
}
