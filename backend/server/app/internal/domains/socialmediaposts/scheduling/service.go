package scheduling

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"

	platformsdomain "github.com/woragis/backend/server/app/internal/domains/socialmediaposts/platforms"
)

// Platform represents the social media platform (duplicated from parent to avoid import cycle).
type Platform string

const (
	PlatformLinkedIn  Platform = "linkedin"
	PlatformTwitter   Platform = "twitter"
	PlatformInstagram Platform = "instagram"
	PlatformMedium    Platform = "medium"
	PlatformSubstack  Platform = "substack"
	PlatformValete    Platform = "valete"
	PlatformWebsite   Platform = "website"
)

// SocialMediaPostRepository is an interface to avoid import cycle with parent package.
type SocialMediaPostRepository interface {
	GetPost(ctx context.Context, postID uuid.UUID) (*SocialMediaPost, error)
	CreatePost(ctx context.Context, post *SocialMediaPost) error
	UpdatePost(ctx context.Context, post *SocialMediaPost) error
}

// SocialMediaPostService is an interface to avoid import cycle with parent package.
type SocialMediaPostService interface {
	GetPost(ctx context.Context, postID uuid.UUID) (*SocialMediaPost, error)
	UpdatePost(ctx context.Context, req UpdateSocialMediaPostRequest) (*SocialMediaPost, error)
}

// SocialMediaPost represents a social media post (duplicated from parent to avoid import cycle).
type SocialMediaPost struct {
	ID            uuid.UUID
	Platform      Platform
	Format        string
	Title         string
	Content       string
	Status        string
	ContentPostID *uuid.UUID
	ScheduledAt   *time.Time
}

// UpdateSocialMediaPostRequest is a request to update a social media post (duplicated from parent to avoid import cycle).
type UpdateSocialMediaPostRequest struct {
	PostID  uuid.UUID
	Title   *string
	Content *string
	Status  *string
}

// Service orchestrates post scheduling workflows.
type Service interface {
	SchedulePost(ctx context.Context, req SchedulePostRequest) (*ScheduledPost, error)
	GetSchedule(ctx context.Context, scheduleID uuid.UUID) (*ScheduledPost, error)
	GetScheduleForDateRange(ctx context.Context, startDate, endDate time.Time) ([]ScheduledPost, error)
	GetUpcomingPosts(ctx context.Context, limit int) ([]ScheduledPost, error)
	UpdateSchedule(ctx context.Context, req UpdateScheduleRequest) (*ScheduledPost, error)
	CancelSchedule(ctx context.Context, scheduleID uuid.UUID) error
	AutoSchedule(ctx context.Context, socialPostID uuid.UUID, platform Platform) (*ScheduledPost, error)
	CheckConflicts(ctx context.Context, scheduledAt time.Time, excludeScheduleID *uuid.UUID) (bool, error)
}

type service struct {
	repo              Repository
	socialPostsRepo   SocialMediaPostRepository
	socialPostsService SocialMediaPostService
	platformsService  platformsdomain.Service
	logger            *slog.Logger
}

var _ Service = (*service)(nil)

// NewService constructs a Service.
func NewService(
	repo Repository,
	socialPostsRepo SocialMediaPostRepository,
	socialPostsService SocialMediaPostService,
	platformsService platformsdomain.Service,
	logger *slog.Logger,
) Service {
	return &service{
		repo:              repo,
		socialPostsRepo:   socialPostsRepo,
		socialPostsService: socialPostsService,
		platformsService:  platformsService,
		logger:            logger,
	}
}

// Request payloads

type SchedulePostRequest struct {
	SocialPostID uuid.UUID `json:"socialPostId"`
	ScheduledAt  time.Time `json:"scheduledAt"`
	PlatformID   *uuid.UUID `json:"platformId,omitempty"`
}

type UpdateScheduleRequest struct {
	ScheduleID  uuid.UUID  `json:"-"`
	ScheduledAt *time.Time `json:"scheduledAt,omitempty"`
	Status      *ScheduledPostStatus `json:"status,omitempty"`
}

// SchedulePost schedules a social media post.
func (s *service) SchedulePost(ctx context.Context, req SchedulePostRequest) (*ScheduledPost, error) {
	// Verify social post exists
	_, err := s.socialPostsRepo.GetPost(ctx, req.SocialPostID)
	if err != nil {
		return nil, err
	}

	// Check for conflicts
	hasConflict, err := s.repo.CheckConflict(ctx, req.ScheduledAt, nil)
	if err != nil {
		return nil, err
	}
	if hasConflict {
		return nil, NewDomainError(ErrCodeConflict, ErrScheduleConflict)
	}

	// Create schedule
	schedule, err := NewScheduledPost(req.SocialPostID, req.ScheduledAt)
	if err != nil {
		return nil, err
	}

	if req.PlatformID != nil {
		schedule.SetPlatformID(*req.PlatformID)
	}

	if err := schedule.UpdateStatus(ScheduledPostStatusScheduled); err != nil {
		return nil, err
	}

	if err := s.repo.CreateSchedule(ctx, schedule); err != nil {
		return nil, err
	}

	// Update social post scheduled time (if needed, update via service)
	// Note: This would require calling socialPostsService.UpdatePost, but we skip it for now
	// to avoid circular dependency. The scheduled time is tracked in ScheduledPost entity.

	return schedule, nil
}

// GetSchedule retrieves a scheduled post.
func (s *service) GetSchedule(ctx context.Context, scheduleID uuid.UUID) (*ScheduledPost, error) {
	return s.repo.GetSchedule(ctx, scheduleID)
}

// GetScheduleForDateRange retrieves scheduled posts for a date range.
func (s *service) GetScheduleForDateRange(ctx context.Context, startDate, endDate time.Time) ([]ScheduledPost, error) {
	return s.repo.GetScheduleForDateRange(ctx, startDate, endDate)
}

// GetUpcomingPosts retrieves upcoming scheduled posts.
func (s *service) GetUpcomingPosts(ctx context.Context, limit int) ([]ScheduledPost, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.repo.GetUpcomingPosts(ctx, limit)
}

// UpdateSchedule updates a scheduled post.
func (s *service) UpdateSchedule(ctx context.Context, req UpdateScheduleRequest) (*ScheduledPost, error) {
	schedule, err := s.repo.GetSchedule(ctx, req.ScheduleID)
	if err != nil {
		return nil, err
	}

	if req.ScheduledAt != nil {
		// Check for conflicts (excluding current schedule)
		hasConflict, err := s.repo.CheckConflict(ctx, *req.ScheduledAt, &req.ScheduleID)
		if err != nil {
			return nil, err
		}
		if hasConflict {
			return nil, NewDomainError(ErrCodeConflict, ErrScheduleConflict)
		}

		schedule.ScheduledAt = req.ScheduledAt.UTC()
		schedule.ScheduledDate = req.ScheduledAt.Truncate(24 * time.Hour)
		schedule.ScheduledTime = time.Date(0, 1, 1, req.ScheduledAt.Hour(), req.ScheduledAt.Minute(), req.ScheduledAt.Second(), 0, time.UTC)
		schedule.UpdatedAt = time.Now().UTC()
	}

	if req.Status != nil {
		if err := schedule.UpdateStatus(*req.Status); err != nil {
			return nil, err
		}
	}

	if err := s.repo.UpdateSchedule(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

// CancelSchedule cancels a scheduled post.
func (s *service) CancelSchedule(ctx context.Context, scheduleID uuid.UUID) error {
	schedule, err := s.repo.GetSchedule(ctx, scheduleID)
	if err != nil {
		return err
	}

	if err := schedule.Cancel(); err != nil {
		return err
	}

	return s.repo.UpdateSchedule(ctx, schedule)
}

// AutoSchedule automatically schedules a post at an optimal time based on platform config.
func (s *service) AutoSchedule(ctx context.Context, socialPostID uuid.UUID, platform Platform) (*ScheduledPost, error) {
	// Get platform config for optimal times (convert Platform to platforms.Platform)
	platformConfig, err := s.platformsService.GetConfigByName(ctx, platformsdomain.Platform(platform))
	if err != nil {
		s.logger.Warn("Failed to get platform config, using default scheduling", "platform", platform, "error", err)
		// Fallback to scheduling 24 hours from now
		scheduledAt := time.Now().UTC().Add(24 * time.Hour)
		return s.SchedulePost(ctx, SchedulePostRequest{
			SocialPostID: socialPostID,
			ScheduledAt:  scheduledAt,
		})
	}

	// Find next optimal time slot
	optimalTime := s.findNextOptimalTime(ctx, platformConfig)
	if optimalTime == nil {
		// Fallback to 24 hours from now
		optimalTime = timePtr(time.Now().UTC().Add(24 * time.Hour))
	}

	return s.SchedulePost(ctx, SchedulePostRequest{
		SocialPostID: socialPostID,
		ScheduledAt:  *optimalTime,
		PlatformID:   &platformConfig.ID,
	})
}

// CheckConflicts checks if a time slot has conflicts.
func (s *service) CheckConflicts(ctx context.Context, scheduledAt time.Time, excludeScheduleID *uuid.UUID) (bool, error) {
	return s.repo.CheckConflict(ctx, scheduledAt, excludeScheduleID)
}

// Helper functions

func (s *service) findNextOptimalTime(ctx context.Context, platformConfig *platformsdomain.PlatformConfig) *time.Time {
	// Simple implementation: find next available slot in the next 7 days
	now := time.Now().UTC()
	for i := 0; i < 7; i++ {
		candidate := now.AddDate(0, 0, i)
		// Try morning (9 AM), afternoon (2 PM), and evening (6 PM)
		for _, hour := range []int{9, 14, 18} {
			candidateTime := time.Date(candidate.Year(), candidate.Month(), candidate.Day(), hour, 0, 0, 0, time.UTC)
			if candidateTime.Before(now) {
				continue
			}

			hasConflict, err := s.repo.CheckConflict(ctx, candidateTime, nil)
			if err != nil {
				s.logger.Warn("Error checking conflict", "error", err)
				continue
			}

			if !hasConflict {
				return &candidateTime
			}
		}
	}

	return nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}
