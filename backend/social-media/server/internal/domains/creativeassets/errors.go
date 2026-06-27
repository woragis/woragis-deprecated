package creativeassets

// Domain error codes
const (
	ErrCodeInvalidPayload = 10300
	ErrCodeInvalidType    = 10301
	ErrCodeNotFound       = 10302
	ErrCodeUnauthorized   = 10303
	ErrCodeNotImplemented = 10304
)

// Domain error messages
const (
	ErrNilCreativeAsset    = "creative asset is nil"
	ErrEmptyAssetID        = "asset ID is empty"
	ErrEmptyUserID         = "user ID is empty"
	ErrEmptyEntityID       = "entity ID is empty"
	ErrUnsupportedEntityType = "unsupported entity type"
	ErrUnsupportedAssetType  = "unsupported asset type"
	ErrUnsupportedPurpose    = "unsupported purpose"
	ErrAssetNotFound         = "creative asset not found"
	ErrUnauthorizedAccess    = "unauthorized access to creative asset"
)

// DomainError represents a domain-level error
type DomainError struct {
	Code    int
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a new domain error
func NewDomainError(code int, message string) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
	}
}

