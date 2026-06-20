package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/container"
	"github.com/rodatboat/crong/internal/handlers"
	"github.com/rodatboat/crong/internal/middleware"
)

func JobsRoutes(app fiber.Router, serviceContainer *container.Container) {
	jobs := app.Group("/jobs")
	handler := handlers.NewJobHandler(serviceContainer.JobService)

	jobs.Post("/", middleware.Protected(), handler.CreateJob)
	jobs.Get("/", middleware.Protected(), handler.ReadJobs)
	jobs.Put("/:id", middleware.Protected(), handler.UpdateJob)
	jobs.Delete("/:id", middleware.Protected(), handler.DeleteJob)
}
