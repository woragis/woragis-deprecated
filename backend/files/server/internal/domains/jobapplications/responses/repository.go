package responses

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for responses.
type Repository interface {
	CreateResponse(ctx context.Context, response *Response) error
	UpdateResponse(ctx context.Context, response *Response) error
	GetResponse(ctx context.Context, responseID uuid.UUID) (*Response, error)
	ListResponses(ctx context.Context, filters ResponseFilters) ([]Response, error)
	DeleteResponse(ctx context.Context, responseID uuid.UUID) error
	GetResponsesByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]Response, error)
}

// ResponseFilters represents filtering options for listing responses.
type ResponseFilters struct {
	JobApplicationID *uuid.UUID
	ResponseType     *ResponseType
	Limit            int
	Offset           int
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateResponse(ctx context.Context, response *Response) error {
	if err := response.Validate(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(response).Error
}

func (r *gormRepository) UpdateResponse(ctx context.Context, response *Response) error {
	if err := response.Validate(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(response).Error
}

func (r *gormRepository) GetResponse(ctx context.Context, responseID uuid.UUID) (*Response, error) {
	var response Response
	if err := r.db.WithContext(ctx).Where("id = ?", responseID).First(&response).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrResponseNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &response, nil
}

func (r *gormRepository) ListResponses(ctx context.Context, filters ResponseFilters) ([]Response, error) {
	var responses []Response
	query := r.db.WithContext(ctx).Model(&Response{})

	if filters.JobApplicationID != nil {
		query = query.Where("job_application_id = ?", *filters.JobApplicationID)
	}
	if filters.ResponseType != nil {
		query = query.Where("response_type = ?", *filters.ResponseType)
	}

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	query = query.Order("response_date DESC")

	if err := query.Find(&responses).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return responses, nil
}

func (r *gormRepository) DeleteResponse(ctx context.Context, responseID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&Response{}, responseID)
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrResponseNotFound)
	}
	return nil
}

func (r *gormRepository) GetResponsesByApplicationID(ctx context.Context, applicationID uuid.UUID) ([]Response, error) {
	return r.ListResponses(ctx, ResponseFilters{
		JobApplicationID: &applicationID,
	})
}

