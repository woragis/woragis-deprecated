package projectcasestudies

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers project case study endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Case study routes
	api.Post("/:id/case-studies", handler.CreateCaseStudy)
	api.Get("/:id/case-studies", handler.ListCaseStudies)
	api.Get("/:id/case-studies/current", handler.GetCaseStudyByProjectID)
	api.Get("/:id/case-studies/public", handler.GetPublicCaseStudyByProjectID) // Public access
	api.Get("/case-studies/:caseStudyID", handler.GetCaseStudy)
	api.Patch("/case-studies/:caseStudyID", handler.UpdateCaseStudy)
	api.Delete("/case-studies/:caseStudyID", handler.DeleteCaseStudy)
}
