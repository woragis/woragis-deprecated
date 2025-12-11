package resumes

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ResumeMetrics holds calculated metrics for a resume.
type ResumeMetrics struct {
	ApplicationsUsed int
	InterviewCount   int
	OfferCount       int
	InterviewRate    float64
	OfferRate        float64
}

// Repository defines persistence operations for resumes.
type Repository interface {
	CreateResume(ctx context.Context, resume *Resume) error
	UpdateResume(ctx context.Context, resume *Resume) error
	DeleteResume(ctx context.Context, resumeID uuid.UUID, userID uuid.UUID) error
	GetResume(ctx context.Context, resumeID uuid.UUID, userID uuid.UUID) (*Resume, error)
	ListResumes(ctx context.Context, userID uuid.UUID) ([]Resume, error)
	ListResumesByTags(ctx context.Context, userID uuid.UUID, tags []string) ([]Resume, error)
	GetMainResume(ctx context.Context, userID uuid.UUID) (*Resume, error)
	GetFeaturedResume(ctx context.Context, userID uuid.UUID) (*Resume, error)
	UnmarkAllAsMain(ctx context.Context, userID uuid.UUID) error
	CalculateResumeMetrics(ctx context.Context, resumeID uuid.UUID) (*ResumeMetrics, error)
	UpdateResumeMetrics(ctx context.Context, resumeID uuid.UUID, metrics *ResumeMetrics) error
}

// gormRepository implements Repository using GORM.
type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GORM-based repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

// CreateResume creates a new resume.
func (r *gormRepository) CreateResume(ctx context.Context, resume *Resume) error {
	return r.db.WithContext(ctx).Create(resume).Error
}

// UpdateResume updates an existing resume.
func (r *gormRepository) UpdateResume(ctx context.Context, resume *Resume) error {
	return r.db.WithContext(ctx).Save(resume).Error
}

// DeleteResume deletes a resume.
func (r *gormRepository) DeleteResume(ctx context.Context, resumeID uuid.UUID, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", resumeID, userID).
		Delete(&Resume{})
	
	if result.Error != nil {
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrResumeNotFound)
	}
	
	return nil
}

// GetResume retrieves a resume by ID.
func (r *gormRepository) GetResume(ctx context.Context, resumeID uuid.UUID, userID uuid.UUID) (*Resume, error) {
	var resume Resume
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", resumeID, userID).
		First(&resume).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrResumeNotFound)
		}
		return nil, err
	}
	
	return &resume, nil
}

// ListResumes lists all resumes for a user.
func (r *gormRepository) ListResumes(ctx context.Context, userID uuid.UUID) ([]Resume, error) {
	var resumes []Resume
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_main DESC, is_featured DESC, created_at DESC").
		Find(&resumes).Error
	
	return resumes, err
}

// ListResumesByTags lists resumes filtered by tags (resume must have at least one matching tag).
func (r *gormRepository) ListResumesByTags(ctx context.Context, userID uuid.UUID, tags []string) ([]Resume, error) {
	if len(tags) == 0 {
		return r.ListResumes(ctx, userID)
	}
	
	var resumes []Resume
	query := r.db.WithContext(ctx).
		Where("user_id = ?", userID)
	
	// For each tag, check if it's contained in the tags JSON array
	for _, tag := range tags {
		normalizedTag := strings.ToLower(strings.TrimSpace(tag))
		if normalizedTag != "" {
			query = query.Where("tags @> ?", `["`+normalizedTag+`"]`)
		}
	}
	
	err := query.
		Order("is_main DESC, is_featured DESC, created_at DESC").
		Find(&resumes).Error
	
	return resumes, err
}

// GetMainResume retrieves the main resume for a user.
func (r *gormRepository) GetMainResume(ctx context.Context, userID uuid.UUID) (*Resume, error) {
	var resume Resume
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_main = ?", userID, true).
		First(&resume).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrNoMainResume)
		}
		return nil, err
	}
	
	return &resume, nil
}

// GetFeaturedResume retrieves a featured resume for a user (fallback if no main).
func (r *gormRepository) GetFeaturedResume(ctx context.Context, userID uuid.UUID) (*Resume, error) {
	var resume Resume
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_featured = ?", userID, true).
		Order("created_at DESC").
		First(&resume).Error
	
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, NewDomainError(ErrCodeNotFound, ErrResumeNotFound)
		}
		return nil, err
	}
	
	return &resume, nil
}

// UnmarkAllAsMain unmarks all resumes as main for a user.
func (r *gormRepository) UnmarkAllAsMain(ctx context.Context, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&Resume{}).
		Where("user_id = ?", userID).
		Update("is_main", false).Error
}

// CalculateResumeMetrics calculates metrics for a resume.
func (r *gormRepository) CalculateResumeMetrics(ctx context.Context, resumeID uuid.UUID) (*ResumeMetrics, error) {
	var metrics ResumeMetrics
	
	// Count total applications using this resume
	var count int64
	err := r.db.WithContext(ctx).
		Table("job_applications").
		Where("resume_id = ?", resumeID).
		Count(&count).Error
	if err != nil {
		return nil, err
	}
	metrics.ApplicationsUsed = int(count)
	
	if metrics.ApplicationsUsed == 0 {
		metrics.InterviewRate = 0.0
		metrics.OfferRate = 0.0
		return &metrics, nil
	}
	
	// Count applications with completed interviews
	// A completed interview is one where completed_date is not null
	var interviewCount int64
	err = r.db.WithContext(ctx).
		Table("job_applications").
		Where("resume_id = ? AND id IN (SELECT DISTINCT job_application_id FROM job_application_stages WHERE completed_date IS NOT NULL)", resumeID).
		Count(&interviewCount).Error
	if err != nil {
		return nil, err
	}
	metrics.InterviewCount = int(interviewCount)
	
	// Calculate interview rate
	if metrics.ApplicationsUsed > 0 {
		metrics.InterviewRate = (float64(metrics.InterviewCount) / float64(metrics.ApplicationsUsed)) * 100.0
	}
	
	// Count applications with offers
	// An offer is: status = 'accepted' OR has a response with response_type = 'offer'
	var offerCount int64
	err = r.db.WithContext(ctx).
		Table("job_applications").
		Where("resume_id = ? AND (status = 'accepted' OR id IN (SELECT DISTINCT job_application_id FROM job_application_responses WHERE response_type = 'offer'))", resumeID).
		Count(&offerCount).Error
	if err != nil {
		return nil, err
	}
	metrics.OfferCount = int(offerCount)
	
	// Calculate offer rate
	if metrics.ApplicationsUsed > 0 {
		metrics.OfferRate = (float64(metrics.OfferCount) / float64(metrics.ApplicationsUsed)) * 100.0
	}
	
	return &metrics, nil
}

// UpdateResumeMetrics updates the cached metrics for a resume.
func (r *gormRepository) UpdateResumeMetrics(ctx context.Context, resumeID uuid.UUID, metrics *ResumeMetrics) error {
	return r.db.WithContext(ctx).
		Model(&Resume{}).
		Where("id = ?", resumeID).
		Updates(map[string]interface{}{
			"applications_used": metrics.ApplicationsUsed,
			"interview_rate":    metrics.InterviewRate,
			"offer_rate":        metrics.OfferRate,
		}).Error
}

