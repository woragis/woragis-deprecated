package files

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// FileType represents the type of file.
type FileType string

const (
	FileTypeDocument FileType = "document"
	FileTypeImage    FileType = "image"
	FileTypeVideo    FileType = "video"
	FileTypeAudio    FileType = "audio"
	FileTypeArchive  FileType = "archive"
	FileTypeOther    FileType = "other"
)

// FileStatus represents the status of a file.
type FileStatus string

const (
	FileStatusUploading FileStatus = "uploading"
	FileStatusProcessing FileStatus = "processing"
	FileStatusReady     FileStatus = "ready"
	FileStatusFailed    FileStatus = "failed"
	FileStatusDeleted   FileStatus = "deleted"
)

// File represents a stored file with metadata.
type File struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID  `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Name        string     `gorm:"column:name;size:255;not null" json:"name"`
	OriginalName string    `gorm:"column:original_name;size:255;not null" json:"originalName"`
	Path        string     `gorm:"column:path;size:512;not null" json:"path"`
	URL         string     `gorm:"column:url;size:512" json:"url,omitempty"`
	MimeType    string     `gorm:"column:mime_type;size:128;not null" json:"mimeType"`
	Size        int64      `gorm:"column:size;not null" json:"size"`
	FileType    FileType   `gorm:"column:file_type;type:varchar(32);not null;index" json:"fileType"`
	Status      FileStatus `gorm:"column:status;type:varchar(32);not null;default:'uploading';index" json:"status"`
	Hash        string     `gorm:"column:hash;size:64;index" json:"hash,omitempty"` // SHA256 hash for deduplication
	Metadata    map[string]interface{} `gorm:"column:metadata;type:jsonb" json:"metadata,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

// TableName specifies the table name for File.
func (File) TableName() string {
	return "files"
}

// NewFile creates a new file entity.
func NewFile(userID uuid.UUID, originalName, path, mimeType string, size int64) (*File, error) {
	file := &File{
		ID:          uuid.New(),
		UserID:      userID,
		OriginalName: strings.TrimSpace(originalName),
		Path:        strings.TrimSpace(path),
		MimeType:    strings.TrimSpace(mimeType),
		Size:        size,
		Status:      FileStatusUploading,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Extract name without extension for display
	nameParts := strings.Split(originalName, ".")
	if len(nameParts) > 0 {
		file.Name = strings.Join(nameParts[:len(nameParts)-1], ".")
	} else {
		file.Name = originalName
	}

	// Determine file type from mime type
	file.FileType = determineFileType(mimeType)

	return file, file.Validate()
}

// Validate ensures file invariants.
func (f *File) Validate() error {
	if f == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilFile)
	}

	if f.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyFileID)
	}

	if f.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if strings.TrimSpace(f.OriginalName) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyFileName)
	}

	if strings.TrimSpace(f.Path) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyFilePath)
	}

	if strings.TrimSpace(f.MimeType) == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyMimeType)
	}

	if f.Size < 0 {
		return NewDomainError(ErrCodeInvalidPayload, ErrInvalidFileSize)
	}

	return nil
}

// determineFileType determines the file type from mime type.
func determineFileType(mimeType string) FileType {
	mimeType = strings.ToLower(mimeType)
	
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return FileTypeImage
	case strings.HasPrefix(mimeType, "video/"):
		return FileTypeVideo
	case strings.HasPrefix(mimeType, "audio/"):
		return FileTypeAudio
	case strings.HasPrefix(mimeType, "application/pdf") ||
		 strings.HasPrefix(mimeType, "application/msword") ||
		 strings.HasPrefix(mimeType, "application/vnd.openxmlformats"):
		return FileTypeDocument
	case strings.HasPrefix(mimeType, "application/zip") ||
		 strings.HasPrefix(mimeType, "application/x-tar") ||
		 strings.HasPrefix(mimeType, "application/x-rar"):
		return FileTypeArchive
	default:
		return FileTypeOther
	}
}

// MarkReady marks the file as ready.
func (f *File) MarkReady() {
	f.Status = FileStatusReady
	f.UpdatedAt = time.Now().UTC()
}

// MarkFailed marks the file as failed.
func (f *File) MarkFailed() {
	f.Status = FileStatusFailed
	f.UpdatedAt = time.Now().UTC()
}

// MarkProcessing marks the file as processing.
func (f *File) MarkProcessing() {
	f.Status = FileStatusProcessing
	f.UpdatedAt = time.Now().UTC()
}

