package jobapplications

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	jobapplicationsdomain "github.com/woragis/backend/server/app/internal/domains/jobapplications"
	jobwebsitesdomain "github.com/woragis/backend/server/app/internal/domains/jobwebsites"
	aiservice "github.com/woragis/backend/server/app/internal/services/ai"
	playwrightservice "github.com/woragis/backend/server/app/internal/services/playwright"
	"github.com/woragis/backend/server/app/pkg/validation"
)

// Worker processes job application jobs from Redis queue.
type Worker struct {
	queue            jobapplicationsdomain.Queue
	applicationRepo  jobapplicationsdomain.Repository
	websiteService   jobwebsitesdomain.Service
	orchestrator     *Orchestrator
	scraper          *playwrightservice.Scraper
	coverLetterService *aiservice.CoverLetterService
	db               *gorm.DB
	logger           *slog.Logger
	stopChan         chan struct{}
}

// NewWorker creates a new job application worker.
func NewWorker(
	queue jobapplicationsdomain.Queue,
	applicationRepo jobapplicationsdomain.Repository,
	websiteService jobwebsitesdomain.Service,
	orchestrator *Orchestrator,
	scraper *playwrightservice.Scraper,
	coverLetterService *aiservice.CoverLetterService,
	db *gorm.DB,
	logger *slog.Logger,
) *Worker {
	return &Worker{
		queue:              queue,
		applicationRepo:    applicationRepo,
		websiteService:     websiteService,
		orchestrator:       orchestrator,
		scraper:            scraper,
		coverLetterService: coverLetterService,
		db:                 db,
		logger:              logger,
		stopChan:           make(chan struct{}),
	}
}

// Start begins processing job application jobs.
func (w *Worker) Start(ctx context.Context) {
	w.logger.Info("job application worker started")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("job application worker stopping (context cancelled)")
			return
		case <-w.stopChan:
			w.logger.Info("job application worker stopping (stop signal)")
			return
		default:
			w.processJob(ctx)
		}
	}
}

// Stop signals the worker to stop.
func (w *Worker) Stop() {
	close(w.stopChan)
}

func (w *Worker) processJob(ctx context.Context) {
	// Dequeue job with 5 second timeout
	job, err := w.queue.DequeueJob(ctx, 5*time.Second)
	if err != nil {
		if err.Error() != "jobapplications: job not found" {
			w.logger.Error("failed to dequeue job", slog.Any("error", err))
		}
		return
	}

	if job == nil {
		// No job available, continue polling
		return
	}

	w.logger.Info("processing job application job",
		slog.String("jobId", job.ID),
		slog.String("company", job.CompanyName),
		slog.String("website", job.Website),
	)

	// Validate job data
	if err := ValidateJobApplicationJob(job); err != nil {
		w.logger.Error("invalid job application job",
			slog.String("jobId", job.ID),
			slog.Any("error", err),
		)
		_ = w.queue.MarkJobFailed(ctx, job.ID, fmt.Sprintf("validation failed: %v", err))
		return
	}

	// Check if we should process this website (rate limit check)
	shouldProcess, err := w.orchestrator.ShouldProcessWebsite(ctx, job.Website)
	if err != nil {
		w.logger.Error("failed to check website availability",
			slog.String("website", job.Website),
			slog.Any("error", err),
		)
		// Re-enqueue job for retry
		_ = w.queue.EnqueueJob(ctx, job)
		return
	}

	if !shouldProcess {
		w.logger.Info("website limit reached, re-enqueuing job",
			slog.String("website", job.Website),
		)
		// Re-enqueue for later
		_ = w.queue.EnqueueJob(ctx, job)
		time.Sleep(1 * time.Hour) // Wait before checking again
		return
	}

	// Process the job
	if err := w.processApplication(ctx, job); err != nil {
		w.logger.Error("failed to process job application",
			slog.String("jobId", job.ID),
			slog.Any("error", err),
		)
		// Mark job as failed
		_ = w.queue.MarkJobFailed(ctx, job.ID, err.Error())
		return
	}

	// Increment website count
	if err := w.orchestrator.IncrementWebsiteCount(ctx, job.Website); err != nil {
		w.logger.Warn("failed to increment website count",
			slog.String("website", job.Website),
			slog.Any("error", err),
		)
	}

	// Mark job as complete
	if err := w.queue.MarkJobComplete(ctx, job.ID); err != nil {
		w.logger.Warn("failed to mark job as complete",
			slog.String("jobId", job.ID),
			slog.Any("error", err),
		)
	}

	w.logger.Info("job application job completed",
		slog.String("jobId", job.ID),
	)
}

func (w *Worker) processApplication(ctx context.Context, job *jobapplicationsdomain.JobApplicationJob) error {
	// Job validation already done in processJob, but validate again for safety
	if err := ValidateJobApplicationJob(job); err != nil {
		return fmt.Errorf("job validation failed: %w", err)
	}

	userID, err := uuid.Parse(job.UserID)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	// Find or create application record
	filters := jobapplicationsdomain.JobApplicationFilters{
		UserID: &userID,
	}
	applications, err := w.applicationRepo.ListJobApplications(ctx, filters)
	if err != nil {
		return fmt.Errorf("failed to list applications: %w", err)
	}

	var application *jobapplicationsdomain.JobApplication
	for i := range applications {
		if applications[i].JobURL == job.JobURL && applications[i].Website == job.Website {
			application = &applications[i]
			break
		}
	}

	if application == nil {
		// Create new application
		application, err = jobapplicationsdomain.NewJobApplication(
			userID,
			job.CompanyName,
			job.Location,
			job.JobTitle,
			job.JobURL,
			job.Website,
		)
		if err != nil {
			return fmt.Errorf("failed to create application: %w", err)
		}
		if err := w.applicationRepo.CreateJobApplication(ctx, application); err != nil {
			return fmt.Errorf("failed to save application: %w", err)
		}
	}

	// Mark as processing
	application.MarkProcessing()
	if err := w.applicationRepo.UpdateJobApplication(ctx, application); err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	// Fetch user profile data
	profile, err := w.fetchUserProfile(ctx, userID)
	if err != nil {
		w.logger.Warn("failed to fetch user profile, using empty profile",
			slog.String("userId", userID.String()),
			slog.Any("error", err),
		)
		profile = aiservice.UserProfile{} // Use empty profile
	}

	// Validate fetched profile data
	if err := aiservice.ValidateUserProfile(profile); err != nil {
		w.logger.Warn("fetched user profile validation failed, using empty profile",
			slog.String("userId", userID.String()),
			slog.Any("error", err),
		)
		profile = aiservice.UserProfile{} // Use empty profile on validation failure
	}

	// Generate cover letter
	jobInfo := aiservice.JobInfo{
		CompanyName:    job.CompanyName,
		JobTitle:       job.JobTitle,
		JobDescription: "", // TODO: Fetch from job URL if possible
		Location:       job.Location,
		Requirements:   []string{}, // TODO: Extract from job description
	}

	// Validate job info before generating cover letter
	if err := aiservice.ValidateJobInfo(jobInfo); err != nil {
		return fmt.Errorf("job info validation failed: %w", err)
	}

	coverLetter, err := w.coverLetterService.GenerateCoverLetter(ctx, profile, jobInfo)
	if err != nil {
		return fmt.Errorf("failed to generate cover letter: %w", err)
	}

	// Validate generated cover letter before using
	if err := validation.ValidateString(coverLetter, 100, 10000, "cover_letter"); err != nil {
		return fmt.Errorf("generated cover letter validation failed: %w", err)
	}

	// Apply to job using Playwright
	if err := w.scraper.ApplyToJob(ctx, job, coverLetter); err != nil {
		application.MarkFailed(fmt.Sprintf("Application failed: %v", err))
		w.applicationRepo.UpdateJobApplication(ctx, application)
		return fmt.Errorf("failed to apply to job: %w", err)
	}

	// Mark as applied
	application.MarkApplied(coverLetter)
	if err := w.applicationRepo.UpdateJobApplication(ctx, application); err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	return nil
}

// fetchUserProfile fetches user profile data from database.
func (w *Worker) fetchUserProfile(ctx context.Context, userID uuid.UUID) (aiservice.UserProfile, error) {
	profile := aiservice.UserProfile{}

	// TODO: Fetch projects, posts, technical writings, skills, etc. from database
	// This is a placeholder - actual implementation will query the respective domains

	// Example structure (to be implemented):
	// projects := fetchProjects(ctx, userID)
	// posts := fetchPosts(ctx, userID)
	// technicalWritings := fetchTechnicalWritings(ctx, userID)
	// skills := fetchSkills(ctx, userID)
	// etc.

	return profile, nil
}

