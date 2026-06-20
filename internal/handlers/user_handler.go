package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/response"
	"github.com/rodatboat/crong/internal/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) LoginUser(c fiber.Ctx) error {
	newJob := new(models.Job)
	if err := c.Bind().Body(newJob); err != nil {
		return err
	}

	// TODO: Call repository to create job in database

	return response.Success(c, newJob)
}

func (h *UserHandler) RegisterUser(c fiber.Ctx) error {
	newUser := new(models.User)
	if err := c.Bind().Body(newUser); err != nil {
		return err
	}

	// TODO: Call repository to create user in database

	return response.Success(c, newUser)
}
