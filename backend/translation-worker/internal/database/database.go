package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Translation represents a translation record in the database.
// This matches the Translation entity from the server.
type Translation struct {
	ID           uuid.UUID         `gorm:"column:id;type:uuid;primaryKey"`
	EntityType   string            `gorm:"column:entity_type;type:varchar(50);not null;index:idx_translation_lookup"`
	EntityID     uuid.UUID         `gorm:"column:entity_id;type:uuid;not null;index:idx_translation_lookup"`
	Language     string            `gorm:"column:language;type:varchar(10);not null;index:idx_translation_lookup"`
	Fields       datatypes.JSON    `gorm:"column:fields;type:jsonb;not null"`
	Status       string            `gorm:"column:status;type:varchar(20);not null;default:'pending';index"`
	ErrorMessage string            `gorm:"column:error_message;type:text"`
	CreatedAt    time.Time         `gorm:"column:created_at"`
	UpdatedAt    time.Time         `gorm:"column:updated_at"`
}

// TableName specifies the table name for Translation.
func (Translation) TableName() string {
	return "translations"
}

// Repository provides database operations for translations.
type Repository interface {
	GetTranslationByEntity(ctx context.Context, entityType string, entityID uuid.UUID, language string) (*Translation, error)
	CreateTranslation(ctx context.Context, translation *Translation) error
	UpdateTranslation(ctx context.Context, translation *Translation) error
	FetchSourceTextFromEntity(ctx context.Context, entityType string, entityID uuid.UUID, fields []string) (map[string]string, error)
}

type repository struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewRepository creates a new translation repository.
func NewRepository(dbURL string, logger *slog.Logger) (Repository, error) {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &repository{
		db:     db,
		logger: logger,
	}, nil
}

// GetTranslationByEntity retrieves a translation by entity type, ID, and language.
func (r *repository) GetTranslationByEntity(ctx context.Context, entityType string, entityID uuid.UUID, language string) (*Translation, error) {
	var translation Translation
	err := r.db.WithContext(ctx).
		Where("entity_type = ? AND entity_id = ? AND language = ?", entityType, entityID, language).
		First(&translation).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get translation: %w", err)
	}

	return &translation, nil
}

// CreateTranslation creates a new translation record.
func (r *repository) CreateTranslation(ctx context.Context, translation *Translation) error {
	if translation.ID == uuid.Nil {
		translation.ID = uuid.New()
	}
	if translation.CreatedAt.IsZero() {
		translation.CreatedAt = time.Now().UTC()
	}
	if translation.UpdatedAt.IsZero() {
		translation.UpdatedAt = time.Now().UTC()
	}

	if err := r.db.WithContext(ctx).Create(translation).Error; err != nil {
		return fmt.Errorf("failed to create translation: %w", err)
	}

	return nil
}

// UpdateTranslation updates an existing translation record.
func (r *repository) UpdateTranslation(ctx context.Context, translation *Translation) error {
	translation.UpdatedAt = time.Now().UTC()

	if err := r.db.WithContext(ctx).Save(translation).Error; err != nil {
		return fmt.Errorf("failed to update translation: %w", err)
	}

	return nil
}

// FetchSourceTextFromEntity fetches source text from the entity table.
// This is used when sourceText is not provided in the job.
func (r *repository) FetchSourceTextFromEntity(ctx context.Context, entityType string, entityID uuid.UUID, fields []string) (map[string]string, error) {
	sourceText := make(map[string]string)

	// Map entity types to table names and field mappings
	// This is a simplified version - you may need to adjust based on your schema
	switch entityType {
	case "testimonial":
		type Testimonial struct {
			Content       string `gorm:"column:content"`
			AuthorRole    string `gorm:"column:author_role"`
			AuthorCompany string `gorm:"column:author_company"`
		}
		var testimonial Testimonial
		if err := r.db.WithContext(ctx).Table("testimonials").Where("id = ?", entityID).First(&testimonial).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "content":
				sourceText["content"] = testimonial.Content
			case "authorRole":
				sourceText["authorRole"] = testimonial.AuthorRole
			case "authorCompany":
				sourceText["authorCompany"] = testimonial.AuthorCompany
			}
		}
	case "project":
		type Project struct {
			Name        string `gorm:"column:name"`
			Description string `gorm:"column:description"`
		}
		var project Project
		if err := r.db.WithContext(ctx).Table("projects").Where("id = ?", entityID).First(&project).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "name":
				sourceText["name"] = project.Name
			case "description":
				sourceText["description"] = project.Description
			}
		}
	case "certification":
		type Certification struct {
			Name        string `gorm:"column:name"`
			Issuer      string `gorm:"column:issuer"`
			Description string `gorm:"column:description"`
		}
		var cert Certification
		if err := r.db.WithContext(ctx).Table("certifications").Where("id = ?", entityID).First(&cert).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "name":
				sourceText["name"] = cert.Name
			case "issuer":
				sourceText["issuer"] = cert.Issuer
			case "description":
				sourceText["description"] = cert.Description
			}
		}
	case "skill":
		type Skill struct {
			Name        string `gorm:"column:name"`
			Description string `gorm:"column:description"`
		}
		var skill Skill
		if err := r.db.WithContext(ctx).Table("skills").Where("id = ?", entityID).First(&skill).Error; err != nil {
			return nil, err
		}
		for _, field := range fields {
			switch field {
			case "name":
				sourceText["name"] = skill.Name
			case "description":
				sourceText["description"] = skill.Description
			}
		}
	// Add more entity types as needed
	default:
		return nil, fmt.Errorf("unsupported entity type for auto-fetch: %s", entityType)
	}

	return sourceText, nil
}

// SetFields sets the translated fields on a Translation.
func (t *Translation) SetFields(fields map[string]string) error {
	fieldsJSON, err := json.Marshal(fields)
	if err != nil {
		return fmt.Errorf("failed to marshal fields: %w", err)
	}
	t.Fields = datatypes.JSON(fieldsJSON)
	t.UpdatedAt = time.Now().UTC()
	return nil
}

// GetFields returns the translated fields as a map.
func (t *Translation) GetFields() (map[string]string, error) {
	var fields map[string]string
	if err := json.Unmarshal(t.Fields, &fields); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fields: %w", err)
	}
	return fields, nil
}
