package certifications

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Repository defines persistence operations for certifications.
type Repository interface {
	CreateCertification(ctx context.Context, certification *Certification) error
	UpdateCertification(ctx context.Context, certification *Certification) error
	GetCertification(ctx context.Context, certificationID uuid.UUID, userID uuid.UUID) (*Certification, error)
	GetCertificationPublic(ctx context.Context, certificationID uuid.UUID) (*Certification, error)
	ListCertifications(ctx context.Context, filters CertificationFilters) ([]Certification, error)
	ListFeaturedCertifications(ctx context.Context) ([]Certification, error)
	GetCertificationsBySkill(ctx context.Context, skillID uuid.UUID) ([]Certification, error)
	DeleteCertification(ctx context.Context, certificationID uuid.UUID, userID uuid.UUID) error
	// Skill relationship methods
	AddCertificationSkill(ctx context.Context, certificationID, skillID uuid.UUID) error
	RemoveCertificationSkill(ctx context.Context, certificationID, skillID uuid.UUID) error
	GetCertificationSkills(ctx context.Context, certificationID uuid.UUID) ([]uuid.UUID, error)
	// Entity link methods (for projects, etc.)
	CreateCertificationEntityLink(ctx context.Context, link *CertificationEntityLink) error
	GetCertificationEntityLinks(ctx context.Context, certificationID uuid.UUID) ([]CertificationEntityLink, error)
	GetEntityCertifications(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]Certification, error)
	DeleteCertificationEntityLink(ctx context.Context, linkID uuid.UUID) error
	DeleteCertificationEntityLinks(ctx context.Context, certificationID uuid.UUID) error
}

// CertificationFilters represents filtering options for listing certifications.
type CertificationFilters struct {
	UserID     *uuid.UUID
	Status     *CertificationStatus
	Category   *CertificationCategory
	Issuer     *string
	Featured   *bool
	ExpiringSoon *bool // Certifications expiring within 90 days
	SkillID    *uuid.UUID
	Limit      int
	Offset     int
	OrderBy    string // "issue_date", "expiry_date", "display_order", "created_at"
	Order      string // "asc", "desc"
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository returns a GORM-backed repository.
func NewGormRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateCertification(ctx context.Context, certification *Certification) error {
	if certification == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilCertification)
	}

	if err := certification.Validate(); err != nil {
		return err
	}

	now := time.Now().UTC()
	certification.CreatedAt = now
	certification.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(certification).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" { // unique_violation
				return NewDomainError(ErrCodeConflict, ErrCertificationAlreadyExists)
			}
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) UpdateCertification(ctx context.Context, certification *Certification) error {
	if certification == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilCertification)
	}

	if certification.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}

	if err := certification.Validate(); err != nil {
		return err
	}

	certification.UpdatedAt = time.Now().UTC()

	result := r.db.WithContext(ctx).Model(&Certification{}).
		Where("id = ?", certification.ID).
		Updates(certification)

	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToUpdate)
	}

	if result.RowsAffected == 0 {
		return NewDomainError(ErrCodeNotFound, ErrCertificationNotFound)
	}

	return nil
}

func (r *gormRepository) GetCertification(ctx context.Context, certificationID uuid.UUID, userID uuid.UUID) (*Certification, error) {
	if certificationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}

	var certification Certification
	err := r.db.WithContext(ctx).
		Preload("Skills").
		Where("id = ? AND user_id = ?", certificationID, userID).
		First(&certification).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrCertificationNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &certification, nil
}

func (r *gormRepository) GetCertificationPublic(ctx context.Context, certificationID uuid.UUID) (*Certification, error) {
	if certificationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}

	var certification Certification
	err := r.db.WithContext(ctx).
		Preload("Skills").
		Where("id = ? AND status = ?", certificationID, CertificationStatusActive).
		First(&certification).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, ErrCertificationNotFound)
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &certification, nil
}

func (r *gormRepository) ListCertifications(ctx context.Context, filters CertificationFilters) ([]Certification, error) {
	query := r.db.WithContext(ctx).Model(&Certification{}).Preload("Skills")

	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}

	if filters.Category != nil {
		query = query.Where("category = ?", *filters.Category)
	}

	if filters.Issuer != nil {
		query = query.Where("issuer ILIKE ?", "%"+*filters.Issuer+"%")
	}

	if filters.Featured != nil {
		query = query.Where("featured = ?", *filters.Featured)
	}

	if filters.ExpiringSoon != nil && *filters.ExpiringSoon {
		// Certifications expiring within 90 days
		expiryThreshold := time.Now().UTC().AddDate(0, 0, 90)
		query = query.Where("expiry_date IS NOT NULL AND expiry_date <= ? AND expiry_date > ?", expiryThreshold, time.Now().UTC())
	}

	if filters.SkillID != nil {
		query = query.Joins("INNER JOIN certification_skills ON certifications.id = certification_skills.certification_id").
			Where("certification_skills.skill_id = ?", *filters.SkillID)
	}

	// Default ordering
	orderBy := filters.OrderBy
	if orderBy == "" {
		orderBy = "issue_date"
	}
	order := filters.Order
	if order == "" {
		order = "desc"
	}
	query = query.Order(orderBy + " " + order)

	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	var certifications []Certification
	if err := query.Find(&certifications).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return certifications, nil
}

func (r *gormRepository) ListFeaturedCertifications(ctx context.Context) ([]Certification, error) {
	var certifications []Certification
	err := r.db.WithContext(ctx).
		Preload("Skills").
		Where("featured = ? AND status = ?", true, CertificationStatusActive).
		Order("display_order ASC, issue_date DESC").
		Find(&certifications).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return certifications, nil
}

func (r *gormRepository) GetCertificationsBySkill(ctx context.Context, skillID uuid.UUID) ([]Certification, error) {
	var certifications []Certification
	err := r.db.WithContext(ctx).
		Preload("Skills").
		Joins("INNER JOIN certification_skills ON certifications.id = certification_skills.certification_id").
		Where("certification_skills.skill_id = ? AND certifications.status = ?", skillID, CertificationStatusActive).
		Order("issue_date DESC").
		Find(&certifications).Error

	if err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return certifications, nil
}

func (r *gormRepository) DeleteCertification(ctx context.Context, certificationID uuid.UUID, userID uuid.UUID) error {
	if certificationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}

	// Verify ownership
	var certification Certification
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", certificationID, userID).
		First(&certification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewDomainError(ErrCodeNotFound, ErrCertificationNotFound)
		}
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	// Delete skill relationships first
	if err := r.db.WithContext(ctx).
		Where("certification_id = ?", certificationID).
		Delete(&CertificationSkill{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	// Delete entity links
	if err := r.db.WithContext(ctx).
		Where("certification_id = ?", certificationID).
		Delete(&CertificationEntityLink{}).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	// Delete certification
	if err := r.db.WithContext(ctx).Delete(&certification).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

// Skill relationship methods

func (r *gormRepository) AddCertificationSkill(ctx context.Context, certificationID, skillID uuid.UUID) error {
	if certificationID == uuid.Nil || skillID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, "certification or skill id cannot be empty")
	}

	// Check if relationship already exists
	var existing CertificationSkill
	err := r.db.WithContext(ctx).
		Where("certification_id = ? AND skill_id = ?", certificationID, skillID).
		First(&existing).Error

	if err == nil {
		// Relationship already exists
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	// Create new relationship
	certSkill := CertificationSkill{
		CertificationID: certificationID,
		SkillID:         skillID,
		CreatedAt:       time.Now().UTC(),
	}

	if err := r.db.WithContext(ctx).Create(&certSkill).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) RemoveCertificationSkill(ctx context.Context, certificationID, skillID uuid.UUID) error {
	if certificationID == uuid.Nil || skillID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, "certification or skill id cannot be empty")
	}

	result := r.db.WithContext(ctx).
		Where("certification_id = ? AND skill_id = ?", certificationID, skillID).
		Delete(&CertificationSkill{})

	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) GetCertificationSkills(ctx context.Context, certificationID uuid.UUID) ([]uuid.UUID, error) {
	if certificationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}

	var certSkills []CertificationSkill
	if err := r.db.WithContext(ctx).
		Where("certification_id = ?", certificationID).
		Find(&certSkills).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	skillIDs := make([]uuid.UUID, len(certSkills))
	for i, cs := range certSkills {
		skillIDs[i] = cs.SkillID
	}

	return skillIDs, nil
}

// Entity link methods

func (r *gormRepository) CreateCertificationEntityLink(ctx context.Context, link *CertificationEntityLink) error {
	if link == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilLink)
	}

	if err := link.Validate(); err != nil {
		return err
	}

	now := time.Now().UTC()
	link.CreatedAt = now
	link.UpdatedAt = now

	if err := r.db.WithContext(ctx).Create(link).Error; err != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) GetCertificationEntityLinks(ctx context.Context, certificationID uuid.UUID) ([]CertificationEntityLink, error) {
	if certificationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}

	var links []CertificationEntityLink
	if err := r.db.WithContext(ctx).Where("certification_id = ?", certificationID).Find(&links).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return links, nil
}

func (r *gormRepository) GetEntityCertifications(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]Certification, error) {
	if entityID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyEntityID)
	}

	var certifications []Certification
	if err := r.db.WithContext(ctx).
		Preload("Skills").
		Table("certifications").
		Joins("INNER JOIN certification_entity_links ON certifications.id = certification_entity_links.certification_id").
		Where("certification_entity_links.entity_type = ? AND certification_entity_links.entity_id = ?", entityType, entityID).
		Find(&certifications).Error; err != nil {
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return certifications, nil
}

func (r *gormRepository) GetCertificationEntityLink(ctx context.Context, linkID uuid.UUID) (*CertificationEntityLink, error) {
	if linkID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyLinkID)
	}

	var link CertificationEntityLink
	err := r.db.WithContext(ctx).Where("id = ?", linkID).First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewDomainError(ErrCodeNotFound, "certification entity link not found")
		}
		return nil, NewDomainError(ErrCodeRepositoryFailure, ErrUnableToFetch)
	}

	return &link, nil
}

func (r *gormRepository) DeleteCertificationEntityLink(ctx context.Context, linkID uuid.UUID) error {
	if linkID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyLinkID)
	}

	result := r.db.WithContext(ctx).Where("id = ?", linkID).Delete(&CertificationEntityLink{})
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

func (r *gormRepository) DeleteCertificationEntityLinks(ctx context.Context, certificationID uuid.UUID) error {
	if certificationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}

	result := r.db.WithContext(ctx).Where("certification_id = ?", certificationID).Delete(&CertificationEntityLink{})
	if result.Error != nil {
		return NewDomainError(ErrCodeRepositoryFailure, ErrUnableToPersist)
	}

	return nil
}

