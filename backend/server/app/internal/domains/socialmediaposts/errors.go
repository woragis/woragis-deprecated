package socialmediaposts

import "errors"

const (
	ErrCodeInvalidPayload          = 8000
	ErrCodeInvalidPlatform         = 8001
	ErrCodeInvalidEntityType       = 8002
	ErrCodeInvalidRelationshipType = 8003
	ErrCodeInvalidStatus          = 8004
	ErrCodeRepositoryFailure       = 8005
	ErrCodeNotFound                = 8006
	ErrCodeConflict                = 8007
)

const (
	ErrNilPost                  = "socialmediaposts: post entity is nil"
	ErrNilLink                  = "socialmediaposts: link entity is nil"
	ErrEmptyPostID              = "socialmediaposts: post id cannot be empty"
	ErrEmptyLinkID              = "socialmediaposts: link id cannot be empty"
	ErrEmptyURL                 = "socialmediaposts: url cannot be empty"
	ErrEmptyEntityID            = "socialmediaposts: entity id cannot be empty"
	ErrPostNotFound             = "socialmediaposts: post not found"
	ErrLinkNotFound             = "socialmediaposts: link not found"
	ErrPostAlreadyExists        = "socialmediaposts: post with this url already exists"
	ErrLinkAlreadyExists        = "socialmediaposts: link already exists"
	ErrUnsupportedPlatform      = "socialmediaposts: unsupported platform"
	ErrUnsupportedEntityType    = "socialmediaposts: unsupported entity type"
	ErrUnsupportedRelationshipType = "socialmediaposts: unsupported relationship type"
	ErrUnsupportedStatus        = "socialmediaposts: unsupported status"
	ErrUnableToPersist          = "socialmediaposts: unable to persist data"
	ErrUnableToFetch            = "socialmediaposts: unable to fetch data"
	ErrUnableToUpdate           = "socialmediaposts: unable to update data"
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

