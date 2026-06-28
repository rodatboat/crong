package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/rodatboat/crong/internal/entities"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/repositories"
	"github.com/rodatboat/crong/internal/utils"
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

func (s *JobExecutionService) ExecuteJob(jobEntity entities.Job) (*entities.JobExecution, error) {
	log.Infof("Executing job %v", jobEntity.ID)

	job := utils.MapJobEntityToJobModel(&jobEntity)
	executionStartTs := time.Now()

	jobBatchId := fmt.Sprintf("%v-%v-%v-%v-%v",
		jobEntity.ID,
		executionStartTs.Year(),
		executionStartTs.Month(),
		executionStartTs.Day(),
		executionStartTs.Unix(),
	)
	// Initialize job execution entity
	jobExecution := &entities.JobExecution{
		BatchIdentifier: jobBatchId,
		JobID:           job.ID,
		ExecutionStatus: entities.EXECUTING,
		Url:             job.Url,
	}

	// Initialize http client
	client := &http.Client{
		Timeout: time.Duration(job.Timeout) * time.Second,
	}

	// Create request
	httpMethod := utils.ReqMethodToHTTPMethod(job.Method)
	req, err := http.NewRequest(
		httpMethod,
		job.Url,
		strings.NewReader(job.Body),
	)
	if err != nil {
		return nil, err
	}

	// Set headers
	for _, header := range job.Headers {
		req.Header.Set(header.Key, header.Value)
	}

	// Set auth
	if job.Auth.Enabled && job.Auth.Username != "" && job.Auth.Password != "" {
		req.SetBasicAuth(job.Auth.Username, job.Auth.Password)
	}

	// Send request
	jobExecution.ExecutedAt = &executionStartTs
	resp, err := client.Do(req)
	finishedAt := time.Now()
	if err != nil {
		log.Errorf("Error executing job: %v", err)
		jobExecution.ExecutionStatus = entities.FAILED
		jobExecution.Error = err.Error()
	} else {
		// Update job execution
		jobExecution.ExecutionStatus = entities.COMPLETED

		// Update job execution
		jobExecution.StatusCode = resp.StatusCode
		jobExecution.StatusText = resp.Status
		jobExecution.DurationMs = int(finishedAt.Sub(executionStartTs).Milliseconds())

		// Read response
		respBody, err := json.Marshal(resp.Body)
		if err != nil {
			log.Errorf("Error reading response body: %v", err)
			return nil, err
		}
		jobExecution.ResponseBody = string(respBody)

		// Read headers
		respHeaders, err := json.Marshal(resp.Header)
		if err != nil {
			log.Errorf("Error reading response headers: %v", err)
			return nil, err
		}
		jobExecution.ResponseHeaders = string(respHeaders)
	}
	log.Infof("Finished executing job %v, status: %v", jobEntity.ID, resp.StatusCode)
	defer resp.Body.Close()

	if err := s.jobExecutionRepo.Create(jobExecution); err != nil {
		log.Errorf("Error creating job execution: %v", err)
		return nil, err
	}

	return jobExecution, nil
}

func (s *JobExecutionService) GetJobExecutionsByJobID(jobID uint) ([]*models.JobExecution, error) {
	// TODO: Retrieve all job executions for a given job ID from the job_executions table
	return nil, nil
}

func (s *JobExecutionService) CreateJobExecution(jobID uint, jobExecution models.JobExecution) error {
	// TODO: Insert a new record into the job_executions table with the provided execution details
	return nil
}
