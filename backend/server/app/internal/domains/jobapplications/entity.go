package jobapplications

import (
	"time"

	"github.com/google/uuid"
)

// ApplicationStatus represents the status of a job application.
type ApplicationStatus string

const (
	ApplicationStatusPending   ApplicationStatus = "pending"
	ApplicationStatusProcessing ApplicationStatus = "processing"
	ApplicationStatusApplied    ApplicationStatus = "applied"
	ApplicationStatusContacted  ApplicationStatus = "contacted"
	ApplicationStatusRejected   ApplicationStatus = "rejected"
	ApplicationStatusAccepted   ApplicationStatus = "accepted"
	ApplicationStatusFailed     ApplicationStatus = "failed"
)

// JobApplication represents a job application record.
type JobApplication struct {
	ID              uuid.UUID         `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID          uuid.UUID         `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	CompanyName     string            `gorm:"column:company_name;size:255;not null" json:"companyName"`
	Location        string            `gorm:"column:location;size:255" json:"location"`
	JobTitle        string            `gorm:"column:job_title;size:255;not null" json:"jobTitle"`
	JobURL          string            `gorm:"column:job_url;size:512;not null" json:"jobUrl"`
	Website         string            `gorm:"column:website;size:50;not null;index" json:"website"` // "linkedin", "glassdoor", etc.
	AppliedAt       *time.Time        `gorm:"column:applied_at" json:"appliedAt,omitempty"`
	CoverLetter     string            `gorm:"column:cover_letter;type:text" json:"coverLetter"`
	LinkedInContact bool              `gorm:"column:linkedin_contact;not null;default:false" json:"linkedInContact"`
	Status          ApplicationStatus `gorm:"column:status;type:varchar(20);not null;default:'pending';index" json:"status"`
	ErrorMessage    string            `gorm:"column:error_message;type:text" json:"errorMessage,omitempty"`
	CreatedAt       time.Time         `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt       time.Time         `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for JobApplication.
func (JobApplication) TableName() string {
	return "job_applications"
}

// NewJobApplication creates a new job application entity.
func NewJobApplication(userID uuid.UUID, companyName, location, jobTitle, jobURL, website string) (*JobApplication, error) {
	app := &JobApplication{
		ID:          uuid.New(),
		UserID:      userID,
		CompanyName: companyName,
		Location:    location,
		JobTitle:    jobTitle,
		JobURL:      jobURL,
		Website:     website,
		Status:      ApplicationStatusPending,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	return app, app.Validate()
}

// Validate ensures job application invariants hold.
func (j *JobApplication) Validate() error {
	if j.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyApplicationID)
	}
	if j.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}
	if j.CompanyName == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyCompanyName)
	}
	if j.JobTitle == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyJobTitle)
	}
	if j.JobURL == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyJobURL)
	}
	if j.Website == "" {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyWebsite)
	}
	if !isValidStatus(j.Status) {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}
	return nil
}

func isValidStatus(status ApplicationStatus) bool {
	switch status {
	case ApplicationStatusPending, ApplicationStatusProcessing, ApplicationStatusApplied,
		ApplicationStatusContacted, ApplicationStatusRejected, ApplicationStatusAccepted, ApplicationStatusFailed:
		return true
	}
	return false
}

// MarkApplied updates the application status to applied and sets the applied timestamp.
func (j *JobApplication) MarkApplied(coverLetter string) {
	now := time.Now().UTC()
	j.Status = ApplicationStatusApplied
	j.AppliedAt = &now
	j.CoverLetter = coverLetter
	j.UpdatedAt = now
}

// MarkContacted updates the application status to contacted.
func (j *JobApplication) MarkContacted() {
	j.Status = ApplicationStatusContacted
	j.LinkedInContact = true
	j.UpdatedAt = time.Now().UTC()
}

// MarkFailed updates the application status to failed with an error message.
func (j *JobApplication) MarkFailed(errorMessage string) {
	j.Status = ApplicationStatusFailed
	j.ErrorMessage = errorMessage
	j.UpdatedAt = time.Now().UTC()
}

// MarkProcessing updates the application status to processing.
func (j *JobApplication) MarkProcessing() {
	j.Status = ApplicationStatusProcessing
	j.UpdatedAt = time.Now().UTC()
}

// UpdateStatus updates the application status.
func (j *JobApplication) UpdateStatus(status ApplicationStatus) error {
	if !isValidStatus(status) {
		return NewDomainError(ErrCodeInvalidStatus, ErrUnsupportedStatus)
	}
	j.Status = status
	j.UpdatedAt = time.Now().UTC()
	return nil
}

// JobApplicationJob represents a job in the Redis queue.
type JobApplicationJob struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	CompanyName string  `json:"companyName"`
	Location  string    `json:"location"`
	JobTitle  string    `json:"jobTitle"`
	JobURL    string    `json:"jobUrl"`
	Website   string    `json:"website"`
}

