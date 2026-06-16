package entities

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Title    string    `gorm:"column:title;not null"`
	Url      string    `gorm:"column:url;not null"`
	FolderID uint      `gorm:"column:folder_id"`
	UserID   uint      `gorm:"column:user_id;not null"`
	Method   ReqMethod `gorm:"column:method;type:smallint;not null"`

	Headers datatypes.JSON `gorm:"column:headers"`
	Auth    datatypes.JSON `gorm:"column:auth"`
	Body    string         `gorm:"column:body"`

	Cron     string `gorm:"column:cron;not null"`
	Timezone string `gorm:"column:timezone"`

	Timeout int  `gorm:"column:timeout"`
	Enabled bool `gorm:"column:enabled;not null"`

	LastExecution time.Time `gorm:"column:last_execution"`
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
