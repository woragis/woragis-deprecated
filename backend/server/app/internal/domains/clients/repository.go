package clients

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for clients.
type Repository interface {
	Create(ctx context.Context, client *Client) error
	Update(ctx context.Context, client *Client) error
	Get(ctx context.Context, userID, id uuid.UUID) (*Client, error)
	List(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]Client, error)
	Delete(ctx context.Context, userID, id uuid.UUID) error
	SetArchived(ctx context.Context, userID, id uuid.UUID, archived bool) error
	FindByPhoneNumber(ctx context.Context, userID uuid.UUID, phoneNumber string) (*Client, error)
}

// GormRepository implements Repository using GORM.
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM-based repository.
func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// Create inserts a new client.
func (r *GormRepository) Create(ctx context.Context, client *Client) error {
	if err := r.db.WithContext(ctx).Create(client).Error; err != nil {
		return fmt.Errorf("%w: %v", NewDomainError(ErrCodeRepositoryFailure, ErrRepositoryFailure), err)
	}
	return nil
}

// Update saves changes to an existing client.
func (r *GormRepository) Update(ctx context.Context, client *Client) error {
	if err := r.db.WithContext(ctx).Save(client).Error; err != nil {
		return fmt.Errorf("%w: %v", NewDomainError(ErrCodeRepositoryFailure, ErrRepositoryFailure), err)
	}
	return nil
}

// Get retrieves a client by ID for a specific user.
func (r *GormRepository) Get(ctx context.Context, userID, id uuid.UUID) (*Client, error) {
	var client Client
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrClientNotFound)
		}
		return nil, fmt.Errorf("%w: %v", NewDomainError(ErrCodeRepositoryFailure, ErrRepositoryFailure), err)
	}
	return &client, nil
}

// List retrieves all clients for a user, optionally including archived ones.
func (r *GormRepository) List(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]Client, error) {
	var clients []Client
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	
	if !includeArchived {
		query = query.Where("is_archived = ?", false)
	}
	
	if err := query.Order("created_at DESC").Find(&clients).Error; err != nil {
		return nil, fmt.Errorf("%w: %v", NewDomainError(ErrCodeRepositoryFailure, ErrRepositoryFailure), err)
	}
	return clients, nil
}

// Delete removes a client permanently.
func (r *GormRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&Client{})
	
	if result.Error != nil {
		return fmt.Errorf("%w: %v", NewDomainError(ErrCodeRepositoryFailure, ErrRepositoryFailure), result.Error)
	}
	
	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrClientNotFound)
	}
	
	return nil
}

// SetArchived updates the archived status of a client.
func (r *GormRepository) SetArchived(ctx context.Context, userID, id uuid.UUID, archived bool) error {
	result := r.db.WithContext(ctx).
		Model(&Client{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_archived", archived)
	
	if result.Error != nil {
		return fmt.Errorf("%w: %v", NewDomainError(ErrCodeRepositoryFailure, ErrRepositoryFailure), result.Error)
	}
	
	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrClientNotFound)
	}
	
	return nil
}

// FindByPhoneNumber finds a client by phone number for a specific user.
func (r *GormRepository) FindByPhoneNumber(ctx context.Context, userID uuid.UUID, phoneNumber string) (*Client, error) {
	var client Client
	normalized := normalizePhoneNumber(phoneNumber)
	
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND phone_number = ?", userID, normalized).
		First(&client).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrClientNotFound)
		}
		return nil, fmt.Errorf("%w: %v", NewDomainError(ErrCodeRepositoryFailure, ErrRepositoryFailure), err)
	}
	return &client, nil
}

