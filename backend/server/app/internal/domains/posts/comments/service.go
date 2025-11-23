package comments

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service orchestrates comment workflows.
type Service interface {
	CreateComment(ctx context.Context, req CreateCommentRequest) (*Comment, error)
	UpdateComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID, req UpdateCommentRequest) (*Comment, error)
	GetComment(ctx context.Context, commentID uuid.UUID) (*Comment, error)
	DeleteComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) error
	ListComments(ctx context.Context, filters CommentFilters) ([]Comment, error)
	ApproveComment(ctx context.Context, commentID uuid.UUID) error
	RejectComment(ctx context.Context, commentID uuid.UUID) error
	MarkCommentAsSpam(ctx context.Context, commentID uuid.UUID) error
	GetCommentCount(ctx context.Context, postID uuid.UUID, status *CommentStatus) (int64, error)
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

type CreateCommentRequest struct {
	PostID      uuid.UUID  `json:"postId"`
	Content     string     `json:"content"`
	AuthorName  string     `json:"authorName"`
	AuthorEmail string     `json:"authorEmail,omitempty"`
	AuthorURL   string     `json:"authorUrl,omitempty"`
	ParentID    *uuid.UUID `json:"parentId,omitempty"` // For reply comments
	UserID      *uuid.UUID `json:"-"` // Set from auth context if available
	IPAddress   string     `json:"-"` // Set from request
	UserAgent   string     `json:"-"` // Set from request
}

type UpdateCommentRequest struct {
	Content *string `json:"content,omitempty"`
}

// Comment operations

func (s *service) CreateComment(ctx context.Context, req CreateCommentRequest) (*Comment, error) {
	comment, err := NewComment(req.PostID, req.Content, req.AuthorName, req.UserID)
	if err != nil {
		return nil, err
	}

	if req.AuthorEmail != "" {
		comment.AuthorEmail = req.AuthorEmail
	}
	if req.AuthorURL != "" {
		comment.AuthorURL = req.AuthorURL
	}
	if req.ParentID != nil {
		comment.SetReplyTo(*req.ParentID)
	}
	if req.IPAddress != "" {
		comment.IPAddress = req.IPAddress
	}
	if req.UserAgent != "" {
		comment.UserAgent = req.UserAgent
	}

	// Auto-approve if user is authenticated
	if req.UserID != nil {
		comment.Approve()
	}

	if err := s.repo.CreateComment(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *service) UpdateComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID, req UpdateCommentRequest) (*Comment, error) {
	comment, err := s.repo.GetComment(ctx, commentID)
	if err != nil {
		return nil, err
	}

	// Verify ownership (only if comment has a user_id)
	if comment.UserID != nil && *comment.UserID != userID {
		return nil, NewDomainError(ErrCodeUnauthorized, ErrUnauthorized)
	}

	if req.Content != nil {
		if err := comment.UpdateContent(*req.Content); err != nil {
			return nil, err
		}
	}

	if err := s.repo.UpdateComment(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *service) GetComment(ctx context.Context, commentID uuid.UUID) (*Comment, error) {
	return s.repo.GetComment(ctx, commentID)
}

func (s *service) DeleteComment(ctx context.Context, userID uuid.UUID, commentID uuid.UUID) error {
	return s.repo.DeleteComment(ctx, commentID, userID)
}

func (s *service) ListComments(ctx context.Context, filters CommentFilters) ([]Comment, error) {
	return s.repo.ListComments(ctx, filters)
}

func (s *service) ApproveComment(ctx context.Context, commentID uuid.UUID) error {
	comment, err := s.repo.GetComment(ctx, commentID)
	if err != nil {
		return err
	}

	comment.Approve()
	return s.repo.UpdateComment(ctx, comment)
}

func (s *service) RejectComment(ctx context.Context, commentID uuid.UUID) error {
	comment, err := s.repo.GetComment(ctx, commentID)
	if err != nil {
		return err
	}

	comment.Reject()
	return s.repo.UpdateComment(ctx, comment)
}

func (s *service) MarkCommentAsSpam(ctx context.Context, commentID uuid.UUID) error {
	comment, err := s.repo.GetComment(ctx, commentID)
	if err != nil {
		return err
	}

	comment.MarkAsSpam()
	return s.repo.UpdateComment(ctx, comment)
}

func (s *service) GetCommentCount(ctx context.Context, postID uuid.UUID, status *CommentStatus) (int64, error) {
	return s.repo.GetCommentCount(ctx, postID, status)
}

