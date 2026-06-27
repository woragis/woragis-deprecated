package files

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for files.
type Repository interface {
	CreateFile(ctx context.Context, file *File) error
	UpdateFile(ctx context.Context, file *File) error
	GetFile(ctx context.Context, fileID uuid.UUID) (*File, error)
	GetFileByUser(ctx context.Context, fileID, userID uuid.UUID) (*File, error)
	DeleteFile(ctx context.Context, fileID, userID uuid.UUID) error
	ListFiles(ctx context.Context, filters FileFilters) ([]File, error)
	GetFileByHash(ctx context.Context, hash string) (*File, error)
}

// FileFilters represents filtering options for listing files.
type FileFilters struct {
	UserID   *uuid.UUID
	FileType *FileType
	Status   *FileStatus
	Limit    int
	Offset   int
	OrderBy  string // "created_at", "updated_at", "size", "name"
	Order    string // "asc", "desc"
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// isUniqueConstraintError checks if the error is a unique constraint violation.
func isUniqueConstraintError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}

func (r *gormRepository) CreateFile(ctx context.Context, file *File) error {
	if err := file.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		if isUniqueConstraintError(err) {
			return NewDomainError(ErrCodeConflict, ErrFileExists)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateFile(ctx context.Context, file *File) error {
	if err := file.Validate(); err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).Save(file).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetFile(ctx context.Context, fileID uuid.UUID) (*File, error) {
	var file File
	err := r.db.WithContext(ctx).Where("id = ?", fileID).First(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrFileNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &file, nil
}

func (r *gormRepository) GetFileByUser(ctx context.Context, fileID, userID uuid.UUID) (*File, error) {
	var file File
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrFileNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &file, nil
}

func (r *gormRepository) DeleteFile(ctx context.Context, fileID, userID uuid.UUID) error {
	// First verify ownership
	var file File
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return NewDomainError(ErrCodeNotFound, ErrFileNotFound)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	// Soft delete
	if err := r.db.WithContext(ctx).Delete(&file).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToDelete)
	}
	return nil
}

func (r *gormRepository) ListFiles(ctx context.Context, filters FileFilters) ([]File, error) {
	var files []File
	query := r.db.WithContext(ctx).Model(&File{})

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.FileType != nil {
		query = query.Where("file_type = ?", *filters.FileType)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	// Ordering
	orderBy := filters.OrderBy
	if orderBy == "" {
		orderBy = "created_at"
	}
	order := filters.Order
	if order == "" {
		order = "desc"
	}
	query = query.Order(orderBy + " " + order)

	// Pagination
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	if err := query.Find(&files).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return files, nil
}

func (r *gormRepository) GetFileByHash(ctx context.Context, hash string) (*File, error) {
	var file File
	err := r.db.WithContext(ctx).Where("hash = ?", hash).First(&file).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrFileNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &file, nil
}

