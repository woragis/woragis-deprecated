package ideas

import (
	"fmt"

	"github.com/woragis/backend/server/app/pkg/validation"
)

// ValidateCreateIdeaPayload validates create idea payload
func ValidateCreateIdeaPayload(payload *createIdeaPayload) error {
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

	// Validate position coordinates (optional, but if provided, validate range)
	if payload.PosX < -10000 || payload.PosX > 10000 {
		return fmt.Errorf("pos_x: must be between -10000 and 10000")
	}
	if payload.PosY < -10000 || payload.PosY > 10000 {
		return fmt.Errorf("pos_y: must be between -10000 and 10000")
	}

	// Validate color (optional, but if provided, validate format - hex color)
	if payload.Color != "" {
		if len(payload.Color) != 7 || payload.Color[0] != '#' {
			return fmt.Errorf("color: must be a valid hex color (e.g., #FF0000)")
		}
	}

	// Validate project ID (optional, but if provided, validate UUID)
	if payload.ProjectID != "" {
		if err := validation.ValidateUUID(payload.ProjectID); err != nil {
			return fmt.Errorf("project_id: %w", err)
		}
	}

	return nil
}

// ValidateUpdateIdeaPayload validates update idea payload
func ValidateUpdateIdeaPayload(payload *updateIdeaPayload) error {
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

	// Validate color (optional, but if provided, validate format)
	if payload.Color != "" {
		if len(payload.Color) != 7 || payload.Color[0] != '#' {
			return fmt.Errorf("color: must be a valid hex color (e.g., #FF0000)")
		}
	}

	// Validate project ID (optional, but if provided, validate UUID)
	if payload.ProjectID != "" {
		if err := validation.ValidateUUID(payload.ProjectID); err != nil {
			return fmt.Errorf("project_id: %w", err)
		}
	}

	return nil
}

// ValidateMoveIdeaPayload validates move idea payload
func ValidateMoveIdeaPayload(payload *moveIdeaPayload) error {
	// Validate position coordinates
	if payload.PosX < -10000 || payload.PosX > 10000 {
		return fmt.Errorf("pos_x: must be between -10000 and 10000")
	}
	if payload.PosY < -10000 || payload.PosY > 10000 {
		return fmt.Errorf("pos_y: must be between -10000 and 10000")
	}
	return nil
}

// ValidateBulkMovePayload validates bulk move payload
func ValidateBulkMovePayload(payload *bulkMovePayload) error {
	// Validate items (required, at least one)
	if len(payload.Items) == 0 {
		return fmt.Errorf("items: at least one item is required")
	}
	if len(payload.Items) > 100 {
		return fmt.Errorf("items: too many items (maximum 100)")
	}

	// Validate each item
	for i, item := range payload.Items {
		// Validate idea ID
		if err := validation.ValidateUUID(item.IdeaID); err != nil {
			return fmt.Errorf("items[%d].idea_id: %w", i, err)
		}

		// Validate position coordinates
		if item.PosX < -10000 || item.PosX > 10000 {
			return fmt.Errorf("items[%d].pos_x: must be between -10000 and 10000", i)
		}
		if item.PosY < -10000 || item.PosY > 10000 {
			return fmt.Errorf("items[%d].pos_y: must be between -10000 and 10000", i)
		}
	}

	return nil
}

// ValidateBulkUpdatePayload validates bulk update payload
func ValidateBulkUpdatePayload(payload *bulkUpdatePayload) error {
	// Validate items (required, at least one)
	if len(payload.Items) == 0 {
		return fmt.Errorf("items: at least one item is required")
	}
	if len(payload.Items) > 100 {
		return fmt.Errorf("items: too many items (maximum 100)")
	}

	// Validate each item
	for i, item := range payload.Items {
		// Validate idea ID
		if err := validation.ValidateUUID(item.IdeaID); err != nil {
			return fmt.Errorf("items[%d].idea_id: %w", i, err)
		}

		// Validate title (required, 1-200 chars)
		if err := validation.ValidateString(item.Title, 1, 200, fmt.Sprintf("items[%d].title", i)); err != nil {
			return fmt.Errorf("items[%d].title: %w", i, err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(item.Title); err != nil {
			return fmt.Errorf("items[%d].title: %w", i, err)
		}
		if err := validation.ValidateNoXSS(item.Title); err != nil {
			return fmt.Errorf("items[%d].title: %w", i, err)
		}

		// Validate description (optional, but if provided, validate)
		if item.Description != "" {
			if err := validation.ValidateString(item.Description, 1, 5000, fmt.Sprintf("items[%d].description", i)); err != nil {
				return fmt.Errorf("items[%d].description: %w", i, err)
			}
			// Check for SQL injection and XSS
			if err := validation.ValidateNoSQLInjection(item.Description); err != nil {
				return fmt.Errorf("items[%d].description: %w", i, err)
			}
			if err := validation.ValidateNoXSS(item.Description); err != nil {
				return fmt.Errorf("items[%d].description: %w", i, err)
			}
		}

		// Validate color (optional, but if provided, validate format)
		if item.Color != "" {
			if len(item.Color) != 7 || item.Color[0] != '#' {
				return fmt.Errorf("items[%d].color: must be a valid hex color (e.g., #FF0000)", i)
			}
		}

		// Validate project ID (optional, but if provided, validate UUID)
		if item.ProjectID != "" {
			if err := validation.ValidateUUID(item.ProjectID); err != nil {
				return fmt.Errorf("items[%d].project_id: %w", i, err)
			}
		}
	}

	return nil
}

// ValidateBulkIDsPayload validates bulk IDs payload
func ValidateBulkIDsPayload(payload *bulkIDsPayload) error {
	// Validate idea IDs (required, at least one)
	if len(payload.IdeaIDs) == 0 {
		return fmt.Errorf("idea_ids: at least one idea ID is required")
	}
	if len(payload.IdeaIDs) > 100 {
		return fmt.Errorf("idea_ids: too many idea IDs (maximum 100)")
	}

	// Validate each idea ID
	for i, idStr := range payload.IdeaIDs {
		if err := validation.ValidateUUID(idStr); err != nil {
			return fmt.Errorf("idea_ids[%d]: %w", i, err)
		}
	}

	return nil
}

// ValidateCreateLinkPayload validates create link payload
func ValidateCreateLinkPayload(payload *createLinkPayload) error {
	// Validate source idea ID (required, UUID)
	if payload.SourceIdeaID == "" {
		return fmt.Errorf("source_idea_id is required")
	}
	if err := validation.ValidateUUID(payload.SourceIdeaID); err != nil {
		return fmt.Errorf("source_idea_id: %w", err)
	}

	// Validate target idea ID (required, UUID)
	if payload.TargetIdeaID == "" {
		return fmt.Errorf("target_idea_id is required")
	}
	if err := validation.ValidateUUID(payload.TargetIdeaID); err != nil {
		return fmt.Errorf("target_idea_id: %w", err)
	}

	// Validate relation (required)
	if payload.Relation == "" {
		return fmt.Errorf("relation is required")
	}
	validRelations := []string{"related", "depends_on", "blocks", "duplicate"}
	isValid := false
	for _, validRelation := range validRelations {
		if payload.Relation == validRelation {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("relation: must be one of: related, depends_on, blocks, duplicate")
	}

	// Validate weight (optional, but if provided, validate range)
	if payload.Weight < 0 || payload.Weight > 1 {
		return fmt.Errorf("weight: must be between 0 and 1")
	}

	return nil
}

// ValidateCollaboratorPayload validates collaborator payload
func ValidateCollaboratorPayload(payload *collaboratorPayload) error {
	// Validate owner ID (required, UUID)
	if payload.OwnerID == "" {
		return fmt.Errorf("owner_id is required")
	}
	if err := validation.ValidateUUID(payload.OwnerID); err != nil {
		return fmt.Errorf("owner_id: %w", err)
	}

	// Validate collaborator ID (required, UUID)
	if payload.CollaboratorID == "" {
		return fmt.Errorf("collaborator_id is required")
	}
	if err := validation.ValidateUUID(payload.CollaboratorID); err != nil {
		return fmt.Errorf("collaborator_id: %w", err)
	}

	// Validate role (required)
	if payload.Role == "" {
		return fmt.Errorf("role is required")
	}
	validRoles := []string{"viewer", "editor", "owner"}
	isValid := false
	for _, validRole := range validRoles {
		if payload.Role == validRole {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("role: must be one of: viewer, editor, owner")
	}

	return nil
}

