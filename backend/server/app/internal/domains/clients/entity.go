package clients

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Client represents a client or contact that can receive WhatsApp messages.
type Client struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;index;not null" json:"userId"`
	Name        string    `gorm:"column:name;size:120;not null" json:"name"`
	Email       string    `gorm:"column:email;size:255;index" json:"email,omitempty"`
	PhoneNumber string    `gorm:"column:phone_number;size:20;index;not null" json:"phoneNumber"`
	Company     string    `gorm:"column:company;size:120" json:"company,omitempty"`
	Notes       string    `gorm:"column:notes;type:text" json:"notes,omitempty"`
	IsArchived  bool      `gorm:"column:is_archived;not null;default:false;index" json:"isArchived"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

// NewClient creates a new Client with the provided fields.
func NewClient(userID uuid.UUID, name, phoneNumber string) (*Client, error) {
	now := time.Now().UTC()
	client := &Client{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        strings.TrimSpace(name),
		PhoneNumber: normalizePhoneNumber(phoneNumber),
		IsArchived:  false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return client, client.Validate()
}

// Validate enforces domain invariants for the client entity.
func (c *Client) Validate() error {
	if c == nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrNilClient)
	}

	if c.ID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyClientID)
	}

	if c.UserID == uuid.Nil {
		return NewDomainError(ErrCodeInvalidPayload, ErrEmptyUserID)
	}

	if c.Name == "" {
		return NewDomainError(ErrCodeInvalidName, ErrEmptyName)
	}

	if len(c.Name) > 120 {
		return NewDomainError(ErrCodeInvalidName, ErrNameTooLong)
	}

	if c.PhoneNumber == "" {
		return NewDomainError(ErrCodeInvalidPhoneNumber, ErrEmptyPhoneNumber)
	}

	if len(c.PhoneNumber) > 20 {
		return NewDomainError(ErrCodeInvalidPhoneNumber, ErrPhoneNumberTooLong)
	}

	if c.Email != "" {
		if len(c.Email) > 255 {
			return NewDomainError(ErrCodeInvalidEmail, ErrEmailTooLong)
		}
	}

	if c.Company != "" && len(c.Company) > 120 {
		return NewDomainError(ErrCodeInvalidPayload, ErrCompanyTooLong)
	}

	return nil
}

// UpdateMutableFields updates mutable attributes and preserves invariants.
func (c *Client) UpdateMutableFields(name, email, phoneNumber, company, notes *string) error {
	if name != nil {
		c.Name = strings.TrimSpace(*name)
	}
	if email != nil {
		c.Email = strings.TrimSpace(*email)
	}
	if phoneNumber != nil {
		c.PhoneNumber = normalizePhoneNumber(*phoneNumber)
	}
	if company != nil {
		c.Company = strings.TrimSpace(*company)
	}
	if notes != nil {
		c.Notes = strings.TrimSpace(*notes)
	}
	c.UpdatedAt = time.Now().UTC()
	return c.Validate()
}

// ToggleArchived sets archived flag.
func (c *Client) ToggleArchived(archived bool) {
	c.IsArchived = archived
	c.UpdatedAt = time.Now().UTC()
}

// normalizePhoneNumber normalizes phone number format.
// Removes spaces, dashes, and parentheses, but keeps + prefix if present.
func normalizePhoneNumber(phone string) string {
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return ""
	}

	// Remove common separators but keep + prefix
	normalized := strings.Builder{}
	for i, r := range phone {
		if i == 0 && r == '+' {
			normalized.WriteRune(r)
		} else if r >= '0' && r <= '9' {
			normalized.WriteRune(r)
		}
	}

	return normalized.String()
}

