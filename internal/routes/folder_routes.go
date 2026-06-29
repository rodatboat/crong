package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/container"
	"github.com/rodatboat/crong/internal/handlers"
	"github.com/rodatboat/crong/internal/middleware"
)

func FolderRoutes(app fiber.Router, serviceContainer *container.Container) {
	folders := app.Group("/folders")
	handler := handlers.NewFolderHandler(serviceContainer.FolderService)

	folders.Post("/", middleware.Protected(serviceContainer.UserRepository), handler.CreateFolder)
	folders.Get("/", middleware.Protected(serviceContainer.UserRepository), handler.ReadFolders)
	folders.Get("/:id", middleware.Protected(serviceContainer.UserRepository), handler.GetFoldersDetailsByID)
	folders.Put("/:id", middleware.Protected(serviceContainer.UserRepository), handler.UpdateFolder)
	folders.Delete("/:id", middleware.Protected(serviceContainer.UserRepository), handler.DeleteFolder)
}
