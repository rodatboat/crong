package models

import (
	"time"
)

type Job struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	FolderID uint   `json:"folder_id"`

	Method  ReqMethod         `json:"method"`
	Headers map[string]string `json:"headers"`
	Auth    *JobAuth          `json:"auth"`
	Body    string            `json:"body"`

	Timezone string `json:"timezone"`
	Timeout  int    `json:"timeout"`
	Enabled  bool   `json:"enabled"`

	LastExecution time.Time `json:"last_execution"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type JobUpdateRequest struct {
	Title    string `json:"title" validate:"required"`
	Url      string `json:"url" validate:"required"`
	FolderID uint   `json:"folder_id"`

	Method  ReqMethod         `json:"method"`
	Headers map[string]string `json:"headers"`
	Auth    *JobAuth          `json:"auth"`
	Body    string            `json:"body"`

	Timezone string `json:"timezone"`
	Timeout  int    `json:"timeout" validate:"required,max=30"`
	Enabled  bool   `json:"enabled"`
}

type JobAuth struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"username"`
	Password string `json:"password"`
}

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
