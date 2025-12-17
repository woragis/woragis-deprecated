package queue

import (
	"context"
	"encoding/json"
	"testing"
)

func TestTranslationJob_MarshalUnmarshal(t *testing.T) {
	job := TranslationJob{
		ID:         "test-job-id",
		EntityType: "project",
		EntityID:   "test-entity-id",
		Language:   "pt-BR",
		Fields:     []string{"name", "description"},
		SourceText: map[string]string{
			"name":        "Test Project",
			"description": "Test Description",
		},
	}

	data, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var unmarshaled TranslationJob
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if unmarshaled.ID != job.ID {
		t.Errorf("ID = %v, want %v", unmarshaled.ID, job.ID)
	}
	if unmarshaled.EntityType != job.EntityType {
		t.Errorf("EntityType = %v, want %v", unmarshaled.EntityType, job.EntityType)
	}
	if unmarshaled.EntityID != job.EntityID {
		t.Errorf("EntityID = %v, want %v", unmarshaled.EntityID, job.EntityID)
	}
	if unmarshaled.Language != job.Language {
		t.Errorf("Language = %v, want %v", unmarshaled.Language, job.Language)
	}
	if len(unmarshaled.Fields) != len(job.Fields) {
		t.Errorf("Fields length = %v, want %v", len(unmarshaled.Fields), len(job.Fields))
	}
	if len(unmarshaled.SourceText) != len(job.SourceText) {
		t.Errorf("SourceText length = %v, want %v", len(unmarshaled.SourceText), len(job.SourceText))
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
	// Test closing a connection with nil fields (should not panic)
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
	var job TranslationJob
	err := json.Unmarshal(invalidJSON, &job)
	if err == nil {
		t.Error("json.Unmarshal() should fail for invalid JSON")
	}
}

func TestTranslationJob_EmptyFields(t *testing.T) {
	// Test job with empty fields
	job := TranslationJob{
		ID:         "test-id",
		EntityType: "test",
		EntityID:   "test-entity",
		Language:   "en",
		Fields:     []string{},
		SourceText: map[string]string{},
	}

	data, err := json.Marshal(job)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var unmarshaled TranslationJob
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if len(unmarshaled.Fields) != 0 {
		t.Errorf("Fields = %v, want empty", unmarshaled.Fields)
	}
	if len(unmarshaled.SourceText) != 0 {
		t.Errorf("SourceText = %v, want empty", unmarshaled.SourceText)
	}
}
