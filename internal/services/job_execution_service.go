package services

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3/log"
	"github.com/rodatboat/crong/internal/entities"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/repositories"
	"github.com/rodatboat/crong/internal/resp"
	"github.com/rodatboat/crong/internal/utils"
	"gorm.io/gorm"
)

type JobExecutionService struct {
	jobExecutionRepo *repositories.JobExecutionRepository
	jobRepo          *repositories.JobRepository
}

func NewJobExecutionService(
	jobExecutionRepo *repositories.JobExecutionRepository,
	jobRepo *repositories.JobRepository,
) *JobExecutionService {
	return &JobExecutionService{
		jobExecutionRepo: jobExecutionRepo,
		jobRepo:          jobRepo,
	}
}

// captureTraceStats returns a context and JobExecutionStats struct that captures HTTP timing metrics
func captureTraceStats() (context.Context, *entities.JobExecutionStats) {
	stats := &entities.JobExecutionStats{}

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) {
			stats.DNSLookupMs = -int(time.Now().UnixMilli())
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			stats.DNSLookupMs += int(time.Now().UnixMilli())
		},
		ConnectStart: func(_, _ string) {
			stats.TCPConnectMs = -int(time.Now().UnixMilli())
		},
		ConnectDone: func(_, _ string, _ error) {
			stats.TCPConnectMs += int(time.Now().UnixMilli())
		},
		TLSHandshakeStart: func() {
			stats.TLSHandshakeMs = -int(time.Now().UnixMilli())
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, _ error) {
			stats.TLSHandshakeMs += int(time.Now().UnixMilli())
		},
		WroteRequest: func(_ httptrace.WroteRequestInfo) {
			stats.RequestWriteMs = -int(time.Now().UnixMilli())
		},
		GotFirstResponseByte: func() {
			stats.TimeToFirstByteMs = -int(time.Now().UnixMilli())
			stats.RequestWriteMs += int(time.Now().UnixMilli())
			stats.TimeToFirstByteMs += int(time.Now().UnixMilli())
		},
	}

	ctx := httptrace.WithClientTrace(context.Background(), trace)
	return ctx, stats
}

func (s *JobExecutionService) ExecuteJob(jobEntity entities.Job) (*entities.JobExecution, error) {
	log.Infof("Executing job %v", jobEntity.ID)

	job := utils.MapJobEntityToJobModel(&jobEntity)
	executionStartTs := time.Now()

	jobBatchId := fmt.Sprintf("%v-%v-%v-%v-%v",
		jobEntity.ID,
		executionStartTs.Year(),
		int(executionStartTs.Month()),
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

	// Attach httptrace to capture HTTP timing metrics
	traceCtx, traceStats := captureTraceStats()
	req = req.WithContext(traceCtx)
	_ = traceStats // Available for future implementation

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
	jobExecution.PlannedFor = &executionStartTs
	resp, err := client.Do(req)
	finishedAt := time.Now()
	log.Infof("Job %v finished at %v, with stats %+v", job.ID, finishedAt, traceStats)

	if err != nil {
		log.Errorf("Error executing job: %v", err)
		jobExecution.ExecutionStatus = entities.FAILED
		jobExecution.Error = err.Error()
		jobExecution.DurationMs = int(finishedAt.Sub(executionStartTs).Milliseconds())
	} else {
		defer resp.Body.Close()

		// Update job execution
		jobExecution.ExecutionStatus = entities.COMPLETED
		jobExecution.StatusCode = resp.StatusCode
		jobExecution.StatusText = resp.Status
		jobExecution.DurationMs = int(finishedAt.Sub(executionStartTs).Milliseconds())

		// Read response body
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Error reading response body: %v", err)
			jobExecution.Error = err.Error()
			jobExecution.ExecutionStatus = entities.FAILED
		} else {
			jobExecution.ResponseBody = string(bodyBytes)
		}

		// Read headers
		respHeaders, err := json.Marshal(resp.Header)
		if err != nil {
			log.Errorf("Error reading response headers: %v", err)
		} else {
			jobExecution.ResponseHeaders = string(respHeaders)
		}

		log.Infof("Finished executing job %v, status: %v", jobEntity.ID, resp.StatusCode)
	}

	if err := s.jobExecutionRepo.Create(jobExecution); err != nil {
		log.Errorf("Error creating job execution: %v", err)
		return nil, err
	}

	return jobExecution, nil
}

func (s *JobExecutionService) GetJobExecutionsByJobID(jobID uint, userID uint, page int, limit int) ([]*models.JobExecution, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	jobExecutions, err := s.jobExecutionRepo.ListByJobID(jobID, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var result []*models.JobExecution
	for _, je := range jobExecutions {
		result = append(result, utils.MapJobExecutionEntityToJobExecutionModel(je))
	}

	return result, nil
}

func (s *JobExecutionService) CleanupOldExecutions(retentionDays int) error {
	log.Infof("Running cleanup job: removing job executions older than %d days", retentionDays)

	if err := s.jobExecutionRepo.DeleteOldExecutions(retentionDays); err != nil {
		log.Errorf("Error cleaning up old job executions: %v", err)
		return err
	}

	log.Infof("Cleanup job completed successfully")
	return nil
}

func (s *JobExecutionService) CreateJobExecution(jobID uint, userID uint) (*models.JobExecution, error) {
	log.Infof("Creating new job execution for job %v, user %v", jobID, userID)

	jobEntity, err := s.jobRepo.FindByJobID(jobID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, resp.ErrNotFound
		}
		return nil, err
	}

	// Run the job
	jobExecution, err := s.ExecuteJob(*jobEntity)
	if err != nil {
		return nil, err
	}
	return utils.MapJobExecutionEntityToJobExecutionModel(jobExecution), nil
}
