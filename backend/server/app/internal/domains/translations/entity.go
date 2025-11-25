package translations

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// EntityType represents the type of entity being translated.
type EntityType string

const (
	EntityTypeTestimonial      EntityType = "testimonial"
	EntityTypePost             EntityType = "post"
	EntityTypeProject          EntityType = "project"
	EntityTypeCaseStudy        EntityType = "case_study"
	EntityTypeProjectCaseStudy EntityType = "project_case_study"
	EntityTypeSystemDesign     EntityType = "system_design"
	EntityTypeProblemSolution  EntityType = "problem_solution"
	EntityTypeCertification    EntityType = "certification"
	EntityTypeAIMLIntegration  EntityType = "aiml_integration"
	EntityTypeImpactMetric     EntityType = "impact_metric"
	EntityTypeSocialMediaPost  EntityType = "social_media_post"
	EntityTypeTechnicalWriting EntityType = "technical_writing"
	EntityTypeInterest         EntityType = "interest"
)

// Language represents supported languages.
type Language string

const (
	LanguageEN    Language = "en"
	LanguagePTBR  Language = "pt-BR"
	LanguageFR    Language = "fr"
	LanguageES    Language = "es"
	LanguageDE    Language = "de"
	LanguageRU    Language = "ru"
	LanguageJA    Language = "ja"
	LanguageKO    Language = "ko"
	LanguageZHCN  Language = "zh-CN"
	LanguageEL    Language = "el"
	LanguageLA    Language = "la"
)

// TranslationStatus represents the status of a translation.
type TranslationStatus string

const (
	TranslationStatusPending   TranslationStatus = "pending"
	TranslationStatusProcessing TranslationStatus = "processing"
	TranslationStatusCompleted  TranslationStatus = "completed"
	TranslationStatusFailed     TranslationStatus = "failed"
)

// Translation represents a translation for an entity.
type Translation struct {
	ID           uuid.UUID         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	EntityType   EntityType        `gorm:"column:entity_type;type:varchar(50);not null;index:idx_translation_lookup" json:"entityType"`
	EntityID     uuid.UUID         `gorm:"column:entity_id;type:uuid;not null;index:idx_translation_lookup" json:"entityId"`
	Language     Language          `gorm:"column:language;type:varchar(10);not null;index:idx_translation_lookup" json:"language"`
	Fields       datatypes.JSON    `gorm:"column:fields;type:jsonb;not null" json:"fields"` // Map of field names to translated values
	Status       TranslationStatus `gorm:"column:status;type:varchar(20);not null;default:'pending';index" json:"status"`
	ErrorMessage string            `gorm:"column:error_message;type:text" json:"errorMessage,omitempty"`
	CreatedAt    time.Time         `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt    time.Time         `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for Translation.
func (Translation) TableName() string {
	return "translations"
}

// NewTranslation creates a new translation entity.
func NewTranslation(entityType EntityType, entityID uuid.UUID, language Language, fields map[string]string) (*Translation, error) {
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	return &Translation{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		Language:   language,
		Fields:     datatypes.JSON(fieldsJSON),
		Status:     TranslationStatusPending,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}, nil
}

// GetFields returns the translated fields as a map.
func (t *Translation) GetFields() (map[string]string, error) {
	var fields map[string]string
	if err := json.Unmarshal(t.Fields, &fields); err != nil {
		return nil, err
	}
	return fields, nil
}

// SetFields sets the translated fields.
func (t *Translation) SetFields(fields map[string]string) error {
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	t.Fields = datatypes.JSON(fieldsJSON)
	t.UpdatedAt = time.Now().UTC()
	return nil
}

// Validate ensures translation invariants hold.
func (t *Translation) Validate() error {
	if t.EntityType == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyEntityType)
	}
	if t.EntityID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyEntityID)
	}
	if t.Language == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyLanguage)
	}
	if !isValidEntityType(t.EntityType) {
		return NewDomainError(ErrCodeInvalidEntityType, ErrUnsupportedEntityType)
	}
	if !isValidLanguage(t.Language) {
		return NewDomainError(ErrCodeInvalidLanguage, ErrUnsupportedLanguage)
	}
	return nil
}

func isValidEntityType(et EntityType) bool {
	switch et {
	case EntityTypeTestimonial, EntityTypePost, EntityTypeProject, EntityTypeCaseStudy, EntityTypeProjectCaseStudy, EntityTypeSystemDesign, EntityTypeProblemSolution, EntityTypeCertification, EntityTypeAIMLIntegration, EntityTypeImpactMetric, EntityTypeSocialMediaPost, EntityTypeTechnicalWriting, EntityTypeInterest:
		return true
	}
	return false
}

func isValidLanguage(lang Language) bool {
	switch lang {
	case LanguageEN, LanguagePTBR, LanguageFR, LanguageES, LanguageDE, LanguageRU, LanguageJA, LanguageKO, LanguageZHCN, LanguageEL, LanguageLA:
		return true
	}
	return false
}

// TranslationJob represents a job in the Redis queue.
type TranslationJob struct {
	ID         string     `json:"id"`
	EntityType EntityType `json:"entityType"`
	EntityID   string     `json:"entityId"`
	Language   Language   `json:"language"`
	Fields     []string   `json:"fields"` // Fields to translate
	SourceText map[string]string `json:"sourceText"` // Original text to translate
}

