package socialmediaposts

import (
	"github.com/gofiber/fiber/v2"
	
	"woragis-social-media-service/internal/domains/socialmediaposts/analytics"
	"woragis-social-media-service/internal/domains/socialmediaposts/assets"
	"woragis-social-media-service/internal/domains/socialmediaposts/content"
	"woragis-social-media-service/internal/domains/socialmediaposts/platforms"
	"woragis-social-media-service/internal/domains/socialmediaposts/scheduling"
)

// SetupRoutes registers social media post endpoints and mounts subdomain routes.
func SetupRoutes(
	api fiber.Router,
	handler Handler,
	platformsHandler platforms.Handler,
	contentHandler content.Handler,
	schedulingHandler scheduling.Handler,
	analyticsHandler analytics.Handler,
	assetsHandler assets.Handler,
) {
	// Post CRUD operations
	api.Post("/", handler.CreatePost)
	api.Get("/", handler.ListPosts)
	api.Get("/by-url", handler.GetPostByURL)
	api.Get("/:id", handler.GetPost)
	api.Patch("/:id", handler.UpdatePost)
	api.Patch("/:id/status", handler.UpdatePostStatus)
	api.Patch("/:id/engagement", handler.UpdatePostEngagement)
	api.Delete("/:id", handler.DeletePost)

	// Link operations
	api.Post("/links", handler.CreateLink)
	api.Patch("/links/:id", handler.UpdateLink)
	api.Delete("/links/:id", handler.DeleteLink)
	api.Get("/posts/:postId/links", handler.GetLinksByPost)
	api.Get("/posts/:postId/entities", handler.GetEntitiesByPost)
	api.Get("/links/by-entity", handler.GetLinksByEntity)
	api.Get("/posts/by-entity", handler.GetPostsByEntity)

	// Subdomain routes
	platforms.SetupRoutes(api.Group("/platforms"), platformsHandler)
	content.SetupRoutes(api.Group("/content"), contentHandler)
	scheduling.SetupRoutes(api.Group("/scheduling"), schedulingHandler)
	analytics.SetupRoutes(api.Group("/analytics"), analyticsHandler)
	assets.SetupRoutes(api.Group("/assets"), assetsHandler)
}

