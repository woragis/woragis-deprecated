package skills

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers skill endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Use the provided router directly (it's already a group with the correct path)

	// Skill CRUD operations
	api.Post("/", handler.CreateSkill)
	api.Get("/", handler.ListSkills)
	api.Get("/with-counts", handler.GetAllSkillsWithProjectCounts)
	api.Get("/search", handler.SearchSkills)
	api.Get("/category", handler.ListSkillsByCategory)
	// Timeline - must be before /:id to avoid route conflict
	api.Get("/timeline", handler.GetSkillsTimeline)
	api.Get("/:id", handler.GetSkill)
	api.Get("/slug/:slug", handler.GetSkillBySlug)
	api.Patch("/:id", handler.UpdateSkill)

	// Project-Skill relationship operations
	// These will be nested under projects in the projects routes
	// But we can also expose them here for convenience
	api.Post("/projects/:projectId/skills/:skillId", handler.AttachSkillToProject)
	api.Delete("/projects/:projectId/skills/:skillId", handler.DetachSkillFromProject)
	api.Get("/projects/:projectId/skills", handler.GetProjectSkills)
	api.Get("/:skillId/projects", handler.GetProjectsBySkill)
}

