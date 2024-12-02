// internal/middleware/jwt_middleware.go
package middleware

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	_ "github.com/golang-jwt/jwt/v4"
)

func JWTAuth() fiber.Handler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(secret),
		ErrorHandler: jwtError,
		TokenLookup:  "header:Authorization",
		AuthScheme:   "Bearer",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error":   "Unauthorized",
		"message": err.Error(),
	})
}
