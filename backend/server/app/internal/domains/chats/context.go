package chats

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	jobapplicationsdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications"
	projectsdomain "github.com/woragis/backend/server/app/internal/domains/projects"
	resumesdomain "github.com/woragis/backend/server/app/internal/domains/resumes"
	skillsdomain "github.com/woragis/backend/server/app/internal/domains/skills"
	userprofilesdomain "github.com/woragis/backend/server/app/internal/domains/userprofiles"
	casestudiesdomain "github.com/woragis/backend/server/app/internal/domains/casestudies"
	technicalwritingsdomain "github.com/woragis/backend/server/app/internal/domains/technicalwritings"
	postsdomain "github.com/woragis/backend/server/app/internal/domains/posts"
	problemsolutionsdomain "github.com/woragis/backend/server/app/internal/domains/problemsolutions"
	experiencesdomain "github.com/woragis/backend/server/app/internal/domains/experiences"
)

// ContextBuilder builds context for chat conversations.
type ContextBuilder struct {
	jobApplicationService jobapplicationsdomain.Service
	resumeService         resumesdomain.Service
	userProfileService    userprofilesdomain.Service
	projectService        projectsdomain.Service
	skillService          skillsdomain.Service
	caseStudyService      casestudiesdomain.Service
	technicalWritingService technicalwritingsdomain.Service
	postService           postsdomain.Service
	problemSolutionService problemsolutionsdomain.Service
	experienceService     experiencesdomain.Service
	logger                *slog.Logger
}

// NewContextBuilder creates a new context builder.
func NewContextBuilder(
	jobApplicationService jobapplicationsdomain.Service,
	resumeService resumesdomain.Service,
	userProfileService userprofilesdomain.Service,
	projectService projectsdomain.Service,
	skillService skillsdomain.Service,
	caseStudyService casestudiesdomain.Service,
	technicalWritingService technicalwritingsdomain.Service,
	postService postsdomain.Service,
	problemSolutionService problemsolutionsdomain.Service,
	experienceService experiencesdomain.Service,
	logger *slog.Logger,
) *ContextBuilder {
	return &ContextBuilder{
		jobApplicationService:  jobApplicationService,
		resumeService:          resumeService,
		userProfileService:     userProfileService,
		projectService:         projectService,
		skillService:           skillService,
		caseStudyService:       caseStudyService,
		technicalWritingService: technicalWritingService,
		postService:            postService,
		problemSolutionService: problemSolutionService,
		experienceService:     experienceService,
		logger:                logger,
	}
}

// BuildContextOptions specifies what context to include.
type BuildContextOptions struct {
	IncludeJobApplication bool
	IncludeResume         bool
	IncludeUserProfile    bool
	IncludeProjects       bool
	IncludeCaseStudies    bool
	IncludeTechnicalWritings bool
	IncludePosts          bool
	IncludeProblemSolutions bool
	IncludeSkills         bool
	IncludeExperiences    bool
}

// BuildContext builds context string for a conversation.
func (cb *ContextBuilder) BuildContext(ctx context.Context, userID uuid.UUID, conv *Conversation, options BuildContextOptions) (string, error) {
	var parts []string

	// User Profile (About Me)
	if options.IncludeUserProfile {
		profile, err := cb.userProfileService.GetProfile(ctx, userID)
		if err == nil && profile != nil && profile.AboutMe != "" {
			parts = append(parts, fmt.Sprintf("## About Me\n%s\n", profile.AboutMe))
		}
	}

	// Job Application Context
	if options.IncludeJobApplication && conv.JobApplicationID != nil {
		app, err := cb.jobApplicationService.GetJobApplication(ctx, *conv.JobApplicationID)
		if err == nil && app != nil {
			parts = append(parts, fmt.Sprintf("## Job Application Context\n"))
			parts = append(parts, fmt.Sprintf("**Company:** %s\n", app.CompanyName))
			parts = append(parts, fmt.Sprintf("**Job Title:** %s\n", app.JobTitle))
			parts = append(parts, fmt.Sprintf("**Location:** %s\n", app.Location))
			parts = append(parts, fmt.Sprintf("**Status:** %s\n", app.Status))
			if app.JobDescription != "" {
				parts = append(parts, fmt.Sprintf("**Job Description:**\n%s\n", app.JobDescription))
			}
			if app.Notes != "" {
				parts = append(parts, fmt.Sprintf("**Notes:** %s\n", app.Notes))
			}
			parts = append(parts, "\n")
		}

		// Resume Context (if linked to job application)
		if options.IncludeResume && app != nil && app.ResumeID != nil {
			resume, err := cb.resumeService.GetResume(ctx, userID, *app.ResumeID)
			if err == nil && resume != nil {
				parts = append(parts, fmt.Sprintf("## Resume Information\n"))
				parts = append(parts, fmt.Sprintf("**Title:** %s\n", resume.Title))
				parts = append(parts, fmt.Sprintf("**File:** %s (%d bytes)\n", resume.FileName, resume.FileSize))
				parts = append(parts, "\n")
			}
		}
	}

	// Projects
	if options.IncludeProjects {
		projects, err := cb.projectService.ListProjects(ctx, userID)
		if err == nil && len(projects) > 0 {
			parts = append(parts, "## Projects\n")
			for i, p := range projects {
				if i >= 10 { // Limit to 10 most recent
					parts = append(parts, fmt.Sprintf("... and %d more projects\n", len(projects)-10))
					break
				}
				parts = append(parts, fmt.Sprintf("- **%s**: %s (Status: %s)\n", p.Name, p.Description, p.Status))
			}
			parts = append(parts, "\n")
		}
	}

	// Skills
	if options.IncludeSkills {
		skills, err := cb.skillService.ListSkills(ctx)
		if err == nil && len(skills) > 0 {
			parts = append(parts, "## Skills\n")
			skillNames := make([]string, 0, len(skills))
			for _, s := range skills {
				if len(skillNames) >= 20 { // Limit to 20 skills
					break
				}
				skillNames = append(skillNames, s.Name)
			}
			parts = append(parts, strings.Join(skillNames, ", "))
			parts = append(parts, "\n\n")
		}
	}

	// Case Studies
	if options.IncludeCaseStudies {
		userIDPtr := &userID
		caseStudies, err := cb.caseStudyService.ListCaseStudies(ctx, casestudiesdomain.ListCaseStudiesFilters{
			UserID: userIDPtr,
			Limit:  10,
		})
		if err == nil && len(caseStudies) > 0 {
			parts = append(parts, "## Case Studies\n")
			for _, cs := range caseStudies {
				// Use Problem as the description since Description field doesn't exist
				description := cs.Problem
				if len(description) > 100 {
					description = description[:100] + "..."
				}
				parts = append(parts, fmt.Sprintf("- **%s**: %s\n", cs.Title, description))
			}
			parts = append(parts, "\n")
		}
	}

	// Technical Writings
	if options.IncludeTechnicalWritings {
		// Note: Would need to add ListTechnicalWritings method
		parts = append(parts, "## Technical Writings\n(Technical writings would be included here)\n\n")
	}

	// Posts
	if options.IncludePosts {
		// Note: Would need to add ListPosts method with user filter
		parts = append(parts, "## Posts\n(Posts would be included here)\n\n")
	}

	// Problem Solutions
	if options.IncludeProblemSolutions {
		// Note: Would need to add ListProblemSolutions method
		parts = append(parts, "## Problem Solutions\n(Problem solutions would be included here)\n\n")
	}

	// Experiences
	if options.IncludeExperiences {
		// Note: Would need to add ListExperiences method with user filter
		parts = append(parts, "## Experiences\n(Experiences would be included here)\n\n")
	}

	return strings.Join(parts, ""), nil
}

// GetDefaultContextOptions returns default context options for job application chats.
func GetDefaultContextOptions() BuildContextOptions {
	return BuildContextOptions{
		IncludeJobApplication:  true,
		IncludeResume:          true,
		IncludeUserProfile:     true,
		IncludeProjects:        true,
		IncludeCaseStudies:     true,
		IncludeTechnicalWritings: true,
		IncludePosts:           true,
		IncludeProblemSolutions: true,
		IncludeSkills:          true,
		IncludeExperiences:     true,
	}
}

