package routes

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Initialize route handlers
	FolderRoutes(api)
	JobsRoutes(api)

	// Catch-all route for undefined endpoints
	app.Use(func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound) // => 404 "Not Found"
	})
}
