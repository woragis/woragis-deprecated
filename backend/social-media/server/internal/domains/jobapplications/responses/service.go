package responses

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates response workflows.
type Service interface {
	CreateResponse(ctx context.Context, jobApplicationID uuid.UUID, responseType ResponseType, responseDate time.Time) (*Response, error)
	GetResponse(ctx context.Context, responseID uuid.UUID) (*Response, error)
	ListResponses(ctx context.Context, filters ResponseFilters) ([]Response, error)
	UpdateResponse(ctx context.Context, responseID uuid.UUID, updates UpdateResponseRequest) (*Response, error)
	DeleteResponse(ctx context.Context, responseID uuid.UUID) error
	GetResponsesByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]Response, error)
}

// UpdateResponseRequest represents fields that can be updated on a response.
type UpdateResponseRequest struct {
	Message         *string
	ContactPerson   *string
	ContactEmail    *string
	ContactPhone    *string
	ResponseChannel *string
}

// ResumeMetricsService is an interface to avoid circular dependencies
type ResumeMetricsService interface {
	RecalculateResumeMetrics(ctx context.Context, resumeID uuid.UUID) error
}

// JobApplicationService is an interface to get job application details
type JobApplicationService interface {
	GetJobApplication(ctx context.Context, applicationID uuid.UUID) (*JobApplication, error)
}

// JobApplication represents a job application (minimal interface)
type JobApplication struct {
	ID       uuid.UUID
	ResumeID *uuid.UUID
}

type service struct {
	repo                 Repository
	jobApplicationService JobApplicationService // Optional: to get resumeId
	resumeMetricsService ResumeMetricsService   // Optional: for updating resume metrics
	logger               *slog.Logger
}

// NewService constructs a Service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// NewServiceWithDependencies constructs a Service with all dependencies.
func NewServiceWithDependencies(repo Repository, jobApplicationService JobApplicationService, resumeMetricsService ResumeMetricsService, logger *slog.Logger) Service {
	return &service{
		repo:                 repo,
		jobApplicationService: jobApplicationService,
		resumeMetricsService: resumeMetricsService,
		logger:               logger,
	}
}

func (s *service) CreateResponse(ctx context.Context, jobApplicationID uuid.UUID, responseType ResponseType, responseDate time.Time) (*Response, error) {
	response, err := NewResponse(jobApplicationID, responseType, responseDate)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateResponse(ctx, response); err != nil {
		return nil, err
	}

	// Recalculate resume metrics when an offer response is created
	if s.resumeMetricsService != nil && s.jobApplicationService != nil && responseType == ResponseTypeOffer {
		application, err := s.jobApplicationService.GetJobApplication(ctx, jobApplicationID)
		if err == nil && application != nil && application.ResumeID != nil {
			if err := s.resumeMetricsService.RecalculateResumeMetrics(ctx, *application.ResumeID); err != nil {
				s.logger.Warn("failed to recalculate resume metrics after creating offer response", "resume_id", application.ResumeID.String(), "error", err)
				// Don't fail the request if metric recalculation fails
			}
		}
	}

	return response, nil
}

func (s *service) GetResponse(ctx context.Context, responseID uuid.UUID) (*Response, error) {
	return s.repo.GetResponse(ctx, responseID)
}

func (s *service) ListResponses(ctx context.Context, filters ResponseFilters) ([]Response, error) {
	return s.repo.ListResponses(ctx, filters)
}

func (s *service) UpdateResponse(ctx context.Context, responseID uuid.UUID, updates UpdateResponseRequest) (*Response, error) {
	response, err := s.repo.GetResponse(ctx, responseID)
	if err != nil {
		return nil, err
	}

	if updates.Message != nil {
		response.UpdateMessage(*updates.Message)
	}

	if updates.ContactPerson != nil || updates.ContactEmail != nil || updates.ContactPhone != nil {
		person := ""
		email := ""
		phone := ""
		if updates.ContactPerson != nil {
			person = *updates.ContactPerson
		} else {
			person = response.ContactPerson
		}
		if updates.ContactEmail != nil {
			email = *updates.ContactEmail
		} else {
			email = response.ContactEmail
		}
		if updates.ContactPhone != nil {
			phone = *updates.ContactPhone
		} else {
			phone = response.ContactPhone
		}
		response.UpdateContactInfo(person, email, phone)
	}

	if updates.ResponseChannel != nil {
		response.UpdateResponseChannel(*updates.ResponseChannel)
	}

	if err := s.repo.UpdateResponse(ctx, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (s *service) DeleteResponse(ctx context.Context, responseID uuid.UUID) error {
	return s.repo.DeleteResponse(ctx, responseID)
}

func (s *service) GetResponsesByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]Response, error) {
	return s.repo.GetResponsesByApplicationID(ctx, applicationID)
}

