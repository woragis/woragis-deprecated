package jobapplications

import (
	"github.com/gofiber/fiber/v2"
	
	"woragis-files-service/internal/domains/jobapplications/responses"
	"woragis-files-service/internal/domains/jobapplications/interviewstages"
)

// SetupRoutes registers job application endpoints and subdomain routes.
func SetupRoutes(api fiber.Router, handler Handler, responseHandler responses.Handler, stageHandler interviewstages.Handler) {
	// Main job application routes
	api.Post("/", handler.CreateJobApplication)
	api.Get("/", handler.ListJobApplications)
	api.Get("/:id", handler.GetJobApplication)
	api.Patch("/:id/status", handler.UpdateJobApplicationStatus)
	api.Patch("/:id", handler.UpdateJobApplication)
	api.Delete("/:id", handler.DeleteJobApplication)
	api.Post("/:id/generate-cover-letter", handler.GenerateCoverLetter)
	
	// Subdomain routes
	responses.SetupRoutes(api.Group("/:applicationId/responses"), responseHandler)
	interviewstages.SetupRoutes(api.Group("/:applicationId/interview-stages"), stageHandler)
}

