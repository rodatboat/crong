package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/container"
	"github.com/rodatboat/crong/internal/handlers"
	"github.com/rodatboat/crong/internal/middleware"
)

func UserRoutes(app fiber.Router, svc *container.Container) {
	users := app.Group("/users")

	users.Post("/login", middleware.Protected(), handlers.LoginUser)
	users.Get("/register", middleware.Protected(), handlers.RegisterUser)
}
