package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/handlers"
	"github.com/rodatboat/crong/internal/middleware"
)

func JobsRoutes(app fiber.Router) {
	jobs := app.Group("/jobs")

	jobs.Post("/", middleware.Protected(), handlers.CreateJob)
	jobs.Get("/", middleware.Protected(), handlers.ReadJobs)
	jobs.Put("/:id", middleware.Protected(), handlers.UpdateJob)
	jobs.Delete("/:id", middleware.Protected(), handlers.DeleteJob)
}
