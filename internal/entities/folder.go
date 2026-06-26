package entities

import "time"

type Folder struct {
	ID     uint   `gorm:"column:id;primaryKey"`
	Name   string `gorm:"column:name;not null"`
	UserID uint   `gorm:"column:user_id;not null"`

	CreatedAt time.Time `gorm:"column:created_at;default:now()"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:now()"`
}
