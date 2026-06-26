package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/middleware"
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

	auth := c.Locals(middleware.AuthContextKey).(*middleware.AuthContext)

	// Validate request
	if validationErrors, err := utils.ValidateStruct(&req); err != nil {
		return resp.HandleValidationError(c, err, validationErrors)
	}

	folder, err := h.folderService.CreateFolder(auth.UserID, &req)
	if err != nil {
		return resp.HandleError(c, err)
	}

	return resp.Send(c, resp.Success(folder))
}

func (h *FolderHandler) ReadFolders(c fiber.Ctx) error {
	auth := c.Locals(middleware.AuthContextKey).(*middleware.AuthContext)

	folders, err := h.folderService.GetFoldersByUser(auth.UserID)
	if err != nil {
		return resp.HandleError(c, err)
	}

	return resp.Send(c, resp.Success(folders))
}

func (h *FolderHandler) GetFoldersDetailsByID(c fiber.Ctx) error {
	folderIDStr := c.Params("id")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil || folderID == 0 {
		return resp.Send(c, resp.BadRequest())
	}

	auth := c.Locals(middleware.AuthContextKey).(*middleware.AuthContext)

	folders, err := h.folderService.GetFolderDetailsByID(uint(folderID), auth.UserID)
	if err != nil {
		return resp.HandleError(c, err)
	}

	return resp.Send(c, resp.Success(folders))
}

func (h *FolderHandler) UpdateFolder(c fiber.Ctx) error {
	folderIDStr := c.Params("id")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil || folderID == 0 {
		return resp.Send(c, resp.BadRequest())
	}

	auth := c.Locals(middleware.AuthContextKey).(*middleware.AuthContext)

	var req models.FolderUpdate
	if err := c.Bind().Body(&req); err != nil {
		return resp.Send(c, resp.BadRequest())
	}

	// Validate request
	if validationErrors, err := utils.ValidateStruct(&req); err != nil {
		return resp.HandleValidationError(c, err, validationErrors)
	}

	folder, err := h.folderService.UpdateFolder(uint(folderID), auth.UserID, &req)
	if err != nil {
		return resp.HandleError(c, err)
	}

	return resp.Send(c, resp.Success(folder))
}

func (h *FolderHandler) DeleteFolder(c fiber.Ctx) error {
	folderIDStr := c.Params("id")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil || folderID == 0 {
		return resp.Send(c, resp.BadRequest())
	}

	auth := c.Locals(middleware.AuthContextKey).(*middleware.AuthContext)

	err = h.folderService.DeleteFolder(uint(folderID), auth.UserID)
	if err != nil {
		return resp.HandleError(c, err)
	}

	return resp.Send(c, resp.Success(nil))
}
