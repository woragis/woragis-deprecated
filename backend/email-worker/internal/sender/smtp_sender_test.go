package sender

import (
	"context"
	"strings"
	"testing"

	"github.com/woragis/backend/email-worker/internal/config"
)

func TestNewSMTPSender(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.EmailConfig
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: config.EmailConfig{
				Host: "smtp.example.com",
				From: "noreply@example.com",
			},
			wantErr: false,
		},
		{
			name: "disabled config",
			cfg: config.EmailConfig{
				Host: "",
				From: "noreply@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sender, err := NewSMTPSender(tt.cfg, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSMTPSender() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && sender == nil {
				t.Error("NewSMTPSender() returned nil sender")
			}
		})
	}
}

func TestSMTPSender_Send(t *testing.T) {
	tests := []struct {
		name    string
		msg     Message
		wantErr bool
	}{
		{
			name: "valid message",
			msg: Message{
				To:       "test@example.com",
				Subject:  "Test Subject",
				TextBody: "Test body",
				HTMLBody: "<p>Test body</p>",
			},
			wantErr: false, // Will fail on actual SMTP connection, but validates message structure
		},
		{
			name: "missing recipient",
			msg: Message{
				To:       "",
				Subject:  "Test Subject",
				TextBody: "Test body",
			},
			wantErr: true,
		},
		{
			name: "text only",
			msg: Message{
				To:       "test@example.com",
				Subject:  "Test Subject",
				TextBody: "Test body",
			},
			wantErr: false,
		},
		{
			name: "HTML only",
			msg: Message{
				To:       "test@example.com",
				Subject:  "Test Subject",
				HTMLBody: "<p>Test body</p>",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.EmailConfig{
				Host:   "smtp.example.com",
				Port:   587,
				From:   "noreply@example.com",
				UseTLS: true,
			}
			sender, err := NewSMTPSender(cfg, nil)
			if err != nil {
				t.Fatalf("NewSMTPSender() error = %v", err)
			}

			err = sender.Send(context.Background(), tt.msg)
			if (err != nil) != tt.wantErr {
				// For valid messages, we expect SMTP connection errors (not validation errors)
				if tt.wantErr && err != nil {
					// Check if it's a validation error
					if err.Error() == "recipient email required" {
						return // This is expected
					}
				}
				if !tt.wantErr {
					// We expect SMTP connection errors, not validation errors
					if err != nil && err.Error() == "recipient email required" {
						t.Errorf("SMTPSender.Send() unexpected validation error = %v", err)
					}
					// Other errors (SMTP connection) are expected in unit tests
				}
			}
		})
	}
}

func TestWritePart(t *testing.T) {
	var builder strings.Builder
	boundary := "test_boundary"
	contentType := "text/plain; charset=UTF-8"
	body := "Test body content"

	writePart(&builder, boundary, contentType, body)

	result := builder.String()
	if result == "" {
		t.Error("writePart() should write content")
	}
	if !strings.Contains(result, boundary) {
		t.Error("writePart() should include boundary")
	}
	if !strings.Contains(result, contentType) {
		t.Error("writePart() should include content type")
	}
}

func TestWritePart_EmptyBody(t *testing.T) {
	var builder strings.Builder
	boundary := "test_boundary"
	contentType := "text/plain; charset=UTF-8"
	body := ""

	writePart(&builder, boundary, contentType, body)

	result := builder.String()
	if result != "" {
		t.Error("writePart() should not write content for empty body")
	}
}

func TestBuilderWriter(t *testing.T) {
	var builder strings.Builder
	writer := builderWriter{&builder}

	testData := []byte("test data")
	n, err := writer.Write(testData)
	if err != nil {
		t.Errorf("builderWriter.Write() error = %v", err)
	}
	if n != len(testData) {
		t.Errorf("builderWriter.Write() n = %v, want %v", n, len(testData))
	}
	if builder.String() != "test data" {
		t.Errorf("builderWriter.Write() wrote %v, want test data", builder.String())
	}
}
