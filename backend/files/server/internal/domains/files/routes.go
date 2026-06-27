package files

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes registers file endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// File CRUD operations
	api.Post("/", handler.UploadFile)
	api.Get("/", handler.ListFiles)
	api.Get("/:id", handler.GetFile)
	api.Get("/:id/download", handler.DownloadFile)
	api.Get("/:id/url", handler.GetFileURL)
	api.Delete("/:id", handler.DeleteFile)
}

