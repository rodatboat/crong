package entities

type Folder struct {
	Name   string `gorm:"column:name;not null"`
	UserID uint   `gorm:"column:user_id;not null"`
}
