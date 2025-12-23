package validation

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid email", "test@example.com", false},
		{"valid email with subdomain", "user@mail.example.com", false},
		{"empty email", "", true},
		{"invalid format", "notanemail", true},
		{"missing @", "testexample.com", true},
		{"missing domain", "test@", true},
		{"too long", string(make([]byte, 255)) + "@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name    string
		uuid    string
		wantErr bool
	}{
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", false},
		{"valid UUID uppercase", "550E8400-E29B-41D4-A716-446655440000", false},
		{"empty UUID", "", true},
		{"invalid format", "not-a-uuid", true},
		{"missing dashes", "550e8400e29b41d4a716446655440000", true},
		{"wrong length", "550e8400-e29b-41d4-a716", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUUID(tt.uuid)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUUID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		min      int
		max      int
		fieldName string
		wantErr  bool
	}{
		{"valid string", "hello", 3, 10, "name", false},
		{"too short", "hi", 3, 10, "name", true},
		{"too long", "this is too long", 3, 10, "name", true},
		{"empty with min 0", "", 0, 10, "name", false},
		{"empty with min > 0", "", 1, 10, "name", true},
		{"exact min", "abc", 3, 10, "name", false},
		{"exact max", "1234567890", 3, 10, "name", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateString(tt.input, tt.min, tt.max, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid http", "http://example.com", false},
		{"valid https", "https://example.com", false},
		{"valid with path", "https://example.com/path", false},
		{"valid with query", "https://example.com?key=value", false},
		{"empty URL", "", true},
		{"invalid scheme", "ftp://example.com", true},
		{"no scheme", "example.com", true},
		{"no host", "https://", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNoSQLInjection(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"safe input", "normal text", false},
		{"SQL injection 1", "' OR '1'='1", true},
		{"SQL injection 2", "'; DROP TABLE users;--", true},
		{"SQL injection 3", "'; DELETE FROM users;--", true},
		{"SQL injection 4", "UNION SELECT * FROM users", true},
		{"safe with quotes", "user's name", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNoSQLInjection(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNoSQLInjection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNoXSS(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"safe input", "normal text", false},
		{"XSS script tag", "<script>alert('xss')</script>", true},
		{"XSS javascript", "javascript:alert('xss')", true},
		{"XSS onerror", "onerror=alert('xss')", true},
		{"XSS iframe", "<iframe src='evil.com'></iframe>", true},
		{"safe HTML", "<p>Hello</p>", false}, // Note: This might need stricter validation
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNoXSS(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNoXSS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"normal string", "hello", "hello"},
		{"with null byte", "hello\x00world", "helloworld"},
		{"with whitespace", "  hello  ", "hello"},
		{"with newlines", "\nhello\n", "hello"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeString(tt.input)
			if got != tt.expected {
				t.Errorf("SanitizeString() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestValidateInt(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		min      int
		max      int
		fieldName string
		wantErr  bool
	}{
		{"valid in range", 5, 1, 10, "count", false},
		{"at min", 1, 1, 10, "count", false},
		{"at max", 10, 1, 10, "count", false},
		{"below min", 0, 1, 10, "count", true},
		{"above max", 11, 1, 10, "count", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateInt(tt.value, tt.min, tt.max, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
