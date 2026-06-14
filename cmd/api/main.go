package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/rodatboat/crong/internal/database"
	"github.com/rodatboat/crong/internal/routes"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database connection
	database.InitDb()

	// Initialize repositories

	// Initialize routes
	app := fiber.New()
	app.Use(cors.New())

	routes.RegisterRoutes(app)

	// Start the server
	app.Listen(":3000")
}
