package utils

import (
	"github.com/rodatboat/crong/internal/entities"
	"github.com/rodatboat/crong/internal/models"
)

func MapUserEntityToUserModel(userEntity *entities.User) *models.User {
	return &models.User{
		ID:        userEntity.ID,
		FirstName: userEntity.FirstName,
		LastName:  userEntity.LastName,
		Email:     userEntity.Email,
		CreatedAt: userEntity.CreatedAt,
		UpdatedAt: userEntity.UpdatedAt,
	}
}
