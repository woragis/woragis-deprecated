package skills

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers skill endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	group := api.Group("/skills")

	// Skill CRUD operations
	group.Post("/", handler.CreateSkill)
	group.Get("/", handler.ListSkills)
	group.Get("/with-counts", handler.GetAllSkillsWithProjectCounts)
	group.Get("/search", handler.SearchSkills)
	group.Get("/category", handler.ListSkillsByCategory)
	group.Get("/:id", handler.GetSkill)
	group.Get("/slug/:slug", handler.GetSkillBySlug)
	group.Patch("/:id", handler.UpdateSkill)

	// Project-Skill relationship operations
	// These will be nested under projects in the projects routes
	// But we can also expose them here for convenience
	group.Post("/projects/:projectId/skills/:skillId", handler.AttachSkillToProject)
	group.Delete("/projects/:projectId/skills/:skillId", handler.DetachSkillFromProject)
	group.Get("/projects/:projectId/skills", handler.GetProjectSkills)
	group.Get("/:skillId/projects", handler.GetProjectsBySkill)
}

