package models

import (
	"time"
)

type Folder struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`

	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FolderCreate struct {
	Name string `json:"name" validate:"required"`
}

type FolderUpdate struct {
	Name string `json:"name" validate:"required"`
}
