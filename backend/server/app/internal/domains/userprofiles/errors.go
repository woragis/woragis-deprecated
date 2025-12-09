package userprofiles

// Domain error codes.
const (
	ErrCodeInvalidPayload = "INVALID_PAYLOAD"
	ErrCodeNotFound       = "NOT_FOUND"
)

// Domain error messages.
const (
	ErrNilProfile      = "userprofiles: profile is nil"
	ErrEmptyProfileID  = "userprofiles: profile ID cannot be empty"
	ErrEmptyUserID     = "userprofiles: user ID cannot be empty"
	ErrProfileNotFound = "userprofiles: profile not found"
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

