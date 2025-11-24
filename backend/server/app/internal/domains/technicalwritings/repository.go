package technicalwritings

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for technical writings.
type Repository interface {
	CreateTechnicalWriting(ctx context.Context, writing *TechnicalWriting) error
	UpdateTechnicalWriting(ctx context.Context, writing *TechnicalWriting) error
	GetTechnicalWriting(ctx context.Context, writingID uuid.UUID, userID uuid.UUID) (*TechnicalWriting, error)
	GetTechnicalWritingPublic(ctx context.Context, writingID uuid.UUID) (*TechnicalWriting, error)
	ListTechnicalWritings(ctx context.Context, filters TechnicalWritingFilters) ([]TechnicalWriting, error)
	ListFeaturedTechnicalWritings(ctx context.Context) ([]TechnicalWriting, error)
	GetWritingsByProject(ctx context.Context, projectID uuid.UUID) ([]TechnicalWriting, error)
	GetWritingsByType(ctx context.Context, writingType WritingType) ([]TechnicalWriting, error)
	GetWritingsByPlatform(ctx context.Context, platform PublicationPlatform) ([]TechnicalWriting, error)
	SearchTechnicalWritings(ctx context.Context, query string) ([]TechnicalWriting, error)
	DeleteTechnicalWriting(ctx context.Context, writingID uuid.UUID, userID uuid.UUID) error
}

// TechnicalWritingFilters represents filtering options for listing writings.
type TechnicalWritingFilters struct {
	UserID      *uuid.UUID
	Type        *WritingType
	Platform    *PublicationPlatform
	ProjectID   *uuid.UUID
	Featured    *bool
	Limit       int
	Offset      int
	OrderBy     string // "created_at", "published_at", "display_order", "views", "likes"
	Order       string // "asc", "desc"
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateTechnicalWriting(ctx context.Context, writing *TechnicalWriting) error {
	if writing == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilWriting)
	}

	if err := writing.Validate(); err != nil {
		return err
	}

	now := time.Now().UTC()
	writing.CreatedAt = now
	writing.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(writing).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return NewDomainError(ErrCodeConflict, ErrWritingAlreadyExists)
			}
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) UpdateTechnicalWriting(ctx context.Context, writing *TechnicalWriting) error {
	if writing == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilWriting)
	}

	if writing.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyWritingID)
	}

	if err := writing.Validate(); err != nil {
		return err
	}

	writing.UpdatedAt = time.Now().UTC()

	result := r.db.WithContext(ctx).Model(&TechnicalWriting{}).
		Where("id = ?", writing.ID).
		Updates(writing)

	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrWritingNotFound)
	}

	return nil
}

func (r *gormRepository) GetTechnicalWriting(ctx context.Context, writingID uuid.UUID, userID uuid.UUID) (*TechnicalWriting, error) {
	if writingID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyWritingID)
	}

	var writing TechnicalWriting
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", writingID, userID).
		First(&writing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrWritingNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &writing, nil
}

func (r *gormRepository) GetTechnicalWritingPublic(ctx context.Context, writingID uuid.UUID) (*TechnicalWriting, error) {
	if writingID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyWritingID)
	}

	var writing TechnicalWriting
	err := r.db.WithContext(ctx).
		Where("id = ? AND featured = ?", writingID, true).
		First(&writing).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrWritingNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &writing, nil
}

func (r *gormRepository) ListTechnicalWritings(ctx context.Context, filters TechnicalWritingFilters) ([]TechnicalWriting, error) {
	query := r.db.WithContext(ctx).Model(&TechnicalWriting{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Type != nil {
		query = query.Where("type = ?", *filters.Type)
	}

	if filters.Platform != nil {
		query = query.Where("platform = ?", *filters.Platform)
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

	var writings []TechnicalWriting
	if err := query.Find(&writings).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return writings, nil
}

func (r *gormRepository) ListFeaturedTechnicalWritings(ctx context.Context) ([]TechnicalWriting, error) {
	var writings []TechnicalWriting
	err := r.db.WithContext(ctx).
		Where("featured = ?", true).
		Order("display_order ASC, published_at DESC").
		Find(&writings).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return writings, nil
}

func (r *gormRepository) GetWritingsByProject(ctx context.Context, projectID uuid.UUID) ([]TechnicalWriting, error) {
	var writings []TechnicalWriting
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("display_order ASC, published_at DESC").
		Find(&writings).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return writings, nil
}

func (r *gormRepository) GetWritingsByType(ctx context.Context, writingType WritingType) ([]TechnicalWriting, error) {
	var writings []TechnicalWriting
	err := r.db.WithContext(ctx).
		Where("type = ?", writingType).
		Order("display_order ASC, published_at DESC").
		Find(&writings).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return writings, nil
}

func (r *gormRepository) GetWritingsByPlatform(ctx context.Context, platform PublicationPlatform) ([]TechnicalWriting, error) {
	var writings []TechnicalWriting
	err := r.db.WithContext(ctx).
		Where("platform = ?", platform).
		Order("display_order ASC, published_at DESC").
		Find(&writings).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return writings, nil
}

func (r *gormRepository) SearchTechnicalWritings(ctx context.Context, query string) ([]TechnicalWriting, error) {
	var writings []TechnicalWriting
	err := r.db.WithContext(ctx).
		Where("title ILIKE ? OR description ILIKE ? OR content ILIKE ?", "%"+query+"%", "%"+query+"%", "%"+query+"%").
		Order("published_at DESC").
		Find(&writings).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return writings, nil
}

func (r *gormRepository) DeleteTechnicalWriting(ctx context.Context, writingID uuid.UUID, userID uuid.UUID) error {
	if writingID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyWritingID)
	}

	// Verify ownership
	var writing TechnicalWriting
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", writingID, userID).
		First(&writing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewDomainError(ErrCodeNotFound, ErrWritingNotFound)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	if err := r.db.WithContext(ctx).Delete(&writing).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

