package testimonials

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers testimonial endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Testimonial CRUD operations
	api.Post("/", handler.CreateTestimonial)
	api.Get("/", handler.ListTestimonials)
	api.Get("/:id", handler.GetTestimonial)
	api.Patch("/:id", handler.UpdateTestimonial)
	api.Delete("/:id", handler.DeleteTestimonial)

	// Moderation operations
	api.Post("/:id/approve", handler.ApproveTestimonial)
	api.Post("/:id/reject", handler.RejectTestimonial)
	api.Post("/:id/hide", handler.HideTestimonial)
}

