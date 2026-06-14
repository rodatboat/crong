package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/response"
)

func LoginUser(c fiber.Ctx) error {
	newJob := new(models.Job)
	if err := c.Bind().Body(newJob); err != nil {
		return err
	}

	// TODO: Call repository to create job in database

	return response.Success(c, newJob)
}

func RegisterUser(c fiber.Ctx) error {
	newUser := new(models.User)
	if err := c.Bind().Body(newUser); err != nil {
		return err
	}

	// TODO: Call repository to create user in database

	return response.Success(c, newUser)
}
