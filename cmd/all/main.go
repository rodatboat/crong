package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/rodatboat/crong/internal/config"
	"github.com/rodatboat/crong/internal/container"
	"github.com/rodatboat/crong/internal/database"
	"github.com/rodatboat/crong/internal/routes"
	"github.com/rodatboat/crong/internal/scheduler"
)

const (
	workerCount = 5
	queueSize   = 100
)

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load() // Ignore error if .env doesn't exist

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db := database.InitDb(cfg)

	// Initialize dependency container (repositories, services)
	container := container.NewContainer(db)

	// Shared cancellation context (OS shutdown signal)
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	// Fiber API setup
	app := fiber.New()
	app.Use(cors.New())
	routes.RegisterRoutes(app, container)

	// Scheduler + Workers
	sched := scheduler.New(container.ScheduleRepository, queueSize)
	pool := scheduler.NewWorkerPool(
		sched.JobQueue(),
		container.JobExecutionService,
		workerCount,
	)

	// Run background system
	go func() {
		log.Info("Scheduler + worker pool starting")
		pool.Start(ctx)
		log.Info("Worker pool stopped")
	}()

	go func() {
		log.Info("Scheduler starting")
		sched.Start(ctx)
		log.Info("Scheduler stopped")
	}()

	// Run HTTP server (blocking)
	go func() {
		log.Infof("API starting on port %s", cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Errorf("API error: %v", err)
			cancel()
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	log.Info("Shutdown signal received")

	_ = app.Shutdown()
}
