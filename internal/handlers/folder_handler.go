package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/models"
	"github.com/rodatboat/crong/internal/resp"
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
		return resp.Send(c, resp.BadRequest())
	}

	// Validate request
	if validationErrors, err := utils.ValidateStruct(&req); err != nil {
		if validationErrors == nil {
			return resp.Send(c, resp.InternalServerError())
		}
		return resp.Send(c, resp.ValidationError(validationErrors))
	}

	// TODO: Call repository to create folder in database

	return resp.Send(c, resp.Success(req))
}

func (h *FolderHandler) ReadFolders(c fiber.Ctx) error {

	// TODO: Call repository to read folders from database

	folders := []models.Folder{}

	return resp.Send(c, resp.Success(folders))
}

func (h *FolderHandler) UpdateFolder(c fiber.Ctx) error {
	folderId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return resp.Send(c, resp.BadRequest())
	}

	var req models.FolderUpdate
	if err := c.Bind().Body(&req); err != nil {
		return resp.Send(c, resp.BadRequest())
	}

	// Validate request
	if validationErrors, err := utils.ValidateStruct(&req); err != nil {
		if validationErrors == nil {
			return resp.Send(c, resp.InternalServerError())
		}
		return resp.Send(c, resp.ValidationError(validationErrors))
	}

	// TODO: Call repository to update folder in database

	folder := models.Folder{
		ID: uint(folderId),
	}

	return resp.Send(c, resp.Success(folder))
}

func (h *FolderHandler) DeleteFolder(c fiber.Ctx) error {
	_, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return resp.Send(c, resp.BadRequest())
	}

	// TODO: Call repository to delete folder from database

	return resp.Send(c, resp.Success(nil))
}
