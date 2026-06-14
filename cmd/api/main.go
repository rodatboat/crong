package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/routes"
)

func main() {
	// Initialize database connection
	// Initialize repositories

	// Initialize routes
	app := fiber.New()
	routes.RegisterRoutes(app)

	// Start the server
	app.Listen(":3000")
}
