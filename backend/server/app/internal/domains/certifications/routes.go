package certifications

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers certification endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Certification routes
	api.Post("/", handler.CreateCertification)
	api.Get("/", handler.ListCertifications)
	api.Get("/featured", handler.ListFeaturedCertifications) // Public access
	api.Get("/skill/:skillId", handler.GetCertificationsBySkill) // Public access - Get certifications by skill
	api.Get("/entities/:entityType/:entityId", handler.GetEntityCertifications) // Public access - Get certifications by entity (must be before /:id)
	api.Get("/:id", handler.GetCertification)
	api.Get("/:id/public", handler.GetCertificationPublic) // Public access
	api.Patch("/:id", handler.UpdateCertification)
	api.Delete("/:id", handler.DeleteCertification)
	
	// Skill relationship routes
	api.Post("/:id/skills/:skillId", handler.AddCertificationSkill)
	api.Delete("/:id/skills/:skillId", handler.RemoveCertificationSkill)
	api.Get("/:id/skills", handler.GetCertificationSkills)
	
	// Entity link routes (for projects, etc.)
	api.Post("/:id/entities", handler.CreateCertificationEntityLink)
	api.Get("/:id/entities", handler.GetCertificationEntityLinks)
	api.Delete("/:id/entities", handler.DeleteCertificationEntityLinks)
	api.Delete("/:id/entities/:linkId", handler.DeleteCertificationEntityLink)
}

