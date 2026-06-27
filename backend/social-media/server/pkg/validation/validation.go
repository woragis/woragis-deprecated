package validation

import (
	"fmt"
	"regexp"
	"strings"
)

// ValidateString validates a string length
func ValidateString(s string, min, max int, fieldName string) error {
	if len(s) < min {
		return fmt.Errorf("%s must be at least %d characters", fieldName, min)
	}
	if len(s) > max {
		return fmt.Errorf("%s must be at most %d characters", fieldName, max)
	}
	return nil
}

// ValidateNoSQLInjection checks for SQL injection patterns
func ValidateNoSQLInjection(s string) error {
	sqlPatterns := []string{
		"(?i)(union|select|insert|update|delete|drop|create|alter|exec|execute)",
		"(?i)(--|;|'|\"|/\\*|\\*/)",
		"(?i)(or|and)\\s+\\d+\\s*=\\s*\\d+",
	}
	
	for _, pattern := range sqlPatterns {
		matched, err := regexp.MatchString(pattern, s)
		if err != nil {
			continue
		}
		if matched {
			return fmt.Errorf("potential SQL injection detected")
		}
	}
	return nil
}

// ValidateNoXSS checks for XSS patterns
func ValidateNoXSS(s string) error {
	xssPatterns := []string{
		"(?i)<script[^>]*>.*?</script>",
		"(?i)javascript:",
		"(?i)on\\w+\\s*=",
		"(?i)<iframe[^>]*>",
		"(?i)<object[^>]*>",
		"(?i)<embed[^>]*>",
	}
	
	for _, pattern := range xssPatterns {
		matched, err := regexp.MatchString(pattern, s)
		if err != nil {
			continue
		}
		if matched {
			return fmt.Errorf("potential XSS detected")
		}
	}
	return nil
}

// ValidateURL validates a URL format
func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("URL cannot be empty")
	}
	
	urlPattern := `^https?://[^\s/$.?#].[^\s]*$`
	matched, err := regexp.MatchString(urlPattern, url)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}
	if !matched {
		return fmt.Errorf("invalid URL format")
	}
	return nil
}

// ValidateUUID validates a UUID format
func ValidateUUID(uuid string) error {
	if uuid == "" {
		return fmt.Errorf("UUID cannot be empty")
	}
	
	uuidPattern := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	matched, err := regexp.MatchString(uuidPattern, strings.ToLower(uuid))
	if err != nil {
		return fmt.Errorf("invalid UUID format")
	}
	if !matched {
		return fmt.Errorf("invalid UUID format")
	}
	return nil
}

// ValidateFileExtension validates file extension
func ValidateFileExtension(filename string, allowedExtensions []string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}
	
	ext := strings.ToLower(strings.TrimPrefix(strings.ToLower(filename[strings.LastIndex(filename, "."):]), "."))
	
	for _, allowed := range allowedExtensions {
		if strings.ToLower(allowed) == ext {
			return nil
		}
	}
	
	return fmt.Errorf("file extension %s is not allowed. Allowed extensions: %v", ext, allowedExtensions)
}

// ValidateFileSize validates file size
func ValidateFileSize(sizeBytes int64, maxBytes int64) error {
	if sizeBytes < 0 {
		return fmt.Errorf("file size cannot be negative")
	}
	if sizeBytes > maxBytes {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size of %d bytes", sizeBytes, maxBytes)
	}
	return nil
}
