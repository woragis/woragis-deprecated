package socialmediaposts

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates social media post workflows.
type Service interface {
	// Post operations
	CreatePost(ctx context.Context, req CreatePostRequest) (*SocialMediaPost, error)
	UpdatePost(ctx context.Context, req UpdatePostRequest) (*SocialMediaPost, error)
	GetPost(ctx context.Context, postID uuid.UUID) (*SocialMediaPost, error)
	GetPostByURL(ctx context.Context, url string) (*SocialMediaPost, error)
	ListPosts(ctx context.Context, filters PostFilters) ([]SocialMediaPost, error)
	DeletePost(ctx context.Context, postID uuid.UUID) error
	UpdatePostEngagement(ctx context.Context, postID uuid.UUID, req UpdateEngagementRequest) (*SocialMediaPost, error)

	// Link operations
	CreateLink(ctx context.Context, req CreateLinkRequest) (*SocialMediaEntityLink, error)
	UpdateLink(ctx context.Context, req UpdateLinkRequest) (*SocialMediaEntityLink, error)
	DeleteLink(ctx context.Context, linkID uuid.UUID) error
	GetLinksByPost(ctx context.Context, postID uuid.UUID) ([]SocialMediaEntityLink, error)
	GetLinksByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]SocialMediaEntityLink, error)
	GetPostsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID, relationshipType *RelationshipType) ([]SocialMediaPostWithLink, error)
	GetEntitiesByPost(ctx context.Context, postID uuid.UUID) ([]EntityLinkInfo, error)
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

var _ Service = (*service)(nil)

// NewService constructs a Service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// Request payloads

type CreatePostRequest struct {
	URL            string     `json:"url"`
	Platform       Platform   `json:"platform"`
	Title          string     `json:"title,omitempty"`
	ContentPreview string     `json:"contentPreview,omitempty"`
	PublishedDate  *time.Time `json:"publishedDate,omitempty"`
}

type UpdatePostRequest struct {
	PostID         uuid.UUID  `json:"-"`
	Title          string     `json:"title,omitempty"`
	ContentPreview string     `json:"contentPreview,omitempty"`
	PublishedDate  *time.Time `json:"publishedDate,omitempty"`
	Status         *PostStatus `json:"status,omitempty"`
}

type UpdateEngagementRequest struct {
	Likes    *int64 `json:"likes,omitempty"`
	Shares   *int64 `json:"shares,omitempty"`
	Comments *int64 `json:"comments,omitempty"`
	Views    *int64 `json:"views,omitempty"`
}

type CreateLinkRequest struct {
	PostID           uuid.UUID        `json:"postId"`
	EntityType       EntityType       `json:"entityType"`
	EntityID         uuid.UUID        `json:"entityId"`
	RelationshipType RelationshipType `json:"relationshipType"`
}

type UpdateLinkRequest struct {
	LinkID           uuid.UUID        `json:"-"`
	RelationshipType RelationshipType `json:"relationshipType"`
}

// Post operations

func (s *service) CreatePost(ctx context.Context, req CreatePostRequest) (*SocialMediaPost, error) {
	// Check if post with same URL already exists
	existing, _ := s.repo.GetPostByURL(ctx, req.URL)
	if existing != nil {
		return nil, NewDomainError(ErrCodeConflict, ErrPostAlreadyExists)
	}

	post, err := NewSocialMediaPost(req.URL, req.Platform, req.Title, req.ContentPreview, req.PublishedDate)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreatePost(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *service) UpdatePost(ctx context.Context, req UpdatePostRequest) (*SocialMediaPost, error) {
	post, err := s.repo.GetPost(ctx, req.PostID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.ContentPreview != "" {
		post.ContentPreview = req.ContentPreview
	}
	if req.PublishedDate != nil {
		post.PublishedDate = req.PublishedDate
	}
	if req.Status != nil {
		if err := post.UpdateStatus(*req.Status); err != nil {
			return nil, err
		}
	}
	post.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdatePost(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *service) GetPost(ctx context.Context, postID uuid.UUID) (*SocialMediaPost, error) {
	return s.repo.GetPost(ctx, postID)
}

func (s *service) GetPostByURL(ctx context.Context, url string) (*SocialMediaPost, error) {
	return s.repo.GetPostByURL(ctx, url)
}

func (s *service) ListPosts(ctx context.Context, filters PostFilters) ([]SocialMediaPost, error) {
	return s.repo.ListPosts(ctx, filters)
}

func (s *service) DeletePost(ctx context.Context, postID uuid.UUID) error {
	return s.repo.DeletePost(ctx, postID)
}

func (s *service) UpdatePostEngagement(ctx context.Context, postID uuid.UUID, req UpdateEngagementRequest) (*SocialMediaPost, error) {
	post, err := s.repo.GetPost(ctx, postID)
	if err != nil {
		return nil, err
	}

	post.UpdateEngagement(req.Likes, req.Shares, req.Comments, req.Views)

	if err := s.repo.UpdatePost(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

// Link operations

func (s *service) CreateLink(ctx context.Context, req CreateLinkRequest) (*SocialMediaEntityLink, error) {
	// Verify post exists
	if _, err := s.repo.GetPost(ctx, req.PostID); err != nil {
		return nil, err
	}

	link, err := NewSocialMediaEntityLink(req.PostID, req.EntityType, req.EntityID, req.RelationshipType)
	if err != nil {
		return nil, err
	}

	if err := s.repo.CreateLink(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *service) UpdateLink(ctx context.Context, req UpdateLinkRequest) (*SocialMediaEntityLink, error) {
	link, err := s.repo.GetLink(ctx, req.LinkID)
	if err != nil {
		return nil, err
	}

	link.RelationshipType = req.RelationshipType
	link.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdateLink(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

func (s *service) DeleteLink(ctx context.Context, linkID uuid.UUID) error {
	return s.repo.DeleteLink(ctx, linkID)
}

func (s *service) GetLinksByPost(ctx context.Context, postID uuid.UUID) ([]SocialMediaEntityLink, error) {
	return s.repo.GetLinksByPost(ctx, postID)
}

func (s *service) GetLinksByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]SocialMediaEntityLink, error) {
	return s.repo.GetLinksByEntity(ctx, entityType, entityID)
}

func (s *service) GetPostsByEntity(ctx context.Context, entityType EntityType, entityID uuid.UUID, relationshipType *RelationshipType) ([]SocialMediaPostWithLink, error) {
	return s.repo.GetPostsByEntity(ctx, entityType, entityID, relationshipType)
}

func (s *service) GetEntitiesByPost(ctx context.Context, postID uuid.UUID) ([]EntityLinkInfo, error) {
	return s.repo.GetEntitiesByPost(ctx, postID)
}

