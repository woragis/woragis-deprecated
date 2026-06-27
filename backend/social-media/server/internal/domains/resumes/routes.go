package resumes

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers resume endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Post("/upload", handler.UploadResume) // File upload endpoint (must be before /:id routes)
	api.Post("/generate", handler.GenerateResume) // Generate resume endpoint (must be before /:id routes)
	api.Post("/", handler.CreateResume)
	api.Get("/tags", handler.ListResumeTags) // Get all tags for autocomplete (must be before /:id routes)
	api.Get("/", handler.ListResumes) // Supports ?tags=tag1,tag2 query parameter
	api.Get("/:id/download", handler.DownloadResumeByID) // Download resume by ID (must be before /:id)
	api.Get("/:id", handler.GetResume)
	api.Patch("/:id", handler.UpdateResume)
	api.Delete("/:id", handler.DeleteResume)
	api.Patch("/:id/main", handler.MarkAsMain)
	api.Patch("/:id/featured", handler.MarkAsFeatured)
	api.Delete("/:id/main", handler.UnmarkAsMain)
	api.Delete("/:id/featured", handler.UnmarkAsFeatured)
	api.Post("/:id/recalculate-metrics", handler.RecalculateMetrics) // Manually recalculate metrics
	api.Get("/jobs/:id", handler.GetJobStatus) // Get job status
	api.Post("/jobs/:id/retry", handler.RetryJob) // Retry failed job
	api.Post("/jobs/:id/cancel", handler.CancelJob) // Cancel pending/processing job
}

// SetupPublicRoutes registers public resume endpoints.
func SetupPublicRoutes(api fiber.Router, handler Handler) {
	api.Get("/resume/download", handler.DownloadResume)
	api.Get("/resume/preview", handler.PreviewResume)
}

// SetupInternalRoutes registers internal resume endpoints (no auth middleware).
func SetupInternalRoutes(api fiber.Router, handler Handler) {
	api.Post("/resumes/complete", handler.CompleteResumeGeneration) // Internal callback for resume worker
}

