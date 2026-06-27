package content

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SocialMediaPostRepository is an interface to avoid import cycle with parent package.
type SocialMediaPostRepository interface {
	GetPost(ctx context.Context, postID uuid.UUID) (*SocialMediaPost, error)
	CreatePost(ctx context.Context, post *SocialMediaPost) error
}

// SocialMediaPostService is an interface to avoid import cycle with parent package.
type SocialMediaPostService interface {
	CreatePost(ctx context.Context, req CreateSocialMediaPostRequest) (*SocialMediaPost, error)
	GetPost(ctx context.Context, postID uuid.UUID) (*SocialMediaPost, error)
}

// SocialMediaPost represents a social media post (duplicated from parent to avoid import cycle).
type SocialMediaPost struct {
	ID            uuid.UUID
	Platform      Platform
	Format        ContentFormat
	Title         string
	Content       string
	Status        string
	ContentPostID *uuid.UUID
}

// CreateSocialMediaPostRequest is a request to create a social media post (duplicated from parent to avoid import cycle).
type CreateSocialMediaPostRequest struct {
	Platform      Platform
	Format        ContentFormat
	Title         string
	Content       string
	ContentPostID *uuid.UUID
}

// Service orchestrates content post repurposing workflows.
type Service interface {
	CreateContentPostFromBackend(ctx context.Context, req CreateContentPostRequest) (*ContentPost, error)
	GetContentPost(ctx context.Context, contentPostID uuid.UUID) (*ContentPostWithSocialPosts, error)
	ListContentPosts(ctx context.Context, filters ContentPostFilters) ([]ContentPost, error)
	GetContentBacklog(ctx context.Context) ([]ContentPost, error)
	UpdateContentPostPriority(ctx context.Context, contentPostID uuid.UUID, priority ContentPostPriority) (*ContentPost, error)
	RepurposeToPlatforms(ctx context.Context, contentPostID uuid.UUID, req RepurposeRequest) ([]*SocialMediaPost, error)
	GetRepurposingHistory(ctx context.Context, contentPostID uuid.UUID) ([]RepurposingHistoryItem, error)
}

type service struct {
	repo              Repository
	socialPostsRepo   SocialMediaPostRepository
	socialPostsService SocialMediaPostService
	logger            *slog.Logger
}

var _ Service = (*service)(nil)

// NewService constructs a Service.
func NewService(repo Repository, socialPostsRepo SocialMediaPostRepository, socialPostsService SocialMediaPostService, logger *slog.Logger) Service {
	return &service{
		repo:              repo,
		socialPostsRepo:   socialPostsRepo,
		socialPostsService: socialPostsService,
		logger:            logger,
	}
}

// Request payloads

type CreateContentPostRequest struct {
	PostID      uuid.UUID `json:"postId"`
	ContentType string    `json:"contentType,omitempty"`
	Project     *string   `json:"project,omitempty"`
	Priority    ContentPostPriority `json:"priority,omitempty"`
}

type RepurposeRequest struct {
	Platforms []RepurposePlatform `json:"platforms"`
}

type RepurposePlatform struct {
	Platform Platform      `json:"platform"`
	Format   ContentFormat `json:"format"`
	Title    string        `json:"title"`
	Content  string        `json:"content"`
}

type ContentPostWithSocialPosts struct {
	ContentPost
	SocialPosts []SocialMediaPost `json:"socialPosts"`
}

type RepurposingHistoryItem struct {
	SocialPost  SocialMediaPost `json:"socialPost"`
	RepurposedAt time.Time      `json:"repurposedAt"`
}

// CreateContentPostFromBackend creates a content post from an existing backend post.
func (s *service) CreateContentPostFromBackend(ctx context.Context, req CreateContentPostRequest) (*ContentPost, error) {
	// Check if content post already exists for this post
	existing, _ := s.repo.GetContentPostByPostID(ctx, req.PostID)
	if existing != nil {
		return nil, NewDomainError(ErrCodeConflict, ErrContentPostExists)
	}

	priority := req.Priority
	if priority == "" {
		priority = ContentPostPriorityMedium
	}

	post, err := NewContentPost(req.PostID, req.ContentType)
	if err != nil {
		return nil, err
	}

	if err := post.UpdatePriority(priority); err != nil {
		return nil, err
	}

	if req.Project != nil {
		post.SetProject(*req.Project)
	}

	if err := s.repo.CreateContentPost(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

// GetContentPost retrieves a content post with its associated social media posts.
func (s *service) GetContentPost(ctx context.Context, contentPostID uuid.UUID) (*ContentPostWithSocialPosts, error) {
	contentPost, err := s.repo.GetContentPost(ctx, contentPostID)
	if err != nil {
		return nil, err
	}

	// Get repurposing relationships
	repurposings, err := s.repo.GetRepurposingsByContentPost(ctx, contentPostID)
	if err != nil {
		return nil, err
	}

	// Fetch social posts
	socialPosts := make([]SocialMediaPost, 0, len(repurposings))
	for _, rep := range repurposings {
		socialPost, err := s.socialPostsRepo.GetPost(ctx, rep.SocialPostID)
		if err != nil {
			s.logger.Warn("Failed to fetch social post", "socialPostID", rep.SocialPostID, "error", err)
			continue
		}
		socialPosts = append(socialPosts, *socialPost)
	}

	return &ContentPostWithSocialPosts{
		ContentPost: *contentPost,
		SocialPosts: socialPosts,
	}, nil
}

// ListContentPosts lists content posts with optional filters.
func (s *service) ListContentPosts(ctx context.Context, filters ContentPostFilters) ([]ContentPost, error) {
	return s.repo.ListContentPosts(ctx, filters)
}

// GetContentBacklog retrieves content posts ready for repurposing.
func (s *service) GetContentBacklog(ctx context.Context) ([]ContentPost, error) {
	return s.repo.GetContentBacklog(ctx)
}

// UpdateContentPostPriority updates the priority of a content post.
func (s *service) UpdateContentPostPriority(ctx context.Context, contentPostID uuid.UUID, priority ContentPostPriority) (*ContentPost, error) {
	post, err := s.repo.GetContentPost(ctx, contentPostID)
	if err != nil {
		return nil, err
	}

	if err := post.UpdatePriority(priority); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateContentPost(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

// RepurposeToPlatforms creates social media posts for multiple platforms from a content post.
func (s *service) RepurposeToPlatforms(ctx context.Context, contentPostID uuid.UUID, req RepurposeRequest) ([]*SocialMediaPost, error) {
	contentPost, err := s.repo.GetContentPost(ctx, contentPostID)
	if err != nil {
		return nil, err
	}

	// Update status to in progress
	if contentPost.Status == ContentPostStatusPending {
		if err := contentPost.UpdateStatus(ContentPostStatusInProgress); err != nil {
			return nil, err
		}
		if err := s.repo.UpdateContentPost(ctx, contentPost); err != nil {
			return nil, err
		}
	}

	createdPosts := make([]*SocialMediaPost, 0, len(req.Platforms))
	contentPostIDPtr := &contentPostID

	for _, platform := range req.Platforms {
		// Create social media post
		socialPost, err := s.socialPostsService.CreatePost(ctx, CreateSocialMediaPostRequest{
			Platform:      platform.Platform,
			Format:        platform.Format,
			Title:         platform.Title,
			Content:       platform.Content,
			ContentPostID: contentPostIDPtr,
		})
		if err != nil {
			s.logger.Error("Failed to create social post", "platform", platform.Platform, "error", err)
			continue
		}

		// Create repurposing relationship
		repurposing, err := NewContentRepurposing(contentPostID, socialPost.ID)
		if err != nil {
			s.logger.Error("Failed to create repurposing", "error", err)
			continue
		}

		if err := s.repo.CreateRepurposing(ctx, repurposing); err != nil {
			s.logger.Error("Failed to persist repurposing", "error", err)
			continue
		}

		createdPosts = append(createdPosts, socialPost)
	}

	// Update status to completed if all posts were created successfully
	if len(createdPosts) == len(req.Platforms) {
		if err := contentPost.UpdateStatus(ContentPostStatusCompleted); err != nil {
			s.logger.Warn("Failed to update status to completed", "error", err)
		} else {
			if err := s.repo.UpdateContentPost(ctx, contentPost); err != nil {
				s.logger.Warn("Failed to persist status update", "error", err)
			}
		}
	}

	return createdPosts, nil
}

// GetRepurposingHistory retrieves the repurposing history for a content post.
func (s *service) GetRepurposingHistory(ctx context.Context, contentPostID uuid.UUID) ([]RepurposingHistoryItem, error) {
	repurposings, err := s.repo.GetRepurposingsByContentPost(ctx, contentPostID)
	if err != nil {
		return nil, err
	}

	history := make([]RepurposingHistoryItem, 0, len(repurposings))
	for _, rep := range repurposings {
		socialPost, err := s.socialPostsRepo.GetPost(ctx, rep.SocialPostID)
		if err != nil {
			s.logger.Warn("Failed to fetch social post", "socialPostID", rep.SocialPostID, "error", err)
			continue
		}

		history = append(history, RepurposingHistoryItem{
			SocialPost:  *socialPost,
			RepurposedAt: rep.RepurposedAt,
		})
	}

	return history, nil
}

// AdaptContentForPlatform adapts content for a specific platform and format.
// This is a placeholder - in a real implementation, this would use AI/ML to adapt content.
func AdaptContentForPlatform(originalContent string, platform Platform, format ContentFormat) (title, content string) {
	// Simple adaptation logic - in production, this would be more sophisticated
	title = extractTitle(originalContent)
	content = adaptContentLength(originalContent, platform, format)
	return title, content
}

// Helper functions for content adaptation

func extractTitle(content string) string {
	// Extract first line or first sentence as title
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) != "" {
		title := strings.TrimSpace(lines[0])
		// Limit title length
		if len(title) > 100 {
			title = title[:97] + "..."
		}
		return title
	}

	// Fallback: use first sentence
	sentences := strings.Split(content, ".")
	if len(sentences) > 0 {
		title := strings.TrimSpace(sentences[0])
		if len(title) > 100 {
			title = title[:97] + "..."
		}
		return title
	}

	return "Untitled"
}

func adaptContentLength(content string, platform Platform, format ContentFormat) string {
	// Simple truncation based on platform/format
	maxLength := getMaxLengthForPlatform(platform, format)
	if len(content) <= maxLength {
		return content
	}

	// Truncate at word boundary
	truncated := content[:maxLength]
	lastSpace := strings.LastIndex(truncated, " ")
	if lastSpace > 0 {
		truncated = truncated[:lastSpace]
	}
	return truncated + "..."
}

func getMaxLengthForPlatform(platform Platform, format ContentFormat) int {
	// Character limits by platform and format
	switch platform {
	case PlatformTwitter:
		if format == FormatThread {
			return 280 * 10 // Thread can have multiple tweets
		}
		return 280
	case PlatformLinkedIn:
		if format == FormatLongForm {
			return 3000
		}
		return 1300
	case PlatformInstagram:
		return 2200
	case PlatformMedium:
		return 10000
	case PlatformSubstack:
		return 50000
	default:
		return 2000
	}
}
