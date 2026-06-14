package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/handlers"
)

func JobsRoutes(app fiber.Router) {
	jobs := app.Group("/jobs")

	jobs.Post("/", handlers.CreateJob)
	jobs.Get("/", handlers.ReadJobs)
	jobs.Put("/:id", handlers.UpdateJob)
	jobs.Delete("/:id", handlers.DeleteJob)
}
