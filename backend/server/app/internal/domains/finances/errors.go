package finances

import "errors"

const (
	ErrCodeInvalidPayload    = 2000
	ErrCodeInvalidType       = 2001
	ErrCodeInvalidCategory   = 2002
	ErrCodeInvalidAmount     = 2003
	ErrCodeInvalidCurrency   = 2004
	ErrCodeRepositoryFailure = 2005
	ErrCodeSummaryFailure    = 2006
)

const (
	ErrNilTransaction             = "finances: transaction entity is nil"
	ErrEmptyTransactionID         = "finances: transaction id cannot be empty"
	ErrEmptyUserID                = "finances: user id cannot be empty"
	ErrUnsupportedTransactionType = "finances: unsupported transaction type"
	ErrEmptyCategory              = "finances: category cannot be empty"
	ErrAmountMustBePositive       = "finances: amount must be positive"
	ErrEmptyCurrency              = "finances: currency cannot be empty"
	ErrCurrencyMustBeISO          = "finances: currency must be a 3-letter ISO code"
	ErrUnableToPersist            = "finances: unable to persist transaction"
	ErrUnableToFetch              = "finances: unable to fetch transactions"
	ErrUnableToSummarize          = "finances: unable to build summary"
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
