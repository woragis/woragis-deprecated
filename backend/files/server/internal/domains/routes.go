package files

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"woragis-files-service/internal/domains/files"
	"woragis-files-service/pkg/authservice"
	"woragis-files-service/pkg/middleware"
)

// SetupRoutes sets up all files service routes
func SetupRoutes(api fiber.Router, db *gorm.DB, authServiceURL string, storageBasePath string, logger *slog.Logger) {
	// Initialize Auth Service client
	authClient := authservice.NewClient(authServiceURL)

	// Apply auth validation middleware to all routes
	api.Use(middleware.AuthValidationMiddleware(middleware.DefaultAuthValidationConfig(authClient)))

	// Initialize repositories
	fileRepo := files.NewGormRepository(db)

	// Initialize services (storage provider nil for now - can be S3, local filesystem, etc.)
	fileService := files.NewService(fileRepo, nil, storageBasePath, logger)

	// Initialize handlers
	fileHandler := files.NewHandler(fileService, logger)

	// Setup routes
	files.SetupRoutes(api.Group("/files"), fileHandler)
}
