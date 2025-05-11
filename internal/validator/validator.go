package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string
	Message string
	Tag     string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a slice of ValidationError
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, err := range e {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Error())
	}
	return sb.String()
}

// ValidateJSON validates a JSON request body against a struct type
func ValidateJSON(data []byte, s any) error {
	// First, try to unmarshal into a map to check for type mismatches
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return ValidationErrors{
			{
				Field:   "body",
				Message: "Invalid JSON format",
				Tag:     "json_format",
			},
		}
	}

	// Now try to unmarshal into the actual struct
	if err := json.Unmarshal(data, s); err != nil {
		// Handle type mismatch errors
		if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
			return ValidationErrors{
				{
					Field:   typeErr.Field,
					Message: fmt.Sprintf("expected %s, got %s", typeErr.Type, typeErr.Value),
					Tag:     "type_mismatch",
				},
			}
		}
		return ValidationErrors{
			{
				Field:   "body",
				Message: "Invalid request structure",
				Tag:     "invalid_structure",
			},
		}
	}

	// If JSON parsing succeeded, validate the struct
	return ValidateStruct(s)
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(s any) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("validator: expected struct, got %s", v.Kind())
	}

	var errors ValidationErrors
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Check required fields
		if tag := field.Tag.Get("validate"); tag != "" {
			if strings.Contains(tag, "required") && isEmpty(value) {
				errors = append(errors, ValidationError{
					Field:   field.Name,
					Message: "this field is required",
					Tag:     "required",
				})
			}
		}

		// Validate nested structs
		if value.Kind() == reflect.Struct {
			if err := ValidateStruct(value.Interface()); err != nil {
				if ve, ok := err.(ValidationErrors); ok {
					for _, e := range ve {
						errors = append(errors, ValidationError{
							Field:   fmt.Sprintf("%s.%s", field.Name, e.Field),
							Message: e.Message,
							Tag:     e.Tag,
						})
					}
				}
			}
		}

		// Validate slices
		if value.Kind() == reflect.Slice {
			for j := 0; j < value.Len(); j++ {
				item := value.Index(j)
				if item.Kind() == reflect.Struct {
					if err := ValidateStruct(item.Interface()); err != nil {
						if ve, ok := err.(ValidationErrors); ok {
							for _, e := range ve {
								errors = append(errors, ValidationError{
									Field:   fmt.Sprintf("%s[%d].%s", field.Name, j, e.Field),
									Message: e.Message,
									Tag:     e.Tag,
								})
							}
						}
					}
				}
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// isEmpty checks if a value is empty
func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}
	return false
}
