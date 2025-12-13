package ai

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	langchainservice "github.com/woragis/backend/server/app/internal/services/langchain"
)

// CoverLetterService generates personalized cover letters using AI.
type CoverLetterService struct {
	aiClient *langchainservice.Client
	logger   *slog.Logger
}

// NewCoverLetterService creates a new cover letter service.
func NewCoverLetterService(aiClient *langchainservice.Client, logger *slog.Logger) *CoverLetterService {
	return &CoverLetterService{
		aiClient: aiClient,
		logger:   logger,
	}
}

// UserProfile represents user data for cover letter generation.
type UserProfile struct {
	Projects          []ProjectInfo
	Posts             []PostInfo
	TechnicalWritings []TechnicalWritingInfo
	Skills            []string
	Interests         []string
	Certifications    []string
}

// ProjectInfo represents a project.
type ProjectInfo struct {
	Name        string
	Description string
	TechStack   []string
}

// PostInfo represents a post.
type PostInfo struct {
	Title   string
	Content string
}

// TechnicalWritingInfo represents technical writing.
type TechnicalWritingInfo struct {
	Title   string
	Content string
}

// JobInfo represents job information.
type JobInfo struct {
	CompanyName    string
	JobTitle       string
	JobDescription string
	Location       string
	Requirements   []string
}

// GenerateCoverLetter generates a personalized cover letter.
func (cls *CoverLetterService) GenerateCoverLetter(ctx context.Context, profile UserProfile, job JobInfo) (string, error) {
	return cls.GenerateCoverLetterWithContext(ctx, profile, job, "")
}

// GenerateCoverLetterWithContext generates a personalized cover letter with additional context.
func (cls *CoverLetterService) GenerateCoverLetterWithContext(ctx context.Context, profile UserProfile, job JobInfo, additionalContext string) (string, error) {
	cls.logger.Info("generating cover letter",
		slog.String("company", job.CompanyName),
		slog.String("jobTitle", job.JobTitle),
	)

	// Build prompt
	prompt := cls.buildPrompt(profile, job)
	
	// Add additional context if provided (e.g., from chat message)
	if additionalContext != "" {
		prompt += fmt.Sprintf("\n\nAdditional Context from Conversation:\n%s\n\nUse this context to further personalize the cover letter.", additionalContext)
	}

	// Call AI service
	req := langchainservice.ChatCompletionRequest{
		Provider:    langchainservice.ProviderOpenAI,
		Model:       "gpt-4o-mini",
		Temperature: 0.7, // Slightly higher for more creative cover letters
		Messages: []langchainservice.ChatMessage{
			{
				Role:      "user",
				Content:   prompt,
				Timestamp: time.Now(),
			},
		},
		MaxTokens: 1500,
	}

	resp, err := cls.aiClient.GenerateCompletion(ctx, req)
	if err != nil {
		cls.logger.Error("failed to generate cover letter", slog.Any("error", err))
		return "", fmt.Errorf("failed to generate cover letter: %w", err)
	}

	coverLetter := resp.Message.Content
	cls.logger.Info("cover letter generated", slog.Int("length", len(coverLetter)))

	return coverLetter, nil
}

// buildPrompt builds the prompt for cover letter generation.
func (cls *CoverLetterService) buildPrompt(profile UserProfile, job JobInfo) string {
	prompt := fmt.Sprintf(`You are a professional cover letter writer. Write a personalized cover letter for the following job application.

Job Information:
- Company: %s
- Position: %s
- Location: %s
- Job Description: %s

Applicant Profile:
`, job.CompanyName, job.JobTitle, job.Location, job.JobDescription)

	// Add projects
	if len(profile.Projects) > 0 {
		prompt += "\nProjects:\n"
		for _, project := range profile.Projects {
			prompt += fmt.Sprintf("- %s: %s (Tech: %v)\n", project.Name, project.Description, project.TechStack)
		}
	}

	// Add skills
	if len(profile.Skills) > 0 {
		prompt += fmt.Sprintf("\nSkills: %v\n", profile.Skills)
	}

	// Add technical writings
	if len(profile.TechnicalWritings) > 0 {
		prompt += "\nTechnical Writings:\n"
		for _, writing := range profile.TechnicalWritings {
			prompt += fmt.Sprintf("- %s: %s\n", writing.Title, writing.Content[:min(200, len(writing.Content))])
		}
	}

	prompt += `
Instructions:
1. Write a professional, engaging cover letter
2. Highlight relevant experience and skills from the applicant's profile
3. Show enthusiasm for the specific role and company
4. Keep it concise (3-4 paragraphs)
5. Use a professional but personable tone
6. Do not include placeholders or generic statements

Write the cover letter now:`

	return prompt
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

