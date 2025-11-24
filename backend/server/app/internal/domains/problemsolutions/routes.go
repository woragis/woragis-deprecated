package problemsolutions

import "github.com/gofiber/fiber/v2"

// SetupRoutes registers problem solution endpoints.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Problem solution routes
	api.Post("/", handler.CreateProblemSolution)
	api.Get("/", handler.ListProblemSolutions)
	api.Get("/featured", handler.ListFeaturedProblemSolutions) // Public access
	api.Get("/matrix", handler.GetProblemSolutionMatrix)        // Public access - Problem-Solution Matrix
	api.Get("/:id", handler.GetProblemSolution)
	api.Get("/:id/public", handler.GetProblemSolutionPublic) // Public access
	api.Patch("/:id", handler.UpdateProblemSolution)
	api.Delete("/:id", handler.DeleteProblemSolution)
}

