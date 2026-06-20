package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/response"
	"github.com/rodatboat/crong/internal/services"
	"github.com/rodatboat/crong/internal/utils"
)

type FolderHandler struct {
	folderService *services.FolderService
}

func NewFolderHandler(folderService *services.FolderService) *FolderHandler {
	return &FolderHandler{
		folderService: folderService,
	}
}

func (h *FolderHandler) CreateFolder(c fiber.Ctx) error {
	var req models.FolderCreate
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	// TODO: Call repository to create folder in database

	return response.Success(c, &req)
}

func (h *FolderHandler) ReadFolders(c fiber.Ctx) error {

	// TODO: Call repository to read folders from database

	folders := []models.Folder{}

	return response.Success(c, folders)
}

func (h *FolderHandler) UpdateFolder(c fiber.Ctx) error {
	folderId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid folder ID")
	}

	var req models.FolderUpdate
	if err := c.Bind().Body(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, err.Error())
	}

	// TODO: Call repository to update folder in database

	folder := models.Folder{
		ID: uint(folderId),
	}

	return response.Success(c, &folder)
}

func (h *FolderHandler) DeleteFolder(c fiber.Ctx) error {
	_, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid folder ID")
	}

	// TODO: Call repository to delete folder from database

	return response.Success(c, nil)
}
