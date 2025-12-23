package validation

import (
	"github.com/gofiber/fiber/v2"
)

// ValidateUUIDParam validates a UUID path parameter
func ValidateUUIDParam(c *fiber.Ctx, paramName string) error {
	uuid := c.Params(paramName)
	if uuid == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Missing required parameter",
			"field":   paramName,
			"message": "UUID parameter is required",
		})
	}
	
	if err := ValidateUUID(uuid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid parameter format",
			"field":   paramName,
			"message": err.Error(),
		})
	}
	
	return nil
}

// ValidateUUIDParamMiddleware creates middleware to validate UUID path parameters
func ValidateUUIDParamMiddleware(paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := ValidateUUIDParam(c, paramName); err != nil {
			return err
		}
		return c.Next()
	}
}

// ValidateEmailQuery validates an email query parameter
func ValidateEmailQuery(c *fiber.Ctx, paramName string, required bool) error {
	email := c.Query(paramName)
	if email == "" {
		if required {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Missing required query parameter",
				"field":   paramName,
				"message": "Email parameter is required",
			})
		}
		return nil // Optional parameter
	}
	
	if err := ValidateEmail(email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid query parameter format",
			"field":   paramName,
			"message": err.Error(),
		})
	}
	
	return nil
}

// ValidateStringQuery validates a string query parameter with length constraints
func ValidateStringQuery(c *fiber.Ctx, paramName string, min, max int, required bool) error {
	value := c.Query(paramName)
	if value == "" {
		if required {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Missing required query parameter",
				"field":   paramName,
				"message": "Parameter is required",
			})
		}
		return nil // Optional parameter
	}
	
	if err := ValidateString(value, min, max, paramName); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid query parameter",
			"field":   paramName,
			"message": err.Error(),
		})
	}
	
	// Check for security issues
	if err := ValidateNoSQLInjection(value); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid query parameter",
			"field":   paramName,
			"message": err.Error(),
		})
	}
	
	return nil
}
