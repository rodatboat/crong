package entities

import (
	"gorm.io/gorm"
)

type Folder struct {
	gorm.Model
	Name   string `gorm:"column:name;not null"`
	UserID uint   `gorm:"column:user_id;not null"`
}
