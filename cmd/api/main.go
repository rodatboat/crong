package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/rodatboat/crong/internal/config"
	"github.com/rodatboat/crong/internal/container"
	"github.com/rodatboat/crong/internal/database"
	"github.com/rodatboat/crong/internal/routes"
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db := database.InitDb(cfg)

	// Initialize dependency container (repositories, services)
	serviceContainer := container.NewContainer(db)

	// Initialize routes with container
	app := fiber.New()
	app.Use(cors.New())

	routes.RegisterRoutes(app, serviceContainer)

	// Start the server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
