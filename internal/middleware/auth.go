package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/rodatboat/crong/internal/resp"
)

func Protected() fiber.Handler {
	log.Info("Initializing authentication middleware")
	secret := os.Getenv("AUTH_SECRET")

	if secret == "" {
		log.Error("AUTH_SECRET environment variable is not set, failed to initialize authentication middleware")
		panic("AUTH_SECRET environment variable is not set, failed to initialize authentication middleware")
	}

	return func(c fiber.Ctx) error {
		auth, err := authenticate(c, secret)
		if err != nil {
			log.Warn(err.Error())
			return resp.Send(c, resp.Unauthorized())
		}

		c.Locals(AuthContextKey, auth)
		log.Info("Route authentication successful")
		return c.Next()
	}
}

func authenticate(c fiber.Ctx, secret string) (*AuthContext, error) {
	auth := c.Get("Authorization")
	if auth == "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	if token != secret {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid bearer token")
	}

	// TODO: Remove TEMP & replace with JWT subject/claims.
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

	return &AuthContext{
		UserID: 1,
	}, nil
}

func AuthDetails(c fiber.Ctx) *AuthContext {
	auth, ok := c.Locals(AuthContextKey).(*AuthContext)
	if !ok {
		return nil
	}
	return auth
}

type AuthContext struct {
	UserID uint `json:"user_id"`
}

const AuthContextKey = "auth"
