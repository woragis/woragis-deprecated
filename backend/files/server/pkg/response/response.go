package response

import (
	"github.com/gofiber/fiber/v2"
)

// Success sends a successful JSON response
func Success(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

// Error sends an error JSON response
func Error(c *fiber.Ctx, statusCode int, code int, data interface{}) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success": false,
		"code":    code,
		"data":    data,
	})
}

