package platforms

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for platform configurations.
type Repository interface {
	CreateConfig(ctx context.Context, config *PlatformConfig) error
	UpdateConfig(ctx context.Context, config *PlatformConfig) error
	GetConfig(ctx context.Context, configID uuid.UUID) (*PlatformConfig, error)
	GetConfigByName(ctx context.Context, name Platform) (*PlatformConfig, error)
	ListConfigs(ctx context.Context, activeOnly bool) ([]PlatformConfig, error)
	DeleteConfig(ctx context.Context, configID uuid.UUID) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateConfig(ctx context.Context, config *PlatformConfig) error {
	if err := config.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(config).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrPlatformConfigExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateConfig(ctx context.Context, config *PlatformConfig) error {
	if err := config.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(config).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrPlatformConfigExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetConfig(ctx context.Context, configID uuid.UUID) (*PlatformConfig, error) {
	var config PlatformConfig
	err := r.db.WithContext(ctx).Where("id = ?", configID).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrPlatformConfigNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &config, nil
}

func (r *gormRepository) GetConfigByName(ctx context.Context, name Platform) (*PlatformConfig, error) {
	var config PlatformConfig
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrPlatformConfigNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &config, nil
}

func (r *gormRepository) ListConfigs(ctx context.Context, activeOnly bool) ([]PlatformConfig, error) {
	var configs []PlatformConfig
	query := r.db.WithContext(ctx)

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Order("display_name asc").Find(&configs).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return configs, nil
}

func (r *gormRepository) DeleteConfig(ctx context.Context, configID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&PlatformConfig{}, "id = ?", configID).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
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
