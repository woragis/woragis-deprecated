package testimonials

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

// Repository defines persistence operations for testimonials.
type Repository interface {
	CreateTestimonial(ctx context.Context, testimonial *Testimonial) error
	UpdateTestimonial(ctx context.Context, testimonial *Testimonial) error
	GetTestimonial(ctx context.Context, testimonialID uuid.UUID) (*Testimonial, error)
	ListTestimonials(ctx context.Context, filters TestimonialFilters) ([]Testimonial, error)
	DeleteTestimonial(ctx context.Context, testimonialID uuid.UUID) error
	ApproveTestimonial(ctx context.Context, testimonialID uuid.UUID) error
	RejectTestimonial(ctx context.Context, testimonialID uuid.UUID) error
	HideTestimonial(ctx context.Context, testimonialID uuid.UUID) error
	// Entity link methods
	CreateTestimonialEntityLink(ctx context.Context, link *TestimonialEntityLink) error
	GetTestimonialEntityLinks(ctx context.Context, testimonialID uuid.UUID) ([]TestimonialEntityLink, error)
	GetEntityTestimonials(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]Testimonial, error)
	DeleteTestimonialEntityLink(ctx context.Context, linkID uuid.UUID) error
	DeleteTestimonialEntityLinks(ctx context.Context, testimonialID uuid.UUID) error
}

// TestimonialFilters represents filtering options for listing testimonials.
type TestimonialFilters struct {
	UserID *uuid.UUID
	Status *TestimonialStatus
	Limit  int
	Offset int
	OrderBy string // "created_at", "updated_at", "display_order"
	Order   string // "asc", "desc"
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateTestimonial(ctx context.Context, testimonial *Testimonial) error {
	if testimonial == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilTestimonial)
	}

	if err := testimonial.Validate(); err != nil {
		return err
	}

	now := time.Now()
	testimonial.CreatedAt = now
	testimonial.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(testimonial).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return NewDomainError(ErrCodeConflict, ErrTestimonialAlreadyExists)
			}
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) UpdateTestimonial(ctx context.Context, testimonial *Testimonial) error {
	if testimonial == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilTestimonial)
	}

	if testimonial.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	if err := testimonial.Validate(); err != nil {
		return err
	}

	testimonial.UpdatedAt = time.Now()

	result := r.db.WithContext(ctx).Model(&Testimonial{}).
		Where("id = ?", testimonial.ID).
		Updates(testimonial)

	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrTestimonialNotFound)
	}

	return nil
}

func (r *gormRepository) GetTestimonial(ctx context.Context, testimonialID uuid.UUID) (*Testimonial, error) {
	if testimonialID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	var testimonial Testimonial
	if err := r.db.WithContext(ctx).Where("id = ?", testimonialID).First(&testimonial).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrTestimonialNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &testimonial, nil
}

func (r *gormRepository) ListTestimonials(ctx context.Context, filters TestimonialFilters) ([]Testimonial, error) {
	query := r.db.WithContext(ctx).Model(&Testimonial{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	// Default ordering
	orderBy := normalizeOrderBy(filters.OrderBy)
	if orderBy == "" {
		orderBy = "display_order"
	}
	order := filters.Order
	if order == "" {
		order = "asc"
	}
	// Validate order direction
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	query = query.Order(orderBy + " " + order)

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var testimonials []Testimonial
	if err := query.Find(&testimonials).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return testimonials, nil
}

func (r *gormRepository) DeleteTestimonial(ctx context.Context, testimonialID uuid.UUID) error {
	if testimonialID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	result := r.db.WithContext(ctx).Where("id = ?", testimonialID).Delete(&Testimonial{})
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrTestimonialNotFound)
	}

	return nil
}

func (r *gormRepository) ApproveTestimonial(ctx context.Context, testimonialID uuid.UUID) error {
	if testimonialID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	result := r.db.WithContext(ctx).Model(&Testimonial{}).
		Where("id = ?", testimonialID).
		Updates(map[string]interface{}{
			"status":     TestimonialStatusApproved,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrTestimonialNotFound)
	}

	return nil
}

func (r *gormRepository) RejectTestimonial(ctx context.Context, testimonialID uuid.UUID) error {
	if testimonialID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	result := r.db.WithContext(ctx).Model(&Testimonial{}).
		Where("id = ?", testimonialID).
		Updates(map[string]interface{}{
			"status":     TestimonialStatusRejected,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrTestimonialNotFound)
	}

	return nil
}

func (r *gormRepository) HideTestimonial(ctx context.Context, testimonialID uuid.UUID) error {
	if testimonialID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	result := r.db.WithContext(ctx).Model(&Testimonial{}).
		Where("id = ?", testimonialID).
		Updates(map[string]interface{}{
			"status":     TestimonialStatusHidden,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrTestimonialNotFound)
	}

	return nil
}

// Entity link methods

func (r *gormRepository) CreateTestimonialEntityLink(ctx context.Context, link *TestimonialEntityLink) error {
	if link == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilLink)
	}

	if err := link.Validate(); err != nil {
		return err
	}

	now := time.Now()
	link.CreatedAt = now
	link.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(link).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) GetTestimonialEntityLinks(ctx context.Context, testimonialID uuid.UUID) ([]TestimonialEntityLink, error) {
	if testimonialID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	var links []TestimonialEntityLink
	if err := r.db.WithContext(ctx).Where("testimonial_id = ?", testimonialID).Find(&links).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return links, nil
}

func (r *gormRepository) GetEntityTestimonials(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]Testimonial, error) {
	if entityID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyEntityID)
	}

	var testimonials []Testimonial
	if err := r.db.WithContext(ctx).
		Table("testimonials").
		Joins("INNER JOIN testimonial_entity_links ON testimonials.id = testimonial_entity_links.testimonial_id").
		Where("testimonial_entity_links.entity_type = ? AND testimonial_entity_links.entity_id = ?", entityType, entityID).
		Find(&testimonials).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return testimonials, nil
}

func (r *gormRepository) DeleteTestimonialEntityLink(ctx context.Context, linkID uuid.UUID) error {
	if linkID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyLinkID)
	}

	result := r.db.WithContext(ctx).Where("id = ?", linkID).Delete(&TestimonialEntityLink{})
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) DeleteTestimonialEntityLinks(ctx context.Context, testimonialID uuid.UUID) error {
	if testimonialID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTestimonialID)
	}

	result := r.db.WithContext(ctx).Where("testimonial_id = ?", testimonialID).Delete(&TestimonialEntityLink{})
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
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
		"createdAt":    "created_at",
		"updatedAt":    "updated_at",
		"displayOrder": "display_order",
		"created_at":   "created_at",
		"updated_at":   "updated_at",
		"display_order": "display_order",
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

