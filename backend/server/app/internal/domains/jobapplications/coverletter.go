package jobapplications

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	aiservice "github.com/woragis/backend/server/app/internal/services/ai"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	"github.com/woragis/backend/server/app/pkg/response"
)

// MessageService is an interface to fetch message content from chat service.
// This will be implemented later when we inject chat service.
type MessageService interface {
	GetMessageContent(ctx context.Context, messageID uuid.UUID, userID uuid.UUID) (string, error)
}

type generateCoverLetterPayload struct {
	MessageID *string `json:"messageId,omitempty"` // Optional: message ID from chat to use as additional context
}

func (h *handler) GenerateCoverLetter(c *fiber.Ctx) error {
	applicationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid application id",
		})
	}

	// Get user ID from context
	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
			"message": "authentication required",
		})
	}

	// Get job application
	application, err := h.service.GetJobApplication(c.Context(), applicationID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Verify ownership
	if application.UserID != userID {
		return response.Error(c, fiber.StatusForbidden, 403, fiber.Map{
			"message": "access denied",
		})
	}

	// Parse optional payload
	var payload generateCoverLetterPayload
	if err := c.BodyParser(&payload); err != nil {
		// Payload is optional, so we ignore parsing errors
		payload = generateCoverLetterPayload{}
	}

	// Get additional context from message if provided
	additionalContext := ""
	if payload.MessageID != nil && *payload.MessageID != "" {
		messageID, err := uuid.Parse(*payload.MessageID)
		if err == nil {
			// TODO: Inject MessageService to fetch message content
			// For now, we'll leave this empty and implement it later
			// This would require injecting a chat service to fetch the message
			// The message content would be added to the prompt as additional context
			_ = messageID // Suppress unused variable warning
		}
	}

	// Check if cover letter generator is available
	if h.coverLetterGenerator == nil {
		return response.Error(c, fiber.StatusNotImplemented, 501, fiber.Map{
			"message": "cover letter generation not available",
		})
	}

	// Build user profile for cover letter generation
	profile := aiservice.UserProfile{
		Projects:          []aiservice.ProjectInfo{},
		Posts:             []aiservice.PostInfo{},
		TechnicalWritings: []aiservice.TechnicalWritingInfo{},
		Skills:            []string{},
		Interests:         []string{},
		Certifications:    []string{},
	}

	// Note: We would need to inject services to fetch user data
	// For now, we'll use a simplified approach and fetch what we can
	// This is a placeholder - in production, you'd inject these services

	// Build job info
	jobInfo := aiservice.JobInfo{
		CompanyName:    application.CompanyName,
		JobTitle:       application.JobTitle,
		JobDescription: application.JobDescription,
		Location:       application.Location,
		Requirements:   []string{}, // Could parse from job description in the future
	}

	// Generate cover letter
	coverLetter, err := h.coverLetterGenerator.GenerateCoverLetterWithContext(
		c.Context(),
		profile,
		jobInfo,
		additionalContext,
	)
	if err != nil {
		h.logger.Error("failed to generate cover letter", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to generate cover letter",
		})
	}

	// Update job application with generated cover letter
	updates := UpdateJobApplicationRequest{
		CoverLetter: &coverLetter,
	}
	updatedApplication, err := h.service.UpdateJobApplication(c.Context(), applicationID, updates)
	if err != nil {
		h.logger.Error("failed to update cover letter", slog.Any("error", err))
		return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
			"message": "failed to update cover letter",
		})
	}

	return response.Success(c, fiber.StatusOK, updatedApplication)
}
