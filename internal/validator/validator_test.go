package validator

import (
	"testing"
)

type Address struct {
	City    string `json:"city" validate:"required"`
	ZipCode int    `json:"zip_code" validate:"required"`
}

type User struct {
	Name    string   `json:"name" validate:"required"`
	Age     int      `json:"age"`
	Address Address  `json:"address" validate:"required"`
	Tags    []string `json:"tags"`
}

func TestValidateJSON_ValidInput(t *testing.T) {
	jsonData := []byte(`{
		"name": "Alice",
		"age": 30,
		"address": {
			"city": "Wonderland",
			"zip_code": 12345
		},
		"tags": ["admin", "user"]
	}`)

	var user User
	err := ValidateJSON(jsonData, &user)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateJSON_MissingRequiredFields(t *testing.T) {
	jsonData := []byte(`{
		"age": 30,
		"address": {
			"zip_code": 12345
		}
	}`)

	var user User
	err := ValidateJSON(jsonData, &user)
	if err == nil {
		t.Fatal("expected validation errors, got nil")
	}

	ve, ok := err.(ValidationErrors)
	if !ok {
		t.Fatalf("expected ValidationErrors, got %T", err)
	}

	if len(ve) < 2 {
		t.Errorf("expected at least 2 validation errors, got %d", len(ve))
	}
}

func TestValidateJSON_InvalidJSON(t *testing.T) {
	jsonData := []byte(`{ invalid json }`)

	var user User
	err := ValidateJSON(jsonData, &user)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}

	ve, ok := err.(ValidationErrors)
	if !ok || len(ve) != 1 || ve[0].Tag != "json_format" {
		t.Errorf("expected json_format error, got: %v", err)
	}
}

func TestValidateJSON_TypeMismatch(t *testing.T) {
	jsonData := []byte(`{
		"name": "Alice",
		"age": "thirty",
		"address": {
			"city": "Wonderland",
			"zip_code": 12345
		}
	}`)

	var user User
	err := ValidateJSON(jsonData, &user)
	if err == nil {
		t.Fatal("expected type mismatch error, got nil")
	}

	ve, ok := err.(ValidationErrors)
	if !ok || len(ve) == 0 {
		t.Errorf("expected ValidationErrors, got: %v", err)
	}
}
