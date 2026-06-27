package entities

import (
	"time"

	"gorm.io/datatypes"
)

type ReqMethod int

const (
	GET ReqMethod = iota
	POST
	PUT
	PATCH
	DELETE

	OPTIONS
	HEAD
	TRACE
	CONNECT
)

func (Job) TableName() string {
	return "jobs"
}

type Job struct {
	ID       uint      `gorm:"column:id;primaryKey"`
	Title    string    `gorm:"column:title;not null"`
	Url      string    `gorm:"column:url;not null"`
	FolderID uint      `gorm:"column:folder_id"`
	UserID   uint      `gorm:"column:user_id;not null"`
	Method   ReqMethod `gorm:"column:method;type:smallint;not null"`

	Headers datatypes.JSON `gorm:"column:headers"`
	Auth    datatypes.JSON `gorm:"column:auth"`
	Body    string         `gorm:"column:body"`
	Cron    string         `gorm:"column:cron;not null"`

	Timezone string `gorm:"column:timezone;not null;default:'UTC'"`

	Timeout int  `gorm:"column:timeout;default:30"`
	Enabled bool `gorm:"column:enabled;not null;default:true"`

	LastExecution time.Time `gorm:"column:last_execution"`
	CreatedAt     time.Time `gorm:"column:created_at;default:now()"`
	UpdatedAt     time.Time `gorm:"column:updated_at;default:now()"`
}

func (JobExecution) TableName() string {
	return "job_executions"
}

type JobExecution struct {
	JobID uint `gorm:"index;not null"`

	ExecutionSuccess bool   `gorm:"column:exec_success"`
	StatusCode       int    `gorm:"column:status_code"`
	DurationMs       int    `gorm:"column:duration_ms"`
	Url              string `gorm:"column:url"`
	BatchIdentifier  string `gorm:"column:batch_identifier"`

	ResponseHeaders string    `gorm:"column:response_headers;type:text"`
	ResponseBody    string    `gorm:"column:response_body;type:text"`
	Error           string    `gorm:"column:error;type:text"`
	ExecutedAt      time.Time `gorm:"column:executed_at;"`
}
