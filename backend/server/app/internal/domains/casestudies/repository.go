package casestudies

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for case studies.
type Repository interface {
	CreateCaseStudy(ctx context.Context, caseStudy *CaseStudy) error
	UpdateCaseStudy(ctx context.Context, caseStudy *CaseStudy) error
	GetCaseStudy(ctx context.Context, caseStudyID uuid.UUID) (*CaseStudy, error)
	GetCaseStudyByProjectSlug(ctx context.Context, projectSlug string) (*CaseStudy, error)
	GetCaseStudyByProjectID(ctx context.Context, projectID uuid.UUID) (*CaseStudy, error)
	ListCaseStudies(ctx context.Context, filters CaseStudyFilters) ([]CaseStudy, error)
	DeleteCaseStudy(ctx context.Context, caseStudyID uuid.UUID) error
}

// CaseStudyFilters represents filtering options for listing case studies.
type CaseStudyFilters struct {
	UserID      *uuid.UUID
	ProjectID   *uuid.UUID
	ProjectSlug *string
	Featured    *bool
	Limit       int
	Offset      int
	OrderBy     string // "created_at", "updated_at", "title"
	Order       string // "asc", "desc"
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateCaseStudy(ctx context.Context, caseStudy *CaseStudy) error {
	if caseStudy == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilCaseStudy)
	}

	if err := caseStudy.Validate(); err != nil {
		return err
	}

	now := time.Now().UTC()
	caseStudy.CreatedAt = now
	caseStudy.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(caseStudy).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return NewDomainError(ErrCodeConflict, ErrCaseStudyAlreadyExists)
			}
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) UpdateCaseStudy(ctx context.Context, caseStudy *CaseStudy) error {
	if caseStudy == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilCaseStudy)
	}

	if caseStudy.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCaseStudyID)
	}

	if err := caseStudy.Validate(); err != nil {
		return err
	}

	caseStudy.UpdatedAt = time.Now().UTC()

	result := r.db.WithContext(ctx).Model(&CaseStudy{}).
		Where("id = ?", caseStudy.ID).
		Updates(caseStudy)

	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrCaseStudyNotFound)
	}

	return nil
}

func (r *gormRepository) GetCaseStudy(ctx context.Context, caseStudyID uuid.UUID) (*CaseStudy, error) {
	if caseStudyID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyCaseStudyID)
	}

	var caseStudy CaseStudy
	if err := r.db.WithContext(ctx).Where("id = ?", caseStudyID).First(&caseStudy).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrCaseStudyNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &caseStudy, nil
}

func (r *gormRepository) GetCaseStudyByProjectSlug(ctx context.Context, projectSlug string) (*CaseStudy, error) {
	if projectSlug == "" {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectSlug)
	}

	var caseStudy CaseStudy
	if err := r.db.WithContext(ctx).Where("project_slug = ?", projectSlug).First(&caseStudy).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrCaseStudyNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &caseStudy, nil
}

func (r *gormRepository) GetCaseStudyByProjectID(ctx context.Context, projectID uuid.UUID) (*CaseStudy, error) {
	if projectID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyProjectID)
	}

	var caseStudy CaseStudy
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).First(&caseStudy).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrCaseStudyNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &caseStudy, nil
}

func (r *gormRepository) ListCaseStudies(ctx context.Context, filters CaseStudyFilters) ([]CaseStudy, error) {
	query := r.db.WithContext(ctx).Model(&CaseStudy{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.ProjectID != nil {
		query = query.Where("project_id = ?", *filters.ProjectID)
	}

	if filters.ProjectSlug != nil {
		query = query.Where("project_slug = ?", *filters.ProjectSlug)
	}

	if filters.Featured != nil {
		query = query.Where("featured = ?", *filters.Featured)
	}

	// Default ordering
	orderBy := normalizeOrderBy(filters.OrderBy)
	if orderBy == "" {
		orderBy = "created_at"
	}
	order := filters.Order
	if order == "" {
		order = "desc"
	}
	// Validate order direction
	if order != "asc" && order != "desc" {
		order = "desc"
	}
	query = query.Order(orderBy + " " + order)

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var caseStudies []CaseStudy
	if err := query.Find(&caseStudies).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return caseStudies, nil
}

func (r *gormRepository) DeleteCaseStudy(ctx context.Context, caseStudyID uuid.UUID) error {
	if caseStudyID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCaseStudyID)
	}

	result := r.db.WithContext(ctx).Where("id = ?", caseStudyID).Delete(&CaseStudy{})
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrCaseStudyNotFound)
	}

	return nil
}

// normalizeOrderBy converts camelCase orderBy values to snake_case database column names
// and validates that the column is allowed for ordering
func normalizeOrderBy(orderBy string) string {
	if orderBy == "" {
		return ""
	}

	// Map of allowed camelCase to snake_case conversions
	allowedColumns := map[string]string{
		"createdAt":   "created_at",
		"updatedAt":   "updated_at",
		"title":       "title",
		"created_at":  "created_at",
		"updated_at":  "updated_at",
	}

	// Check if it's already in the map
	if normalized, ok := allowedColumns[orderBy]; ok {
		return normalized
	}

	// Convert camelCase to snake_case
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(orderBy, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ToLower(snake)

	// Validate the converted value is allowed
	if _, ok := allowedColumns[snake]; ok {
		return snake
	}

	// If not in allowed list, return empty string to use default
	return ""
}

