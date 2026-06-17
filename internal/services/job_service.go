package services

import "github.com/rodatboat/crong/internal/models"

/**
 * This file contains the core business logic for managing jobs, including:
 * - Creating, updating, deleting jobs
 * - Validating cron expressions and converting them to schedule structs
 * - Interacting with the database layer to persist job data
 * - Handling any complex logic related to job schedule syncing based on cron expression
 */
func CronExpressionToSchedule(cronExpr string) (models.Schedule, error) {
	// TODO: Validate cron expression and convert to Schedule struct
	return models.Schedule{}, nil
}

func UpdateJobSchedule(jobID uint, jobSchedule models.Schedule) error {
	// TODO: Update the schedule tables (minute, hour, mday, wday, month) based on the provided Schedule struct
	return nil
}

func CreateJobExecution(jobID uint, jobExecution models.JobExecution) error {
	// TODO: Insert a new record into the job_executions table with the provided execution details
	return nil
}
