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

	users.Post("/login", handler.LoginUser)
	users.Post("/register", handler.RegisterUser)
	users.Get("/me", middleware.Protected(serviceContainer.UserRepository), handler.GetUserDetails)
	users.Put("/me", middleware.Protected(serviceContainer.UserRepository), handler.UpdateUser)
}
