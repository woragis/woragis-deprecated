package content

import "errors"

const (
	ErrCodeInvalidPayload          = 8200
	ErrCodeInvalidStatus           = 8201
	ErrCodeInvalidPriority         = 8202
	ErrCodeRepositoryFailure       = 8203
	ErrCodeNotFound                = 8204
	ErrCodeConflict                = 8205
)

const (
	ErrNilContentPost        = "content: content post entity is nil"
	ErrNilRepurposing        = "content: repurposing entity is nil"
	ErrEmptyContentPostID    = "content: content post id cannot be empty"
	ErrEmptyRepurposingID    = "content: repurposing id cannot be empty"
	ErrEmptyPostID           = "content: post id cannot be empty"
	ErrEmptySocialPostID     = "content: social post id cannot be empty"
	ErrContentPostNotFound   = "content: content post not found"
	ErrRepurposingNotFound   = "content: repurposing not found"
	ErrContentPostExists     = "content: content post already exists"
	ErrUnsupportedStatus     = "content: unsupported status"
	ErrUnsupportedPriority   = "content: unsupported priority"
	ErrUnableToPersist       = "content: unable to persist data"
	ErrUnableToFetch         = "content: unable to fetch data"
	ErrUnableToUpdate        = "content: unable to update data"
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
