package entities

import (
	"gorm.io/gorm"
)

type Folder struct {
	gorm.Model
	Name   string `gorm:"not null" json:"name"`
	UserID uint   `json:"user_id"`
}
