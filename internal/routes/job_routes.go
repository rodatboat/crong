package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/container"
	"github.com/rodatboat/crong/internal/handlers"
	"github.com/rodatboat/crong/internal/middleware"
)

func JobsRoutes(app fiber.Router, serviceContainer *container.Container) {
	jobs := app.Group("/jobs")

	jobs.Post("/", middleware.Protected(), handlers.NewJobHandler(serviceContainer.JobService).CreateJob)
	jobs.Get("/", middleware.Protected(), handlers.NewJobHandler(serviceContainer.JobService).ReadJobs)
	jobs.Put("/:id", middleware.Protected(), handlers.NewJobHandler(serviceContainer.JobService).UpdateJob)
	jobs.Delete("/:id", middleware.Protected(), handlers.NewJobHandler(serviceContainer.JobService).DeleteJob)
}
