package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rodatboat/crong/internal/config"
	"github.com/rodatboat/crong/internal/container"
	"github.com/rodatboat/crong/internal/database"
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

	// Initialize database connection
	db := database.InitDb(cfg)

	// Initialize dependency container (repositories, services)
	serviceContainer := container.NewContainer(db)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	sched := scheduler.New(serviceContainer.ScheduleRepository, queueSize)
	pool := scheduler.NewWorkerPool(sched.JobQueue(), serviceContainer.JobExecutionService, workerCount)

	// Start workers in background; they drain the queue until it is closed.
	go pool.Start(ctx)

	log.Println("Scheduler process started")
	sched.Start(ctx) // blocks until ctx is cancelled, then closes the queue
	log.Println("Scheduler process exiting")
}
