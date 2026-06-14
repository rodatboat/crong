package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/handlers"
)

func FolderRoutes(app fiber.Router) {
	folders := app.Group("/folders")

	folders.Post("/", handlers.CreateFolder)
	folders.Get("/", handlers.ReadFolders)
	folders.Put("/:id", handlers.UpdateFolder)
	folders.Delete("/:id", handlers.DeleteFolder)
}
