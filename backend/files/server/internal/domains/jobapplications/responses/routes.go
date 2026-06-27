package responses

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers response endpoints.
// The routes are nested under /job-applications/:applicationId/responses
// so applicationId is available in the route params.
func SetupRoutes(api fiber.Router, handler Handler) {
	api.Post("/", handler.CreateResponse)
	api.Get("/", handler.ListResponses)
	api.Get("/:id", handler.GetResponse)
	api.Patch("/:id", handler.UpdateResponse)
	api.Delete("/:id", handler.DeleteResponse)
}

