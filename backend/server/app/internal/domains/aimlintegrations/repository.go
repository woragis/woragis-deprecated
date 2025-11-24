package aimlintegrations

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for AI/ML integrations.
type Repository interface {
	CreateAIMLIntegration(ctx context.Context, integration *AIMLIntegration) error
	UpdateAIMLIntegration(ctx context.Context, integration *AIMLIntegration) error
	GetAIMLIntegration(ctx context.Context, integrationID uuid.UUID, userID uuid.UUID) (*AIMLIntegration, error)
	GetAIMLIntegrationPublic(ctx context.Context, integrationID uuid.UUID) (*AIMLIntegration, error)
	ListAIMLIntegrations(ctx context.Context, filters AIMLIntegrationFilters) ([]AIMLIntegration, error)
	ListFeaturedAIMLIntegrations(ctx context.Context) ([]AIMLIntegration, error)
	GetIntegrationsByProject(ctx context.Context, projectID uuid.UUID) ([]AIMLIntegration, error)
	GetIntegrationsByType(ctx context.Context, integrationType IntegrationType) ([]AIMLIntegration, error)
	GetIntegrationsByFramework(ctx context.Context, framework Framework) ([]AIMLIntegration, error)
	DeleteAIMLIntegration(ctx context.Context, integrationID uuid.UUID, userID uuid.UUID) error
}

// AIMLIntegrationFilters represents filtering options for listing integrations.
type AIMLIntegrationFilters struct {
	UserID      *uuid.UUID
	Type        *IntegrationType
	Framework   *Framework
	ProjectID   *uuid.UUID
	Featured    *bool
	Limit       int
	Offset      int
	OrderBy     string // "created_at", "display_order", "title"
	Order       string // "asc", "desc"
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateAIMLIntegration(ctx context.Context, integration *AIMLIntegration) error {
	if integration == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilIntegration)
	}

	if err := integration.Validate(); err != nil {
		return err
	}

	now := time.Now().UTC()
	integration.CreatedAt = now
	integration.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(integration).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return NewDomainError(ErrCodeConflict, ErrIntegrationAlreadyExists)
			}
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) UpdateAIMLIntegration(ctx context.Context, integration *AIMLIntegration) error {
	if integration == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilIntegration)
	}

	if integration.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyIntegrationID)
	}

	if err := integration.Validate(); err != nil {
		return err
	}

	integration.UpdatedAt = time.Now().UTC()

	result := r.db.WithContext(ctx).Model(&AIMLIntegration{}).
		Where("id = ?", integration.ID).
		Updates(integration)

	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrIntegrationNotFound)
	}

	return nil
}

func (r *gormRepository) GetAIMLIntegration(ctx context.Context, integrationID uuid.UUID, userID uuid.UUID) (*AIMLIntegration, error) {
	if integrationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyIntegrationID)
	}

	var integration AIMLIntegration
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", integrationID, userID).
		First(&integration).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrIntegrationNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &integration, nil
}

func (r *gormRepository) GetAIMLIntegrationPublic(ctx context.Context, integrationID uuid.UUID) (*AIMLIntegration, error) {
	if integrationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyIntegrationID)
	}

	var integration AIMLIntegration
	err := r.db.WithContext(ctx).
		Where("id = ? AND featured = ?", integrationID, true).
		First(&integration).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrIntegrationNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &integration, nil
}

func (r *gormRepository) ListAIMLIntegrations(ctx context.Context, filters AIMLIntegrationFilters) ([]AIMLIntegration, error) {
	query := r.db.WithContext(ctx).Model(&AIMLIntegration{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.Framework != nil {
		query = query.Where("framework = ?", *filters.Framework)
	}

	if filters.ProjectID != nil {
		query = query.Where("project_id = ?", *filters.ProjectID)
	}

	if filters.Featured != nil {
		query = query.Where("featured = ?", *filters.Featured)
	}

	// Default ordering
	orderBy := filters.OrderBy
	if orderBy == "" {
		orderBy = "display_order"
	}
	order := filters.Order
	if order == "" {
		order = "asc"
	}
	query = query.Order(orderBy + " " + order)

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var integrations []AIMLIntegration
	if err := query.Find(&integrations).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return integrations, nil
}

func (r *gormRepository) ListFeaturedAIMLIntegrations(ctx context.Context) ([]AIMLIntegration, error) {
	var integrations []AIMLIntegration
	err := r.db.WithContext(ctx).
		Where("featured = ?", true).
		Order("display_order ASC, created_at DESC").
		Find(&integrations).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return integrations, nil
}

func (r *gormRepository) GetIntegrationsByProject(ctx context.Context, projectID uuid.UUID) ([]AIMLIntegration, error) {
	var integrations []AIMLIntegration
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("display_order ASC, created_at DESC").
		Find(&integrations).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return integrations, nil
}

func (r *gormRepository) GetIntegrationsByType(ctx context.Context, integrationType IntegrationType) ([]AIMLIntegration, error) {
	var integrations []AIMLIntegration
	err := r.db.WithContext(ctx).
		Where("type = ?", integrationType).
		Order("display_order ASC, created_at DESC").
		Find(&integrations).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return integrations, nil
}

func (r *gormRepository) GetIntegrationsByFramework(ctx context.Context, framework Framework) ([]AIMLIntegration, error) {
	var integrations []AIMLIntegration
	err := r.db.WithContext(ctx).
		Where("framework = ?", framework).
		Order("display_order ASC, created_at DESC").
		Find(&integrations).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return integrations, nil
}

func (r *gormRepository) DeleteAIMLIntegration(ctx context.Context, integrationID uuid.UUID, userID uuid.UUID) error {
	if integrationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyIntegrationID)
	}

	// Verify ownership
	var integration AIMLIntegration
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", integrationID, userID).
		First(&integration).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewDomainError(ErrCodeNotFound, ErrIntegrationNotFound)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	if err := r.db.WithContext(ctx).Delete(&integration).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

