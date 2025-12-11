package resumes

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

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

// Resume represents a generated resume PDF file.
type Resume struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Title      string    `gorm:"column:title;size:255;not null" json:"title"`
	IsMain     bool      `gorm:"column:is_main;not null;default:false;index" json:"isMain"`
	IsFeatured bool      `gorm:"column:is_featured;not null;default:false;index" json:"isFeatured"`
	FilePath          string    `gorm:"column:file_path;size:512;not null" json:"filePath"`
	FileName          string    `gorm:"column:file_name;size:255;not null" json:"fileName"`
	FileSize          int64     `gorm:"column:file_size;default:0" json:"fileSize"`
	Tags              JSONArray `gorm:"column:tags;type:jsonb;default:'[]'" json:"tags"`
	ApplicationsUsed  int       `gorm:"column:applications_used;default:0" json:"applicationsUsed"`
	InterviewRate     float64   `gorm:"column:interview_rate;default:0" json:"interviewRate"` // Percentage (0-100)
	OfferRate         float64   `gorm:"column:offer_rate;default:0" json:"offerRate"`       // Percentage (0-100)
	CreatedAt         time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for Resume.
func (Resume) TableName() string {
	return "resumes"
}

// NewResume creates a new resume entity.
func NewResume(userID uuid.UUID, title, filePath, fileName string, fileSize int64, tags JSONArray) (*Resume, error) {
	// Normalize tags: lowercase, trim, remove duplicates, limit to 10
	normalizedTags := normalizeTags(tags)
	
	resume := &Resume{
		ID:               uuid.New(),
		UserID:           userID,
		Title:            strings.TrimSpace(title),
		FilePath:         strings.TrimSpace(filePath),
		FileName:         strings.TrimSpace(fileName),
		FileSize:         fileSize,
		Tags:             normalizedTags,
		ApplicationsUsed: 0,
		InterviewRate:    0.0,
		OfferRate:        0.0,
		IsMain:           false,
		IsFeatured:       false,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	return resume, resume.Validate()
}

// normalizeTags normalizes tags: lowercase, trim, remove duplicates, limit to 10.
func normalizeTags(tags JSONArray) JSONArray {
	if tags == nil {
		return JSONArray{}
	}
	
	seen := make(map[string]bool)
	result := JSONArray{}
	
	for _, tag := range tags {
		normalized := strings.ToLower(strings.TrimSpace(tag))
		if normalized != "" && !seen[normalized] && len(result) < 10 {
			seen[normalized] = true
			result = append(result, normalized)
		}
	}
	
	return result
}

// Validate ensures resume invariants hold.
func (r *Resume) Validate() error {
	if r == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilResume)
	}

	if r.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyResumeID)
	}

	if r.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if r.Title == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyResumeTitle)
	}

	if r.FilePath == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyFilePath)
	}

	if r.FileName == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyFileName)
	}

	if r.FileSize < 0 {
		return NewDomainError(ErrCodeInvalidPayload, ErrInvalidFileSize)
	}

	return nil
}

// MarkAsMain sets the resume as the main resume and unmarks others.
func (r *Resume) MarkAsMain() {
	r.IsMain = true
	r.UpdatedAt = time.Now().UTC()
}

// UnmarkAsMain removes the main flag.
func (r *Resume) UnmarkAsMain() {
	r.IsMain = false
	r.UpdatedAt = time.Now().UTC()
}

// MarkAsFeatured sets the resume as featured.
func (r *Resume) MarkAsFeatured() {
	r.IsFeatured = true
	r.UpdatedAt = time.Now().UTC()
}

// UnmarkAsFeatured removes the featured flag.
func (r *Resume) UnmarkAsFeatured() {
	r.IsFeatured = false
	r.UpdatedAt = time.Now().UTC()
}

// UpdateTitle updates the resume title.
func (r *Resume) UpdateTitle(title string) error {
	r.Title = strings.TrimSpace(title)
	r.UpdatedAt = time.Now().UTC()
	return r.Validate()
}

// UpdateTags updates the resume tags.
func (r *Resume) UpdateTags(tags JSONArray) error {
	r.Tags = normalizeTags(tags)
	r.UpdatedAt = time.Now().UTC()
	return r.Validate()
}

// UpdateFilePath updates the file path and name.
func (r *Resume) UpdateFilePath(filePath, fileName string, fileSize int64) error {
	r.FilePath = strings.TrimSpace(filePath)
	r.FileName = strings.TrimSpace(fileName)
	r.FileSize = fileSize
	r.UpdatedAt = time.Now().UTC()
	return r.Validate()
}

