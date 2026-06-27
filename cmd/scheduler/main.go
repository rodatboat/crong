package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/rodatboat/crong/internal/config"
	"github.com/rodatboat/crong/internal/container"
	"github.com/rodatboat/crong/internal/database"
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
	log.Printf("Service container initialized: %+v", serviceContainer)
}
