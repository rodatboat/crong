package models

import (
	"time"

	"github.com/rodatboat/crong/internal/entities"
)

type Job struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Url      string `json:"url"`
	FolderID *uint  `json:"folder_id"`

	Method  entities.ReqMethod `json:"method"`
	Headers []JobHeaders       `json:"headers"`
	Auth    JobAuth            `json:"auth"`
	Body    string             `json:"body"`
	Cron    string             `json:"cron"`

	Timezone string `json:"timezone"`
	Timeout  int    `json:"timeout"`
	Enabled  bool   `json:"enabled"`

	LastExecution *time.Time `json:"last_execution"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type JobCreateRequest struct {
	Title    string `json:"title" validate:"required,min=1,max=255"`
	Url      string `json:"url" validate:"required,url,min=1,max=2048"`
	FolderID *uint  `json:"folder_id"`

	Method  entities.ReqMethod `json:"method" validate:"validmethod"`
	Headers []JobHeaders       `json:"headers"`
	Auth    JobAuth            `json:"auth"`
	Body    string             `json:"body" validate:"max=10000"`
	Cron    string             `json:"cron" validate:"required,min=1,validcron,max=100"`

	Timezone string `json:"timezone" validate:"validtimezone,max=50" default:"America/Chicago"`
	Timeout  int    `json:"timeout" validate:"required,min=0,max=30"`
	Enabled  bool   `json:"enabled"`
}

type JobUpdateRequest struct {
	JobCreateRequest
}

type JobAuth struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"username" validate:"required_if=Enabled true,max=255"`
	Password string `json:"password" validate:"required_if=Enabled true,max=255"`
}

type JobHeaders struct {
	Key   string `json:"key" validate:"required,min=1,max=255"`
	Value string `json:"value" validate:"required,min=1,max=1024"`
}

type JobExecution struct {
	ID              uint   `json:"id"`
	JobID           uint   `json:"job_id"`
	BatchIdentifier string `json:"batch_identifier"`

	ExecutionStatus entities.ExecutionStatus `json:"execution_status"`
	StatusCode      int                      `json:"status_code"`
	StatusText      string                   `json:"status_text"`
	DurationMs      int                      `json:"duration_ms"`
	Url             string                   `json:"url"`

	ResponseHeaders string `json:"response_headers"`
	ResponseBody    string `json:"response_body"`
	Error           string `json:"error"`

	ExecutedAt *time.Time `json:"executed_at"`
	PlannedFor *time.Time `json:"planned_for"`
	CreatedAt  time.Time  `json:"created_at"`
}

type JobExecutionStats struct {
	JobExecutionID    uint `json:"job_execution_id"`
	DNSLookupMs       int  `json:"dns_lookup_ms"`
	TCPConnectMs      int  `json:"tcp_connect_ms"`
	TLSHandshakeMs    int  `json:"tls_handshake_ms"`
	TimeToFirstByteMs int  `json:"time_to_first_byte_ms"`
	RequestWriteMs    int  `json:"request_write_ms"`
}
