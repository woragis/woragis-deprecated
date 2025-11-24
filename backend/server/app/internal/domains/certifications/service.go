package certifications

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service orchestrates certification workflows.
type Service interface {
	CreateCertification(ctx context.Context, userID uuid.UUID, req CreateCertificationRequest) (*Certification, error)
	UpdateCertification(ctx context.Context, userID, certificationID uuid.UUID, req UpdateCertificationRequest) (*Certification, error)
	GetCertification(ctx context.Context, certificationID uuid.UUID, userID uuid.UUID) (*Certification, error)
	GetCertificationPublic(ctx context.Context, certificationID uuid.UUID) (*Certification, error)
	ListCertifications(ctx context.Context, filters ListCertificationsFilters) ([]Certification, error)
	ListFeaturedCertifications(ctx context.Context) ([]Certification, error)
	GetCertificationsBySkill(ctx context.Context, skillID uuid.UUID) ([]Certification, error)
	DeleteCertification(ctx context.Context, userID, certificationID uuid.UUID) error
	// Skill relationship methods
	AddCertificationSkill(ctx context.Context, userID, certificationID, skillID uuid.UUID) error
	RemoveCertificationSkill(ctx context.Context, userID, certificationID, skillID uuid.UUID) error
	GetCertificationSkills(ctx context.Context, certificationID uuid.UUID, userID uuid.UUID) ([]uuid.UUID, error)
	// Entity link methods (for projects, etc.)
	CreateCertificationEntityLink(ctx context.Context, userID, certificationID uuid.UUID, entityType EntityType, entityID uuid.UUID) error
	GetCertificationEntityLinks(ctx context.Context, certificationID uuid.UUID, userID uuid.UUID) ([]CertificationEntityLink, error)
	GetEntityCertifications(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]Certification, error)
	DeleteCertificationEntityLink(ctx context.Context, userID, linkID uuid.UUID) error
	DeleteCertificationEntityLinks(ctx context.Context, userID, certificationID uuid.UUID) error
}

type service struct {
	repo   Repository
	logger *slog.Logger
}

var _ Service = (*service)(nil)

// NewService constructs a Service.
func NewService(repo Repository, logger *slog.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

// Request payloads

type CreateCertificationRequest struct {
	Name            string                `json:"name"`
	Issuer          string                `json:"issuer"`
	IssueDate       time.Time             `json:"issueDate"`
	ExpiryDate      *time.Time            `json:"expiryDate,omitempty"`
	CredentialID    string                `json:"credentialId,omitempty"`
	VerificationURL string                `json:"verificationUrl,omitempty"`
	CertificateURL  string                `json:"certificateUrl,omitempty"`
	Description     string                `json:"description,omitempty"`
	Status          CertificationStatus  `json:"status,omitempty"`
	Category        CertificationCategory `json:"category"`
	Featured        bool                  `json:"featured,omitempty"`
	DisplayOrder    int                   `json:"displayOrder,omitempty"`
	SkillIDs        []uuid.UUID           `json:"skillIds,omitempty"` // Skills to link
}

type UpdateCertificationRequest struct {
	Name            *string                `json:"name,omitempty"`
	Issuer          *string                `json:"issuer,omitempty"`
	IssueDate       *time.Time              `json:"issueDate,omitempty"`
	ExpiryDate      *time.Time              `json:"expiryDate,omitempty"`
	CredentialID    *string                 `json:"credentialId,omitempty"`
	VerificationURL *string                 `json:"verificationUrl,omitempty"`
	CertificateURL  *string                 `json:"certificateUrl,omitempty"`
	Description     *string                 `json:"description,omitempty"`
	Status          *CertificationStatus    `json:"status,omitempty"`
	Category        *CertificationCategory  `json:"category,omitempty"`
	Featured        *bool                   `json:"featured,omitempty"`
	DisplayOrder    *int                    `json:"displayOrder,omitempty"`
	SkillIDs        []uuid.UUID             `json:"skillIds,omitempty"` // Skills to link (replaces existing)
}

type ListCertificationsFilters struct {
	UserID      *uuid.UUID
	Status      *CertificationStatus
	Category    *CertificationCategory
	Issuer      *string
	Featured    *bool
	ExpiringSoon *bool
	SkillID     *uuid.UUID
	Limit       int
	Offset      int
	OrderBy     string
	Order       string
}

// Service methods

func (s *service) CreateCertification(ctx context.Context, userID uuid.UUID, req CreateCertificationRequest) (*Certification, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	status := req.Status
	if status == "" {
		status = CertificationStatusActive
	}

	certification, err := NewCertification(userID, req.Name, req.Issuer, req.IssueDate)
	if err != nil {
		return nil, err
	}

	// Update optional fields
	if req.ExpiryDate != nil {
		certification.ExpiryDate = req.ExpiryDate
	}
	if req.CredentialID != "" {
		certification.CredentialID = req.CredentialID
	}
	if req.VerificationURL != "" {
		certification.VerificationURL = req.VerificationURL
	}
	if req.CertificateURL != "" {
		certification.CertificateURL = req.CertificateURL
	}
	if req.Description != "" {
		certification.Description = req.Description
	}
	if status != "" {
		if err := certification.SetStatus(status); err != nil {
			return nil, err
		}
	}
	if req.Category != "" {
		certification.Category = req.Category
	}
	certification.Featured = req.Featured
	certification.DisplayOrder = req.DisplayOrder

	if err := certification.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateCertification(ctx, certification); err != nil {
		return nil, err
	}

	// Link skills if provided
	if len(req.SkillIDs) > 0 {
		for _, skillID := range req.SkillIDs {
			if err := s.repo.AddCertificationSkill(ctx, certification.ID, skillID); err != nil {
				s.logger.Warn("Failed to link skill to certification",
					slog.String("certificationId", certification.ID.String()),
					slog.String("skillId", skillID.String()),
					slog.Any("error", err),
				)
			}
		}
		// Reload to get skills
		certification, err = s.repo.GetCertification(ctx, certification.ID, userID)
		if err != nil {
			return nil, err
		}
	}

	return certification, nil
}

func (s *service) UpdateCertification(ctx context.Context, userID, certificationID uuid.UUID, req UpdateCertificationRequest) (*Certification, error) {
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if certificationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}

	// Get existing certification
	certification, err := s.repo.GetCertification(ctx, certificationID, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	name := ""
	if req.Name != nil {
		name = *req.Name
	}
	issuer := ""
	if req.Issuer != nil {
		issuer = *req.Issuer
	}
	credentialID := ""
	if req.CredentialID != nil {
		credentialID = *req.CredentialID
	}
	verificationURL := ""
	if req.VerificationURL != nil {
		verificationURL = *req.VerificationURL
	}
	certificateURL := ""
	if req.CertificateURL != nil {
		certificateURL = *req.CertificateURL
	}
	description := ""
	if req.Description != nil {
		description = *req.Description
	}
	issueDate := certification.IssueDate
	if req.IssueDate != nil {
		issueDate = *req.IssueDate
	}
	expiryDate := certification.ExpiryDate
	if req.ExpiryDate != nil {
		expiryDate = req.ExpiryDate
	}
	category := certification.Category
	if req.Category != nil {
		category = *req.Category
	}

	if err := certification.UpdateDetails(name, issuer, credentialID, verificationURL, certificateURL, description, issueDate, expiryDate, category); err != nil {
		return nil, err
	}

	if req.Status != nil {
		if err := certification.SetStatus(*req.Status); err != nil {
			return nil, err
		}
	}
	if req.Featured != nil {
		certification.SetFeatured(*req.Featured)
	}
	if req.DisplayOrder != nil {
		certification.SetDisplayOrder(*req.DisplayOrder)
	}

	if err := s.repo.UpdateCertification(ctx, certification); err != nil {
		return nil, err
	}

	// Update skill relationships if provided
	if req.SkillIDs != nil {
		// Get existing skills
		existingSkills, err := s.repo.GetCertificationSkills(ctx, certificationID)
		if err != nil {
			return nil, err
		}

		// Remove skills that are not in the new list
		skillMap := make(map[uuid.UUID]bool)
		for _, skillID := range req.SkillIDs {
			skillMap[skillID] = true
		}

		for _, existingSkillID := range existingSkills {
			if !skillMap[existingSkillID] {
				if err := s.repo.RemoveCertificationSkill(ctx, certificationID, existingSkillID); err != nil {
					s.logger.Warn("Failed to remove skill from certification",
						slog.String("certificationId", certificationID.String()),
						slog.String("skillId", existingSkillID.String()),
						slog.Any("error", err),
					)
				}
			}
		}

		// Add new skills
		for _, skillID := range req.SkillIDs {
			if err := s.repo.AddCertificationSkill(ctx, certificationID, skillID); err != nil {
				s.logger.Warn("Failed to add skill to certification",
					slog.String("certificationId", certificationID.String()),
					slog.String("skillId", skillID.String()),
					slog.Any("error", err),
				)
			}
		}

		// Reload to get updated skills
		certification, err = s.repo.GetCertification(ctx, certificationID, userID)
		if err != nil {
			return nil, err
		}
	}

	return certification, nil
}

func (s *service) GetCertification(ctx context.Context, certificationID uuid.UUID, userID uuid.UUID) (*Certification, error) {
	if certificationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}

	return s.repo.GetCertification(ctx, certificationID, userID)
}

func (s *service) GetCertificationPublic(ctx context.Context, certificationID uuid.UUID) (*Certification, error) {
	if certificationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}

	return s.repo.GetCertificationPublic(ctx, certificationID)
}

func (s *service) ListCertifications(ctx context.Context, filters ListCertificationsFilters) ([]Certification, error) {
	repoFilters := CertificationFilters{
		UserID:      filters.UserID,
		Status:      filters.Status,
		Category:    filters.Category,
		Issuer:      filters.Issuer,
		Featured:    filters.Featured,
		ExpiringSoon: filters.ExpiringSoon,
		SkillID:     filters.SkillID,
		Limit:       filters.Limit,
		Offset:      filters.Offset,
		OrderBy:     filters.OrderBy,
		Order:       filters.Order,
	}

	return s.repo.ListCertifications(ctx, repoFilters)
}

func (s *service) ListFeaturedCertifications(ctx context.Context) ([]Certification, error) {
	return s.repo.ListFeaturedCertifications(ctx)
}

func (s *service) GetCertificationsBySkill(ctx context.Context, skillID uuid.UUID) ([]Certification, error) {
	if skillID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, "skill id cannot be empty")
	}

	return s.repo.GetCertificationsBySkill(ctx, skillID)
}

func (s *service) DeleteCertification(ctx context.Context, userID, certificationID uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if certificationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}

	return s.repo.DeleteCertification(ctx, certificationID, userID)
}

// Skill relationship methods

func (s *service) AddCertificationSkill(ctx context.Context, userID, certificationID, skillID uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if certificationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}
	if skillID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, "skill id cannot be empty")
	}

	// Verify ownership
	_, err := s.repo.GetCertification(ctx, certificationID, userID)
	if err != nil {
		return err
	}

	return s.repo.AddCertificationSkill(ctx, certificationID, skillID)
}

func (s *service) RemoveCertificationSkill(ctx context.Context, userID, certificationID, skillID uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if certificationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}
	if skillID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, "skill id cannot be empty")
	}

	// Verify ownership
	_, err := s.repo.GetCertification(ctx, certificationID, userID)
	if err != nil {
		return err
	}

	return s.repo.RemoveCertificationSkill(ctx, certificationID, skillID)
}

func (s *service) GetCertificationSkills(ctx context.Context, certificationID uuid.UUID, userID uuid.UUID) ([]uuid.UUID, error) {
	if certificationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	// Verify ownership
	_, err := s.repo.GetCertification(ctx, certificationID, userID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetCertificationSkills(ctx, certificationID)
}

// Entity link methods

func (s *service) CreateCertificationEntityLink(ctx context.Context, userID, certificationID uuid.UUID, entityType EntityType, entityID uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if certificationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}
	if entityID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyEntityID)
	}

	// Verify ownership
	_, err := s.repo.GetCertification(ctx, certificationID, userID)
	if err != nil {
		return err
	}

	link, err := NewCertificationEntityLink(certificationID, entityType, entityID)
	if err != nil {
		return err
	}

	return s.repo.CreateCertificationEntityLink(ctx, link)
}

func (s *service) GetCertificationEntityLinks(ctx context.Context, certificationID uuid.UUID, userID uuid.UUID) ([]CertificationEntityLink, error) {
	if certificationID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}
	if userID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	// Verify ownership
	_, err := s.repo.GetCertification(ctx, certificationID, userID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetCertificationEntityLinks(ctx, certificationID)
}

func (s *service) GetEntityCertifications(ctx context.Context, entityType EntityType, entityID uuid.UUID) ([]Certification, error) {
	if entityID == uuid.Nil {
		return nil, NewDomainError(ErrCodeInvalidPayload, ErrEmptyEntityID)
	}

	return s.repo.GetEntityCertifications(ctx, entityType, entityID)
}

func (s *service) DeleteCertificationEntityLink(ctx context.Context, userID, linkID uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if linkID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyLinkID)
	}

	// Get link to find certification ID
	link, err := s.repo.GetCertificationEntityLink(ctx, linkID)
	if err != nil {
		return err
	}

	// Verify ownership through certification
	_, err = s.repo.GetCertification(ctx, link.CertificationID, userID)
	if err != nil {
		return err
	}

	return s.repo.DeleteCertificationEntityLink(ctx, linkID)
}

func (s *service) DeleteCertificationEntityLinks(ctx context.Context, userID, certificationID uuid.UUID) error {
	if userID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if certificationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}

	// Verify ownership
	_, err := s.repo.GetCertification(ctx, certificationID, userID)
	if err != nil {
		return err
	}

	return s.repo.DeleteCertificationEntityLinks(ctx, certificationID)
}

