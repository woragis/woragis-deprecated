package socialmedia

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	// "woragis-social-media-service/internal/domains/socialmediaposts"
	// "woragis-social-media-service/internal/domains/socialmediaposts/analytics"
	// "woragis-social-media-service/internal/domains/socialmediaposts/assets"
	// "woragis-social-media-service/internal/domains/socialmediaposts/content"
	// "woragis-social-media-service/internal/domains/socialmediaposts/platforms"
	// "woragis-social-media-service/internal/domains/socialmediaposts/scheduling"
	"woragis-social-media-service/internal/domains/creativeassets"
	"woragis-social-media-service/pkg/authservice"
	"woragis-social-media-service/pkg/middleware"
)

// SetupRoutes sets up all social media service routes
func SetupRoutes(api fiber.Router, db *gorm.DB, authServiceURL string, logger *slog.Logger) {
	// Initialize Auth Service client
	authClient := authservice.NewClient(authServiceURL)

	// Apply auth validation middleware to all routes
	api.Use(middleware.AuthValidationMiddleware(middleware.DefaultAuthValidationConfig(authClient)))

	// Initialize repositories
	// TODO: Re-enable once type incompatibilities are resolved
	// socialMediaPostRepo := socialmediaposts.NewGormRepository(db)
	creativeAssetRepo := creativeassets.NewRepository(db)
	
	// Initialize subdomain repositories for social media posts
	// TODO: Re-enable once type incompatibilities are resolved
	// analyticsRepo := analytics.NewGormRepository(db)
	// assetsRepo := assets.NewGormRepository(db)
	// contentRepo := content.NewGormRepository(db)
	// platformsRepo := platforms.NewGormRepository(db)
	// schedulingRepo := scheduling.NewGormRepository(db)

	// Initialize services
	// TODO: Re-enable once type incompatibilities are resolved
	// socialMediaPostService := socialmediaposts.NewService(socialMediaPostRepo, logger)
	creativeAssetService := creativeassets.NewService(creativeAssetRepo, nil) // creativeClient nil for now
	
	// Initialize subdomain services for social media posts
	// TODO: Re-enable once type incompatibilities are resolved
	// analyticsService := analytics.NewService(analyticsRepo, logger)
	// assetsService := assets.NewService(assetsRepo, logger)
	// TODO: Fix type incompatibility between content/scheduling and main socialmediaposts types
	// contentService := content.NewService(contentRepo, socialMediaPostRepo, socialMediaPostService, logger)
	// platformsService := platforms.NewService(platformsRepo, logger)
	// schedulingService := scheduling.NewService(schedulingRepo, socialMediaPostRepo, socialMediaPostService, platformsService, logger)

	// Initialize handlers (simplified - without translation enricher for now)
	// TODO: Re-enable once type incompatibilities are resolved
	// socialMediaPostHandler := socialmediaposts.NewHandler(socialMediaPostService, nil, nil, logger) // enricher, translationService
	creativeAssetHandler := creativeassets.NewHandler(creativeAssetService, logger)
	
	// Initialize subdomain handlers for social media posts
	// TODO: Re-enable once type incompatibilities are resolved
	// analyticsHandler := analytics.NewHandler(analyticsService, logger)
	// assetsHandler := assets.NewHandler(assetsService, logger)
	// contentHandler := content.NewHandler(contentService, logger)
	// platformsHandler := platforms.NewHandler(platformsService, logger)
	// schedulingHandler := scheduling.NewHandler(schedulingService, logger)

	// Setup routes
	// TODO: Fix SetupRoutes call once handlers are properly initialized and type incompatibilities are resolved
	// socialMediaPostsGroup := api.Group("/social-media-posts")
	// socialmediaposts.SetupRoutes(socialMediaPostsGroup, socialMediaPostHandler, platformsHandler, contentHandler, schedulingHandler, analyticsHandler, assetsHandler)
	
	// Setup creative assets routes - need to pass app and auth middleware
	// For now, we'll use a simpler approach - register routes directly
	creativeAssetsGroup := api.Group("/creative-assets")
	creativeAssetsGroup.Get("/:id/data", creativeAssetHandler.GetAssetData)
	creativeAssetsGroup.Post("/", creativeAssetHandler.CreateAsset)
	creativeAssetsGroup.Get("/:id", creativeAssetHandler.GetAsset)
	creativeAssetsGroup.Delete("/:id", creativeAssetHandler.DeleteAsset)
	creativeAssetsGroup.Get("/entity/:entityType/:entityId", creativeAssetHandler.GetAssetsByEntity)
	creativeAssetsGroup.Get("/entity/:entityType/:entityId/purpose", creativeAssetHandler.GetAssetByEntityAndPurpose)
	creativeAssetsGroup.Post("/generate/image", creativeAssetHandler.GenerateImage)
	creativeAssetsGroup.Post("/generate/thumbnail", creativeAssetHandler.GenerateThumbnail)
	creativeAssetsGroup.Post("/generate/diagram", creativeAssetHandler.GenerateDiagram)
}
