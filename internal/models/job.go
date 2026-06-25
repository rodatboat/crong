package models

import (
	"time"

	"github.com/rodatboat/crong/internal/entities"
)

type Job struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	FolderID uint   `json:"folder_id"`

	Method  entities.ReqMethod `json:"method"`
	Headers []JobHeaders       `json:"headers"`
	Auth    JobAuth            `json:"auth"`
	Body    string             `json:"body"`
	Cron    string             `json:"cron"`

	Timezone string `json:"timezone"`
	Timeout  int    `json:"timeout"`
	Enabled  bool   `json:"enabled"`

	LastExecution time.Time `json:"last_execution"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type JobCreateRequest struct {
	Title    string `json:"title" validate:"required"`
	Url      string `json:"url" validate:"required,url"`
	FolderID uint   `json:"folder_id"`

	Method  entities.ReqMethod `json:"method" validate:"validmethod"`
	Headers []JobHeaders       `json:"headers"`
	Auth    JobAuth            `json:"auth"`
	Body    string             `json:"body"`
	Cron    string             `json:"cron" validate:"required,min=9"`

	Timezone string `json:"timezone"`
	Timeout  int    `json:"timeout" validate:"required,max=30"`
	Enabled  bool   `json:"enabled"`
}

type JobUpdateRequest struct {
	JobCreateRequest
}

type JobAuth struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"username" validate:"required_if=Enabled true"`
	Password string `json:"password" validate:"required_if=Enabled true"`
}

type JobHeaders struct {
	Key   string `json:"key" validate:"max=255"`
	Value string `json:"value" validate:"max=1024"`
}

type JobExecution struct {
	ID              uint      `json:"id"`
	JobID           uint      `json:"job_id"`
	Success         bool      `json:"success"`
	StatusCode      int       `json:"status_code"`
	DurationMs      int       `json:"duration_ms"`
	Url             string    `json:"url"`
	ResponseBody    string    `json:"response_body"`
	Error           string    `json:"error"`
	ExecutedAt      time.Time `json:"executed_at"`
	BatchIdentifier string    `json:"batch_identifier"`
}
