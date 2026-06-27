package resumes

// Error codes for resume domain.
const (
	ErrCodeInvalidPayload  = "INVALID_PAYLOAD"
	ErrCodeInvalidName     = "INVALID_NAME"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeFileNotFound    = "FILE_NOT_FOUND"
	ErrCodeFileReadError   = "FILE_READ_ERROR"
	ErrCodeInvalidFileSize = "INVALID_FILE_SIZE"
)

// Error messages.
const (
	ErrNilResume       = "resumes: resume cannot be nil"
	ErrEmptyResumeID   = "resumes: resume ID cannot be empty"
	ErrEmptyUserID     = "resumes: user ID cannot be empty"
	ErrEmptyResumeTitle = "resumes: resume title cannot be empty"
	ErrEmptyFilePath   = "resumes: file path cannot be empty"
	ErrEmptyFileName   = "resumes: file name cannot be empty"
	ErrInvalidFileSize = "resumes: file size cannot be negative"
	ErrResumeNotFound  = "resumes: resume not found"
	ErrFileNotFound    = "resumes: resume file not found"
	ErrFileReadError   = "resumes: error reading resume file"
	ErrNoMainResume    = "resumes: no main resume found"
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

