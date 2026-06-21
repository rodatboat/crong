package entities

type User struct {
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Email     string `gorm:"not null;unique"`
	Password  string `gorm:"not null"`
	Verified  bool   `gorm:"not null;default:false"`
}
