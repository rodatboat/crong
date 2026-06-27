package middleware

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rodatboat/crong/internal/resp"
)

const AuthContextKey = "auth"

type AuthContext struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
}

// JWTClaims represents the JWT claims structure
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

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
			return resp.HandleError(c, err)
		}

		c.Locals(AuthContextKey, auth)
		log.Info("Route authentication successful")
		return c.Next()
	}
}

func authenticate(c fiber.Ctx, secret string) (*AuthContext, error) {
	auth := c.Get("Authorization")
	if auth == "" {
		return nil, resp.ErrUnauthorized
	}

	tokenString := strings.TrimPrefix(auth, "Bearer ")
	if tokenString == auth {
		return nil, resp.ErrUnauthorized
	}

	// Parse and validate JWT token
	claims, err := ParseJWT(tokenString, secret)
	if err != nil {
		return nil, resp.ErrUnauthorized
	}

	return &AuthContext{
		UserID: claims.UserID,
		Email:  claims.Email,
	}, nil
}

func AuthDetails(c fiber.Ctx) *AuthContext {
	auth, ok := c.Locals(AuthContextKey).(*AuthContext)
	if !ok {
		return nil
	}
	return auth
}

// GenerateJWT creates a new JWT token for the user
func GenerateJWT(userID uint, email string) (string, error) {
	secret := os.Getenv("AUTH_SECRET")
	if secret == "" {
		return "", errors.New("AUTH_SECRET environment variable is not set")
	}

	// Set token expiration (e.g., 24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT validates and parses a JWT token
func ParseJWT(tokenString string, secret string) (*JWTClaims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and validate claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
