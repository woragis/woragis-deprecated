package certifications

import "errors"

const (
	ErrCodeInvalidPayload    = 9000
	ErrCodeInvalidName       = 9001
	ErrCodeInvalidIssuer     = 9002
	ErrCodeInvalidDate       = 9003
	ErrCodeInvalidStatus     = 9004
	ErrCodeInvalidCategory   = 9005
	ErrCodeRepositoryFailure = 9006
	ErrCodeNotFound          = 9007
	ErrCodeUnauthorized      = 9008
	ErrCodeConflict          = 9009
	ErrCodeInvalidEntityType = 9010
)

const (
	ErrNilCertification        = "certifications: certification entity is nil"
	ErrEmptyCertificationID    = "certifications: certification id cannot be empty"
	ErrEmptyUserID             = "certifications: user id cannot be empty"
	ErrEmptyName               = "certifications: name cannot be empty"
	ErrEmptyIssuer             = "certifications: issuer cannot be empty"
	ErrEmptyIssueDate          = "certifications: issue date cannot be empty"
	ErrExpiryBeforeIssue       = "certifications: expiry date cannot be before issue date"
	ErrUnsupportedStatus       = "certifications: unsupported certification status"
	ErrUnsupportedCategory     = "certifications: unsupported certification category"
	ErrUnsupportedEntityType    = "certifications: unsupported entity type"
	ErrCertificationNotFound   = "certifications: certification not found"
	ErrUnableToPersist         = "certifications: unable to persist data"
	ErrUnableToFetch           = "certifications: unable to fetch data"
	ErrUnableToUpdate          = "certifications: unable to update data"
	ErrUnauthorized            = "certifications: unauthorized access"
	ErrCertificationAlreadyExists = "certifications: certification already exists"
	ErrNilLink                 = "certifications: certification entity link is nil"
	ErrEmptyLinkID             = "certifications: link id cannot be empty"
	ErrEmptyEntityID           = "certifications: entity id cannot be empty"
)

type DomainError struct {
	Code    int
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

func NewDomainError(code int, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

func AsDomainError(err error) (*DomainError, bool) {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr, true
	}
	return nil, false
}

