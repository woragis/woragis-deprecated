package sender

import "context"

// Message represents an outbound email payload.
type Message struct {
	To       string
	Subject  string
	TextBody string
	HTMLBody string
}

// Sender defines the contract for sending e-mails.
type Sender interface {
	Send(ctx context.Context, msg Message) error
}
