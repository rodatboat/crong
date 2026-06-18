package container

import (
	"github.com/rodatboat/crong/internal/repositories"
	"github.com/rodatboat/crong/internal/services"
	"gorm.io/gorm"
)

// Container holds all repositories and services
type Container struct {
	// Repositories
	JobRepository    *repositories.JobRepository
	UserRepository   *repositories.UserRepository
	FolderRepository *repositories.FolderRepository

	// Services
	JobService    *services.JobService
	UserService   *services.UserService
	FolderService *services.FolderService
}

// NewContainer initializes all dependencies
func NewContainer(db *gorm.DB) *Container {
	// Initialize repositories
	jobRepo := repositories.NewJobRepository(db)
	userRepo := repositories.NewUserRepository(db)
	folderRepo := repositories.NewFolderRepository(db)

	// Initialize services with their dependencies
	jobService := services.NewJobService(jobRepo)
	userService := services.NewUserService(userRepo)
	folderService := services.NewFolderService(folderRepo)

	return &Container{
		JobRepository:    jobRepo,
		UserRepository:   userRepo,
		FolderRepository: folderRepo,
		JobService:       jobService,
		UserService:      userService,
		FolderService:    folderService,
	}
}
