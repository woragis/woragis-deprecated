package comments

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes registers comment-related routes.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Use the provided router directly (it's already a group with the correct path)
	
	// Comment routes
	api.Post("/", handler.CreateComment)
	api.Get("/", handler.ListComments)
	api.Get("/count", handler.GetCommentCount)
	api.Get("/:id", handler.GetComment)
	api.Patch("/:id", handler.UpdateComment)
	api.Delete("/:id", handler.DeleteComment)

	// Moderation routes (require auth)
	api.Post("/:id/approve", handler.ApproveComment)
	api.Post("/:id/reject", handler.RejectComment)
	api.Post("/:id/spam", handler.MarkCommentAsSpam)
}

