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

	folders.Post("/", middleware.Protected(), handler.CreateFolder)
	folders.Get("/", middleware.Protected(), handler.ReadFolders)
	folders.Put("/:id", middleware.Protected(), handler.UpdateFolder)
	folders.Delete("/:id", middleware.Protected(), handler.DeleteFolder)
}
