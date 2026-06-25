package entities

type Folder struct {
	ID     uint   `gorm:"column:id;primaryKey"`
	Name   string `gorm:"column:name;not null"`
	UserID uint   `gorm:"column:user_id;not null"`
}
