package skills

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"woragis-skills-service/internal/domains/skills"
	"woragis-skills-service/internal/domains/interests"
	"woragis-skills-service/pkg/authservice"
	"woragis-skills-service/pkg/middleware"
)

// SetupRoutes sets up all skills service routes
func SetupRoutes(api fiber.Router, db *gorm.DB, authServiceURL string, logger *slog.Logger) {
	// Initialize Auth Service client
	authClient := authservice.NewClient(authServiceURL)

	// Apply auth validation middleware to all routes
	api.Use(middleware.AuthValidationMiddleware(middleware.DefaultAuthValidationConfig(authClient)))

	// Initialize repositories
	skillRepo := skills.NewGormRepository(db)
	interestRepo := interests.NewGormRepository(db)

	// Initialize services
	skillService := skills.NewService(skillRepo, logger)
	interestService := interests.NewService(interestRepo, logger)

	// Initialize handlers (simplified - without translation enricher for now)
	skillHandler := skills.NewHandler(skillService, nil, nil, logger)
	interestHandler := interests.NewHandler(interestService, nil, nil, logger)

	// Setup routes
	skills.SetupRoutes(api.Group("/skills"), skillHandler)
	interests.SetupRoutes(api.Group("/interests"), interestHandler)
}
