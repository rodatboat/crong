package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/container"
	"github.com/rodatboat/crong/internal/handlers"
	"github.com/rodatboat/crong/internal/middleware"
)

func UserRoutes(app fiber.Router, serviceContainer *container.Container) {
	users := app.Group("/users")
	handler := handlers.NewUserHandler(serviceContainer.UserService)

	users.Post("/login", middleware.Protected(), handler.LoginUser)
	users.Get("/register", middleware.Protected(), handler.RegisterUser)
}
