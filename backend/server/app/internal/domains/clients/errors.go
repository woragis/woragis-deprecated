package clients

// Domain error codes.
const (
	ErrCodeInvalidPayload     = 1000
	ErrCodeInvalidName        = 1001
	ErrCodeInvalidPhoneNumber  = 1002
	ErrCodeInvalidEmail        = 1003
	ErrCodeNotFound           = 1004
	ErrCodeRepositoryFailure  = 1005
)

// Domain error messages.
const (
	ErrNilClient           = "clients: client cannot be nil"
	ErrEmptyClientID       = "clients: client ID cannot be empty"
	ErrEmptyUserID         = "clients: user ID cannot be empty"
	ErrEmptyName           = "clients: name cannot be empty"
	ErrNameTooLong         = "clients: name cannot exceed 120 characters"
	ErrEmptyPhoneNumber    = "clients: phone number cannot be empty"
	ErrPhoneNumberTooLong  = "clients: phone number cannot exceed 20 characters"
	ErrEmailTooLong        = "clients: email cannot exceed 255 characters"
	ErrCompanyTooLong      = "clients: company name cannot exceed 120 characters"
	ErrClientNotFound      = "clients: client not found"
	ErrRepositoryFailure   = "clients: repository operation failed"
)

// DomainError represents a domain-level error.
type DomainError struct {
	Code    int
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new domain error.
func NewDomainError(code int, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

