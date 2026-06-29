package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/container"
	"github.com/rodatboat/crong/internal/handlers"
	"github.com/rodatboat/crong/internal/middleware"
)

func JobsRoutes(app fiber.Router, serviceContainer *container.Container) {
	jobs := app.Group("/jobs")
	handler := handlers.NewJobHandler(serviceContainer.JobService, serviceContainer.JobExecutionService)

	jobs.Post("/", middleware.Protected(serviceContainer.UserRepository), handler.CreateJob)
	jobs.Get("/:id/run", middleware.Protected(serviceContainer.UserRepository), handler.RunJob)
	jobs.Get("/:id/executions", middleware.Protected(serviceContainer.UserRepository), handler.GetJobExecutions)
	jobs.Get("/", middleware.Protected(serviceContainer.UserRepository), handler.ReadJobs)
	jobs.Get("/:id", middleware.Protected(serviceContainer.UserRepository), handler.GetJobsDetailsByID)
	jobs.Put("/:id", middleware.Protected(serviceContainer.UserRepository), handler.UpdateJob)
	jobs.Delete("/:id", middleware.Protected(serviceContainer.UserRepository), handler.DeleteJob)
}
