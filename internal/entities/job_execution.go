package entities

import "time"

func (JobExecution) TableName() string {
	return "job_executions"
}

type ExecutionStatus int

const (
	SUCCESS ExecutionStatus = iota
	EXECUTING
	FAILED
)

type JobExecution struct {
	ID              uint   `gorm:"column:id;primaryKey"`
	JobID           uint   `gorm:"index;not null"`
	BatchIdentifier string `gorm:"column:batch_identifier;unique"`

	ExecutionStatus ExecutionStatus `gorm:"column:exec_success"`
	StatusCode      int             `gorm:"column:status_code"`
	StatusText      string          `gorm:"column:status_text"`
	DurationMs      int             `gorm:"column:duration_ms"`
	Url             string          `gorm:"column:url"`

	ResponseHeaders string    `gorm:"column:response_headers;type:text"`
	ResponseBody    string    `gorm:"column:response_body;type:text"`
	Error           string    `gorm:"column:error;type:text"`
	ExecutedAt      time.Time `gorm:"column:executed_at;"`
	PlannedFor      time.Time `gorm:"column:planned_for;"`
	CreatedAt       time.Time `gorm:"column:created_at;default:now()"`
	UpdatedAt       time.Time `gorm:"column:updated_at;default:now()"`
}
