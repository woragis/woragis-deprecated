package testimonials

import (
	"context"
	"errors"
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

