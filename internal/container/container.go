package container

import (
	"github.com/rodatboat/crong/internal/repositories"
	"github.com/rodatboat/crong/internal/services"
	"gorm.io/gorm"
)

// Container holds all repositories and services
type Container struct {
	// Repositories
	JobRepository      *repositories.JobRepository
	UserRepository     *repositories.UserRepository
	FolderRepository   *repositories.FolderRepository
	ScheduleRepository *repositories.ScheduleRepository

	// Services
	JobService      *services.JobService
	UserService     *services.UserService
	FolderService   *services.FolderService
	ScheduleService *services.ScheduleService
}

// NewContainer initializes all dependencies
func NewContainer(db *gorm.DB) *Container {
	// Initialize repositories (order matters - ScheduleRepository before JobRepository)
	userRepo := repositories.NewUserRepository(db)
	folderRepo := repositories.NewFolderRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
	jobRepo := repositories.NewJobRepository(db, scheduleRepo)

	// Initialize services with their dependencies
	userService := services.NewUserService(userRepo)
	folderService := services.NewFolderService(folderRepo)
	scheduleService := services.NewScheduleService()
	jobService := services.NewJobService(jobRepo, scheduleRepo, folderService, scheduleService)

	return &Container{
		JobRepository:      jobRepo,
		UserRepository:     userRepo,
		FolderRepository:   folderRepo,
		ScheduleRepository: scheduleRepo,
		JobService:         jobService,
		UserService:        userService,
		FolderService:      folderService,
		ScheduleService:    scheduleService,
	}
}
