package entities

import (
	"time"

	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Title    string `json:"title"`
	Url      string `json:"url"`
	FolderID uint   `json:"folder_id"`
	UserID   uint   `json:"user_id"`

	Method uint `json:"method"`
	// Headers
	// Auth
	Body string `json:"body"`

	Timezone string `json:"timezone"`
	Timeout  int    `json:"timeout"`
	Enabled  bool   `json:"enabled"`

	LastExecution time.Time `json:"last_execution"`
}
