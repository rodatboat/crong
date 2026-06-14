package models

import "time"

type Folder struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
