package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates any struct using go-playground/validator
func ValidateStruct(input any) error {
	if err := validate.Struct(input); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return fmt.Errorf("validation error: %s", validationErrors.Error())
		}
		return err
	}
	return nil
}