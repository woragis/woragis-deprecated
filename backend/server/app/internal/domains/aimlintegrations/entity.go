package aimlintegrations

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

// IntegrationType represents the type of AI/ML integration.
type IntegrationType string

const (
	IntegrationTypeRAG           IntegrationType = "rag"              // Retrieval-Augmented Generation
	IntegrationTypeLLM          IntegrationType = "llm"               // Large Language Model
	IntegrationTypeMLModel       IntegrationType = "ml_model"         // Machine Learning Model
	IntegrationTypeComputerVision IntegrationType = "computer_vision"  // Computer Vision
	IntegrationTypeNLP           IntegrationType = "nlp"              // Natural Language Processing
	IntegrationTypeRecommendation IntegrationType = "recommendation"  // Recommendation System
	IntegrationTypeChatbot       IntegrationType = "chatbot"          // Chatbot/Virtual Assistant
	IntegrationTypeAnomalyDetection IntegrationType = "anomaly_detection" // Anomaly Detection
	IntegrationTypePredictiveAnalytics IntegrationType = "predictive_analytics" // Predictive Analytics
	IntegrationTypeGenerativeAI   IntegrationType = "generative_ai"    // Generative AI
	IntegrationTypeOther         IntegrationType = "other"
)

// Framework represents the AI/ML framework or platform used.
type Framework string

const (
	FrameworkOpenAI      Framework = "openai"
	FrameworkAnthropic   Framework = "anthropic"
	FrameworkHuggingFace Framework = "huggingface"
	FrameworkTensorFlow   Framework = "tensorflow"
	FrameworkPyTorch      Framework = "pytorch"
	FrameworkLangChain    Framework = "langchain"
	FrameworkLlamaIndex   Framework = "llamaindex"
	FrameworkCohere       Framework = "cohere"
	FrameworkGoogleAI     Framework = "google_ai"
	FrameworkAzureAI      Framework = "azure_ai"
	FrameworkAWSBedrock   Framework = "aws_bedrock"
	FrameworkCustom       Framework = "custom"
	FrameworkOther        Framework = "other"
)

// AIMLIntegration represents an AI/ML integration showcase entry.
type AIMLIntegration struct {
	ID          uuid.UUID       `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID       `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Title       string          `gorm:"column:title;size:255;not null" json:"title"`
	Description string          `gorm:"column:description;type:text;not null" json:"description"`
	Type        IntegrationType `gorm:"column:type;type:varchar(50);not null;index" json:"type"`
	Framework   Framework       `gorm:"column:framework;type:varchar(50);not null;index" json:"framework"`
	// Model information
	ModelName   string          `gorm:"column:model_name;size:255" json:"modelName,omitempty"`
	ModelVersion string         `gorm:"column:model_version;size:100" json:"modelVersion,omitempty"`
	// Use case and impact
	UseCase     string          `gorm:"column:use_case;type:text" json:"useCase,omitempty"`
	Impact      string          `gorm:"column:impact;type:text" json:"impact,omitempty"`
	// Technical details
	Technologies JSONArray      `gorm:"column:technologies;type:jsonb" json:"technologies,omitempty"` // Array of technology names
	Architecture string         `gorm:"column:architecture;type:text" json:"architecture,omitempty"`
	// Metrics and results
	Metrics     string          `gorm:"column:metrics;type:text" json:"metrics,omitempty"` // JSON string or description
	// Links to other entities
	ProjectID   *uuid.UUID      `gorm:"column:project_id;type:uuid;index" json:"projectId,omitempty"`
	CaseStudyID *uuid.UUID      `gorm:"column:case_study_id;type:uuid;index" json:"caseStudyId,omitempty"`
	// Display and organization
	Featured    bool            `gorm:"column:featured;not null;default:false;index" json:"featured"`
	DisplayOrder int            `gorm:"column:display_order;not null;default:0;index" json:"displayOrder"`
	// Media and documentation
	DemoURL     string          `gorm:"column:demo_url;size:500" json:"demoUrl,omitempty"`
	DocumentationURL string      `gorm:"column:documentation_url;size:500" json:"documentationUrl,omitempty"`
	GitHubURL   string          `gorm:"column:github_url;size:500" json:"githubUrl,omitempty"`
	// Timestamps
	CreatedAt   time.Time       `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time       `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for AIMLIntegration.
func (AIMLIntegration) TableName() string {
	return "aiml_integrations"
}

// NewAIMLIntegration creates a new AI/ML integration entity.
func NewAIMLIntegration(userID uuid.UUID, title, description string, integrationType IntegrationType, framework Framework) (*AIMLIntegration, error) {
	integration := &AIMLIntegration{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Type:        integrationType,
		Framework:   framework,
		Technologies: JSONArray{},
		Featured:    false,
		DisplayOrder: 0,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	return integration, integration.Validate()
}

// Validate ensures AI/ML integration invariants hold.
func (a *AIMLIntegration) Validate() error {
	if a == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilIntegration)
	}
	if a.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyIntegrationID)
	}
	if a.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if strings.TrimSpace(a.Title) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyTitle)
	}
	if strings.TrimSpace(a.Description) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyDescription)
	}
	if !isValidIntegrationType(a.Type) {
		return NewDomainError(ErrCodeInvalidType, ErrUnsupportedIntegrationType)
	}
	if !isValidFramework(a.Framework) {
		return NewDomainError(ErrCodeInvalidFramework, ErrUnsupportedFramework)
	}
	return nil
}

// UpdateDetails updates integration details.
func (a *AIMLIntegration) UpdateDetails(title, description, useCase, impact, architecture, metrics string) {
	if title != "" {
		a.Title = strings.TrimSpace(title)
	}
	if description != "" {
		a.Description = strings.TrimSpace(description)
	}
	if useCase != "" {
		a.UseCase = strings.TrimSpace(useCase)
	}
	if impact != "" {
		a.Impact = strings.TrimSpace(impact)
	}
	if architecture != "" {
		a.Architecture = strings.TrimSpace(architecture)
	}
	if metrics != "" {
		a.Metrics = strings.TrimSpace(metrics)
	}
	a.UpdatedAt = time.Now().UTC()
}

// SetModelInfo updates model information.
func (a *AIMLIntegration) SetModelInfo(modelName, modelVersion string) {
	a.ModelName = strings.TrimSpace(modelName)
	a.ModelVersion = strings.TrimSpace(modelVersion)
	a.UpdatedAt = time.Now().UTC()
}

// SetTechnologies updates the technologies array.
func (a *AIMLIntegration) SetTechnologies(technologies []string) {
	a.Technologies = JSONArray(technologies)
	a.UpdatedAt = time.Now().UTC()
}

// SetProjectLink updates the project link.
func (a *AIMLIntegration) SetProjectLink(projectID uuid.UUID) {
	a.ProjectID = &projectID
	a.UpdatedAt = time.Now().UTC()
}

// ClearProjectLink removes the project link.
func (a *AIMLIntegration) ClearProjectLink() {
	a.ProjectID = nil
	a.UpdatedAt = time.Now().UTC()
}

// SetCaseStudyLink updates the case study link.
func (a *AIMLIntegration) SetCaseStudyLink(caseStudyID uuid.UUID) {
	a.CaseStudyID = &caseStudyID
	a.UpdatedAt = time.Now().UTC()
}

// ClearCaseStudyLink removes the case study link.
func (a *AIMLIntegration) ClearCaseStudyLink() {
	a.CaseStudyID = nil
	a.UpdatedAt = time.Now().UTC()
}

// SetFeatured updates the featured flag.
func (a *AIMLIntegration) SetFeatured(featured bool) {
	a.Featured = featured
	a.UpdatedAt = time.Now().UTC()
}

// SetDisplayOrder updates the display order.
func (a *AIMLIntegration) SetDisplayOrder(order int) {
	a.DisplayOrder = order
	a.UpdatedAt = time.Now().UTC()
}

// SetURLs updates demo, documentation, and GitHub URLs.
func (a *AIMLIntegration) SetURLs(demoURL, documentationURL, githubURL string) {
	if demoURL != "" {
		a.DemoURL = strings.TrimSpace(demoURL)
	}
	if documentationURL != "" {
		a.DocumentationURL = strings.TrimSpace(documentationURL)
	}
	if githubURL != "" {
		a.GitHubURL = strings.TrimSpace(githubURL)
	}
	a.UpdatedAt = time.Now().UTC()
}

// Validation helpers

func isValidIntegrationType(it IntegrationType) bool {
	switch it {
	case IntegrationTypeRAG, IntegrationTypeLLM, IntegrationTypeMLModel,
		IntegrationTypeComputerVision, IntegrationTypeNLP, IntegrationTypeRecommendation,
		IntegrationTypeChatbot, IntegrationTypeAnomalyDetection, IntegrationTypePredictiveAnalytics,
		IntegrationTypeGenerativeAI, IntegrationTypeOther:
		return true
	}
	return false
}

func isValidFramework(f Framework) bool {
	switch f {
	case FrameworkOpenAI, FrameworkAnthropic, FrameworkHuggingFace,
		FrameworkTensorFlow, FrameworkPyTorch, FrameworkLangChain, FrameworkLlamaIndex,
		FrameworkCohere, FrameworkGoogleAI, FrameworkAzureAI, FrameworkAWSBedrock,
		FrameworkCustom, FrameworkOther:
		return true
	}
	return false
}

// JSONArray is a custom type for storing JSON arrays in PostgreSQL.
type JSONArray []string

// Value implements the driver.Valuer interface.
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface.
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte(value.(string)), j)
	}
	return json.Unmarshal(bytes, j)
}

