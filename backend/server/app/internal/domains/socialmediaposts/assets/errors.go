package assets

import "errors"

const (
	ErrCodeInvalidPayload    = 8500
	ErrCodeInvalidAssetType  = 8501
	ErrCodeRepositoryFailure = 8502
	ErrCodeNotFound          = 8503
	ErrCodeConflict          = 8504
)

const (
	ErrNilContentAsset      = "assets: content asset entity is nil"
	ErrEmptyContentAssetID  = "assets: content asset id cannot be empty"
	ErrEmptyFilePath        = "assets: file path cannot be empty"
	ErrMissingPostID        = "assets: either content post id or social post id must be provided"
	ErrContentAssetNotFound  = "assets: content asset not found"
	ErrContentAssetExists    = "assets: content asset already exists"
	ErrUnsupportedAssetType  = "assets: unsupported asset type"
	ErrUnableToPersist      = "assets: unable to persist data"
	ErrUnableToFetch        = "assets: unable to fetch data"
	ErrUnableToUpdate       = "assets: unable to update data"
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
