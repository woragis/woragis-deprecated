package jobapplications

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	apikeysdomain "github.com/woragis/backend/server/app/internal/domains/apikeys"
	authdomain "github.com/woragis/backend/server/app/internal/domains/auth"
	aiservice "github.com/woragis/backend/server/app/internal/services/ai"
	"github.com/woragis/backend/server/app/pkg/response"
)

// ConversationCreator is an interface for creating conversations.
// This interface breaks the import cycle between jobapplications and chats domains.
type ConversationCreator interface {
	CreateConversation(ctx context.Context, req CreateConversationRequest) (*Conversation, error)
}

// ConversationCreatorFunc is a function type that implements ConversationCreator.
type ConversationCreatorFunc func(ctx context.Context, req CreateConversationRequest) (*Conversation, error)

// CreateConversation implements ConversationCreator interface.
func (f ConversationCreatorFunc) CreateConversation(ctx context.Context, req CreateConversationRequest) (*Conversation, error) {
	return f(ctx, req)
}

// CreateConversationRequest represents a request to create a conversation.
type CreateConversationRequest struct {
	UserID           uuid.UUID
	Title            string
	Description      string
	JobApplicationID *uuid.UUID
}

// Conversation represents a conversation (minimal interface to avoid import cycle).
type Conversation struct {
	ID uuid.UUID
}

// ResumeService is an interface to avoid circular dependencies with resumes domain.
type ResumeService interface {
	GetResume(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error)
}

// Resume represents a resume (minimal interface to avoid import cycle).
type Resume struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userId"`
	Title      string    `json:"title"`
	IsMain     bool      `json:"isMain"`
	IsFeatured bool      `json:"isFeatured"`
	FilePath   string    `json:"filePath"`
	FileName   string    `json:"fileName"`
	FileSize   int64     `json:"fileSize"`
	Tags       []string  `json:"tags"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// CoverLetterGenerator is an interface for generating cover letters.
type CoverLetterGenerator interface {
	GenerateCoverLetterWithContext(ctx context.Context, profile aiservice.UserProfile, job aiservice.JobInfo, additionalContext string) (string, error)
}

// Handler exposes job application endpoints.
type Handler interface {
	CreateJobApplication(c *fiber.Ctx) error
	GetJobApplication(c *fiber.Ctx) error
	ListJobApplications(c *fiber.Ctx) error
	UpdateJobApplicationStatus(c *fiber.Ctx) error
	UpdateJobApplication(c *fiber.Ctx) error
	DeleteJobApplication(c *fiber.Ctx) error
	GenerateCoverLetter(c *fiber.Ctx) error
}

type handler struct {
	service          Service
	conversationCreator ConversationCreator // Optional: for auto-creating conversations
	resumeService    ResumeService          // Optional: for including resume data in responses
	coverLetterGenerator CoverLetterGenerator // Optional: for generating cover letters
	logger          *slog.Logger
}

// NewHandler constructs a job application handler.
func NewHandler(service Service, logger *slog.Logger) Handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

// NewHandlerWithChatService constructs a job application handler with conversation creator.
func NewHandlerWithChatService(service Service, conversationCreator ConversationCreator, logger *slog.Logger) Handler {
	return &handler{
		service:            service,
		conversationCreator: conversationCreator,
		logger:             logger,
	}
}

// NewHandlerWithDependencies constructs a job application handler with all dependencies.
func NewHandlerWithDependencies(service Service, conversationCreator ConversationCreator, resumeService ResumeService, coverLetterGenerator CoverLetterGenerator, logger *slog.Logger) Handler {
	return &handler{
		service:            service,
		conversationCreator: conversationCreator,
		resumeService:      resumeService,
		coverLetterGenerator: coverLetterGenerator,
		logger:             logger,
	}
}

type createJobApplicationPayload struct {
	CompanyName   string   `json:"companyName"`
	Location      string   `json:"location"`
	JobTitle      string   `json:"jobTitle"`
	JobURL        string   `json:"jobUrl"`
	Website       string   `json:"website"`
	InterestLevel string   `json:"interestLevel,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	FollowUpDate  string   `json:"followUpDate,omitempty"`
	Notes         string   `json:"notes,omitempty"`
}

type updateStatusPayload struct {
	Status ApplicationStatus `json:"status"`
}

func (h *handler) CreateJobApplication(c *fiber.Ctx) error {
	// Get user ID from context
	var userID uuid.UUID
	var err error
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		userID = apiKey.UserID
	} else {
		userID, err = authdomain.UserIDFromContext(c)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, 401, fiber.Map{
				"message": "authentication required",
			})
		}
	}

	var payload createJobApplicationPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request payload",
		})
	}

	// Validate payload
	if err := ValidateCreateJobApplicationPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": err.Error(),
		})
	}

	// Normalize website to lowercase
	payload.Website = strings.ToLower(strings.TrimSpace(payload.Website))
	
	// Normalize URL (add https:// prefix if missing)
	payload.JobURL = normalizeURL(payload.JobURL)

	application, err := h.service.RequestJobApplication(
		c.Context(),
		userID,
		payload.CompanyName,
		payload.Location,
		payload.JobTitle,
		payload.JobURL,
		payload.Website,
	)
	if err != nil {
		return h.handleError(c, err)
	}

	// Update additional fields if provided
	updates := UpdateJobApplicationRequest{}
	hasUpdates := false

	if payload.InterestLevel != "" {
		updates.InterestLevel = &payload.InterestLevel
		hasUpdates = true
	}
	if len(payload.Tags) > 0 {
		updates.Tags = JSONArray(payload.Tags)
		hasUpdates = true
	}
	if payload.FollowUpDate != "" {
		followUpDate, err := time.Parse(time.RFC3339, payload.FollowUpDate)
		if err == nil {
			updates.FollowUpDate = &followUpDate
			hasUpdates = true
		}
	}
	if payload.Notes != "" {
		updates.Notes = &payload.Notes
		hasUpdates = true
	}

	if hasUpdates {
		application, err = h.service.UpdateJobApplication(c.Context(), application.ID, updates)
		if err != nil {
			// Log error but don't fail the request
			h.logger.Warn("failed to update application fields", slog.Any("error", err))
		}
	}

	// Auto-create conversation for this job application
	if h.conversationCreator != nil {
		jobAppID := application.ID
		title := fmt.Sprintf("Job Application: %s at %s", payload.JobTitle, payload.CompanyName)
		description := fmt.Sprintf("Chat about the %s position at %s", payload.JobTitle, payload.CompanyName)
		
		_, err := h.conversationCreator.CreateConversation(c.Context(), CreateConversationRequest{
			UserID:           userID,
			Title:            title,
			Description:      description,
			JobApplicationID: &jobAppID,
		})
		if err != nil {
			// Log error but don't fail the request
			h.logger.Warn("failed to auto-create conversation for job application", 
				slog.String("application_id", jobAppID.String()),
				slog.Any("error", err))
		}
	}

	return response.Success(c, fiber.StatusCreated, application)
}

func (h *handler) GetJobApplication(c *fiber.Ctx) error {
	applicationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid application id",
		})
	}

	application, err := h.service.GetJobApplication(c.Context(), applicationID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Include full resume data if resumeId exists
	responseData := fiber.Map{
		"id":                  application.ID,
		"userId":             application.UserID,
		"companyName":        application.CompanyName,
		"location":            application.Location,
		"jobTitle":            application.JobTitle,
		"jobUrl":              application.JobURL,
		"website":             application.Website,
		"appliedAt":           application.AppliedAt,
		"coverLetter":         application.CoverLetter,
		"linkedInContact":     application.LinkedInContact,
		"status":              application.Status,
		"errorMessage":        application.ErrorMessage,
		"resumeId":            application.ResumeID,
		"salaryMin":           application.SalaryMin,
		"salaryMax":           application.SalaryMax,
		"salaryCurrency":     application.SalaryCurrency,
		"jobDescription":      application.JobDescription,
		"deadline":            application.Deadline,
		"interestLevel":      application.InterestLevel,
		"notes":               application.Notes,
		"tags":                application.Tags,
		"followUpDate":        application.FollowUpDate,
		"responseReceivedAt":  application.ResponseReceivedAt,
		"rejectionReason":     application.RejectionReason,
		"interviewCount":      application.InterviewCount,
		"nextInterviewDate":   application.NextInterviewDate,
		"source":              application.Source,
		"applicationMethod":   application.ApplicationMethod,
		"language":            application.Language,
		"createdAt":           application.CreatedAt,
		"updatedAt":           application.UpdatedAt,
	}

	// If resumeId exists and resumeService is available, fetch full resume data
	if application.ResumeID != nil && h.resumeService != nil {
		resume, err := h.resumeService.GetResume(c.Context(), application.UserID, *application.ResumeID)
		if err == nil {
			responseData["resume"] = fiber.Map{
				"id":         resume.ID,
				"userId":     resume.UserID,
				"title":      resume.Title,
				"isMain":     resume.IsMain,
				"isFeatured": resume.IsFeatured,
				"filePath":   resume.FilePath,
				"fileName":   resume.FileName,
				"fileSize":   resume.FileSize,
				"tags":       resume.Tags,
				"createdAt":  resume.CreatedAt,
				"updatedAt":  resume.UpdatedAt,
			}
		} else {
			// Log error but don't fail the request
			h.logger.Warn("failed to fetch resume data",
				slog.String("resume_id", application.ResumeID.String()),
				slog.Any("error", err),
			)
		}
	}

	return response.Success(c, fiber.StatusOK, responseData)
}

func (h *handler) ListJobApplications(c *fiber.Ctx) error {
	// Get user ID from context
	var userID *uuid.UUID
	if apiKey, hasAPIKey := apikeysdomain.APIKeyFromContext(c); hasAPIKey {
		id := apiKey.UserID
		userID = &id
	} else if uid, err := authdomain.UserIDFromContext(c); err == nil {
		userID = &uid
	}

	filters := JobApplicationFilters{
		UserID: userID,
	}

	// Capture query parameters for validation
	website := c.Query("website")
	status := c.Query("status")
	resumeIDStr := c.Query("resumeId")
	interestLevel := c.Query("interestLevel")
	source := c.Query("source")
	applicationMethod := c.Query("applicationMethod")
	language := c.Query("language")
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	// Validate query parameters
	if err := ValidateListJobApplicationsQueryParams(limit, offset, website, status, resumeIDStr, interestLevel, source, applicationMethod, language); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": err.Error(),
		})
	}

	// Optional query parameters
	if website != "" {
		normalizedWebsite := strings.ToLower(strings.TrimSpace(website))
		filters.Website = &normalizedWebsite
	}
	if status != "" {
		appStatus := ApplicationStatus(status)
		filters.Status = &appStatus
	}
	if resumeIDStr != "" {
		if resumeID, err := uuid.Parse(resumeIDStr); err == nil {
			filters.ResumeID = &resumeID
		}
	}
	if interestLevel != "" {
		filters.InterestLevel = &interestLevel
	}
	if source != "" {
		filters.Source = &source
	}
	if applicationMethod != "" {
		filters.ApplicationMethod = &applicationMethod
	}
	if language != "" {
		filters.Language = &language
	}

	// Pagination
	if limit > 0 {
		filters.Limit = limit
	}
	if offset > 0 {
		filters.Offset = offset
	}

	applications, err := h.service.ListJobApplications(c.Context(), filters)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"applications": applications,
		"count":        len(applications),
	})
}

func (h *handler) UpdateJobApplicationStatus(c *fiber.Ctx) error {
	applicationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid application id",
		})
	}

	var payload updateStatusPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request payload",
		})
	}

	// Validate payload
	if err := ValidateUpdateStatusPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": err.Error(),
		})
	}

	if err := h.service.UpdateJobApplicationStatus(c.Context(), applicationID, payload.Status); err != nil {
		return h.handleError(c, err)
	}

	application, err := h.service.GetJobApplication(c.Context(), applicationID)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, application)
}

type updateJobApplicationPayload struct {
	ResumeID          *string   `json:"resumeId,omitempty"`
	SalaryMin         *int      `json:"salaryMin,omitempty"`
	SalaryMax         *int      `json:"salaryMax,omitempty"`
	SalaryCurrency    *string   `json:"salaryCurrency,omitempty"`
	JobDescription    *string   `json:"jobDescription,omitempty"`
	Deadline          *string   `json:"deadline,omitempty"` // ISO 8601 format
	InterestLevel     *string   `json:"interestLevel,omitempty"`
	Notes             *string   `json:"notes,omitempty"`
	Tags              JSONArray `json:"tags,omitempty"`
	FollowUpDate      *string   `json:"followUpDate,omitempty"` // ISO 8601 format
	ResponseReceivedAt *string  `json:"responseReceivedAt,omitempty"` // ISO 8601 format
	RejectionReason   *string   `json:"rejectionReason,omitempty"`
	NextInterviewDate *string   `json:"nextInterviewDate,omitempty"` // ISO 8601 format
	Source            *string   `json:"source,omitempty"`
	ApplicationMethod *string  `json:"applicationMethod,omitempty"`
	Language          *string   `json:"language,omitempty"` // ISO 639-1 language code (2 characters)
}

func (h *handler) UpdateJobApplication(c *fiber.Ctx) error {
	applicationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid application id",
		})
	}

	var payload updateJobApplicationPayload
	if err := c.BodyParser(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid request payload",
		})
	}

	// Validate payload
	if err := ValidateUpdateJobApplicationPayload(&payload); err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": err.Error(),
		})
	}

	updates := UpdateJobApplicationRequest{}

	if payload.ResumeID != nil {
		resumeID, err := uuid.Parse(*payload.ResumeID)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
				"message": "invalid resume id",
			})
		}
		updates.ResumeID = &resumeID
	}
	updates.SalaryMin = payload.SalaryMin
	updates.SalaryMax = payload.SalaryMax
	updates.SalaryCurrency = payload.SalaryCurrency
	updates.JobDescription = payload.JobDescription
	if payload.Deadline != nil {
		deadline, err := time.Parse(time.RFC3339, *payload.Deadline)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
				"message": "invalid deadline format, use ISO 8601",
			})
		}
		updates.Deadline = &deadline
	}
	updates.InterestLevel = payload.InterestLevel
	updates.Notes = payload.Notes
	if payload.Tags != nil {
		updates.Tags = payload.Tags
	}
	if payload.FollowUpDate != nil {
		followUpDate, err := time.Parse(time.RFC3339, *payload.FollowUpDate)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
				"message": "invalid follow-up date format, use ISO 8601",
			})
		}
		updates.FollowUpDate = &followUpDate
	}
	if payload.ResponseReceivedAt != nil {
		responseReceivedAt, err := time.Parse(time.RFC3339, *payload.ResponseReceivedAt)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
				"message": "invalid response received date format, use ISO 8601",
			})
		}
		updates.ResponseReceivedAt = &responseReceivedAt
	}
	updates.RejectionReason = payload.RejectionReason
	if payload.NextInterviewDate != nil {
		nextInterviewDate, err := time.Parse(time.RFC3339, *payload.NextInterviewDate)
		if err != nil {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
				"message": "invalid next interview date format, use ISO 8601",
			})
		}
		updates.NextInterviewDate = &nextInterviewDate
	}
	updates.Source = payload.Source
	updates.ApplicationMethod = payload.ApplicationMethod
	if payload.Language != nil {
		// Validate language code is exactly 2 characters
		if len(*payload.Language) != 2 {
			return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
				"message": "language must be exactly 2 characters (ISO 639-1 code)",
			})
		}
		updates.Language = payload.Language
	}

	application, err := h.service.UpdateJobApplication(c.Context(), applicationID, updates)
	if err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, application)
}

func (h *handler) DeleteJobApplication(c *fiber.Ctx) error {
	applicationID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.Error(c, fiber.StatusBadRequest, ErrCodeInvalidPayload, fiber.Map{
			"message": "invalid application id",
		})
	}

	userID, err := authdomain.UserIDFromContext(c)
	if err != nil {
		return err
	}

	// Verify the application belongs to the user
	application, err := h.service.GetJobApplication(c.Context(), applicationID)
	if err != nil {
		return h.handleError(c, err)
	}

	if application.UserID != userID {
		return response.Error(c, fiber.StatusForbidden, ErrCodeAccessDenied, fiber.Map{
			"message": "access denied",
		})
	}

	// Delete the job application
	// Note: Related conversations will have their job_application_id set to NULL
	// to preserve chat history while unlinking from the deleted application
	if err := h.service.DeleteJobApplication(c.Context(), applicationID); err != nil {
		return h.handleError(c, err)
	}

	return response.Success(c, fiber.StatusOK, fiber.Map{
		"message": "job application deleted successfully",
	})
}

func (h *handler) handleError(c *fiber.Ctx, err error) error {
	if domainErr, ok := AsDomainError(err); ok {
		statusCode := fiber.StatusInternalServerError
		switch domainErr.Code {
		case ErrCodeNotFound:
			statusCode = fiber.StatusNotFound
		case ErrCodeInvalidPayload, ErrCodeInvalidStatus:
			statusCode = fiber.StatusBadRequest
		case ErrCodeJobQueueFailure, ErrCodeAIServiceFailure, ErrCodePlaywrightFailure:
			statusCode = fiber.StatusServiceUnavailable
		case ErrCodeDatabaseConstraint, ErrCodeDatabaseValueTooLong, ErrCodeDatabaseUniqueViolation, ErrCodeDatabaseForeignKeyViolation:
			statusCode = fiber.StatusBadRequest
		case ErrCodeDatabaseConnection:
			statusCode = fiber.StatusServiceUnavailable
		case ErrCodeAccessDenied:
			statusCode = fiber.StatusForbidden
		}

		return response.Error(c, statusCode, domainErr.Code, fiber.Map{
			"message": domainErr.Message,
		})
	}

	h.logger.Error("unhandled error", slog.Any("error", err))
	return response.Error(c, fiber.StatusInternalServerError, 500, fiber.Map{
		"message": "internal server error",
	})
}

