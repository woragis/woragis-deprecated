package playwright

import (
	"fmt"

	"github.com/woragis/backend/server/app/pkg/validation"
	jobapplicationsdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications"
)

// ValidateApplyToJobRequest validates the request to apply to a job
func ValidateApplyToJobRequest(job *jobapplicationsdomain.JobApplicationJob, coverLetter string) error {
	// Validate job URL (required, must be valid URL)
	if job.JobURL == "" {
		return fmt.Errorf("job_url is required")
	}
	if err := validation.ValidateURL(job.JobURL); err != nil {
		return fmt.Errorf("job_url: %w", err)
	}

	// Validate website (required)
	if job.Website == "" {
		return fmt.Errorf("website is required")
	}
	if err := validation.ValidateString(job.Website, 1, 255, "website"); err != nil {
		return fmt.Errorf("website: %w", err)
	}

	// Validate company name (required)
	if err := validation.ValidateString(job.CompanyName, 1, 200, "company_name"); err != nil {
		return fmt.Errorf("company_name: %w", err)
	}

	// Validate cover letter (required, 100-10000 chars)
	if err := validation.ValidateString(coverLetter, 100, 10000, "cover_letter"); err != nil {
		return fmt.Errorf("cover_letter: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(coverLetter); err != nil {
		return fmt.Errorf("cover_letter: %w", err)
	}
	if err := validation.ValidateNoXSS(coverLetter); err != nil {
		return fmt.Errorf("cover_letter: %w", err)
	}

	return nil
}

