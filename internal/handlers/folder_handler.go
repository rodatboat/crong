package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/response"
)

func CreateFolder(c fiber.Ctx) error {
	newFolder := new(models.Folder)
	if err := c.Bind().Body(newFolder); err != nil {
		return err
	}

	// TODO: Call repository to create folder in database

	return response.Success(c, newFolder)
}

func ReadFolders(c fiber.Ctx) error {

	// TODO: Call repository to read folders from database

	folders := []models.Folder{}

	return response.Success(c, folders)
}

func UpdateFolder(c fiber.Ctx) error {

	folderId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid folder ID")
	}

	// TODO: Call repository to update folder in database

	folder := models.Folder{
		ID: uint(folderId),
	}

	return response.Success(c, &folder)
}

func DeleteFolder(c fiber.Ctx) error {
	_, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid folder ID")
	}

	// TODO: Call repository to delete folder from database

	return response.Success(c, nil)
}
