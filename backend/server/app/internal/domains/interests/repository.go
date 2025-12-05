package interests

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for interests.
type Repository interface {
	CreateInterest(ctx context.Context, interest *Interest) error
	UpdateInterest(ctx context.Context, interest *Interest) error
	DeleteInterest(ctx context.Context, interestID uuid.UUID) error
	GetInterest(ctx context.Context, interestID uuid.UUID) (*Interest, error)
	GetInterestBySlug(ctx context.Context, slug string) (*Interest, error)
	GetInterestByTitle(ctx context.Context, title string) (*Interest, error)
	ListInterests(ctx context.Context) ([]Interest, error)
	ListFeaturedInterests(ctx context.Context) ([]Interest, error)
	SearchInterests(ctx context.Context, query string) ([]Interest, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateInterest(ctx context.Context, interest *Interest) error {
	if err := interest.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(interest).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrInterestAlreadyExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateInterest(ctx context.Context, interest *Interest) error {
	if err := interest.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(interest).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrInterestAlreadyExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) DeleteInterest(ctx context.Context, interestID uuid.UUID) error {
	// Verify interest exists
	if _, err := r.GetInterest(ctx, interestID); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Where("id = ?", interestID).Delete(&Interest{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
	}
	return nil
}

func (r *gormRepository) GetInterest(ctx context.Context, interestID uuid.UUID) (*Interest, error) {
	var interest Interest
	err := r.db.WithContext(ctx).Where("id = ?", interestID).First(&interest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrInterestNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &interest, nil
}

func (r *gormRepository) GetInterestBySlug(ctx context.Context, slug string) (*Interest, error) {
	var interest Interest
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&interest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrInterestNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &interest, nil
}

func (r *gormRepository) GetInterestByTitle(ctx context.Context, title string) (*Interest, error) {
	var interest Interest
	err := r.db.WithContext(ctx).Where("title = ?", title).First(&interest).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrInterestNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &interest, nil
}

func (r *gormRepository) ListInterests(ctx context.Context) ([]Interest, error) {
	var interests []Interest
	if err := r.db.WithContext(ctx).
		Order("title asc").
		Find(&interests).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return interests, nil
}

func (r *gormRepository) ListFeaturedInterests(ctx context.Context) ([]Interest, error) {
	var interests []Interest
	if err := r.db.WithContext(ctx).
		Where("featured = ?", true).
		Order("title asc").
		Find(&interests).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return interests, nil
}

func (r *gormRepository) SearchInterests(ctx context.Context, query string) ([]Interest, error) {
	var interests []Interest
	pattern := "%" + query + "%"
	if err := r.db.WithContext(ctx).
		Where("LOWER(title) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)", pattern, pattern).
		Order("title asc").
		Find(&interests).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return interests, nil
}

// isUniqueConstraintError checks if the error is a unique constraint violation.
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// Check for PostgreSQL unique constraint violation (code 23505)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}

	// Check for SQLite unique constraint violation
	if err.Error() == "UNIQUE constraint failed" {
		return true
	}

	return false
}

