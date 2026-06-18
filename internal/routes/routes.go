package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/container"
)

func RegisterRoutes(app *fiber.App, svc *container.Container) {
	api := app.Group("/api")

	// Initialize route handlers
	UserRoutes(api, svc)
	FolderRoutes(api, svc)
	JobsRoutes(api, svc)

	// Catch-all route for undefined endpoints
	app.Use(func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound) // => 404 "Not Found"
	})
}
