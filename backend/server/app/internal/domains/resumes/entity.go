package resumes

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Resume represents a generated resume PDF file.
type Resume struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Title      string    `gorm:"column:title;size:255;not null" json:"title"`
	IsMain     bool      `gorm:"column:is_main;not null;default:false;index" json:"isMain"`
	IsFeatured bool      `gorm:"column:is_featured;not null;default:false;index" json:"isFeatured"`
	FilePath   string    `gorm:"column:file_path;size:512;not null" json:"filePath"`
	FileName   string    `gorm:"column:file_name;size:255;not null" json:"fileName"`
	FileSize   int64     `gorm:"column:file_size;default:0" json:"fileSize"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for Resume.
func (Resume) TableName() string {
	return "resumes"
}

// NewResume creates a new resume entity.
func NewResume(userID uuid.UUID, title, filePath, fileName string, fileSize int64) (*Resume, error) {
	resume := &Resume{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     strings.TrimSpace(title),
		FilePath:  strings.TrimSpace(filePath),
		FileName:  strings.TrimSpace(fileName),
		FileSize:  fileSize,
		IsMain:    false,
		IsFeatured: false,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return resume, resume.Validate()
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

// UpdateFilePath updates the file path and name.
func (r *Resume) UpdateFilePath(filePath, fileName string, fileSize int64) error {
	r.FilePath = strings.TrimSpace(filePath)
	r.FileName = strings.TrimSpace(fileName)
	r.FileSize = fileSize
	r.UpdatedAt = time.Now().UTC()
	return r.Validate()
}

