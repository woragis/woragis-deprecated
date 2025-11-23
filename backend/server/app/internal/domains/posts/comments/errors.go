package comments

const (
	ErrCodeInvalidPayload        = 2101
	ErrCodeInvalidContent        = 2102
	ErrCodeInvalidStatus         = 2103
	ErrCodeInvalidName           = 2104
	ErrCodeCommentNotFound       = 2105
	ErrCodeUnauthorized          = 2106
	ErrCodeRepositoryFailure     = 2107
	ErrCodeUnsupportedCommentStatus = 2108
)

const (
	ErrNilComment              = "comments: comment entity is nil"
	ErrEmptyCommentID          = "comments: comment id cannot be empty"
	ErrEmptyPostID             = "comments: post id cannot be empty"
	ErrEmptyAuthorName         = "comments: author name cannot be empty"
	ErrEmptyCommentContent     = "comments: comment content cannot be empty"
	ErrCommentNotFound         = "comments: comment not found"
	ErrUnsupportedCommentStatus = "comments: unsupported comment status"
	ErrUnauthorized            = "comments: unauthorized to perform this action"

	ErrUnableToPersist = "comments: unable to persist data"
	ErrUnableToFetch   = "comments: unable to fetch data"
	ErrUnableToUpdate  = "comments: unable to update data"
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

