package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupCORS configures CORS middleware for the Fiber app
// Note: Logger middleware is now handled by pkg/logger
func SetupCORS(app *fiber.App, corsCfg *CORSConfig) {
	if corsCfg.Enabled {
		app.Use(cors.New(cors.Config{
			AllowOrigins:     corsCfg.AllowedOrigins,
			AllowMethods:     corsCfg.AllowedMethods,
			AllowHeaders:     corsCfg.AllowedHeaders,
			ExposeHeaders:    corsCfg.ExposedHeaders,
			AllowCredentials: corsCfg.AllowCredentials,
			MaxAge:           corsCfg.MaxAge,
		}))
	}
}

// SetupMiddleware is deprecated - use SetupCORS instead
// Kept for backward compatibility
func SetupMiddleware(app *fiber.App, corsCfg *CORSConfig) {
	SetupCORS(app, corsCfg)
}
