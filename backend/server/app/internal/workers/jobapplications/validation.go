package jobapplications

import (
	"fmt"

	"github.com/woragis/backend/server/app/pkg/validation"
	jobapplicationsdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications"
)

// ValidateJobApplicationJob validates a job application job before processing
func ValidateJobApplicationJob(job *jobapplicationsdomain.JobApplicationJob) error {
	// Validate user ID (required, UUID)
	if job.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if err := validation.ValidateUUID(job.UserID); err != nil {
		return fmt.Errorf("user_id: %w", err)
	}

	// Validate company name (required, 1-200 chars)
	if err := validation.ValidateString(job.CompanyName, 1, 200, "company_name"); err != nil {
		return fmt.Errorf("company_name: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(job.CompanyName); err != nil {
		return fmt.Errorf("company_name: %w", err)
	}
	if err := validation.ValidateNoXSS(job.CompanyName); err != nil {
		return fmt.Errorf("company_name: %w", err)
	}

	// Validate job title (required, 1-200 chars)
	if err := validation.ValidateString(job.JobTitle, 1, 200, "job_title"); err != nil {
		return fmt.Errorf("job_title: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(job.JobTitle); err != nil {
		return fmt.Errorf("job_title: %w", err)
	}
	if err := validation.ValidateNoXSS(job.JobTitle); err != nil {
		return fmt.Errorf("job_title: %w", err)
	}

	// Validate job URL (required, must be valid URL)
	if job.JobURL == "" {
		return fmt.Errorf("job_url is required")
	}
	if err := validation.ValidateURL(job.JobURL); err != nil {
		return fmt.Errorf("job_url: %w", err)
	}

	// Validate website (required, 1-255 chars)
	if job.Website == "" {
		return fmt.Errorf("website is required")
	}
	if err := validation.ValidateString(job.Website, 1, 255, "website"); err != nil {
		return fmt.Errorf("website: %w", err)
	}

	// Validate location (optional, but if provided, validate)
	if job.Location != "" {
		if err := validation.ValidateString(job.Location, 1, 200, "location"); err != nil {
			return fmt.Errorf("location: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(job.Location); err != nil {
			return fmt.Errorf("location: %w", err)
		}
		if err := validation.ValidateNoXSS(job.Location); err != nil {
			return fmt.Errorf("location: %w", err)
		}
	}

	// Validate job ID (required, UUID)
	if job.ID == "" {
		return fmt.Errorf("job id is required")
	}
	if err := validation.ValidateUUID(job.ID); err != nil {
		return fmt.Errorf("job id: %w", err)
	}

	return nil
}

