package middleware

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/rodatboat/crong/internal/response"
)

func Protected() fiber.Handler {
	log.Info("Initializing authentication middleware")
	secret := os.Getenv("AUTH_SECRET")

	if secret == "" {
		log.Error("AUTH_SECRET environment variable is not set, failed to initialize authentication middleware")
		panic("AUTH_SECRET environment variable is not set, failed to initialize authentication middleware")
	}

	// TEMP START
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Warn("Missing Authorization header")
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
		}
		if authHeader != "Bearer "+secret {
			log.Warn("Invalid Authorization header")
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		log.Info("Route authentication successful")
		return c.Next()
	}
	// TEMP END

	// TODO: Remove TEMP once ready to accept JWT
	// return jwtware.New(jwtware.Config{
	// 	SigningKey: jwtware.SigningKey{
	// 		Key: []byte(secret),
	// 	},
	// 	SuccessHandler: func(c fiber.Ctx) error {
	// 		log.Info("Route authentication successful")
	// 		return c.Next()
	// 	},
	// 	ErrorHandler: func(c fiber.Ctx, err error) error {
	// 		log.Warn("Route authentication failed")
	// 		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
	// 	},
	// })
}
