package services

import (
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/repositories"
)

type JobExecutionService struct {
	jobRepo         *repositories.JobRepository
	scheduleRepo    *repositories.ScheduleRepository
	folderService   *FolderService
	scheduleService *ScheduleService
}

func NewJobExecutionService(
	jobRepo *repositories.JobRepository,
	scheduleRepo *repositories.ScheduleRepository,
	folderService *FolderService,
	scheduleService *ScheduleService,
) *JobExecutionService {
	return &JobExecutionService{
		jobRepo:         jobRepo,
		scheduleRepo:    scheduleRepo,
		folderService:   folderService,
		scheduleService: scheduleService,
	}
}

func (s *JobExecutionService) CreateJobExecution(jobID uint, jobExecution models.JobExecution) error {
	// TODO: Insert a new record into the job_executions table with the provided execution details
	return nil
}
