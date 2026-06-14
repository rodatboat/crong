package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/handlers"
	"github.com/rodatboat/crong/internal/middleware"
)

func FolderRoutes(app fiber.Router) {
	folders := app.Group("/folders")

	folders.Post("/", middleware.Protected(), handlers.CreateFolder)
	folders.Get("/", middleware.Protected(), handlers.ReadFolders)
	folders.Put("/:id", middleware.Protected(), handlers.UpdateFolder)
	folders.Delete("/:id", middleware.Protected(), handlers.DeleteFolder)
}
