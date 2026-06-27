package resumes

import (
	"fmt"
	"strings"

	"woragis-social-media-service/pkg/validation"
)

// createResumePayload represents the payload for CreateResume
type createResumePayload struct {
	Title    string   `json:"title"`
	FilePath string   `json:"filePath"`
	FileName string   `json:"fileName"`
	FileSize int64    `json:"fileSize"`
	Tags     []string `json:"tags"`
}

// updateResumePayload represents the payload for UpdateResume
type updateResumePayload struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags"`
}

// generateResumePayload represents the payload for GenerateResume
type generateResumePayload struct {
	JobApplicationID string `json:"jobApplicationId"`
	Language         string `json:"language"`
	Template         string `json:"template,omitempty"`
}

// ValidateCreateResumePayload validates create resume payload
func ValidateCreateResumePayload(payload *createResumePayload) error {
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

	// Validate file path (optional, but if provided, validate)
	if payload.FilePath != "" {
		if err := validation.ValidateString(payload.FilePath, 1, 500, "filePath"); err != nil {
			return fmt.Errorf("filePath: %w", err)
		}
		// Check for path traversal
		if strings.Contains(payload.FilePath, "..") {
			return fmt.Errorf("filePath: invalid path")
		}
	}

	// Validate file name (optional, but if provided, validate)
	if payload.FileName != "" {
		if err := validation.ValidateString(payload.FileName, 1, 255, "fileName"); err != nil {
			return fmt.Errorf("fileName: %w", err)
		}
		// Validate file extension
		allowedExts := []string{".pdf", ".doc", ".docx"}
		if err := validation.ValidateFileExtension(payload.FileName, allowedExts); err != nil {
			return fmt.Errorf("fileName: %w", err)
		}
	}

	// Validate file size (optional, but if provided, validate)
	if payload.FileSize > 0 {
		maxSize := int64(10 * 1024 * 1024) // 10MB
		if err := validation.ValidateFileSize(payload.FileSize, maxSize); err != nil {
			return fmt.Errorf("fileSize: %w", err)
		}
	}

	// Validate tags (optional, but if provided, validate each tag)
	if len(payload.Tags) > 0 {
		if len(payload.Tags) > 20 {
			return fmt.Errorf("tags: too many tags (maximum 20)")
		}
		for i, tag := range payload.Tags {
			if err := validation.ValidateString(tag, 1, 50, fmt.Sprintf("tags[%d]", i)); err != nil {
				return fmt.Errorf("tags[%d]: %w", i, err)
			}
			// Check for SQL injection and XSS
			if err := validation.ValidateNoSQLInjection(tag); err != nil {
				return fmt.Errorf("tags[%d]: %w", i, err)
			}
			if err := validation.ValidateNoXSS(tag); err != nil {
				return fmt.Errorf("tags[%d]: %w", i, err)
			}
		}
	}

	return nil
}

// ValidateUpdateResumePayload validates update resume payload
func ValidateUpdateResumePayload(payload *updateResumePayload) error {
	// Validate title (optional, but if provided, validate)
	if payload.Title != "" {
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
	}

	// Validate tags (optional, but if provided, validate each tag)
	if len(payload.Tags) > 0 {
		if len(payload.Tags) > 20 {
			return fmt.Errorf("tags: too many tags (maximum 20)")
		}
		for i, tag := range payload.Tags {
			if err := validation.ValidateString(tag, 1, 50, fmt.Sprintf("tags[%d]", i)); err != nil {
				return fmt.Errorf("tags[%d]: %w", i, err)
			}
			// Check for SQL injection and XSS
			if err := validation.ValidateNoSQLInjection(tag); err != nil {
				return fmt.Errorf("tags[%d]: %w", i, err)
			}
			if err := validation.ValidateNoXSS(tag); err != nil {
				return fmt.Errorf("tags[%d]: %w", i, err)
			}
		}
	}

	return nil
}

// ValidateGenerateResumePayload validates generate resume payload
func ValidateGenerateResumePayload(payload *generateResumePayload) error {
	// Validate job application ID (required)
	if payload.JobApplicationID == "" {
		return fmt.Errorf("jobApplicationId is required")
	}
	if err := validation.ValidateUUID(payload.JobApplicationID); err != nil {
		return fmt.Errorf("jobApplicationId: %w", err)
	}

	// Validate language (required, ISO 639-1 format)
	if payload.Language == "" {
		return fmt.Errorf("language is required")
	}
	if len(payload.Language) != 2 {
		return fmt.Errorf("language: must be exactly 2 characters (ISO 639-1 code)")
	}
	// Should be lowercase
	if payload.Language != strings.ToLower(payload.Language) {
		return fmt.Errorf("language: must be lowercase ISO 639-1 code")
	}

	// Validate template (optional, but if provided, validate)
	if payload.Template != "" {
		if err := validation.ValidateString(payload.Template, 1, 100, "template"); err != nil {
			return fmt.Errorf("template: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(payload.Template); err != nil {
			return fmt.Errorf("template: %w", err)
		}
		if err := validation.ValidateNoXSS(payload.Template); err != nil {
			return fmt.Errorf("template: %w", err)
		}
	}

	return nil
}

// ValidateListResumesQueryParams validates query parameters for ListResumes
func ValidateListResumesQueryParams(limit, offset int, search string) error {
	// Validate limit
	if limit < 1 {
		return fmt.Errorf("limit: must be at least 1")
	}
	if limit > 200 {
		return fmt.Errorf("limit: must be at most 200")
	}

	// Validate offset
	if offset < 0 {
		return fmt.Errorf("offset: must be at least 0")
	}

	// Validate search (optional, but if provided, validate length)
	if search != "" {
		if err := validation.ValidateString(search, 1, 200, "search"); err != nil {
			return fmt.Errorf("search: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(search); err != nil {
			return fmt.Errorf("search: %w", err)
		}
		if err := validation.ValidateNoXSS(search); err != nil {
			return fmt.Errorf("search: %w", err)
		}
	}

	return nil
}

// ValidateUploadResumeFile validates uploaded file
func ValidateUploadResumeFile(filename string, size int64, contentType string) error {
	// Validate file extension
	allowedExts := []string{".pdf"}
	if err := validation.ValidateFileExtension(filename, allowedExts); err != nil {
		return fmt.Errorf("file: %w", err)
	}

	// Validate file size (max 10MB)
	maxSize := int64(10 * 1024 * 1024)
	if err := validation.ValidateFileSize(size, maxSize); err != nil {
		return fmt.Errorf("file: %w", err)
	}

	// Validate content type
	if contentType != "" && contentType != "application/pdf" {
		return fmt.Errorf("file: only PDF files are allowed")
	}

	return nil
}

