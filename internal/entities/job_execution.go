package entities

import "gorm.io/gorm"

type JobExecution struct {
	gorm.Model

	JobID uint `gorm:"index;not null"`

	Success bool `gorm:"column:success"`

	StatusCode int `gorm:"column:status_code"`
	DurationMs int `gorm:"column:duration_ms"`

	ResponseBody string `gorm:"column:response_body;type:text"`
	Error        string `gorm:"column:error;type:text"`
}
