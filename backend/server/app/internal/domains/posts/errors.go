package posts

import "errors"

// Domain error codes.
const (
	ErrCodeInvalidPayload        = 2001
	ErrCodeInvalidTitle          = 2002
	ErrCodeInvalidContent        = 2003
	ErrCodeInvalidStatus         = 2004
	ErrCodeInvalidSlug           = 2005
	ErrCodeInvalidName           = 2006
	ErrCodePostNotFound          = 2007
	ErrCodeCategoryNotFound      = 2008
	ErrCodeTagNotFound           = 2009
	ErrCodeSkillNotFound         = 2010
	ErrCodeDuplicateSlug         = 2011
	ErrCodeUnauthorized          = 2012
	ErrCodeRepositoryFailure     = 2013
	ErrCodeUnsupportedPostStatus = 2014
)

// Domain error messages.
const (
	ErrNilPost              = "posts: post entity is nil"
	ErrEmptyPostID          = "posts: post id cannot be empty"
	ErrEmptyPostTitle       = "posts: post title cannot be empty"
	ErrEmptyPostContent     = "posts: post content cannot be empty"
	ErrEmptyPostSlug        = "posts: post slug cannot be empty"
	ErrPostNotFound         = "posts: post not found"
	ErrPostSlugTaken        = "posts: post slug already taken"
	ErrUnsupportedPostStatus = "posts: unsupported post status"

	ErrNilCategory       = "posts: category entity is nil"
	ErrEmptyCategoryID   = "posts: category id cannot be empty"
	ErrEmptyCategoryName = "posts: category name cannot be empty"
	ErrEmptyCategorySlug = "posts: category slug cannot be empty"
	ErrCategoryNotFound  = "posts: category not found"
	ErrCategorySlugTaken = "posts: category slug already taken"

	ErrNilTag       = "posts: tag entity is nil"
	ErrEmptyTagID   = "posts: tag id cannot be empty"
	ErrEmptyTagName = "posts: tag name cannot be empty"
	ErrEmptyTagSlug = "posts: tag slug cannot be empty"
	ErrTagNotFound  = "posts: tag not found"
	ErrTagSlugTaken = "posts: tag slug already taken"

	ErrEmptyUserID  = "posts: user id cannot be empty"
	ErrUnauthorized = "posts: unauthorized to perform this action"

	ErrUnableToPersist = "posts: unable to persist data"
	ErrUnableToFetch  = "posts: unable to fetch data"
	ErrUnableToUpdate  = "posts: unable to update data"
)

// DomainError represents a domain-specific error.
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

// AsDomainError checks if an error is a domain error.
func AsDomainError(err error) (*DomainError, bool) {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr, true
	}
	return nil, false
}

