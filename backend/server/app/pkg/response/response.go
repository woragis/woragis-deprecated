package response

import "github.com/gofiber/fiber/v2"

// Success sends a JSON success response with the provided payload.
func Success(c *fiber.Ctx, status int, payload any) error {
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"data":    payload,
	})
}

// Error sends a JSON error response exposing only the numeric code.
func Error(c *fiber.Ctx, status int, code int, details any) error {
	errorBody := fiber.Map{
		"code": code,
	}
	if details != nil {
		errorBody["details"] = details
	}

	return c.Status(status).JSON(fiber.Map{
		"success": false,
		"error":   errorBody,
	})
}
