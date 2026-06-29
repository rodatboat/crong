package utils

import (
	"encoding/json"

	"github.com/rodatboat/crong/internal/entities"
	"github.com/rodatboat/crong/internal/models"
	"gorm.io/datatypes"
)

/**
 * mapJobCreateRequestToEntity converts a JobCreateRequest to an entities.Job
 */
func MapJobCreateRequestToEntity(userID uint, req *models.JobCreateRequest) *entities.Job {
	return &entities.Job{
		Title:    req.Title,
		Url:      req.Url,
		FolderID: req.FolderID,
		UserID:   userID,

		Method:  req.Method,
		Headers: ConvertHeadersToJSON(req.Headers),
		Auth:    ConvertAuthToJSON(req.Auth),
		Body:    req.Body,
		Cron:    req.Cron,

		Timezone: req.Timezone,
		Timeout:  req.Timeout,
		Enabled:  req.Enabled,
	}
}

func MapJobEntityToJobModel(jobEntity *entities.Job) *models.Job {
	return &models.Job{
		ID:       jobEntity.ID,
		Title:    jobEntity.Title,
		Url:      jobEntity.Url,
		FolderID: jobEntity.FolderID,

		Method:  jobEntity.Method,
		Headers: ConvertHeadersJSONToHeadersModel(jobEntity.Headers),
		Auth:    ConvertAuthJSONToAuthModel(jobEntity.Auth),
		Body:    jobEntity.Body,
		Cron:    jobEntity.Cron,

		Timezone: jobEntity.Timezone,
		Timeout:  jobEntity.Timeout,
		Enabled:  jobEntity.Enabled,

		LastExecution: jobEntity.LastExecution,
		CreatedAt:     jobEntity.CreatedAt,
		UpdatedAt:     jobEntity.UpdatedAt,
	}
}

func MapJobExecutionEntityToJobExecutionModel(jobExecutionEntity *entities.JobExecution) *models.JobExecution {
	return &models.JobExecution{
		ID:              jobExecutionEntity.ID,
		JobID:           jobExecutionEntity.JobID,
		BatchIdentifier: jobExecutionEntity.BatchIdentifier,

		ExecutionStatus: jobExecutionEntity.ExecutionStatus,
		StatusCode:      jobExecutionEntity.StatusCode,
		StatusText:      jobExecutionEntity.StatusText,
		DurationMs:      jobExecutionEntity.DurationMs,
		Url:             jobExecutionEntity.Url,

		ResponseHeaders: jobExecutionEntity.ResponseHeaders,
		ResponseBody:    jobExecutionEntity.ResponseBody,
		Error:           jobExecutionEntity.Error,

		ExecutedAt: jobExecutionEntity.ExecutedAt,
		PlannedFor: jobExecutionEntity.PlannedFor,
	}
}

/**
 * convertHeadersToJSON converts []models.JobHeaders to datatypes.JSON
 */
func ConvertHeadersToJSON(headers []models.JobHeaders) datatypes.JSON {
	if len(headers) == 0 {
		return nil
	}
	data, _ := json.Marshal(headers)
	return datatypes.JSON(data)
}

/**
 * convertHeadersToJSON converts datatypes.JSON to []models.JobHeaders
 */
func ConvertHeadersJSONToHeadersModel(headersJSON datatypes.JSON) []models.JobHeaders {
	var headers []models.JobHeaders
	json.Unmarshal([]byte(headersJSON), &headers)
	return headers
}

/**
 * convertAuthToJSON converts models.JobAuth to datatypes.JSON
 */
func ConvertAuthToJSON(auth models.JobAuth) datatypes.JSON {
	data, _ := json.Marshal(auth)
	return datatypes.JSON(data)
}

/**
 * convertAuthJSONToAuthModel converts datatypes.JSON to models.JobAuth
 */
func ConvertAuthJSONToAuthModel(authJSON datatypes.JSON) models.JobAuth {
	var auth models.JobAuth
	json.Unmarshal([]byte(authJSON), &auth)
	return auth
}
