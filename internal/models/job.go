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
	Headers map[string]string  `json:"headers"`
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
	Url      string `json:"url" validate:"required"`
	FolderID uint   `json:"folder_id"`

	Method  entities.ReqMethod `json:"method" validate:"required"`
	Headers map[string]string  `json:"headers"`
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
	Username string `json:"username"`
	Password string `json:"password"`
}
