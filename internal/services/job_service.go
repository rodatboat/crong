package services

import (
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/repositories"
)

/**
 * This file contains the core business logic for managing jobs, including:
 * - Creating, updating, deleting jobs
 * - Validating cron expressions and converting them to schedule structs
 * - Interacting with the database layer to persist job data
 * - Handling any complex logic related to job schedule syncing based on cron expression
 */

type JobService struct {
	jobRepo *repositories.JobRepository
}

func NewJobService(jobRepo *repositories.JobRepository) *JobService {
	return &JobService{
		jobRepo: jobRepo,
	}
}

func (s *JobService) GetJobsByUser(userID uint) ([]*models.Job, error) {
	// Call repository layer
	_, err := s.jobRepo.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	// TODO: Convert entities to models and apply business logic
	return nil, nil
}

func (s *JobService) CreateJob(userID uint, req *models.JobCreateRequest) (*models.Job, error) {
	// TODO: Validate job request
	// TODO: Convert cron expression to schedule
	// TODO: Create job in repository
	return nil, nil
}

func (s *JobService) UpdateJob(jobID uint, userID uint, req *models.JobUpdateRequest) (*models.Job, error) {
	// TODO: Verify job belongs to user
	// TODO: Update schedule tables if cron changed
	return nil, nil
}

func (s *JobService) DeleteJob(jobID uint, userID uint) error {
	// TODO: Verify job belongs to user
	// TODO: Delete job and associated schedules
	return nil
}

func (s *JobService) CronExpressionToSchedule(cronExpr string) (models.Schedule, error) {
	// TODO: Validate cron expression and convert to Schedule struct
	return models.Schedule{}, nil
}

func (s *JobService) UpdateJobSchedule(jobID uint, jobSchedule models.Schedule) error {
	// TODO: Update the schedule tables (minute, hour, mday, wday, month) based on the provided Schedule struct
	return nil
}

func (s *JobService) CreateJobExecution(jobID uint, jobExecution models.JobExecution) error {
	// TODO: Insert a new record into the job_executions table with the provided execution details
	return nil
}
