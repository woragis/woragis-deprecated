package jobapplications

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates job application workflows.
type Service interface {
	RequestJobApplication(ctx context.Context, userID uuid.UUID, companyName, location, jobTitle, jobURL, website string) (*JobApplication, error)
	GetJobApplication(ctx context.Context, applicationID uuid.UUID) (*JobApplication, error)
	ListJobApplications(ctx context.Context, filters JobApplicationFilters) ([]JobApplication, error)
	UpdateJobApplicationStatus(ctx context.Context, applicationID uuid.UUID, status ApplicationStatus) error
	UpdateJobApplication(ctx context.Context, applicationID uuid.UUID, updates UpdateJobApplicationRequest) (*JobApplication, error)
	DeleteJobApplication(ctx context.Context, applicationID uuid.UUID) error
	ProcessJobApplicationJob(ctx context.Context, job *JobApplicationJob) error
}

// UpdateJobApplicationRequest represents fields that can be updated on a job application.
type UpdateJobApplicationRequest struct {
	ResumeID          *uuid.UUID
	SalaryMin         *int
	SalaryMax         *int
	SalaryCurrency    *string
	JobDescription    *string
	CoverLetter       *string
	Deadline          *time.Time
	InterestLevel     *string
	Notes             *string
	Tags              JSONArray
	FollowUpDate      *time.Time
	ResponseReceivedAt *time.Time
	RejectionReason   *string
	NextInterviewDate *time.Time
	Source            *string
	ApplicationMethod *string
	Language          *string
}

// ResumeMetricsService is an interface to avoid circular dependencies
type ResumeMetricsService interface {
	RecalculateResumeMetrics(ctx context.Context, resumeID uuid.UUID) error
}

type service struct {
	repo                Repository
	queue               Queue
	chatsRepo           ChatsRepository // For unlinking conversations on delete
	preferencesService  UserPreferencesService // For getting user defaults
	resumeMetricsService ResumeMetricsService // Optional: for updating resume metrics
	logger              *slog.Logger
}

// ChatsRepository defines the interface needed from chats domain for unlinking conversations.
type ChatsRepository interface {
	UnlinkFromJobApplication(ctx context.Context, jobApplicationID uuid.UUID) error
}

// UserPreferencesService defines the interface needed from userpreferences domain.
type UserPreferencesService interface {
	GetDefaultLanguage(ctx context.Context, userID uuid.UUID) (string, error)
	GetDefaultCurrency(ctx context.Context, userID uuid.UUID) (string, error)
	GetDefaultWebsite(ctx context.Context, userID uuid.UUID) (string, error)
}

// NewService constructs a Service.
func NewService(repo Repository, queue Queue, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		queue:  queue,
		logger: logger,
	}
}

// NewServiceWithChats constructs a Service with chats repository for cascade operations.
func NewServiceWithChats(repo Repository, queue Queue, chatsRepo ChatsRepository, logger *slog.Logger) Service {
	return &service{
		repo:      repo,
		queue:     queue,
		chatsRepo: chatsRepo,
		logger:    logger,
	}
}

// NewServiceWithDependencies constructs a Service with all dependencies.
func NewServiceWithDependencies(repo Repository, queue Queue, chatsRepo ChatsRepository, preferencesService UserPreferencesService, logger *slog.Logger) Service {
	return &service{
		repo:               repo,
		queue:              queue,
		chatsRepo:          chatsRepo,
		preferencesService: preferencesService,
		logger:             logger,
	}
}

// NewServiceWithResumeMetrics constructs a Service with resume metrics service.
func NewServiceWithResumeMetrics(repo Repository, queue Queue, chatsRepo ChatsRepository, preferencesService UserPreferencesService, resumeMetricsService ResumeMetricsService, logger *slog.Logger) Service {
	return &service{
		repo:                repo,
		queue:               queue,
		chatsRepo:           chatsRepo,
		preferencesService:  preferencesService,
		resumeMetricsService: resumeMetricsService,
		logger:              logger,
	}
}

func (s *service) RequestJobApplication(ctx context.Context, userID uuid.UUID, companyName, location, jobTitle, jobURL, website string) (*JobApplication, error) {
	// Normalize website to lowercase
	website = strings.ToLower(strings.TrimSpace(website))
	
	// If website is blank, try to get from user preferences
	if website == "" && s.preferencesService != nil {
		if defaultWebsite, err := s.preferencesService.GetDefaultWebsite(ctx, userID); err == nil && defaultWebsite != "" {
			website = defaultWebsite
		}
	}
	
	// Normalize URL (add https:// prefix if missing)
	jobURL = normalizeURL(jobURL)
	
	// Create job application record
	application, err := NewJobApplication(userID, companyName, location, jobTitle, jobURL, website)
	if err != nil {
		return nil, err
	}

	// Apply user defaults if not already set
	if s.preferencesService != nil {
		if application.Language == "" {
			if defaultLang, err := s.preferencesService.GetDefaultLanguage(ctx, userID); err == nil {
				application.Language = defaultLang
			}
		}
		if application.SalaryCurrency == "" {
			if defaultCurrency, err := s.preferencesService.GetDefaultCurrency(ctx, userID); err == nil {
				application.SalaryCurrency = defaultCurrency
			}
		}
	}

	// Save to database
	if err := s.repo.CreateJobApplication(ctx, application); err != nil {
		return nil, err
	}

	// Create job for queue
	job := &JobApplicationJob{
		ID:          uuid.New().String(),
		UserID:      userID.String(),
		CompanyName: companyName,
		Location:    location,
		JobTitle:    jobTitle,
		JobURL:      jobURL,
		Website:     website,
	}

	// Enqueue job
	if err := s.queue.EnqueueJob(ctx, job); err != nil {
		// If queue fails, mark application as failed
		application.MarkFailed(fmt.Sprintf("Failed to enqueue job: %v", err))
		s.repo.UpdateJobApplication(ctx, application)
		return nil, err
	}

	// Mark as processing
	application.MarkProcessing()
	if err := s.repo.UpdateJobApplication(ctx, application); err != nil {
		return nil, err
	}

	return application, nil
}

func (s *service) GetJobApplication(ctx context.Context, applicationID uuid.UUID) (*JobApplication, error) {
	return s.repo.GetJobApplication(ctx, applicationID)
}

func (s *service) ListJobApplications(ctx context.Context, filters JobApplicationFilters) ([]JobApplication, error) {
	return s.repo.ListJobApplications(ctx, filters)
}

func (s *service) UpdateJobApplicationStatus(ctx context.Context, applicationID uuid.UUID, status ApplicationStatus) error {
	application, err := s.repo.GetJobApplication(ctx, applicationID)
	if err != nil {
		return err
	}

	oldStatus := application.Status
	if err := application.UpdateStatus(status); err != nil {
		return err
	}

	if err := s.repo.UpdateJobApplication(ctx, application); err != nil {
		return err
	}

	// Recalculate resume metrics if status changed to "accepted" or if resumeId exists
	if (oldStatus != status && status == ApplicationStatusAccepted) || application.ResumeID != nil {
		if s.resumeMetricsService != nil && application.ResumeID != nil {
			if err := s.resumeMetricsService.RecalculateResumeMetrics(ctx, *application.ResumeID); err != nil {
				s.logger.Warn("failed to recalculate resume metrics", "resume_id", application.ResumeID.String(), "error", err)
				// Don't fail the request if metric recalculation fails
			}
		}
	}

	return nil
}

func (s *service) UpdateJobApplication(ctx context.Context, applicationID uuid.UUID, updates UpdateJobApplicationRequest) (*JobApplication, error) {
	application, err := s.repo.GetJobApplication(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	if updates.ResumeID != nil {
		application.ResumeID = updates.ResumeID
	}
	if updates.SalaryMin != nil {
		application.SalaryMin = updates.SalaryMin
	}
	if updates.SalaryMax != nil {
		application.SalaryMax = updates.SalaryMax
	}
	if updates.SalaryCurrency != nil {
		application.SalaryCurrency = *updates.SalaryCurrency
	}
	if updates.JobDescription != nil {
		application.JobDescription = *updates.JobDescription
	}
	if updates.CoverLetter != nil {
		application.CoverLetter = *updates.CoverLetter
	}
	if updates.Deadline != nil {
		application.Deadline = updates.Deadline
	}
	if updates.InterestLevel != nil {
		application.InterestLevel = *updates.InterestLevel
	}
	if updates.Notes != nil {
		application.Notes = *updates.Notes
	}
	if updates.Tags != nil {
		application.Tags = updates.Tags
	}
	if updates.FollowUpDate != nil {
		application.FollowUpDate = updates.FollowUpDate
	}
	if updates.ResponseReceivedAt != nil {
		application.ResponseReceivedAt = updates.ResponseReceivedAt
	}
	if updates.RejectionReason != nil {
		application.RejectionReason = *updates.RejectionReason
	}
	if updates.NextInterviewDate != nil {
		application.NextInterviewDate = updates.NextInterviewDate
	}
	if updates.Source != nil {
		application.Source = *updates.Source
	}
	if updates.ApplicationMethod != nil {
		application.ApplicationMethod = *updates.ApplicationMethod
	}
	if updates.Language != nil {
		application.Language = *updates.Language
	}

	application.UpdatedAt = time.Now().UTC()

	// Track if resumeId changed
	var oldResumeID *uuid.UUID
	if application.ResumeID != nil {
		oldResumeIDCopy := *application.ResumeID
		oldResumeID = &oldResumeIDCopy
	}
	if updates.ResumeID != nil {
		application.ResumeID = updates.ResumeID
	}

	if err := s.repo.UpdateJobApplication(ctx, application); err != nil {
		return nil, err
	}

	// Recalculate metrics for both old and new resume if resumeId changed
	if s.resumeMetricsService != nil {
		if oldResumeID != nil && (application.ResumeID == nil || *oldResumeID != *application.ResumeID) {
			// Old resume metrics need updating
			if err := s.resumeMetricsService.RecalculateResumeMetrics(ctx, *oldResumeID); err != nil {
				s.logger.Warn("failed to recalculate old resume metrics", "resume_id", oldResumeID.String(), "error", err)
			}
		}
		if application.ResumeID != nil {
			// New resume metrics need updating
			if err := s.resumeMetricsService.RecalculateResumeMetrics(ctx, *application.ResumeID); err != nil {
				s.logger.Warn("failed to recalculate new resume metrics", "resume_id", application.ResumeID.String(), "error", err)
			}
		}
	}

	return application, nil
}

func (s *service) DeleteJobApplication(ctx context.Context, applicationID uuid.UUID) error {
	// First, verify the application exists
	application, err := s.repo.GetJobApplication(ctx, applicationID)
	if err != nil {
		return err
	}

	// Unlink conversations from this job application (preserve chat history)
	if s.chatsRepo != nil {
		if err := s.chatsRepo.UnlinkFromJobApplication(ctx, applicationID); err != nil {
			s.logger.Warn("failed to unlink conversations from job application",
				"application_id", applicationID.String(),
				"error", err)
			// Continue with deletion even if unlinking fails
		}
	}

	// Delete the job application
	if err := s.repo.DeleteJobApplication(ctx, applicationID); err != nil {
		return err
	}

	s.logger.Info("job application deleted",
		"application_id", applicationID.String(),
		"company_name", application.CompanyName,
		"job_title", application.JobTitle)

	return nil
}

// normalizeURL adds https:// prefix to URL if it's missing.
// Preserves existing http:// or https:// prefixes.
func normalizeURL(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return url
	}
	
	// Check if URL already has a protocol
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	
	// Add https:// prefix
	return "https://" + url
}

// ProcessJobApplicationJob is called by the worker to process a job.
// This is a placeholder - actual processing will be done in the worker with Playwright.
func (s *service) ProcessJobApplicationJob(ctx context.Context, job *JobApplicationJob) error {
	// This method will be called by the worker, but the actual application logic
	// will be in the worker itself using Playwright.
	// This is here for interface consistency.
	
	// Normalize website to lowercase
	job.Website = strings.ToLower(strings.TrimSpace(job.Website))
	
	// Normalize URL (add https:// prefix if missing)
	job.JobURL = normalizeURL(job.JobURL)
	
	userID, err := uuid.Parse(job.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Find or create application record
	filters := JobApplicationFilters{
		UserID: &userID,
	}
	applications, err := s.repo.ListJobApplications(ctx, filters)
	if err != nil {
		return err
	}

	var application *JobApplication
	for _, app := range applications {
		if app.JobURL == job.JobURL && app.Website == job.Website {
			application = &app
			break
		}
	}

	if application == nil {
		// Create new application
		application, err = NewJobApplication(userID, job.CompanyName, job.Location, job.JobTitle, job.JobURL, job.Website)
		if err != nil {
			return err
		}
		if err := s.repo.CreateJobApplication(ctx, application); err != nil {
			return err
		}
	}

	// Mark as processing
	application.MarkProcessing()
	return s.repo.UpdateJobApplication(ctx, application)
}

