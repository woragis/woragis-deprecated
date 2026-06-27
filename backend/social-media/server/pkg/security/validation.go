package security

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

// RequestSizeLimitMiddleware limits request body size
func RequestSizeLimitMiddleware(maxSize int64) fiber.Handler {
	return func(c *fiber.Ctx) error {
		contentLength := int64(c.Request().Header.ContentLength())
		if contentLength > maxSize {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error": "Request body too large",
				"message": fmt.Sprintf("Maximum request size is %d bytes", maxSize),
			})
		}
		return c.Next()
	}
}

// InputSanitizationMiddleware sanitizes input to prevent injection attacks
func InputSanitizationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Sanitize query parameters
		queries := c.Queries()
		for key, value := range queries {
			sanitized := SanitizeString(value)
			if sanitized != value {
				// Update query parameter if sanitization changed it
				c.Request().URI().QueryArgs().Set(key, sanitized)
			}
		}

		return c.Next()
	}
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(s string) string {
	// Remove null bytes
	s = strings.ReplaceAll(s, "\x00", "")

	// Trim whitespace
	s = strings.TrimSpace(s)

	// Remove control characters (except newline, tab, carriage return)
	var result strings.Builder
	for _, r := range s {
		if unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r' {
			continue
		}
		result.WriteRune(r)
	}

	return result.String()
}

