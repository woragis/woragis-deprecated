package projects

import (
	"fmt"

	"github.com/woragis/backend/server/app/pkg/validation"
)

// ValidateCreateProjectPayload validates create project payload
func ValidateCreateProjectPayload(payload *createProjectPayload) error {
	// Validate name (required, 1-200 chars)
	if err := validation.ValidateString(payload.Name, 1, 200, "name"); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(payload.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := validation.ValidateNoXSS(payload.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}

	// Validate description (optional, but if provided, validate)
	if payload.Description != "" {
		if err := validation.ValidateString(payload.Description, 1, 5000, "description"); err != nil {
			return fmt.Errorf("description: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(payload.Description); err != nil {
			return fmt.Errorf("description: %w", err)
		}
		if err := validation.ValidateNoXSS(payload.Description); err != nil {
			return fmt.Errorf("description: %w", err)
		}
	}

	// Validate status (optional, but if provided, validate)
	if payload.Status != "" {
		validStatuses := []string{"planning", "in_progress", "on_hold", "completed", "cancelled"}
		isValid := false
		for _, validStatus := range validStatuses {
			if payload.Status == validStatus {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("status: must be one of: planning, in_progress, on_hold, completed, cancelled")
		}
	}

	// Validate health score (optional, but if provided, validate range)
	if payload.HealthScore != 0 {
		if payload.HealthScore < 0 || payload.HealthScore > 100 {
			return fmt.Errorf("health_score: must be between 0 and 100")
		}
	}

	// Validate financial metrics (optional, but if provided, validate range)
	if payload.MRR < 0 {
		return fmt.Errorf("mrr: must be non-negative")
	}
	if payload.CAC < 0 {
		return fmt.Errorf("cac: must be non-negative")
	}
	if payload.LTV < 0 {
		return fmt.Errorf("ltv: must be non-negative")
	}
	if payload.ChurnRate < 0 || payload.ChurnRate > 100 {
		return fmt.Errorf("churn_rate: must be between 0 and 100")
	}

	return nil
}

// ValidateUpdateStatusPayload validates update status payload
func ValidateUpdateStatusPayload(payload *updateStatusPayload) error {
	if payload.Status == "" {
		return fmt.Errorf("status is required")
	}
	validStatuses := []string{"planning", "in_progress", "on_hold", "completed", "cancelled"}
	isValid := false
	for _, validStatus := range validStatuses {
		if payload.Status == validStatus {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("status: must be one of: planning, in_progress, on_hold, completed, cancelled")
	}
	return nil
}

// ValidateUpdateMetricsPayload validates update metrics payload
func ValidateUpdateMetricsPayload(payload *updateMetricsPayload) error {
	// Validate health score (optional, but if provided, validate range)
	if payload.HealthScore != 0 {
		if payload.HealthScore < 0 || payload.HealthScore > 100 {
			return fmt.Errorf("health_score: must be between 0 and 100")
		}
	}

	// Validate financial metrics (optional, but if provided, validate range)
	if payload.MRR < 0 {
		return fmt.Errorf("mrr: must be non-negative")
	}
	if payload.CAC < 0 {
		return fmt.Errorf("cac: must be non-negative")
	}
	if payload.LTV < 0 {
		return fmt.Errorf("ltv: must be non-negative")
	}
	if payload.ChurnRate < 0 || payload.ChurnRate > 100 {
		return fmt.Errorf("churn_rate: must be between 0 and 100")
	}

	return nil
}

// ValidateAddMilestonePayload validates add milestone payload
func ValidateAddMilestonePayload(payload *addMilestonePayload) error {
	// Validate title (required, 1-200 chars)
	if err := validation.ValidateString(payload.Title, 1, 200, "title"); err != nil {
		return fmt.Errorf("title: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(payload.Title); err != nil {
		return fmt.Errorf("title: %w", err)
	}
	if err := validation.ValidateNoXSS(payload.Title); err != nil {
		return fmt.Errorf("title: %w", err)
	}

	// Validate description (optional, but if provided, validate)
	if payload.Description != "" {
		if err := validation.ValidateString(payload.Description, 1, 2000, "description"); err != nil {
			return fmt.Errorf("description: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(payload.Description); err != nil {
			return fmt.Errorf("description: %w", err)
		}
		if err := validation.ValidateNoXSS(payload.Description); err != nil {
			return fmt.Errorf("description: %w", err)
		}
	}

	return nil
}

// ValidateCreateColumnPayload validates create column payload
func ValidateCreateColumnPayload(payload *createColumnPayload) error {
	// Validate name (required, 1-100 chars)
	if err := validation.ValidateString(payload.Name, 1, 100, "name"); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(payload.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := validation.ValidateNoXSS(payload.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}

	// Validate WIP limit (optional, but if provided, validate range)
	if payload.WIPLimit != nil {
		if *payload.WIPLimit < 0 {
			return fmt.Errorf("wip_limit: must be non-negative")
		}
		if *payload.WIPLimit > 1000 {
			return fmt.Errorf("wip_limit: must be at most 1000")
		}
	}

	return nil
}

// ValidateCreateCardPayload validates create card payload
func ValidateCreateCardPayload(payload *createCardPayload) error {
	// Validate column ID (required, UUID)
	if payload.ColumnID == "" {
		return fmt.Errorf("column_id is required")
	}
	if err := validation.ValidateUUID(payload.ColumnID); err != nil {
		return fmt.Errorf("column_id: %w", err)
	}

	// Validate title (required, 1-200 chars)
	if err := validation.ValidateString(payload.Title, 1, 200, "title"); err != nil {
		return fmt.Errorf("title: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(payload.Title); err != nil {
		return fmt.Errorf("title: %w", err)
	}
	if err := validation.ValidateNoXSS(payload.Title); err != nil {
		return fmt.Errorf("title: %w", err)
	}

	// Validate description (optional, but if provided, validate)
	if payload.Description != "" {
		if err := validation.ValidateString(payload.Description, 1, 2000, "description"); err != nil {
			return fmt.Errorf("description: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(payload.Description); err != nil {
			return fmt.Errorf("description: %w", err)
		}
		if err := validation.ValidateNoXSS(payload.Description); err != nil {
			return fmt.Errorf("description: %w", err)
		}
	}

	// Validate milestone ID (optional, but if provided, validate UUID)
	if payload.MilestoneID != "" {
		if err := validation.ValidateUUID(payload.MilestoneID); err != nil {
			return fmt.Errorf("milestone_id: %w", err)
		}
	}

	return nil
}

// ValidateMoveCardPayload validates move card payload
func ValidateMoveCardPayload(payload *moveCardPayload) error {
	// Validate target column ID (required, UUID)
	if payload.TargetColumnID == "" {
		return fmt.Errorf("target_column_id is required")
	}
	if err := validation.ValidateUUID(payload.TargetColumnID); err != nil {
		return fmt.Errorf("target_column_id: %w", err)
	}

	// Validate target position (must be non-negative)
	if payload.TargetPosition < 0 {
		return fmt.Errorf("target_position: must be non-negative")
	}

	return nil
}

// ValidateDependencyPayload validates dependency payload
func ValidateDependencyPayload(payload *dependencyPayload) error {
	// Validate depends on project ID (required, UUID)
	if payload.DependsOnProjectID == "" {
		return fmt.Errorf("depends_on_project_id is required")
	}
	if err := validation.ValidateUUID(payload.DependsOnProjectID); err != nil {
		return fmt.Errorf("depends_on_project_id: %w", err)
	}

	// Validate type (required)
	if payload.Type == "" {
		return fmt.Errorf("type is required")
	}
	validTypes := []string{"blocks", "related", "duplicate"}
	isValid := false
	for _, validType := range validTypes {
		if payload.Type == validType {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("type: must be one of: blocks, related, duplicate")
	}

	return nil
}

