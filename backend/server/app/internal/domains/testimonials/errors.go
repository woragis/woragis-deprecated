package testimonials

import "errors"

const (
	ErrCodeInvalidPayload    = 8000
	ErrCodeInvalidAuthorName = 8001
	ErrCodeInvalidContent    = 8002
	ErrCodeInvalidRating     = 8003
	ErrCodeInvalidStatus     = 8004
	ErrCodeRepositoryFailure = 8005
	ErrCodeNotFound          = 8006
	ErrCodeUnauthorized      = 8007
	ErrCodeConflict          = 8008
)

const (
	ErrNilTestimonial        = "testimonials: testimonial entity is nil"
	ErrEmptyTestimonialID    = "testimonials: testimonial id cannot be empty"
	ErrEmptyUserID           = "testimonials: user id cannot be empty"
	ErrEmptyAuthorName       = "testimonials: author name cannot be empty"
	ErrEmptyContent          = "testimonials: content cannot be empty"
	ErrInvalidRating         = "testimonials: rating must be between 1 and 5"
	ErrTestimonialNotFound   = "testimonials: testimonial not found"
	ErrUnsupportedStatus     = "testimonials: unsupported testimonial status"
	ErrUnableToPersist       = "testimonials: unable to persist data"
	ErrUnableToFetch         = "testimonials: unable to fetch data"
	ErrUnableToUpdate        = "testimonials: unable to update data"
	ErrUnauthorized          = "testimonials: unauthorized access"
	ErrTestimonialAlreadyExists = "testimonials: testimonial already exists"
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

