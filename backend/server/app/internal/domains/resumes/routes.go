package resumes

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers resume endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Post("/upload", handler.UploadResume) // File upload endpoint (must be before /:id routes)
	api.Post("/", handler.CreateResume)
	api.Get("/", handler.ListResumes)
	api.Get("/:id", handler.GetResume)
	api.Patch("/:id", handler.UpdateResume)
	api.Delete("/:id", handler.DeleteResume)
	api.Patch("/:id/main", handler.MarkAsMain)
	api.Patch("/:id/featured", handler.MarkAsFeatured)
	api.Delete("/:id/main", handler.UnmarkAsMain)
	api.Delete("/:id/featured", handler.UnmarkAsFeatured)
}

// SetupPublicRoutes registers public resume endpoints.
func SetupPublicRoutes(api fiber.Router, handler Handler) {
	api.Get("/resume/download", handler.DownloadResume)
	api.Get("/resume/preview", handler.PreviewResume)
}

