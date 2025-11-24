package certifications

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// CertificationStatus represents the status of a certification.
type CertificationStatus string

const (
	CertificationStatusActive   CertificationStatus = "active"
	CertificationStatusExpired  CertificationStatus = "expired"
	CertificationStatusRevoked  CertificationStatus = "revoked"
	CertificationStatusRenewed  CertificationStatus = "renewed"
)

// CertificationCategory represents the category of a certification.
type CertificationCategory string

const (
	CertificationCategoryCloud        CertificationCategory = "cloud"
	CertificationCategorySecurity     CertificationCategory = "security"
	CertificationCategoryProgramming  CertificationCategory = "programming"
	CertificationCategoryDatabase     CertificationCategory = "database"
	CertificationCategoryDevOps       CertificationCategory = "devops"
	CertificationCategoryArchitecture CertificationCategory = "architecture"
	CertificationCategoryOther        CertificationCategory = "other"
)

// Certification represents a professional certification or credential.
type Certification struct {
	ID              uuid.UUID            `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID            `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Name            string               `gorm:"column:name;size:255;not null" json:"name"`
	Issuer          string               `gorm:"column:issuer;size:255;not null" json:"issuer"`
	IssueDate       time.Time            `gorm:"column:issue_date;type:date;not null;index" json:"issueDate"`
	ExpiryDate      *time.Time           `gorm:"column:expiry_date;type:date;index" json:"expiryDate,omitempty"`
	CredentialID    string               `gorm:"column:credential_id;size:255" json:"credentialId,omitempty"`
	VerificationURL string               `gorm:"column:verification_url;size:512" json:"verificationUrl,omitempty"`
	CertificateURL  string               `gorm:"column:certificate_url;size:512" json:"certificateUrl,omitempty"` // Image/PDF URL
	Description     string               `gorm:"column:description;type:text" json:"description,omitempty"`
	Status          CertificationStatus  `gorm:"column:status;type:varchar(32);not null;default:'active';index" json:"status"`
	Category        CertificationCategory `gorm:"column:category;type:varchar(32);not null;index" json:"category"`
	Featured        bool                 `gorm:"column:featured;not null;default:false;index" json:"featured"`
	DisplayOrder    int                  `gorm:"column:display_order;not null;default:0;index" json:"displayOrder"`
	CreatedAt       time.Time            `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time            `gorm:"column:updated_at" json:"updatedAt"`
	
	// Relationships
	Skills []Skill `gorm:"many2many:certification_skills;" json:"skills,omitempty"`
}

// TableName specifies the table name for Certification.
func (Certification) TableName() string {
	return "certifications"
}

// Skill represents a skill linked to a certification (for many-to-many relationship).
// This is a reference type, not the full skill entity.
type Skill struct {
	ID   uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name string     `gorm:"column:name;size:120" json:"name"`
}

// TableName specifies the table name for the skills join table reference.
func (Skill) TableName() string {
	return "skills"
}

// CertificationSkill represents the many-to-many relationship between certifications and skills.
type CertificationSkill struct {
	CertificationID uuid.UUID `gorm:"column:certification_id;type:uuid;primaryKey;index:idx_cert_skill" json:"certificationId"`
	SkillID         uuid.UUID `gorm:"column:skill_id;type:uuid;primaryKey;index:idx_cert_skill" json:"skillId"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"createdAt"`
}

// TableName specifies the table name for CertificationSkill.
func (CertificationSkill) TableName() string {
	return "certification_skills"
}

// EntityType represents the type of entity being linked to a certification.
type EntityType string

const (
	EntityTypeProject EntityType = "project"
)

// CertificationEntityLink represents the relationship between a certification and an entity (e.g., project).
type CertificationEntityLink struct {
	ID             uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	CertificationID uuid.UUID `gorm:"column:certification_id;type:uuid;not null;index:idx_cert_entity_link" json:"certificationId"`
	EntityType     EntityType `gorm:"column:entity_type;type:varchar(50);not null;index:idx_cert_entity_link" json:"entityType"`
	EntityID       uuid.UUID  `gorm:"column:entity_id;type:uuid;not null;index:idx_cert_entity_link" json:"entityId"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for CertificationEntityLink.
func (CertificationEntityLink) TableName() string {
	return "certification_entity_links"
}

// NewCertificationEntityLink creates a new link between a certification and an entity.
func NewCertificationEntityLink(certificationID uuid.UUID, entityType EntityType, entityID uuid.UUID) (*CertificationEntityLink, error) {
	link := &CertificationEntityLink{
		ID:             uuid.New(),
		CertificationID: certificationID,
		EntityType:     entityType,
		EntityID:       entityID,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
	return link, link.Validate()
}

// Validate ensures certification entity link invariants hold.
func (l *CertificationEntityLink) Validate() error {
	if l == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilLink)
	}
	if l.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyLinkID)
	}
	if l.CertificationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}
	if l.EntityID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyEntityID)
	}
	if !isValidEntityType(l.EntityType) {
		return NewDomainError(ErrCodeInvalidEntityType, ErrUnsupportedEntityType)
	}
	return nil
}

// NewCertification creates a new certification entity.
func NewCertification(userID uuid.UUID, name, issuer string, issueDate time.Time) (*Certification, error) {
	cert := &Certification{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      strings.TrimSpace(name),
		Issuer:    strings.TrimSpace(issuer),
		IssueDate: issueDate,
		Status:    CertificationStatusActive,
		Category:  CertificationCategoryOther,
		Featured:  false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	return cert, cert.Validate()
}

// Validate ensures certification invariants hold.
func (c *Certification) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilCertification)
	}
	if c.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCertificationID)
	}
	if c.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if strings.TrimSpace(c.Name) == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyName)
	}
	if strings.TrimSpace(c.Issuer) == "" {
		return NewDomainError(ErrCodeInvalidIssuer, ErrEmptyIssuer)
	}
	if c.IssueDate.IsZero() {
		return NewDomainError(ErrCodeInvalidDate, ErrEmptyIssueDate)
	}
	if c.ExpiryDate != nil && !c.ExpiryDate.IsZero() && c.ExpiryDate.Before(c.IssueDate) {
		return NewDomainError(ErrCodeInvalidDate, ErrExpiryBeforeIssue)
	}
	if !isValidStatus(c.Status) {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}
	if !isValidCategory(c.Category) {
		return NewDomainError(ErrCodeInvalidCategory, ErrUnsupportedCategory)
	}
	return nil
}

// UpdateDetails updates certification details.
func (c *Certification) UpdateDetails(name, issuer, credentialID, verificationURL, certificateURL, description string, issueDate time.Time, expiryDate *time.Time, category CertificationCategory) error {
	if name != "" {
		c.Name = strings.TrimSpace(name)
	}
	if issuer != "" {
		c.Issuer = strings.TrimSpace(issuer)
	}
	if credentialID != "" {
		c.CredentialID = strings.TrimSpace(credentialID)
	}
	if verificationURL != "" {
		c.VerificationURL = strings.TrimSpace(verificationURL)
	}
	if certificateURL != "" {
		c.CertificateURL = strings.TrimSpace(certificateURL)
	}
	if description != "" {
		c.Description = strings.TrimSpace(description)
	}
	if !issueDate.IsZero() {
		c.IssueDate = issueDate
	}
	if expiryDate != nil {
		c.ExpiryDate = expiryDate
	}
	if category != "" {
		c.Category = category
	}
	c.UpdatedAt = time.Now().UTC()
	return c.Validate()
}

// SetStatus updates the certification status.
func (c *Certification) SetStatus(status CertificationStatus) error {
	if !isValidStatus(status) {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}
	c.Status = status
	c.UpdatedAt = time.Now().UTC()
	return nil
}

// SetFeatured updates the featured flag.
func (c *Certification) SetFeatured(featured bool) {
	c.Featured = featured
	c.UpdatedAt = time.Now().UTC()
}

// SetDisplayOrder updates the display order.
func (c *Certification) SetDisplayOrder(order int) {
	c.DisplayOrder = order
	c.UpdatedAt = time.Now().UTC()
}

// IsExpired checks if the certification is expired.
func (c *Certification) IsExpired() bool {
	if c.ExpiryDate == nil || c.ExpiryDate.IsZero() {
		return false // No expiry date means it doesn't expire
	}
	return time.Now().After(*c.ExpiryDate)
}

// DaysUntilExpiry returns the number of days until expiry (negative if expired).
func (c *Certification) DaysUntilExpiry() *int {
	if c.ExpiryDate == nil || c.ExpiryDate.IsZero() {
		return nil
	}
	days := int(time.Until(*c.ExpiryDate).Hours() / 24)
	return &days
}

// Validation helpers

func isValidStatus(s CertificationStatus) bool {
	switch s {
	case CertificationStatusActive, CertificationStatusExpired, CertificationStatusRevoked, CertificationStatusRenewed:
		return true
	}
	return false
}

func isValidCategory(c CertificationCategory) bool {
	switch c {
	case CertificationCategoryCloud, CertificationCategorySecurity, CertificationCategoryProgramming,
		CertificationCategoryDatabase, CertificationCategoryDevOps, CertificationCategoryArchitecture,
		CertificationCategoryOther:
		return true
	}
	return false
}

func isValidEntityType(et EntityType) bool {
	switch et {
	case EntityTypeProject:
		return true
	}
	return false
}

