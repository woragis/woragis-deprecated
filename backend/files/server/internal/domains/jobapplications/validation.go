package jobapplications

import (
	"fmt"
	"strings"

	"woragis-files-service/pkg/validation"
)

// ValidateCreateJobApplicationPayload validates create job application payload
func ValidateCreateJobApplicationPayload(payload *createJobApplicationPayload) error {
	// Validate company name (required, 1-200 chars)
	if err := validation.ValidateString(payload.CompanyName, 1, 200, "companyName"); err != nil {
		return fmt.Errorf("companyName: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(payload.CompanyName); err != nil {
		return fmt.Errorf("companyName: %w", err)
	}
	if err := validation.ValidateNoXSS(payload.CompanyName); err != nil {
		return fmt.Errorf("companyName: %w", err)
	}

	// Validate location (optional, but if provided, validate)
	if payload.Location != "" {
		if err := validation.ValidateString(payload.Location, 1, 200, "location"); err != nil {
			return fmt.Errorf("location: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(payload.Location); err != nil {
			return fmt.Errorf("location: %w", err)
		}
		if err := validation.ValidateNoXSS(payload.Location); err != nil {
			return fmt.Errorf("location: %w", err)
		}
	}

	// Validate job title (required, 1-200 chars)
	if err := validation.ValidateString(payload.JobTitle, 1, 200, "jobTitle"); err != nil {
		return fmt.Errorf("jobTitle: %w", err)
	}
	// Check for SQL injection and XSS
	if err := validation.ValidateNoSQLInjection(payload.JobTitle); err != nil {
		return fmt.Errorf("jobTitle: %w", err)
	}
	if err := validation.ValidateNoXSS(payload.JobTitle); err != nil {
		return fmt.Errorf("jobTitle: %w", err)
	}

	// Validate job URL (required, must be valid URL)
	if payload.JobURL != "" {
		if err := validation.ValidateURL(payload.JobURL); err != nil {
			return fmt.Errorf("jobUrl: %w", err)
		}
	}

	// Validate website (optional, but if provided, validate format)
	if payload.Website != "" {
		website := strings.ToLower(strings.TrimSpace(payload.Website))
		// Website should be a domain name, not a full URL
		if len(website) > 255 {
			return fmt.Errorf("website: too long (maximum 255 characters)")
		}
		// Basic domain validation (can be enhanced)
		if strings.Contains(website, " ") {
			return fmt.Errorf("website: invalid format")
		}
	}

	// Validate interest level (optional, but if provided, validate)
	if payload.InterestLevel != "" {
		validLevels := []string{"low", "medium", "high", "very_high"}
		isValid := false
		for _, level := range validLevels {
			if strings.ToLower(payload.InterestLevel) == level {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("interestLevel: must be one of: low, medium, high, very_high")
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

	// Validate notes (optional, but if provided, validate length)
	if payload.Notes != "" {
		if err := validation.ValidateString(payload.Notes, 1, 5000, "notes"); err != nil {
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

// ValidateUpdateStatusPayload validates update status payload
func ValidateUpdateStatusPayload(payload *updateStatusPayload) error {
	// Status is validated by the ApplicationStatus type itself
	// But we can add additional validation if needed
	if payload.Status == "" {
		return fmt.Errorf("status is required")
	}
	return nil
}

// ValidateUpdateJobApplicationPayload validates update job application payload
func ValidateUpdateJobApplicationPayload(payload *updateJobApplicationPayload) error {
	// Validate resume ID (optional, but if provided, validate UUID)
	if payload.ResumeID != nil && *payload.ResumeID != "" {
		if err := validation.ValidateUUID(*payload.ResumeID); err != nil {
			return fmt.Errorf("resumeId: %w", err)
		}
	}

	// Validate salary min (optional, but if provided, validate range)
	if payload.SalaryMin != nil {
		if *payload.SalaryMin < 0 {
			return fmt.Errorf("salaryMin: must be non-negative")
		}
		if *payload.SalaryMin > 10000000 {
			return fmt.Errorf("salaryMin: must be at most 10,000,000")
		}
	}

	// Validate salary max (optional, but if provided, validate range)
	if payload.SalaryMax != nil {
		if *payload.SalaryMax < 0 {
			return fmt.Errorf("salaryMax: must be non-negative")
		}
		if *payload.SalaryMax > 10000000 {
			return fmt.Errorf("salaryMax: must be at most 10,000,000")
		}
	}

	// Validate salary currency (optional, but if provided, validate format)
	if payload.SalaryCurrency != nil && *payload.SalaryCurrency != "" {
		currency := *payload.SalaryCurrency
		if err := validation.ValidateString(currency, 3, 3, "salaryCurrency"); err != nil {
			return fmt.Errorf("salaryCurrency: %w", err)
		}
		// Should be uppercase ISO 4217 code
		if currency != strings.ToUpper(currency) {
			return fmt.Errorf("salaryCurrency: must be uppercase ISO 4217 code")
		}
	}

	// Validate job description (optional, but if provided, validate)
	if payload.JobDescription != nil && *payload.JobDescription != "" {
		desc := *payload.JobDescription
		if err := validation.ValidateString(desc, 1, 10000, "jobDescription"); err != nil {
			return fmt.Errorf("jobDescription: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(desc); err != nil {
			return fmt.Errorf("jobDescription: %w", err)
		}
		if err := validation.ValidateNoXSS(desc); err != nil {
			return fmt.Errorf("jobDescription: %w", err)
		}
	}

	// Validate interest level (optional, but if provided, validate)
	if payload.InterestLevel != nil && *payload.InterestLevel != "" {
		validLevels := []string{"low", "medium", "high", "very_high"}
		isValid := false
		for _, level := range validLevels {
			if strings.ToLower(*payload.InterestLevel) == level {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("interestLevel: must be one of: low, medium, high, very_high")
		}
	}

	// Validate notes (optional, but if provided, validate)
	if payload.Notes != nil && *payload.Notes != "" {
		notes := *payload.Notes
		if err := validation.ValidateString(notes, 1, 5000, "notes"); err != nil {
			return fmt.Errorf("notes: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(notes); err != nil {
			return fmt.Errorf("notes: %w", err)
		}
		if err := validation.ValidateNoXSS(notes); err != nil {
			return fmt.Errorf("notes: %w", err)
		}
	}

	// Validate tags (optional, but if provided, validate each tag)
	if payload.Tags != nil && len(payload.Tags) > 0 {
		if len(payload.Tags) > 20 {
			return fmt.Errorf("tags: too many tags (maximum 20)")
		}
		for i, tag := range payload.Tags {
			tagStr := fmt.Sprintf("%v", tag)
			if err := validation.ValidateString(tagStr, 1, 50, fmt.Sprintf("tags[%d]", i)); err != nil {
				return fmt.Errorf("tags[%d]: %w", i, err)
			}
			// Check for SQL injection and XSS
			if err := validation.ValidateNoSQLInjection(tagStr); err != nil {
				return fmt.Errorf("tags[%d]: %w", i, err)
			}
			if err := validation.ValidateNoXSS(tagStr); err != nil {
				return fmt.Errorf("tags[%d]: %w", i, err)
			}
		}
	}

	// Validate rejection reason (optional, but if provided, validate)
	if payload.RejectionReason != nil && *payload.RejectionReason != "" {
		reason := *payload.RejectionReason
		if err := validation.ValidateString(reason, 1, 500, "rejectionReason"); err != nil {
			return fmt.Errorf("rejectionReason: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(reason); err != nil {
			return fmt.Errorf("rejectionReason: %w", err)
		}
		if err := validation.ValidateNoXSS(reason); err != nil {
			return fmt.Errorf("rejectionReason: %w", err)
		}
	}

	// Validate source (optional, but if provided, validate)
	if payload.Source != nil && *payload.Source != "" {
		source := *payload.Source
		if err := validation.ValidateString(source, 1, 100, "source"); err != nil {
			return fmt.Errorf("source: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(source); err != nil {
			return fmt.Errorf("source: %w", err)
		}
		if err := validation.ValidateNoXSS(source); err != nil {
			return fmt.Errorf("source: %w", err)
		}
	}

	// Validate application method (optional, but if provided, validate)
	if payload.ApplicationMethod != nil && *payload.ApplicationMethod != "" {
		method := *payload.ApplicationMethod
		if err := validation.ValidateString(method, 1, 100, "applicationMethod"); err != nil {
			return fmt.Errorf("applicationMethod: %w", err)
		}
		// Check for SQL injection and XSS
		if err := validation.ValidateNoSQLInjection(method); err != nil {
			return fmt.Errorf("applicationMethod: %w", err)
		}
		if err := validation.ValidateNoXSS(method); err != nil {
			return fmt.Errorf("applicationMethod: %w", err)
		}
	}

	// Validate language (optional, but if provided, validate ISO 639-1 format)
	if payload.Language != nil && *payload.Language != "" {
		lang := *payload.Language
		if len(lang) != 2 {
			return fmt.Errorf("language: must be exactly 2 characters (ISO 639-1 code)")
		}
		// Should be lowercase
		if lang != strings.ToLower(lang) {
			return fmt.Errorf("language: must be lowercase ISO 639-1 code")
		}
	}

	return nil
}

// ValidateListJobApplicationsQueryParams validates query parameters for ListJobApplications
func ValidateListJobApplicationsQueryParams(limit, offset int, website, status, resumeIDStr, interestLevel, source, applicationMethod, language string) error {
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

	// Validate website (optional)
	if website != "" {
		if err := validation.ValidateString(website, 1, 255, "website"); err != nil {
			return fmt.Errorf("website: %w", err)
		}
	}

	// Validate status (optional)
	if status != "" {
		validStatuses := []string{"draft", "applied", "interviewing", "offer", "rejected", "withdrawn", "accepted"}
		isValid := false
		for _, validStatus := range validStatuses {
			if strings.ToLower(status) == validStatus {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("status: must be one of: draft, applied, interviewing, offer, rejected, withdrawn, accepted")
		}
	}

	// Validate resume ID (optional)
	if resumeIDStr != "" {
		if err := validation.ValidateUUID(resumeIDStr); err != nil {
			return fmt.Errorf("resumeId: %w", err)
		}
	}

	// Validate interest level (optional)
	if interestLevel != "" {
		validLevels := []string{"low", "medium", "high", "very_high"}
		isValid := false
		for _, level := range validLevels {
			if strings.ToLower(interestLevel) == level {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("interestLevel: must be one of: low, medium, high, very_high")
		}
	}

	// Validate source (optional)
	if source != "" {
		if err := validation.ValidateString(source, 1, 100, "source"); err != nil {
			return fmt.Errorf("source: %w", err)
		}
	}

	// Validate application method (optional)
	if applicationMethod != "" {
		if err := validation.ValidateString(applicationMethod, 1, 100, "applicationMethod"); err != nil {
			return fmt.Errorf("applicationMethod: %w", err)
		}
	}

	// Validate language (optional)
	if language != "" {
		if len(language) != 2 {
			return fmt.Errorf("language: must be exactly 2 characters (ISO 639-1 code)")
		}
	}

	return nil
}

