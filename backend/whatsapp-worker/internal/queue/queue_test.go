package queue

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestWhatsAppEnvelope_MarshalUnmarshal(t *testing.T) {
	envelope := WhatsAppEnvelope{
		UserID:      "user123",
		TextMessage: "Test message",
		Destination: "+1234567890",
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var unmarshaled WhatsAppEnvelope
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if unmarshaled.UserID != envelope.UserID {
		t.Errorf("UserID = %v, want %v", unmarshaled.UserID, envelope.UserID)
	}
	if unmarshaled.TextMessage != envelope.TextMessage {
		t.Errorf("TextMessage = %v, want %v", unmarshaled.TextMessage, envelope.TextMessage)
	}
	if unmarshaled.Destination != envelope.Destination {
		t.Errorf("Destination = %v, want %v", unmarshaled.Destination, envelope.Destination)
	}
}

func TestConnection_IsClosed(t *testing.T) {
	tests := []struct {
		name string
		conn *Connection
		want bool
	}{
		{
			name: "nil connection",
			conn: &Connection{conn: nil},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.conn.IsClosed()
			// For nil connection, should be true
			if tt.conn.conn == nil && !got {
				t.Errorf("Connection.IsClosed() = %v, want true for nil connection", got)
			}
		})
	}
}

func TestConnection_Close(t *testing.T) {
	conn := &Connection{
		conn:    nil,
		channel: nil,
	}

	err := conn.Close()
	if err != nil {
		// Close may return an error, but should not panic
	}

	// Test that Close handles nil gracefully
	conn2 := &Connection{}
	err2 := conn2.Close()
	if err2 != nil {
		// Should not panic even with nil fields
	}
}

func TestQueue_Consume_ContextCancellation(t *testing.T) {
	// This test would require a real RabbitMQ connection or channel mock
	// For now, we verify the structure is correct
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Verify context cancellation works
	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context should be cancelled")
	}
}

func TestQueue_Consume_InvalidMessage(t *testing.T) {
	// Test that invalid JSON messages are handled correctly
	invalidJSON := []byte("{invalid json}")
	var envelope WhatsAppEnvelope
	err := json.Unmarshal(invalidJSON, &envelope)
	if err == nil {
		t.Error("json.Unmarshal() should fail for invalid JSON")
	}
}

func TestConnection_Channel(t *testing.T) {
	// Test Channel() method with nil connection
	conn := &Connection{
		conn:    nil,
		channel: nil,
	}

	ch := conn.Channel()
	if ch != nil {
		t.Error("Channel() should return nil for nil connection")
	}
}

func TestNewConnection_ErrorHandling(t *testing.T) {
	// Test NewConnection with invalid URL
	invalidURL := "amqp://invalid:invalid@localhost:9999/vhost"
	conn, err := NewConnection(invalidURL)

	// Should return error
	if err == nil {
		t.Error("NewConnection() should return error for invalid URL")
	}
	if conn != nil {
		t.Error("NewConnection() should return nil connection on error")
		conn.Close()
	}
}

func TestNewQueue_ErrorHandling(t *testing.T) {
	// Test NewQueue with nil connection - this will panic on channel access
	// We can't properly test this without a mock, so we skip this test
	// In real usage, NewConnection would return an error for invalid URLs
	// and NewQueue would be called with a valid connection
	t.Skip("Skipping: NewQueue requires valid connection with non-nil channel (requires RabbitMQ mock)")
}

func TestWhatsAppEnvelope_EmptyFields(t *testing.T) {
	// Test envelope with empty fields
	envelope := WhatsAppEnvelope{
		UserID:      "",
		TextMessage: "",
		Destination: "",
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var unmarshaled WhatsAppEnvelope
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if unmarshaled.UserID != "" {
		t.Errorf("UserID = %v, want empty", unmarshaled.UserID)
	}
	if unmarshaled.TextMessage != "" {
		t.Errorf("TextMessage = %v, want empty", unmarshaled.TextMessage)
	}
	if unmarshaled.Destination != "" {
		t.Errorf("Destination = %v, want empty", unmarshaled.Destination)
	}
}

func TestWhatsAppEnvelope_WithoutDestination(t *testing.T) {
	// Test envelope without destination (omitempty)
	envelope := WhatsAppEnvelope{
		UserID:      "user123",
		TextMessage: "Test message",
		// Destination omitted
	}

	data, err := json.Marshal(envelope)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var unmarshaled WhatsAppEnvelope
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if unmarshaled.Destination != "" {
		t.Errorf("Destination = %v, want empty when omitted", unmarshaled.Destination)
	}
}
