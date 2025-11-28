package jobapplications

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// Service orchestrates job application workflows.
type Service interface {
	RequestJobApplication(ctx context.Context, userID uuid.UUID, companyName, location, jobTitle, jobURL, website string) (*JobApplication, error)
	GetJobApplication(ctx context.Context, applicationID uuid.UUID) (*JobApplication, error)
	ListJobApplications(ctx context.Context, filters JobApplicationFilters) ([]JobApplication, error)
	UpdateJobApplicationStatus(ctx context.Context, applicationID uuid.UUID, status ApplicationStatus) error
	ProcessJobApplicationJob(ctx context.Context, job *JobApplicationJob) error
}

type service struct {
	repo  Repository
	queue Queue
	logger *slog.Logger
}

// NewService constructs a Service.
func NewService(repo Repository, queue Queue, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		queue:  queue,
		logger: logger,
	}
}

func (s *service) RequestJobApplication(ctx context.Context, userID uuid.UUID, companyName, location, jobTitle, jobURL, website string) (*JobApplication, error) {
	// Create job application record
	application, err := NewJobApplication(userID, companyName, location, jobTitle, jobURL, website)
	if err != nil {
		return nil, err
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

	return application.UpdateStatus(status)
}

// ProcessJobApplicationJob is called by the worker to process a job.
// This is a placeholder - actual processing will be done in the worker with Playwright.
func (s *service) ProcessJobApplicationJob(ctx context.Context, job *JobApplicationJob) error {
	// This method will be called by the worker, but the actual application logic
	// will be in the worker itself using Playwright.
	// This is here for interface consistency.
	
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

