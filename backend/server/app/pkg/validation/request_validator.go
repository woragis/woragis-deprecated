package validation

import (
	"github.com/gofiber/fiber/v2"
)

// ValidateRequest validates common request fields
type ValidateRequest struct {
	Email    *string `json:"email,omitempty"`
	UUID     *string `json:"uuid,omitempty"`
	URL      *string `json:"url,omitempty"`
	String   *string `json:"string,omitempty"`
	StringMin int    `json:"-"`
	StringMax int    `json:"-"`
	FieldName string `json:"-"`
}

// Validate performs validation on request fields
func (v *ValidateRequest) Validate() error {
	if v.Email != nil {
		if err := ValidateEmail(*v.Email); err != nil {
			return err
		}
	}

	if v.UUID != nil {
		if err := ValidateUUID(*v.UUID); err != nil {
			return err
		}
	}

	if v.URL != nil {
		if err := ValidateURL(*v.URL); err != nil {
			return err
		}
	}

	if v.String != nil {
		min := v.StringMin
		max := v.StringMax
		if max == 0 {
			max = 1000 // Default max length
		}
		fieldName := v.FieldName
		if fieldName == "" {
			fieldName = "string"
		}
		if err := ValidateString(*v.String, min, max, fieldName); err != nil {
			return err
		}
		// Check for SQL injection and XSS
		if err := ValidateNoSQLInjection(*v.String); err != nil {
			return err
		}
		if err := ValidateNoXSS(*v.String); err != nil {
			return err
		}
	}

	return nil
}

// ValidateRequestMiddleware creates a middleware that validates request body
func ValidateRequestMiddleware(validator func(*fiber.Ctx) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := validator(c); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Next()
	}
}

// ValidateQueryParams validates query parameters
func ValidateQueryParams(c *fiber.Ctx, validations map[string]func(string) error) error {
	for key, validate := range validations {
		value := c.Query(key)
		if value != "" {
			if err := validate(value); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Invalid query parameter",
					"field":   key,
					"message": err.Error(),
				})
			}
		}
	}
	return nil
}

// ValidatePathParams validates path parameters
func ValidatePathParams(c *fiber.Ctx, validations map[string]func(string) error) error {
	for key, validate := range validations {
		value := c.Params(key)
		if value != "" {
			if err := validate(value); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Invalid path parameter",
					"field":   key,
					"message": err.Error(),
				})
			}
		}
	}
	return nil
}
