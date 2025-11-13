package projects

import "errors"

const (
	ErrCodeInvalidPayload     = 4000
	ErrCodeInvalidName        = 4001
	ErrCodeInvalidStatus      = 4002
	ErrCodeInvalidHealthScore = 4003
	ErrCodeInvalidMetrics     = 4004
	ErrCodeRepositoryFailure  = 4005
	ErrCodeNotFound           = 4006
	ErrCodeConflict           = 4007
	ErrCodeInvalidDependency  = 4008
)

const (
	ErrNilProject                = "projects: project entity is nil"
	ErrNilMilestone              = "projects: milestone entity is nil"
	ErrEmptyProjectID            = "projects: project id cannot be empty"
	ErrEmptyMilestoneID          = "projects: milestone id cannot be empty"
	ErrEmptyUserID               = "projects: user id cannot be empty"
	ErrEmptyProjectName          = "projects: project name cannot be empty"
	ErrUnsupportedStatus         = "projects: unsupported status transition"
	ErrHealthScoreOutOfRange     = "projects: health score must be between 0 and 100"
	ErrMetricsMustBePositive     = "projects: metrics must be non-negative"
	ErrEmptyMilestoneTitle       = "projects: milestone title cannot be empty"
	ErrProjectNotFound           = "projects: project not found"
	ErrMilestoneNotFound         = "projects: milestone not found"
	ErrUnableToPersist           = "projects: unable to persist data"
	ErrUnableToFetch             = "projects: unable to fetch data"
	ErrUnableToUpdate            = "projects: unable to update data"
	ErrNilKanbanColumn           = "projects: kanban column entity is nil"
	ErrEmptyKanbanColumnID       = "projects: kanban column id cannot be empty"
	ErrKanbanColumnNotFound      = "projects: kanban column not found"
	ErrEmptyKanbanColumnName     = "projects: kanban column name cannot be empty"
	ErrInvalidKanbanPosition     = "projects: kanban position must be zero or positive"
	ErrInvalidWIPLimit           = "projects: kanban WIP limit cannot be negative"
	ErrNilKanbanCard             = "projects: kanban card entity is nil"
	ErrEmptyKanbanCardID         = "projects: kanban card id cannot be empty"
	ErrKanbanCardNotFound        = "projects: kanban card not found"
	ErrEmptyKanbanCardTitle      = "projects: kanban card title cannot be empty"
	ErrSelfDependencyNotAllowed  = "projects: project cannot depend on itself"
	ErrNilDependency             = "projects: dependency entity is nil"
	ErrEmptyDependencyID         = "projects: dependency id cannot be empty"
	ErrUnsupportedDependencyType = "projects: unsupported dependency type"
	ErrDependencyAlreadyExists   = "projects: dependency already exists"
	ErrDependencyNotFound        = "projects: dependency not found"
	ErrInvalidColumnOrder        = "projects: column order payload must include all columns"
	ErrWIPLimitExceeded          = "projects: kanban column WIP limit reached"
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
