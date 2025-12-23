package validation

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	// EmailRegex validates email format
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	
	// UUIDRegex validates UUID format
	UUIDRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
)

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	
	if len(email) > 254 {
		return fmt.Errorf("email is too long (maximum 254 characters)")
	}
	
	if !EmailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	
	return nil
}

// ValidateUUID validates a UUID string
func ValidateUUID(uuid string) error {
	if uuid == "" {
		return fmt.Errorf("UUID is required")
	}
	
	if !UUIDRegex.MatchString(strings.ToLower(uuid)) {
		return fmt.Errorf("invalid UUID format")
	}
	
	return nil
}

// ValidateString validates a string with length constraints
func ValidateString(input string, min, max int, fieldName string) error {
	if input == "" && min > 0 {
		return fmt.Errorf("%s is required", fieldName)
	}
	
	length := utf8.RuneCountInString(input)
	if length < min {
		return fmt.Errorf("%s is too short (minimum %d characters)", fieldName, min)
	}
	if length > max {
		return fmt.Errorf("%s is too long (maximum %d characters)", fieldName, max)
	}
	
	return nil
}

// ValidateURL validates a URL
func ValidateURL(urlString string) error {
	if urlString == "" {
		return fmt.Errorf("URL is required")
	}
	
	u, err := url.Parse(urlString)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}
	
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme")
	}
	
	if u.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	
	return nil
}

// ValidateInt validates an integer within a range
func ValidateInt(value, min, max int, fieldName string) error {
	if value < min {
		return fmt.Errorf("%s is too small (minimum %d)", fieldName, min)
	}
	if value > max {
		return fmt.Errorf("%s is too large (maximum %d)", fieldName, max)
	}
	return nil
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	return input
}

// ValidateNoSQLInjection checks for SQL injection patterns
func ValidateNoSQLInjection(input string) error {
	// Common SQL injection patterns
	sqlPatterns := []string{
		"' OR '1'='1",
		"' OR '1'='1'--",
		"' OR '1'='1'/*",
		"'; DROP TABLE",
		"'; DELETE FROM",
		"'; UPDATE",
		"UNION SELECT",
		"'; INSERT INTO",
	}
	
	lowerInput := strings.ToLower(input)
	for _, pattern := range sqlPatterns {
		if strings.Contains(lowerInput, strings.ToLower(pattern)) {
			return fmt.Errorf("potentially dangerous input detected")
		}
	}
	
	return nil
}

// ValidateNoXSS checks for XSS patterns
func ValidateNoXSS(input string) error {
	// Common XSS patterns
	xssPatterns := []string{
		"<script",
		"</script>",
		"javascript:",
		"onerror=",
		"onload=",
		"onclick=",
		"<iframe",
		"<img",
		"<svg",
	}
	
	lowerInput := strings.ToLower(input)
	for _, pattern := range xssPatterns {
		if strings.Contains(lowerInput, strings.ToLower(pattern)) {
			return fmt.Errorf("potentially dangerous input detected")
		}
	}
	
	return nil
}

// ValidateFileExtension validates file extension
func ValidateFileExtension(filename string, allowedExts []string) error {
	if filename == "" {
		return fmt.Errorf("filename is required")
	}
	
	ext := strings.ToLower(filename[strings.LastIndex(filename, "."):])
	
	for _, allowedExt := range allowedExts {
		if ext == strings.ToLower(allowedExt) {
			return nil
		}
	}
	
	return fmt.Errorf("invalid file extension. Allowed: %v", allowedExts)
}

// ValidateFileSize validates file size
func ValidateFileSize(size, maxSize int64) error {
	if size <= 0 {
		return fmt.Errorf("file size must be greater than 0")
	}
	if size > maxSize {
		return fmt.Errorf("file too large (maximum %d bytes)", maxSize)
	}
	return nil
}
