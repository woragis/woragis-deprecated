package chats

import (
	"fmt"

	"github.com/woragis/backend/server/app/pkg/validation"
)

// ValidateCreateConversationPayload validates create conversation payload
func ValidateCreateConversationPayload(payload *createConversationPayload) error {
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
		if err := validation.ValidateString(payload.Description, 1, 1000, "description"); err != nil {
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

	// Validate idea ID (optional, but if provided, validate UUID)
	if payload.IdeaID != "" {
		if err := validation.ValidateUUID(payload.IdeaID); err != nil {
			return fmt.Errorf("ideaId: %w", err)
		}
	}

	// Validate project ID (optional, but if provided, validate UUID)
	if payload.ProjectID != "" {
		if err := validation.ValidateUUID(payload.ProjectID); err != nil {
			return fmt.Errorf("projectId: %w", err)
		}
	}

	// Validate job application ID (optional, but if provided, validate UUID)
	if payload.JobApplicationID != "" {
		if err := validation.ValidateUUID(payload.JobApplicationID); err != nil {
			return fmt.Errorf("jobApplicationId: %w", err)
		}
	}

	return nil
}

// ValidateAppendMessagePayload validates append message payload
func ValidateAppendMessagePayload(payload *appendMessagePayload) error {
	// Validate role (required)
	if payload.Role == "" {
		return fmt.Errorf("role is required")
	}
	validRoles := []string{"user", "assistant", "system"}
	isValid := false
	for _, validRole := range validRoles {
		if payload.Role == validRole {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("role: must be one of: user, assistant, system")
	}

	// Validate content (required, 1-50000 chars)
	if err := validation.ValidateString(payload.Content, 1, 50000, "content"); err != nil {
		return fmt.Errorf("content: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(payload.Content); err != nil {
		return fmt.Errorf("content: %w", err)
	}
	if err := validation.ValidateNoXSS(payload.Content); err != nil {
		return fmt.Errorf("content: %w", err)
	}

	// Validate agent (optional, but if provided, validate)
	if payload.Agent != "" {
		if err := validation.ValidateString(payload.Agent, 1, 100, "agent"); err != nil {
			return fmt.Errorf("agent: %w", err)
		}
	}

	// Validate provider (optional, but if provided, validate)
	if payload.Provider != "" {
		if err := validation.ValidateString(payload.Provider, 1, 100, "provider"); err != nil {
			return fmt.Errorf("provider: %w", err)
		}
	}

	// Validate model (optional, but if provided, validate)
	if payload.Model != "" {
		if err := validation.ValidateString(payload.Model, 1, 100, "model"); err != nil {
			return fmt.Errorf("model: %w", err)
		}
	}

	// Validate max tokens (optional, but if provided, validate range)
	if payload.MaxTokens > 0 {
		if payload.MaxTokens < 1 {
			return fmt.Errorf("max_tokens: must be at least 1")
		}
		if payload.MaxTokens > 100000 {
			return fmt.Errorf("max_tokens: must be at most 100,000")
		}
	}

	// Validate temperature (optional, but if provided, validate range)
	if payload.Temperature != 0 {
		if payload.Temperature < 0 {
			return fmt.Errorf("temperature: must be at least 0")
		}
		if payload.Temperature > 2 {
			return fmt.Errorf("temperature: must be at most 2")
		}
	}

	return nil
}

// ValidateBulkUpdatePayload validates bulk update payload
func ValidateBulkUpdatePayload(payload *bulkUpdatePayload) error {
	// Validate conversation IDs (required, at least one)
	if len(payload.ConversationIDs) == 0 {
		return fmt.Errorf("conversation_ids: at least one conversation ID is required")
	}
	if len(payload.ConversationIDs) > 100 {
		return fmt.Errorf("conversation_ids: too many conversation IDs (maximum 100)")
	}

	// Validate each conversation ID
	for i, idStr := range payload.ConversationIDs {
		if err := validation.ValidateUUID(idStr); err != nil {
			return fmt.Errorf("conversation_ids[%d]: %w", i, err)
		}
	}

	return nil
}

// ValidateAssignmentPayload validates assignment payload
func ValidateAssignmentPayload(payload *assignmentPayload) error {
	// Validate agent ID (required)
	if payload.AgentID == "" {
		return fmt.Errorf("agent_id is required")
	}
	if err := validation.ValidateUUID(payload.AgentID); err != nil {
		return fmt.Errorf("agent_id: %w", err)
	}

	// Validate agent name (optional, but if provided, validate)
	if payload.AgentName != "" {
		if err := validation.ValidateString(payload.AgentName, 1, 100, "agent_name"); err != nil {
			return fmt.Errorf("agent_name: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(payload.AgentName); err != nil {
			return fmt.Errorf("agent_name: %w", err)
		}
		if err := validation.ValidateNoXSS(payload.AgentName); err != nil {
			return fmt.Errorf("agent_name: %w", err)
		}
	}

	// Validate notes (optional, but if provided, validate)
	if payload.Notes != "" {
		if err := validation.ValidateString(payload.Notes, 1, 1000, "notes"); err != nil {
			return fmt.Errorf("notes: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(payload.Notes); err != nil {
			return fmt.Errorf("notes: %w", err)
		}
		if err := validation.ValidateNoXSS(payload.Notes); err != nil {
			return fmt.Errorf("notes: %w", err)
		}
	}

	return nil
}

// ValidateSearchQueryParams validates search query parameters
func ValidateSearchQueryParams(query string, limit int) error {
	// Validate query (optional, but if provided, validate length)
	if query != "" {
		if err := validation.ValidateString(query, 1, 200, "query"); err != nil {
			return fmt.Errorf("query: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(query); err != nil {
			return fmt.Errorf("query: %w", err)
		}
		if err := validation.ValidateNoXSS(query); err != nil {
			return fmt.Errorf("query: %w", err)
		}
	}

	// Validate limit
	if limit < 1 {
		return fmt.Errorf("limit: must be at least 1")
	}
	if limit > 200 {
		return fmt.Errorf("limit: must be at most 200")
	}

	return nil
}

