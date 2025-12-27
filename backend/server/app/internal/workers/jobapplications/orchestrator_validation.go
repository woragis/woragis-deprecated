package jobapplications

import (
	"fmt"

	"github.com/woragis/backend/server/app/pkg/validation"
)

// ValidateWebsiteName validates a website name
func ValidateWebsiteName(websiteName string) error {
	if websiteName == "" {
		return fmt.Errorf("website_name is required")
	}

	// Validate length
	if err := validation.ValidateString(websiteName, 1, 100, "website_name"); err != nil {
		return fmt.Errorf("website_name: %w", err)
	}

	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(websiteName); err != nil {
		return fmt.Errorf("website_name: %w", err)
	}
	if err := validation.ValidateNoXSS(websiteName); err != nil {
		return fmt.Errorf("website_name: %w", err)
	}

	return nil
}

