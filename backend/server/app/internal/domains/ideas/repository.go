package ideas

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for ideas and links.
type Repository interface {
	CreateIdea(ctx context.Context, idea *Idea) error
	UpdateIdea(ctx context.Context, idea *Idea) error
	GetIdea(ctx context.Context, id, userID uuid.UUID) (*Idea, error)
	GetIdeaByID(ctx context.Context, id uuid.UUID) (*Idea, error)
	GetIdeaBySlug(ctx context.Context, slug string, userID uuid.UUID) (*Idea, error)
	ListIdeas(ctx context.Context, userID uuid.UUID) ([]Idea, error)
	BulkMoveIdeas(ctx context.Context, userID uuid.UUID, updates []IdeaPositionUpdate) error
	BulkUpdateDetails(ctx context.Context, userID uuid.UUID, updates []IdeaDetailUpdate) error
	DeleteIdeas(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
	RestoreIdeas(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error
	CreateVersion(ctx context.Context, version *IdeaVersion) error
	ListVersions(ctx context.Context, ideaID, userID uuid.UUID, limit int) ([]IdeaVersion, error)
	CreateLink(ctx context.Context, link *IdeaLink) error
	ListLinks(ctx context.Context, filters LinkFilters) ([]IdeaLink, error)
	AddCollaborator(ctx context.Context, collaborator *IdeaCollaborator) error
	RemoveCollaborator(ctx context.Context, ownerID, collaboratorID uuid.UUID) error
	ListCollaborators(ctx context.Context, ownerID uuid.UUID) ([]IdeaCollaborator, error)
	HasCollaborator(ctx context.Context, ownerID, collaboratorID uuid.UUID) (bool, error)
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

func (r *gormRepository) GetIdeaByID(ctx context.Context, id uuid.UUID) (*Idea, error) {
	var idea Idea
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&idea).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrIdeaNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &idea, nil
}

func (r *gormRepository) GetIdeaBySlug(ctx context.Context, slug string, userID uuid.UUID) (*Idea, error) {
	var idea Idea
	err := r.db.WithContext(ctx).Where("slug = ? AND user_id = ?", slug, userID).First(&idea).Error
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

// IdeaPositionUpdate holds coordinates for a node.
type IdeaPositionUpdate struct {
	IdeaID uuid.UUID
	PosX   float64
	PosY   float64
}

// IdeaDetailUpdate holds metadata for a node.
type IdeaDetailUpdate struct {
	IdeaID      uuid.UUID
	Title       string
	Description string
	Color       string
	ProjectID   *uuid.UUID
}

func (r *gormRepository) BulkMoveIdeas(ctx context.Context, userID uuid.UUID, updates []IdeaPositionUpdate) error {
	if len(updates) == 0 {
		return nil
	}
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	now := time.Now().UTC()
	for _, upd := range updates {
		if err := tx.Model(&Idea{}).
			Where("user_id = ? AND id = ?", userID, upd.IdeaID).
			Updates(map[string]any{
				"pos_x":      upd.PosX,
				"pos_y":      upd.PosY,
				"version":    gorm.Expr("version + ?", 1),
				"updated_at": now,
			}).Error; err != nil {
			tx.Rollback()
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) BulkUpdateDetails(ctx context.Context, userID uuid.UUID, updates []IdeaDetailUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	now := time.Now().UTC()
	for _, upd := range updates {
		payload := map[string]any{
			"updated_at": now,
			"version":    gorm.Expr("version + ?", 1),
		}
		if strings.TrimSpace(upd.Title) != "" {
			payload["title"] = strings.TrimSpace(upd.Title)
		}
		if upd.Description != "" {
			payload["description"] = strings.TrimSpace(upd.Description)
		}
		if upd.Color != "" {
			payload["color"] = strings.TrimSpace(upd.Color)
		}
		if upd.ProjectID != nil {
			payload["project_id"] = upd.ProjectID
		}
		if err := tx.Model(&Idea{}).
			Where("user_id = ? AND id = ?", userID, upd.IdeaID).
			Updates(payload).Error; err != nil {
			tx.Rollback()
			return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) DeleteIdeas(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND id IN ?", userID, ids).
		Delete(&Idea{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) RestoreIdeas(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).
		Model(&Idea{}).
		Unscoped().
		Where("user_id = ? AND id IN ?", userID, ids).
		Update("deleted_at", nil).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) CreateVersion(ctx context.Context, version *IdeaVersion) error {
	if err := r.db.WithContext(ctx).Create(version).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) ListVersions(ctx context.Context, ideaID, userID uuid.UUID, limit int) ([]IdeaVersion, error) {
	if limit <= 0 {
		limit = 20
	}
	var versions []IdeaVersion
	if err := r.db.WithContext(ctx).
		Where("idea_id = ? AND user_id = ?", ideaID, userID).
		Order("version desc").
		Limit(limit).
		Find(&versions).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return versions, nil
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

// LinkFilters constrains relationship queries.
type LinkFilters struct {
	UserID           uuid.UUID
	IdeaID           uuid.UUID
	Relation         string
	RelationContains string
	MinWeight        *float64
	MaxWeight        *float64
	Bidirectional    *bool
}

func (r *gormRepository) ListLinks(ctx context.Context, filters LinkFilters) ([]IdeaLink, error) {
	var links []IdeaLink
	query := r.db.WithContext(ctx).Where("user_id = ?", filters.UserID)
	if filters.IdeaID != uuid.Nil {
		query = query.Where("source_idea_id = ? OR target_idea_id = ?", filters.IdeaID, filters.IdeaID)
	}
	if strings.TrimSpace(filters.Relation) != "" {
		query = query.Where("LOWER(relation) = ?", strings.ToLower(strings.TrimSpace(filters.Relation)))
	}
	if strings.TrimSpace(filters.RelationContains) != "" {
		query = query.Where("LOWER(relation) LIKE ?", "%"+strings.ToLower(strings.TrimSpace(filters.RelationContains))+"%")
	}
	if filters.MinWeight != nil {
		query = query.Where("weight >= ?", *filters.MinWeight)
	}
	if filters.MaxWeight != nil {
		query = query.Where("weight <= ?", *filters.MaxWeight)
	}
	if filters.Bidirectional != nil {
		query = query.Where("bidirectional = ?", *filters.Bidirectional)
	}

	if err := query.Find(&links).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return links, nil
}

func (r *gormRepository) AddCollaborator(ctx context.Context, collaborator *IdeaCollaborator) error {
	if err := collaborator.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(collaborator).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return NewDomainError(ErrCodeCollaboratorConflict, ErrCollaboratorExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) RemoveCollaborator(ctx context.Context, ownerID, collaboratorID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("owner_id = ? AND collaborator_id = ?", ownerID, collaboratorID).
		Delete(&IdeaCollaborator{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) ListCollaborators(ctx context.Context, ownerID uuid.UUID) ([]IdeaCollaborator, error) {
	var collaborators []IdeaCollaborator
	if err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("created_at desc").
		Find(&collaborators).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return collaborators, nil
}

func (r *gormRepository) HasCollaborator(ctx context.Context, ownerID, collaboratorID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&IdeaCollaborator{}).
		Where("owner_id = ? AND collaborator_id = ?", ownerID, collaboratorID).
		Count(&count).Error; err != nil {
		return false, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return count > 0, nil
}
