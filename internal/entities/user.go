package entities

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Firstname string `gorm:"not null" json:"first_name"`
	Lastname  string `gorm:"not null" json:"last_name"`
	Email     string `gorm:"not null;unique" json:"email"`
	Password  string `gorm:"not null" json:"-"`
}
