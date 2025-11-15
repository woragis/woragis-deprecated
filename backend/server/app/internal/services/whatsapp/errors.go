package whatsapp

import "errors"

// Domain error codes for WhatsApp service.
const (
	ErrCodeInvalidPayload        = 8000
	ErrCodeServiceNotConfigured  = 8001
	ErrCodeNotConnected          = 8002
	ErrCodeClientNotFound        = 8003
	ErrCodeUserNotFound          = 8004
	ErrCodeNoPhoneNumber         = 8005
	ErrCodeSendFailure           = 8006
	ErrCodeAIGenerationFailure   = 8007
	ErrCodeInvalidPhoneNumber    = 8008
	ErrCodeRepositoryFailure     = 8009
)

// Domain error messages.
const (
	ErrServiceNotConfigured     = "whatsapp: service not configured"
	ErrNotConnected             = "whatsapp: WhatsApp is not connected. Please connect via QR code first"
	ErrClientNotFound           = "whatsapp: client not found"
	ErrUserNotFound             = "whatsapp: user not found"
	ErrClientNoPhoneNumber      = "whatsapp: client does not have a phone number configured"
	ErrUserNoPhoneNumber        = "whatsapp: user does not have a phone number configured"
	ErrSendFailure              = "whatsapp: failed to send message"
	ErrAIGenerationFailure      = "whatsapp: failed to generate message with AI"
	ErrInvalidPhoneNumber       = "whatsapp: invalid phone number format"
	ErrEmptyMessage             = "whatsapp: message cannot be empty"
	ErrEmptyClientID            = "whatsapp: client_id is required"
	ErrInvalidClientID          = "whatsapp: invalid client_id format"
	ErrRepositoryFailure        = "whatsapp: repository operation failed"
	ErrAIClientNotConfigured   = "whatsapp: AI client not configured"
	ErrAIEmptyResponse          = "whatsapp: AI returned empty message"
)

// DomainError represents a domain-level error for WhatsApp service.
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

// AsDomainError attempts to cast an error into a DomainError.
func AsDomainError(err error) (*DomainError, bool) {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr, true
	}
	return nil, false
}

