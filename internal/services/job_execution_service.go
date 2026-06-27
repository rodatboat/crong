package services

import (
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/repositories"
)

type JobExecutionService struct {
	jobExecutionRepo *repositories.JobExecutionRepository
}

func NewJobExecutionService(
	jobExecutionRepo *repositories.JobExecutionRepository,
) *JobExecutionService {
	return &JobExecutionService{
		jobExecutionRepo: jobExecutionRepo,
	}
}

func (s *JobExecutionService) RunJob(jobID uint) error {
	return nil
}

func (s *JobExecutionService) GetJobExecutionsByJobID(jobID uint) ([]*models.JobExecution, error) {
	// TODO: Retrieve all job executions for a given job ID from the job_executions table
	return nil, nil
}

func (s *JobExecutionService) CreateJobExecution(jobID uint, jobExecution models.JobExecution) error {
	// TODO: Insert a new record into the job_executions table with the provided execution details
	return nil
}
