package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

// Repository abstraction for persistence.
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// GormRepository implements Repository using PostgreSQL via GORM.
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new repository instance.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// Create persists a new user.
func (r *GormRepository) Create(ctx context.Context, user *User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return mapPersistenceError(err)
	}

	return nil
}

// FindByEmail retrieves a user by email, returning a domain error when not found.
func (r *GormRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeUserNotFound, ErrUserNotFound)
		}
		return nil, mapPersistenceError(err)
	}

	return &user, nil
}

func mapPersistenceError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return NewDomainError(ErrCodeEmailAlreadyExists, ErrUserAlreadyExists)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return NewDomainError(ErrCodeEmailAlreadyExists, ErrUserAlreadyExists)
	}

	// TODO: inspect pq error codes for duplicate violations and map precisely.
	return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
}
