package ideas

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for ideas and links.
type Repository interface {
	CreateIdea(ctx context.Context, idea *Idea) error
	UpdateIdea(ctx context.Context, idea *Idea) error
	GetIdea(ctx context.Context, id, userID uuid.UUID) (*Idea, error)
	ListIdeas(ctx context.Context, userID uuid.UUID) ([]Idea, error)

	CreateLink(ctx context.Context, link *IdeaLink) error
	ListLinks(ctx context.Context, userID uuid.UUID, ideaID uuid.UUID) ([]IdeaLink, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a new repository instance.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateIdea(ctx context.Context, idea *Idea) error {
	if err := idea.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(idea).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) UpdateIdea(ctx context.Context, idea *Idea) error {
	if err := idea.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(idea).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) GetIdea(ctx context.Context, id, userID uuid.UUID) (*Idea, error) {
	var idea Idea
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&idea).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrIdeaNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &idea, nil
}

func (r *gormRepository) ListIdeas(ctx context.Context, userID uuid.UUID) ([]Idea, error) {
	var ideas []Idea
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at desc").
		Find(&ideas).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return ideas, nil
}

func (r *gormRepository) CreateLink(ctx context.Context, link *IdeaLink) error {
	if err := link.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(link).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) ListLinks(ctx context.Context, userID uuid.UUID, ideaID uuid.UUID) ([]IdeaLink, error) {
	var links []IdeaLink
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if ideaID != uuid.Nil {
		query = query.Where("source_idea_id = ? OR target_idea_id = ?", ideaID, ideaID)
	}

	if err := query.Find(&links).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return links, nil
}
