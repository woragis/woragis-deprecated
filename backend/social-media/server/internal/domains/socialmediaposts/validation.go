package socialmediaposts

import (
	"fmt"

	"woragis-social-media-service/pkg/validation"
)

// ValidateCreatePostPayload validates create post payload
func ValidateCreatePostPayload(payload *createPostPayload) error {
	// Validate title (required, 1-500 chars)
	if err := validation.ValidateString(payload.Title, 1, 500, "title"); err != nil {
		return fmt.Errorf("title: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(payload.Title); err != nil {
		return fmt.Errorf("title: %w", err)
	}
	if err := validation.ValidateNoXSS(payload.Title); err != nil {
		return fmt.Errorf("title: %w", err)
	}

	// Validate content (required, 1-10000 chars)
	if err := validation.ValidateString(payload.Content, 1, 10000, "content"); err != nil {
		return fmt.Errorf("content: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(payload.Content); err != nil {
		return fmt.Errorf("content: %w", err)
	}
	if err := validation.ValidateNoXSS(payload.Content); err != nil {
		return fmt.Errorf("content: %w", err)
	}

	// Validate content post ID (optional, but if provided, validate UUID)
	if payload.ContentPostID != nil && *payload.ContentPostID != "" {
		if err := validation.ValidateUUID(*payload.ContentPostID); err != nil {
			return fmt.Errorf("contentPostId: %w", err)
		}
	}

	return nil
}

// ValidateUpdatePostPayload validates update post payload
func ValidateUpdatePostPayload(payload *updatePostPayload) error {
	// Validate title (optional, but if provided, validate)
	if payload.Title != nil && *payload.Title != "" {
		if err := validation.ValidateString(*payload.Title, 1, 500, "title"); err != nil {
			return fmt.Errorf("title: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(*payload.Title); err != nil {
			return fmt.Errorf("title: %w", err)
		}
		if err := validation.ValidateNoXSS(*payload.Title); err != nil {
			return fmt.Errorf("title: %w", err)
		}
	}

	// Validate content (optional, but if provided, validate)
	if payload.Content != nil && *payload.Content != "" {
		if err := validation.ValidateString(*payload.Content, 1, 10000, "content"); err != nil {
			return fmt.Errorf("content: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(*payload.Content); err != nil {
			return fmt.Errorf("content: %w", err)
		}
		if err := validation.ValidateNoXSS(*payload.Content); err != nil {
			return fmt.Errorf("content: %w", err)
		}
	}

	return nil
}

// ValidateUpdateEngagementPayload validates update engagement payload
func ValidateUpdateEngagementPayload(payload *updateEngagementPayload) error {
	// Validate likes (optional, but if provided, validate range)
	if payload.Likes != nil {
		if *payload.Likes < 0 {
			return fmt.Errorf("likes: must be non-negative")
		}
		if *payload.Likes > 1000000000 {
			return fmt.Errorf("likes: must be at most 1,000,000,000")
		}
	}

	// Validate shares (optional, but if provided, validate range)
	if payload.Shares != nil {
		if *payload.Shares < 0 {
			return fmt.Errorf("shares: must be non-negative")
		}
		if *payload.Shares > 1000000000 {
			return fmt.Errorf("shares: must be at most 1,000,000,000")
		}
	}

	// Validate comments (optional, but if provided, validate range)
	if payload.Comments != nil {
		if *payload.Comments < 0 {
			return fmt.Errorf("comments: must be non-negative")
		}
		if *payload.Comments > 1000000000 {
			return fmt.Errorf("comments: must be at most 1,000,000,000")
		}
	}

	// Validate views (optional, but if provided, validate range)
	if payload.Views != nil {
		if *payload.Views < 0 {
			return fmt.Errorf("views: must be non-negative")
		}
		if *payload.Views > 10000000000 {
			return fmt.Errorf("views: must be at most 10,000,000,000")
		}
	}

	return nil
}

// ValidateCreateLinkPayload validates create link payload
func ValidateCreateLinkPayload(payload *createLinkPayload) error {
	// PostID, EntityID are UUIDs and validated by type system
	// EntityType and RelationshipType are enums validated by type system
	// No additional validation needed beyond type checking
	return nil
}

// ValidateUpdateLinkPayload validates update link payload
func ValidateUpdateLinkPayload(payload *updateLinkPayload) error {
	// RelationshipType is an enum validated by type system
	// No additional validation needed
	return nil
}

