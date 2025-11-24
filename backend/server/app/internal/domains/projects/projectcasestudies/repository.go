package projectcasestudies

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines persistence operations for project case studies.
type Repository interface {
	CreateCaseStudy(ctx context.Context, caseStudy *ProjectCaseStudy) error
	UpdateCaseStudy(ctx context.Context, caseStudy *ProjectCaseStudy) error
	GetCaseStudy(ctx context.Context, caseStudyID uuid.UUID, userID uuid.UUID) (*ProjectCaseStudy, error)
	GetCaseStudyByProjectID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*ProjectCaseStudy, error)
	GetCaseStudyByProjectIDPublic(ctx context.Context, projectID uuid.UUID) (*ProjectCaseStudy, error)
	ListCaseStudies(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectCaseStudy, error)
	DeleteCaseStudy(ctx context.Context, caseStudyID uuid.UUID, userID uuid.UUID) error
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateCaseStudy(ctx context.Context, caseStudy *ProjectCaseStudy) error {
	if err := caseStudy.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Create(caseStudy).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}
	return nil
}

func (r *gormRepository) UpdateCaseStudy(ctx context.Context, caseStudy *ProjectCaseStudy) error {
	if err := caseStudy.Validate(); err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Save(caseStudy).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}

func (r *gormRepository) GetCaseStudy(ctx context.Context, caseStudyID uuid.UUID, userID uuid.UUID) (*ProjectCaseStudy, error) {
	var caseStudy ProjectCaseStudy
	err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = project_case_studies.project_id").
		Where("project_case_studies.id = ? AND projects.user_id = ?", caseStudyID, userID).
		First(&caseStudy).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrCaseStudyNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &caseStudy, nil
}

func (r *gormRepository) GetCaseStudyByProjectID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*ProjectCaseStudy, error) {
	var caseStudy ProjectCaseStudy
	err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = project_case_studies.project_id").
		Where("project_case_studies.project_id = ? AND projects.user_id = ?", projectID, userID).
		First(&caseStudy).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrCaseStudyNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &caseStudy, nil
}

func (r *gormRepository) GetCaseStudyByProjectIDPublic(ctx context.Context, projectID uuid.UUID) (*ProjectCaseStudy, error) {
	var caseStudy ProjectCaseStudy
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		First(&caseStudy).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrCaseStudyNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return &caseStudy, nil
}

func (r *gormRepository) ListCaseStudies(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]ProjectCaseStudy, error) {
	var caseStudies []ProjectCaseStudy
	err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = project_case_studies.project_id").
		Where("project_case_studies.project_id = ? AND projects.user_id = ?", projectID, userID).
		Order("created_at DESC").
		Find(&caseStudies).Error
	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}
	return caseStudies, nil
}

func (r *gormRepository) DeleteCaseStudy(ctx context.Context, caseStudyID uuid.UUID, userID uuid.UUID) error {
	var caseStudy ProjectCaseStudy
	if err := r.db.WithContext(ctx).
		Joins("JOIN projects ON projects.id = project_case_studies.project_id").
		Where("project_case_studies.id = ? AND projects.user_id = ?", caseStudyID, userID).
		First(&caseStudy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return NewDomainError(ErrCodeNotFound, ErrCaseStudyNotFound)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	if err := r.db.WithContext(ctx).Delete(&caseStudy).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}
	return nil
}
