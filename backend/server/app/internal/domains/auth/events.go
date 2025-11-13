package auth

import "time"

// DomainEvent represents internal domain happenings for observability.
type DomainEvent struct {
	Name      string
	Timestamp time.Time
	Payload   map[string]any
}

// NewDomainEvent creates a domain event with current timestamp.
func NewDomainEvent(name string, payload map[string]any) DomainEvent {
	return DomainEvent{
		Name:      name,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}
}
