package posts

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes registers post-related routes.
func SetupRoutes(api fiber.Router, handler Handler) {
	// Use the provided router directly (it's already a group with the correct path)
	
	// Post routes
	api.Post("/", handler.CreatePost)
	api.Get("/", handler.ListPosts)
	api.Get("/slug/:slug", handler.GetPostBySlug)
	api.Get("/:id", handler.GetPost)
	api.Patch("/:id", handler.UpdatePost)
	api.Delete("/:id", handler.DeletePost)

	// Post relationship routes
	api.Get("/:id/skills", handler.GetPostSkills)
	api.Post("/:id/skills", handler.AttachSkillToPost)
	api.Delete("/:id/skills/:skillId", handler.DetachSkillFromPost)
	
	api.Get("/:id/categories", handler.GetPostCategories)
	api.Post("/:id/categories", handler.AttachCategoryToPost)
	api.Delete("/:id/categories/:categoryId", handler.DetachCategoryFromPost)
	
	api.Get("/:id/tags", handler.GetPostTags)
	api.Post("/:id/tags", handler.AttachTagToPost)
	api.Delete("/:id/tags/:tagId", handler.DetachTagFromPost)

	// Category routes
	api.Post("/categories", handler.CreateCategory)
	api.Get("/categories", handler.ListCategories)
	api.Get("/categories/slug/:slug", handler.GetCategoryBySlug)
	api.Get("/categories/:id", handler.GetCategory)
	api.Patch("/categories/:id", handler.UpdateCategory)

	// Tag routes
	api.Get("/tags", handler.ListTags)
	api.Get("/tags/slug/:slug", handler.GetTagBySlug)
	api.Get("/tags/:id", handler.GetTag)

	// Creative assets routes
	api.Post("/:id/assets/generate/thumbnail", handler.GeneratePostThumbnail)
	api.Post("/:id/assets/generate/featured-image", handler.GeneratePostFeaturedImage)
	api.Post("/:id/assets/generate/og-image", handler.GeneratePostOGImage)
	api.Get("/:id/assets", handler.GetPostAssets)
}

