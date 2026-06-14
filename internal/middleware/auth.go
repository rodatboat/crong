package middleware

import (
	"os"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/rodatboat/crong/internal/response"
)

func Protected() fiber.Handler {
	secret := os.Getenv("AUTH_SECRET")

	if secret == "" {
		panic("AUTH_SECRET environment variable is not set, failed to initialize authentication middleware")
	}

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(secret),
		},
		SuccessHandler: func(c fiber.Ctx) error {
			return c.Next()
		},
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
		},
	})
}
