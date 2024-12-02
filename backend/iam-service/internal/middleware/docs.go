// backend/iam-service/internal/middleware/docs.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func ServeOpenAPISpec(c *fiber.Ctx) error {
	return c.SendFile("./static/docs/openapi.yaml")
}

func ServeRedoc(c *fiber.Ctx) error {
	return c.SendFile("./static/redoc.html")
}
