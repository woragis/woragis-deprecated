package interviewstages

import (
	"time"

	"github.com/google/uuid"
)

// StageType represents the type of interview stage.
type StageType string

const (
	StageTypePhoneScreen    StageType = "phone-screen"
	StageTypeTechnical      StageType = "technical"
	StageTypeBehavioral      StageType = "behavioral"
	StageTypeSystemDesign   StageType = "system-design"
	StageTypeFinal          StageType = "final"
	StageTypeHR             StageType = "hr"
	StageTypeManager        StageType = "manager"
	StageTypePanel          StageType = "panel"
	StageTypeOther          StageType = "other"
)

// StageOutcome represents the outcome of an interview stage.
type StageOutcome string

const (
	StageOutcomePending StageOutcome = "pending"
	StageOutcomePassed  StageOutcome = "passed"
	StageOutcomeFailed  StageOutcome = "failed"
	StageOutcomeCancelled StageOutcome = "cancelled"
)

// InterviewStage represents an interview stage for a job application.
type InterviewStage struct {
	ID                uuid.UUID    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	JobApplicationID  uuid.UUID    `gorm:"column:job_application_id;type:uuid;index;not null" json:"jobApplicationId"`
	StageType         StageType    `gorm:"column:stage_type;type:varchar(50);not null;index" json:"stageType"`
	ScheduledDate     *time.Time   `gorm:"column:scheduled_date" json:"scheduledDate,omitempty"`
	CompletedDate     *time.Time   `gorm:"column:completed_date" json:"completedDate,omitempty"`
	InterviewerName   string       `gorm:"column:interviewer_name;size:255" json:"interviewerName,omitempty"`
	InterviewerEmail  string       `gorm:"column:interviewer_email;size:255" json:"interviewerEmail,omitempty"`
	Location           string       `gorm:"column:location;size:255" json:"location,omitempty"` // "in-person", "video", "phone", or physical address
	Notes              string       `gorm:"column:notes;type:text" json:"notes,omitempty"`
	Feedback           string       `gorm:"column:feedback;type:text" json:"feedback,omitempty"`
	Outcome            StageOutcome `gorm:"column:outcome;type:varchar(20);not null;default:'pending';index" json:"outcome"`
	CreatedAt          time.Time    `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          time.Time    `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for InterviewStage.
func (InterviewStage) TableName() string {
	return "job_application_interview_stages"
}

// NewInterviewStage creates a new interview stage entity.
func NewInterviewStage(jobApplicationID uuid.UUID, stageType StageType) (*InterviewStage, error) {
	stage := &InterviewStage{
		ID:               uuid.New(),
		JobApplicationID: jobApplicationID,
		StageType:        stageType,
		Outcome:          StageOutcomePending,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	return stage, stage.Validate()
}

// Validate ensures interview stage invariants hold.
func (s *InterviewStage) Validate() error {
	if s.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyStageID)
	}
	if s.JobApplicationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyJobApplicationID)
	}
	if !isValidStageType(s.StageType) {
		return NewDomainError(ErrCodeInvalidStageType, ErrUnsupportedStageType)
	}
	if !isValidOutcome(s.Outcome) {
		return NewDomainError(ErrCodeInvalidOutcome, ErrUnsupportedOutcome)
	}
	return nil
}

func isValidStageType(stageType StageType) bool {
	switch stageType {
	case StageTypePhoneScreen, StageTypeTechnical, StageTypeBehavioral,
		StageTypeSystemDesign, StageTypeFinal, StageTypeHR, StageTypeManager,
		StageTypePanel, StageTypeOther:
		return true
	}
	return false
}

func isValidOutcome(outcome StageOutcome) bool {
	switch outcome {
	case StageOutcomePending, StageOutcomePassed, StageOutcomeFailed, StageOutcomeCancelled:
		return true
	}
	return false
}

// Schedule sets the scheduled date for the interview.
func (s *InterviewStage) Schedule(scheduledDate time.Time) {
	s.ScheduledDate = &scheduledDate
	s.UpdatedAt = time.Now().UTC()
}

// Complete marks the interview as completed with an outcome.
func (s *InterviewStage) Complete(completedDate time.Time, outcome StageOutcome) error {
	if !isValidOutcome(outcome) {
		return NewDomainError(ErrCodeInvalidOutcome, ErrUnsupportedOutcome)
	}
	s.CompletedDate = &completedDate
	s.Outcome = outcome
	s.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateInterviewerInfo updates the interviewer information.
func (s *InterviewStage) UpdateInterviewerInfo(name, email string) {
	s.InterviewerName = name
	s.InterviewerEmail = email
	s.UpdatedAt = time.Now().UTC()
}

// UpdateLocation updates the interview location.
func (s *InterviewStage) UpdateLocation(location string) {
	s.Location = location
	s.UpdatedAt = time.Now().UTC()
}

// UpdateNotes updates the interview notes.
func (s *InterviewStage) UpdateNotes(notes string) {
	s.Notes = notes
	s.UpdatedAt = time.Now().UTC()
}

// UpdateFeedback updates the interview feedback.
func (s *InterviewStage) UpdateFeedback(feedback string) {
	s.Feedback = feedback
	s.UpdatedAt = time.Now().UTC()
}

