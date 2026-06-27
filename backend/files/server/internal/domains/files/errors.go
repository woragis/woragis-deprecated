package files

// Domain error codes
const (
	ErrCodeInvalidPayload = "INVALID_PAYLOAD"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeRepositoryFailure = "REPOSITORY_FAILURE"
	ErrCodeStorageFailure = "STORAGE_FAILURE"
)

// Domain error messages
const (
	ErrNilFile        = "file cannot be nil"
	ErrEmptyFileID    = "file ID cannot be empty"
	ErrEmptyUserID    = "user ID cannot be empty"
	ErrEmptyFileName  = "file name cannot be empty"
	ErrEmptyFilePath  = "file path cannot be empty"
	ErrEmptyMimeType  = "mime type cannot be empty"
	ErrInvalidFileSize = "file size must be non-negative"
	ErrFileNotFound   = "file not found"
	ErrFileExists     = "file already exists"
	ErrUnableToPersist = "unable to persist file"
	ErrUnableToFetch   = "unable to fetch file"
	ErrUnableToUpdate  = "unable to update file"
	ErrUnableToDelete  = "unable to delete file"
	ErrStorageError    = "storage operation failed"
)

// DomainError represents a domain-specific error.
type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new domain error.
func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

