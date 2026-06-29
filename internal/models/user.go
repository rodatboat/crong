package models

import (
	"time"
)

type User struct {
	ID uint `json:"id"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`

	AuthToken string `json:"auth_token,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,max=255"`
	Password string `json:"password" validate:"required,max=255"`
}

type UserRegister struct {
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
	Email     string `json:"email" validate:"required,max=255"`
	Password  string `json:"password" validate:"required,max=255"`
}

type UserUpdate struct {
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"max=100"`
}
