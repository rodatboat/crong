package repositories

import (
	"github.com/rodatboat/crong/internal/entities"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email string, firstName string, lastName string, passwordHash string) (*entities.User, error) {
	user := &entities.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Password:  passwordHash,
	}

	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(id uint, firstName string, lastName string) (*entities.User, error) {
	var user entities.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName

	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
