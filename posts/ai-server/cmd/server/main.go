package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/woragis/posts-ai-service/internal/config"
	"github.com/woragis/posts-ai-service/internal/db"
	"github.com/woragis/posts-ai-service/internal/handlers"
	"github.com/woragis/posts-ai-service/internal/services"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	dbpool, err := db.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Run migrations
	if err := db.RunMigrations(context.Background(), dbpool); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services
	aiSvc := services.NewAIService(cfg.AIServiceURL)
	chatSvc := services.NewChatService(dbpool, aiSvc)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Posts AI Service v1.0.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	}))

	// Health check
	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"service": "posts-ai-service",
		})
	})

	// Routes
	api := app.Group("/api/v1")

	// Chat endpoints
	api.Post("/chats/generate", handlers.GenerateDraft(chatSvc))
	api.Post("/posts/:id/ai/improve", handlers.ImproveContent(chatSvc))
	api.Get("/chats/:id", handlers.GetChat(chatSvc))
	api.Get("/chats", handlers.ListChats(chatSvc))
	api.Get("/usage/stats", handlers.GetUsageStats(chatSvc))

	// WebSocket endpoint
	app.Get("/ws/chats/:id", handlers.ChatWebSocket(chatSvc))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", cfg.Port)); err != nil {
			if err.Error() != "the server closed the listener" {
				log.Fatalf("Server error: %v", err)
			}
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
}
