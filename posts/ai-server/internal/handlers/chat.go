package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/woragis/posts-ai-service/internal/services"
)

func GenerateDraft(chatSvc *services.ChatService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			UserID string `json:"user_id"`
			PostID *string `json:"post_id,omitempty"`
			Prompt string `json:"prompt"`
			Agent  string `json:"agent"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id"})
		}

		var postID *uuid.UUID
		if req.PostID != nil {
			pid, err := uuid.Parse(*req.PostID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post_id"})
			}
			postID = &pid
		}

		// Set streaming headers
		c.Set("Content-Type", "application/x-ndjson")
		c.Set("Transfer-Encoding", "chunked")

		// Stream response
		response, err := chatSvc.GenerateDraft(c.Context(), userID, postID, req.Prompt, req.Agent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		// Send complete response
		return c.SendString(response)
	}
}

func ImproveContent(chatSvc *services.ChatService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			UserID          string `json:"user_id"`
			PostID          string `json:"post_id"`
			ImprovementText string `json:"improvement"`
			Agent           string `json:"agent"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		userID, err := uuid.Parse(req.UserID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id"})
		}

		postID, err := uuid.Parse(req.PostID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid post_id"})
		}

		c.Set("Content-Type", "application/x-ndjson")
		c.Set("Transfer-Encoding", "chunked")

		response, err := chatSvc.ImproveContent(c.Context(), userID, postID, req.ImprovementText, req.Agent)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.SendString(response)
	}
}

func GetChat(chatSvc *services.ChatService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		chatID, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid chat_id"})
		}

		chat, turns, err := chatSvc.GetChat(c.Context(), chatID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"chat":  chat,
			"turns": turns,
		})
	}
}

func ListChats(chatSvc *services.ChatService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := uuid.Parse(c.Query("user_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id"})
		}

		limit := c.QueryInt("limit", 20)
		offset := c.QueryInt("offset", 0)

		if limit > 100 {
			limit = 100
		}

		chats, total, err := chatSvc.ListChats(c.Context(), userID, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"chats": chats,
			"total": total,
			"limit": limit,
			"offset": offset,
		})
	}
}

func GetUsageStats(chatSvc *services.ChatService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := uuid.Parse(c.Query("user_id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id"})
		}

		days := c.QueryInt("days", 30)
		if days > 365 {
			days = 365
		}

		stats, err := chatSvc.GetUsageStats(c.Context(), userID, days)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"stats": stats,
			"days": days,
		})
	}
}
