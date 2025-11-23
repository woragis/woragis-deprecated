package casestudies

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers case study endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Case study CRUD operations
	api.Post("/", handler.CreateCaseStudy)
	api.Get("/", handler.ListCaseStudies)
	api.Get("/project-slug/:projectSlug", handler.GetCaseStudyByProjectSlug)
	api.Get("/:id", handler.GetCaseStudy)
	api.Patch("/:id", handler.UpdateCaseStudy)
	api.Delete("/:id", handler.DeleteCaseStudy)
}

