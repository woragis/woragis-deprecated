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
	ErrCodeInvalidVisibility  = 4009
	ErrCodeInvalidSectionType = 4010
	ErrCodeInvalidTechCategory = 4011
	ErrCodeInvalidDiagramType = 4012
)

const (
	ErrNilProject                = "projects: project entity is nil"
	ErrNilMilestone              = "projects: milestone entity is nil"
	ErrEmptyProjectID            = "projects: project id cannot be empty"
	ErrEmptyMilestoneID          = "projects: milestone id cannot be empty"
	ErrEmptyUserID               = "projects: user id cannot be empty"
	ErrEmptyProjectName          = "projects: project name cannot be empty"
	ErrEmptyProjectSlug          = "projects: project slug cannot be empty"
	ErrUnsupportedStatus         = "projects: unsupported status transition"
	ErrHealthScoreOutOfRange     = "projects: health score must be between 0 and 100"
	ErrMetricsMustBePositive     = "projects: metrics must be non-negative"
	ErrEmptyMilestoneTitle       = "projects: milestone title cannot be empty"
	ErrProjectNotFound           = "projects: project not found"
	ErrMilestoneNotFound         = "projects: milestone not found"
	ErrUnableToPersist           = "projects: unable to persist data"
	ErrUnableToFetch             = "projects: unable to fetch data"
	ErrUnableToUpdate            = "projects: unable to update data"
	ErrUnableToDelete            = "projects: unable to delete data"
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
	ErrNilDocumentation          = "projects: documentation entity is nil"
	ErrEmptyDocumentationID      = "projects: documentation id cannot be empty"
	ErrDocumentationNotFound     = "projects: documentation not found"
	ErrNilDocumentationSection   = "projects: documentation section entity is nil"
	ErrEmptySectionID            = "projects: section id cannot be empty"
	ErrEmptySectionTitle         = "projects: section title cannot be empty"
	ErrSectionNotFound           = "projects: documentation section not found"
	ErrUnsupportedVisibility     = "projects: unsupported visibility setting"
	ErrUnsupportedSectionType    = "projects: unsupported section type"
	ErrNilTechnology             = "projects: technology entity is nil"
	ErrEmptyTechnologyID         = "projects: technology id cannot be empty"
	ErrEmptyTechnologyName       = "projects: technology name cannot be empty"
	ErrTechnologyNotFound        = "projects: technology not found"
	ErrUnsupportedTechCategory   = "projects: unsupported technology category"
	ErrNilFileStructure          = "projects: file structure entity is nil"
	ErrEmptyFileStructureID      = "projects: file structure id cannot be empty"
	ErrEmptyFilePath             = "projects: file path cannot be empty"
	ErrFileStructureNotFound     = "projects: file structure not found"
	ErrNilArchitectureDiagram    = "projects: architecture diagram entity is nil"
	ErrEmptyDiagramID            = "projects: diagram id cannot be empty"
	ErrEmptyDiagramTitle         = "projects: diagram title cannot be empty"
	ErrDiagramNotFound           = "projects: architecture diagram not found"
	ErrUnsupportedDiagramType   = "projects: unsupported diagram type"
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
