package languages

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers language learning endpoints.
func SetupRoutes(api fiber.Router, handler *Handler) {
	languagesGroup := api.Group("/languages")

	languagesGroup.Post("/sessions", handler.PostStudySession)
	languagesGroup.Get("/sessions", handler.GetStudySessions)

	languagesGroup.Post("/vocabulary", handler.PostVocabulary)
	languagesGroup.Get("/vocabulary", handler.GetVocabulary)

	languagesGroup.Get("/summary", handler.GetSummary)
}
