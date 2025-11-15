package whatsapp

import (
	"context"
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

// userRepositoryAdapter adapts auth domain repository to UserRepository interface.
type userRepositoryAdapter struct {
	repo interface{}
}

// NewUserRepositoryAdapter creates an adapter for user repository.
// The repo parameter should be an auth domain repository.
func NewUserRepositoryAdapter(repo interface{}) UserRepository {
	return &userRepositoryAdapter{repo: repo}
}

// GetUserPhoneNumber retrieves a user's phone number by user ID.
func (a *userRepositoryAdapter) GetUserPhoneNumber(ctx context.Context, userID uuid.UUID) (string, error) {
	// Use reflection to call FindByID method
	repoValue := reflect.ValueOf(a.repo)
	findByIDMethod := repoValue.MethodByName("FindByID")
	if !findByIDMethod.IsValid() {
		return "", fmt.Errorf("FindByID method not found on repository")
	}

	// Call FindByID(ctx, userID)
	results := findByIDMethod.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(userID),
	})

	if len(results) != 2 {
		return "", fmt.Errorf("unexpected number of return values from FindByID")
	}

	// Check for error
	if !results[1].IsNil() {
		// Return domain error for user not found
		return "", NewDomainError(ErrCodeUserNotFound, ErrUserNotFound)
	}

	// Get the user object (first return value)
	userValue := results[0]
	if userValue.IsNil() {
		return "", NewDomainError(ErrCodeUserNotFound, ErrUserNotFound)
	}

	// Extract PhoneNumber field using reflection
	return extractPhoneNumberViaReflection(userValue.Interface())
}

// extractPhoneNumberFromUser extracts phone number from user using type assertion
func extractPhoneNumberFromUser(user interface{}) (string, error) {
	// Try common patterns for accessing PhoneNumber
	switch u := user.(type) {
	case interface{ PhoneNumber() string }:
		return u.PhoneNumber(), nil
	case interface{ GetPhoneNumber() string }:
		return u.GetPhoneNumber(), nil
	default:
		// Use reflection as last resort
		return extractPhoneNumberViaReflection(user)
	}
}

// clientRepositoryAdapter adapts clients domain repository to ClientRepository interface.
type clientRepositoryAdapter struct {
	repo interface{}
}

// NewClientRepositoryAdapter creates an adapter for client repository.
// The repo parameter should be a clients domain repository.
func NewClientRepositoryAdapter(repo interface{}) ClientRepository {
	return &clientRepositoryAdapter{repo: repo}
}

// GetClientPhoneNumber retrieves a client's phone number by client ID.
func (a *clientRepositoryAdapter) GetClientPhoneNumber(ctx context.Context, userID, clientID uuid.UUID) (string, error) {
	// Use reflection to call Get method
	repoValue := reflect.ValueOf(a.repo)
	getMethod := repoValue.MethodByName("Get")
	if !getMethod.IsValid() {
		return "", fmt.Errorf("Get method not found on repository")
	}

	// Call Get(ctx, userID, clientID)
	results := getMethod.Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(userID),
		reflect.ValueOf(clientID),
	})

	if len(results) != 2 {
		return "", fmt.Errorf("unexpected number of return values from Get")
	}

	// Check for error
	if !results[1].IsNil() {
		// Return domain error for client not found
		return "", NewDomainError(ErrCodeClientNotFound, ErrClientNotFound)
	}

	// Get the client object (first return value)
	clientValue := results[0]
	if clientValue.IsNil() {
		return "", NewDomainError(ErrCodeClientNotFound, ErrClientNotFound)
	}

	// Extract PhoneNumber field using reflection
	return extractPhoneNumberViaReflection(clientValue.Interface())
}

// extractPhoneNumberFromClient extracts phone number from client using type assertion
func extractPhoneNumberFromClient(client interface{}) (string, error) {
	switch c := client.(type) {
	case interface{ PhoneNumber() string }:
		return c.PhoneNumber(), nil
	case interface{ GetPhoneNumber() string }:
		return c.GetPhoneNumber(), nil
	default:
		return extractPhoneNumberViaReflection(client)
	}
}

// extractPhoneNumberViaReflection uses reflection to extract PhoneNumber field
func extractPhoneNumberViaReflection(obj interface{}) (string, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected struct, got %v", v.Kind())
	}

	field := v.FieldByName("PhoneNumber")
	if !field.IsValid() {
		return "", fmt.Errorf("PhoneNumber field not found")
	}

	if field.Kind() != reflect.String {
		return "", fmt.Errorf("PhoneNumber field is not a string")
	}

	return field.String(), nil
}
