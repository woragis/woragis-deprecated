package ai

import (
	"fmt"

	"github.com/woragis/backend/server/app/pkg/validation"
)

// ValidateUserProfile validates user profile data
func ValidateUserProfile(profile UserProfile) error {
	// Validate projects (optional, but if provided, validate each)
	if len(profile.Projects) > 100 {
		return fmt.Errorf("projects: too many projects (maximum 100)")
	}
	for i, project := range profile.Projects {
		if err := validation.ValidateString(project.Name, 1, 200, fmt.Sprintf("projects[%d].name", i)); err != nil {
			return fmt.Errorf("projects[%d].name: %w", i, err)
		}
		if err := validation.ValidateString(project.Description, 1, 5000, fmt.Sprintf("projects[%d].description", i)); err != nil {
			return fmt.Errorf("projects[%d].description: %w", i, err)
		}
		if len(project.TechStack) > 50 {
			return fmt.Errorf("projects[%d].tech_stack: too many technologies (maximum 50)", i)
		}
		for j, tech := range project.TechStack {
			if err := validation.ValidateString(tech, 1, 100, fmt.Sprintf("projects[%d].tech_stack[%d]", i, j)); err != nil {
				return fmt.Errorf("projects[%d].tech_stack[%d]: %w", i, j, err)
			}
		}
	}

	// Validate posts (optional, but if provided, validate each)
	if len(profile.Posts) > 100 {
		return fmt.Errorf("posts: too many posts (maximum 100)")
	}
	for i, post := range profile.Posts {
		if err := validation.ValidateString(post.Title, 1, 500, fmt.Sprintf("posts[%d].title", i)); err != nil {
			return fmt.Errorf("posts[%d].title: %w", i, err)
		}
		if err := validation.ValidateString(post.Content, 1, 50000, fmt.Sprintf("posts[%d].content", i)); err != nil {
			return fmt.Errorf("posts[%d].content: %w", i, err)
		}
	}

	// Validate technical writings (optional, but if provided, validate each)
	if len(profile.TechnicalWritings) > 100 {
		return fmt.Errorf("technical_writings: too many technical writings (maximum 100)")
	}
	for i, writing := range profile.TechnicalWritings {
		if err := validation.ValidateString(writing.Title, 1, 500, fmt.Sprintf("technical_writings[%d].title", i)); err != nil {
			return fmt.Errorf("technical_writings[%d].title: %w", i, err)
		}
		if err := validation.ValidateString(writing.Content, 1, 100000, fmt.Sprintf("technical_writings[%d].content", i)); err != nil {
			return fmt.Errorf("technical_writings[%d].content: %w", i, err)
		}
	}

	// Validate skills (optional, but if provided, validate each)
	if len(profile.Skills) > 200 {
		return fmt.Errorf("skills: too many skills (maximum 200)")
	}
	for i, skill := range profile.Skills {
		if err := validation.ValidateString(skill, 1, 100, fmt.Sprintf("skills[%d]", i)); err != nil {
			return fmt.Errorf("skills[%d]: %w", i, err)
		}
	}

	// Validate interests (optional, but if provided, validate each)
	if len(profile.Interests) > 200 {
		return fmt.Errorf("interests: too many interests (maximum 200)")
	}
	for i, interest := range profile.Interests {
		if err := validation.ValidateString(interest, 1, 100, fmt.Sprintf("interests[%d]", i)); err != nil {
			return fmt.Errorf("interests[%d]: %w", i, err)
		}
	}

	// Validate certifications (optional, but if provided, validate each)
	if len(profile.Certifications) > 100 {
		return fmt.Errorf("certifications: too many certifications (maximum 100)")
	}
	for i, cert := range profile.Certifications {
		if err := validation.ValidateString(cert, 1, 200, fmt.Sprintf("certifications[%d]", i)); err != nil {
			return fmt.Errorf("certifications[%d]: %w", i, err)
		}
	}

	return nil
}

// ValidateJobInfo validates job information
func ValidateJobInfo(job JobInfo) error {
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

	// Validate job description (optional, but if provided, validate)
	if job.JobDescription != "" {
		if err := validation.ValidateString(job.JobDescription, 1, 50000, "job_description"); err != nil {
			return fmt.Errorf("job_description: %w", err)
		}
	}

	// Validate location (optional, but if provided, validate)
	if job.Location != "" {
		if err := validation.ValidateString(job.Location, 1, 200, "location"); err != nil {
			return fmt.Errorf("location: %w", err)
		}
	}

	// Validate requirements (optional, but if provided, validate each)
	if len(job.Requirements) > 100 {
		return fmt.Errorf("requirements: too many requirements (maximum 100)")
	}
	for i, req := range job.Requirements {
		if err := validation.ValidateString(req, 1, 500, fmt.Sprintf("requirements[%d]", i)); err != nil {
			return fmt.Errorf("requirements[%d]: %w", i, err)
		}
	}

	return nil
}

