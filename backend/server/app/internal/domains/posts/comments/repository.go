package comments

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for comments.
type Repository interface {
	CreateComment(ctx context.Context, comment *Comment) error
	UpdateComment(ctx context.Context, comment *Comment) error
	GetComment(ctx context.Context, commentID uuid.UUID) (*Comment, error)
	DeleteComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error
	ListComments(ctx context.Context, filters CommentFilters) ([]Comment, error)
	GetCommentCount(ctx context.Context, postID uuid.UUID, status *CommentStatus) (int64, error)
}

// CommentFilters represents filtering options for listing comments.
type CommentFilters struct {
	PostID   *uuid.UUID
	UserID   *uuid.UUID
	ParentID *uuid.UUID // nil means top-level comments, uuid.Nil means all, specific ID means replies to that comment
	Status   *CommentStatus
	Limit    int
	Offset   int
	OrderBy  string // "created_at", "updated_at"
	Order    string // "asc", "desc"
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// isUniqueConstraintError checks if the error is a unique constraint violation.
func isUniqueConstraintError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}

func (r *gormRepository) CreateComment(ctx context.Context, comment *Comment) error {
	if err := comment.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(comment).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateComment(ctx context.Context, comment *Comment) error {
	if err := comment.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(comment).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetComment(ctx context.Context, commentID uuid.UUID) (*Comment, error) {
	var comment Comment
	err := r.db.WithContext(ctx).Where("id = ?", commentID).First(&comment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeCommentNotFound, ErrCommentNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &comment, nil
}

func (r *gormRepository) DeleteComment(ctx context.Context, commentID uuid.UUID, userID uuid.UUID) error {
	// First verify ownership
	var comment Comment
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", commentID, userID).First(&comment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return NewDomainError(ErrCodeCommentNotFound, ErrCommentNotFound)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	// Delete the comment
	if err := r.db.WithContext(ctx).Delete(&comment).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) ListComments(ctx context.Context, filters CommentFilters) ([]Comment, error) {
	var comments []Comment
	query := r.db.WithContext(ctx).Model(&Comment{})

	if filters.PostID != nil {
		query = query.Where("post_id = ?", *filters.PostID)
	}

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.ParentID != nil {
		if *filters.ParentID == uuid.Nil {
			// Top-level comments only (parent_id IS NULL)
			query = query.Where("parent_id IS NULL")
		} else {
			// Replies to specific comment
			query = query.Where("parent_id = ?", *filters.ParentID)
		}
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	// Ordering
	orderBy := filters.OrderBy
	if orderBy == "" {
		orderBy = "created_at"
	}
	order := filters.Order
	if order == "" {
		order = "asc" // Comments typically shown oldest first
	}
	query = query.Order(orderBy + " " + order)

	// Pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	if err := query.Find(&comments).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return comments, nil
}

func (r *gormRepository) GetCommentCount(ctx context.Context, postID uuid.UUID, status *CommentStatus) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&Comment{}).Where("post_id = ?", postID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return count, nil
}

