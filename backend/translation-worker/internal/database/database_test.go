package database

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func TestTranslation_SetFields(t *testing.T) {
	translation := &Translation{
		ID:         uuid.New(),
		EntityType: "project",
		EntityID:   uuid.New(),
		Language:   "pt-BR",
		Status:     "pending",
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	fields := map[string]string{
		"name":        "Projeto Teste",
		"description": "Descrição do projeto",
	}

	err := translation.SetFields(fields)
	if err != nil {
		t.Fatalf("SetFields() error = %v", err)
	}

	// Verify fields were set
	retrievedFields, err := translation.GetFields()
	if err != nil {
		t.Fatalf("GetFields() error = %v", err)
	}

	if retrievedFields["name"] != fields["name"] {
		t.Errorf("GetFields() name = %v, want %v", retrievedFields["name"], fields["name"])
	}
	if retrievedFields["description"] != fields["description"] {
		t.Errorf("GetFields() description = %v, want %v", retrievedFields["description"], fields["description"])
	}
}

func TestTranslation_GetFields(t *testing.T) {
	fields := map[string]string{
		"name":        "Test Project",
		"description": "Test Description",
	}

	fieldsJSON, _ := json.Marshal(fields)
	translation := &Translation{
		ID:         uuid.New(),
		EntityType: "project",
		EntityID:   uuid.New(),
		Language:   "en",
		Fields:     datatypes.JSON(fieldsJSON),
		Status:     "completed",
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	retrievedFields, err := translation.GetFields()
	if err != nil {
		t.Fatalf("GetFields() error = %v", err)
	}

	if retrievedFields["name"] != fields["name"] {
		t.Errorf("GetFields() name = %v, want %v", retrievedFields["name"], fields["name"])
	}
	if retrievedFields["description"] != fields["description"] {
		t.Errorf("GetFields() description = %v, want %v", retrievedFields["description"], fields["description"])
	}
}

func TestTranslation_GetFields_Empty(t *testing.T) {
	translation := &Translation{
		ID:         uuid.New(),
		EntityType: "project",
		EntityID:   uuid.New(),
		Language:   "en",
		Fields:     datatypes.JSON("{}"),
		Status:     "pending",
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	fields, err := translation.GetFields()
	if err != nil {
		t.Fatalf("GetFields() error = %v", err)
	}

	if len(fields) != 0 {
		t.Errorf("GetFields() should return empty map, got %v", fields)
	}
}

func TestTranslation_TableName(t *testing.T) {
	var translation Translation
	if translation.TableName() != "translations" {
		t.Errorf("TableName() = %v, want translations", translation.TableName())
	}
}

// Note: Full repository tests would require a test database or more sophisticated mocking
// These tests verify the basic functionality of the Translation entity methods
