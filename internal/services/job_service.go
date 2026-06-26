package services

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v3/log"
	"github.com/rodatboat/crong/internal/entities"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/repositories"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

/**
 * This file contains the core business logic for managing jobs, including:
 * - Creating, updating, deleting jobs
 * - Validating cron expressions and converting them to schedule structs
 * - Interacting with the database layer to persist job data
 * - Handling any complex logic related to job schedule syncing based on cron expression
 */
type JobService struct {
	jobRepo         *repositories.JobRepository
	scheduleRepo    *repositories.ScheduleRepository
	folderService   *FolderService
	scheduleService *ScheduleService
}

func NewJobService(
	jobRepo *repositories.JobRepository,
	scheduleRepo *repositories.ScheduleRepository,
	folderService *FolderService,
	scheduleService *ScheduleService,
) *JobService {
	return &JobService{
		jobRepo:         jobRepo,
		scheduleRepo:    scheduleRepo,
		folderService:   folderService,
		scheduleService: scheduleService,
	}
}

func (s *JobService) GetJobsByUser(userID uint) ([]*models.Job, error) {
	log.Infof("Getting jobs for user: %v", userID)
	// Call repository layer
	jobEntities, err := s.jobRepo.FindByUser(userID)
	if err != nil {
		return nil, err
	}

	jobs := make([]*models.Job, len(jobEntities))
	for _, jobEntity := range jobEntities {
		jobs = append(jobs, s.mapJobEntityToJobModel(jobEntity))
	}
	return jobs, nil
}

func (s *JobService) GetJobsDetailsByID(jobID uint, userID uint) (*models.Job, error) {
	log.Infof("Getting job details: %v", jobID)
	// Call repository layer
	jobEntity, err := s.jobRepo.FindByJobID(jobID, userID)
	if err != nil {
		return nil, err
	}

	return s.mapJobEntityToJobModel(jobEntity), nil
}

func (s *JobService) CreateJob(userID uint, req *models.JobCreateRequest) (*models.Job, error) {
	log.Infof("Creating job: %v", req)

	// Validate user owns the folder, and that folder exists (if provided)
	if req.FolderID > 0 {
		if !s.folderService.FolderExists(req.FolderID, userID) {
			log.Errorf("Folder not found: %v", req.FolderID)
			return nil, fmt.Errorf("Folder not found")
		}
	}

	// Parse cron expression to schedule model
	scheduleModel, err := s.scheduleService.CronExpressionToSchedule(req.Cron)
	if err != nil {
		log.Errorf("Error parsing cron expression: %v", err)
		return nil, err
	}

	// Map request to job entity
	jobEntity := s.mapJobCreateRequestToEntity(userID, req)

	// Create job and schedule atomically using transaction callback
	err = s.jobRepo.WithTransaction(func(tx *gorm.DB) error {
		// Create the job
		if err := s.jobRepo.Create(tx, jobEntity); err != nil {
			return err
		}

		// Now we have the job ID - map schedule model to entity
		scheduleEntity := s.scheduleService.ScheduleModelToEntity(jobEntity.ID, scheduleModel)

		// Create all schedule entries in the same transaction
		if err := s.scheduleRepo.CreateSchedules(tx, scheduleEntity); err != nil {
			return fmt.Errorf("Failed to create job schedule: %v", err)
		}

		return nil
	})

	if err != nil {
		log.Errorf("Error creating job: %v", err)
		return nil, err
	}

	return s.mapJobEntityToJobModel(jobEntity), nil
}

func (s *JobService) UpdateJob(jobID uint, userID uint, req *models.JobUpdateRequest) (*models.Job, error) {
	log.Infof("Updating existing job: %v", req)

	// Validate folder exists if provided
	if req.FolderID > 0 {
		if !s.folderService.FolderExists(req.FolderID, userID) {
			log.Errorf("Folder not found: %v", req.FolderID)
			return nil, fmt.Errorf("Folder not found")
		}
	}

	// Delete old schedules and create new ones atomically
	err := s.jobRepo.WithTransaction(func(tx *gorm.DB) error {
		// TODO: Update job entity

		// TODO: Update schedule tables if cron changed. If job.cron == req.cron, do nothing

		// Parse cron expression to schedule model
		scheduleModel, err := s.scheduleService.CronExpressionToSchedule(req.Cron)
		if err != nil {
			log.Errorf("Error parsing cron expression: %v", err)
			return err
		}

		// Delete existing schedules
		if err := s.scheduleRepo.DeleteSchedulesByJobID(tx, jobID); err != nil {
			return err
		}

		// Create new schedules
		scheduleEntity := s.scheduleService.ScheduleModelToEntity(jobID, scheduleModel)
		if err := s.scheduleRepo.CreateSchedules(tx, scheduleEntity); err != nil {
			log.Errorf("Failed to updating job schedule: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Errorf("Error updating job: %v", err)
		return nil, err
	}

	// TODO: Convert entity back to model and return
	return nil, nil
}

func (s *JobService) DeleteJob(jobID uint, userID uint) error {
	// TODO: Verify job belongs to user
	return s.jobRepo.WithTransaction(func(tx *gorm.DB) error {
		// Delete job
		if err := tx.Delete(&entities.Job{}, jobID).Error; err != nil {
			return err
		}
		// Delete associated schedules
		return s.scheduleRepo.DeleteSchedulesByJobID(tx, jobID)
	})
}

func (s *JobService) CreateJobExecution(jobID uint, jobExecution models.JobExecution) error {
	// TODO: Insert a new record into the job_executions table with the provided execution details
	return nil
}

// ========== UTILITIES ==========

/**
 * mapJobCreateRequestToEntity converts a JobCreateRequest to an entities.Job
 */
func (s *JobService) mapJobCreateRequestToEntity(userID uint, req *models.JobCreateRequest) *entities.Job {
	return &entities.Job{
		Title:    req.Title,
		Url:      req.Url,
		FolderID: req.FolderID,
		UserID:   userID,
		Method:   req.Method,
		Headers:  convertHeadersToJSON(req.Headers),
		Auth:     convertAuthToJSON(req.Auth),
		Body:     req.Body,
		Cron:     req.Cron,
		Timezone: req.Timezone,
		Timeout:  req.Timeout,
		Enabled:  req.Enabled,
	}
}

/**
 * convertHeadersToJSON converts []models.JobHeaders to datatypes.JSON
 */
func convertHeadersToJSON(headers []models.JobHeaders) datatypes.JSON {
	if len(headers) == 0 {
		return nil
	}
	data, _ := json.Marshal(headers)
	return datatypes.JSON(data)
}

/**
 * convertAuthToJSON converts models.JobAuth to datatypes.JSON
 */
func convertAuthToJSON(auth models.JobAuth) datatypes.JSON {
	data, _ := json.Marshal(auth)
	return datatypes.JSON(data)
}

func (s *JobService) mapJobEntityToJobModel(jobEntity *entities.Job) *models.Job {
	return &models.Job{
		ID:            jobEntity.ID,
		Title:         jobEntity.Title,
		Url:           jobEntity.Url,
		FolderID:      jobEntity.FolderID,
		Method:        jobEntity.Method,
		Headers:       convertHeadersJSONToHeadersModel(jobEntity.Headers),
		Auth:          convertAuthJSONToAuthModel(jobEntity.Auth),
		Body:          jobEntity.Body,
		Cron:          jobEntity.Cron,
		Timezone:      jobEntity.Timezone,
		Timeout:       jobEntity.Timeout,
		Enabled:       jobEntity.Enabled,
		LastExecution: jobEntity.LastExecution,
		CreatedAt:     jobEntity.CreatedAt,
		UpdatedAt:     jobEntity.UpdatedAt,
	}
}

/**
 * convertAuthJSONToAuthModel converts datatypes.JSON to models.JobAuth
 */
func convertAuthJSONToAuthModel(authJSON datatypes.JSON) models.JobAuth {
	var auth models.JobAuth
	json.Unmarshal([]byte(authJSON), &auth)
	return auth
}

/**
 * convertHeadersToJSON converts datatypes.JSON to []models.JobHeaders
 */
func convertHeadersJSONToHeadersModel(headersJSON datatypes.JSON) []models.JobHeaders {
	var headers []models.JobHeaders
	json.Unmarshal([]byte(headersJSON), &headers)
	return headers
}
