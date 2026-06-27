package resumes

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service orchestrates resume workflows.
type Service interface {
	CreateResume(ctx context.Context, userID uuid.UUID, title, filePath, fileName string, fileSize int64, tags JSONArray) (*Resume, error)
	UpdateResume(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID, title string, tags JSONArray) (*Resume, error)
	DeleteResume(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) error
	GetResume(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error)
	ListResumes(ctx context.Context, userID uuid.UUID) ([]Resume, error)
	ListResumesByTags(ctx context.Context, userID uuid.UUID, tags []string) ([]Resume, error)
	MarkAsMain(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error)
	MarkAsFeatured(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error)
	UnmarkAsMain(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error)
	UnmarkAsFeatured(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error)
	GetMainResume(ctx context.Context, userID uuid.UUID) (*Resume, error)
	GetFeaturedResume(ctx context.Context, userID uuid.UUID) (*Resume, error)
	GetBestResume(ctx context.Context, userID uuid.UUID) (*Resume, error) // Returns main > featured > most recent
	RecalculateResumeMetrics(ctx context.Context, resumeID uuid.UUID) error
}

// service implements Service.
type service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new resume service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// CreateResume creates a new resume.
func (s *service) CreateResume(ctx context.Context, userID uuid.UUID, title, filePath, fileName string, fileSize int64, tags JSONArray) (*Resume, error) {
	resume, err := NewResume(userID, title, filePath, fileName, fileSize, tags)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateResume(ctx, resume); err != nil {
		return nil, err
	}

	return resume, nil
}

// UpdateResume updates an existing resume.
func (s *service) UpdateResume(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID, title string, tags JSONArray) (*Resume, error) {
	resume, err := s.repo.GetResume(ctx, resumeID, userID)
	if err != nil {
		return nil, err
	}

	if title != "" {
		if err := resume.UpdateTitle(title); err != nil {
			return nil, err
		}
	}

	if tags != nil {
		if err := resume.UpdateTags(tags); err != nil {
			return nil, err
		}
	}

	if err := s.repo.UpdateResume(ctx, resume); err != nil {
		return nil, err
	}

	return resume, nil
}

// DeleteResume deletes a resume.
func (s *service) DeleteResume(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) error {
	return s.repo.DeleteResume(ctx, resumeID, userID)
}

// GetResume retrieves a resume by ID.
func (s *service) GetResume(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error) {
	return s.repo.GetResume(ctx, resumeID, userID)
}

// ListResumes lists all resumes for a user.
func (s *service) ListResumes(ctx context.Context, userID uuid.UUID) ([]Resume, error) {
	return s.repo.ListResumes(ctx, userID)
}

// ListResumesByTags lists resumes filtered by tags.
func (s *service) ListResumesByTags(ctx context.Context, userID uuid.UUID, tags []string) ([]Resume, error) {
	return s.repo.ListResumesByTags(ctx, userID, tags)
}

// MarkAsMain marks a resume as main and unmarks others.
func (s *service) MarkAsMain(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error) {
	// Unmark all other resumes as main
	if err := s.repo.UnmarkAllAsMain(ctx, userID); err != nil {
		return nil, err
	}

	// Mark this resume as main
	resume, err := s.repo.GetResume(ctx, resumeID, userID)
	if err != nil {
		return nil, err
	}

	resume.MarkAsMain()
	if err := s.repo.UpdateResume(ctx, resume); err != nil {
		return nil, err
	}

	return resume, nil
}

// MarkAsFeatured marks a resume as featured.
func (s *service) MarkAsFeatured(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error) {
	resume, err := s.repo.GetResume(ctx, resumeID, userID)
	if err != nil {
		return nil, err
	}

	resume.MarkAsFeatured()
	if err := s.repo.UpdateResume(ctx, resume); err != nil {
		return nil, err
	}

	return resume, nil
}

// UnmarkAsMain removes the main flag from a resume.
func (s *service) UnmarkAsMain(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error) {
	resume, err := s.repo.GetResume(ctx, resumeID, userID)
	if err != nil {
		return nil, err
	}

	resume.UnmarkAsMain()
	if err := s.repo.UpdateResume(ctx, resume); err != nil {
		return nil, err
	}

	return resume, nil
}

// UnmarkAsFeatured removes the featured flag from a resume.
func (s *service) UnmarkAsFeatured(ctx context.Context, userID uuid.UUID, resumeID uuid.UUID) (*Resume, error) {
	resume, err := s.repo.GetResume(ctx, resumeID, userID)
	if err != nil {
		return nil, err
	}

	resume.UnmarkAsFeatured()
	if err := s.repo.UpdateResume(ctx, resume); err != nil {
		return nil, err
	}

	return resume, nil
}

// GetMainResume retrieves the main resume.
func (s *service) GetMainResume(ctx context.Context, userID uuid.UUID) (*Resume, error) {
	return s.repo.GetMainResume(ctx, userID)
}

// GetFeaturedResume retrieves a featured resume.
func (s *service) GetFeaturedResume(ctx context.Context, userID uuid.UUID) (*Resume, error) {
	return s.repo.GetFeaturedResume(ctx, userID)
}

// GetBestResume returns the best resume using priority: main > featured > most recent.
func (s *service) GetBestResume(ctx context.Context, userID uuid.UUID) (*Resume, error) {
	// Try main first
	resume, err := s.repo.GetMainResume(ctx, userID)
	if err == nil {
		return resume, nil
	}

	// Try featured
	resume, err = s.repo.GetFeaturedResume(ctx, userID)
	if err == nil {
		return resume, nil
	}

	// Fallback to most recent
	resumes, err := s.repo.ListResumes(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(resumes) == 0 {
		return nil, NewDomainError(ErrCodeNotFound, ErrResumeNotFound)
	}

	// Return the most recent (first in list since it's ordered by created_at DESC)
	return &resumes[0], nil
}

// RecalculateResumeMetrics recalculates and updates metrics for a resume.
func (s *service) RecalculateResumeMetrics(ctx context.Context, resumeID uuid.UUID) error {
	metrics, err := s.repo.CalculateResumeMetrics(ctx, resumeID)
	if err != nil {
		return err
	}
	
	return s.repo.UpdateResumeMetrics(ctx, resumeID, metrics)
}

