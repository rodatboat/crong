package routes

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Initialize route handlers
	FolderRoutes(api)
	JobsRoutes(api)
}
