package responses

import (
	"time"

	"github.com/google/uuid"
)

// ResponseType represents the type of response received.
type ResponseType string

const (
	ResponseTypeRejection  ResponseType = "rejection"
	ResponseTypeInterview  ResponseType = "interview"
	ResponseTypeOffer      ResponseType = "offer"
	ResponseTypeNoResponse ResponseType = "no-response"
)

// Response represents a response received for a job application.
type Response struct {
	ID                uuid.UUID    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	JobApplicationID  uuid.UUID    `gorm:"column:job_application_id;type:uuid;index;not null" json:"jobApplicationId"`
	ResponseType      ResponseType `gorm:"column:response_type;type:varchar(20);not null;index" json:"responseType"`
	ResponseDate      time.Time    `gorm:"column:response_date;not null" json:"responseDate"`
	Message           string       `gorm:"column:message;type:text" json:"message,omitempty"`
	ContactPerson     string       `gorm:"column:contact_person;size:255" json:"contactPerson,omitempty"`
	ContactEmail      string       `gorm:"column:contact_email;size:255" json:"contactEmail,omitempty"`
	ContactPhone      string       `gorm:"column:contact_phone;size:50" json:"contactPhone,omitempty"`
	ResponseChannel   string       `gorm:"column:response_channel;size:50" json:"responseChannel,omitempty"` // "email", "phone", "linkedin", etc.
	CreatedAt         time.Time    `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt         time.Time    `gorm:"column:updated_at" json:"updatedAt"`
}

// TableName specifies the table name for Response.
func (Response) TableName() string {
	return "job_application_responses"
}

// NewResponse creates a new response entity.
func NewResponse(jobApplicationID uuid.UUID, responseType ResponseType, responseDate time.Time) (*Response, error) {
	resp := &Response{
		ID:               uuid.New(),
		JobApplicationID: jobApplicationID,
		ResponseType:     responseType,
		ResponseDate:     responseDate,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	return resp, resp.Validate()
}

// Validate ensures response invariants hold.
func (r *Response) Validate() error {
	if r.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyResponseID)
	}
	if r.JobApplicationID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyJobApplicationID)
	}
	if !isValidResponseType(r.ResponseType) {
		return NewDomainError(ErrCodeInvalidResponseType, ErrUnsupportedResponseType)
	}
	return nil
}

func isValidResponseType(responseType ResponseType) bool {
	switch responseType {
	case ResponseTypeRejection, ResponseTypeInterview, ResponseTypeOffer, ResponseTypeNoResponse:
		return true
	}
	return false
}

// UpdateMessage updates the response message.
func (r *Response) UpdateMessage(message string) {
	r.Message = message
	r.UpdatedAt = time.Now().UTC()
}

// UpdateContactInfo updates the contact information.
func (r *Response) UpdateContactInfo(name, email, phone string) {
	r.ContactPerson = name
	r.ContactEmail = email
	r.ContactPhone = phone
	r.UpdatedAt = time.Now().UTC()
}

// UpdateResponseChannel updates the response channel.
func (r *Response) UpdateResponseChannel(channel string) {
	r.ResponseChannel = channel
	r.UpdatedAt = time.Now().UTC()
}

