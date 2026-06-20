package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	
	// Use JSON field names in error messages instead of struct field names
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// ValidateStruct validates a struct and returns a user-friendly error message
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	// Type assert to validator.ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fmt.Errorf("validation failed: %v", err)
	}

	// Build user-friendly error messages
	var errMsgs []string
	for _, fieldErr := range validationErrors {
		errMsgs = append(errMsgs, formatFieldError(fieldErr))
	}

	return fmt.Errorf("%s", strings.Join(errMsgs, "; "))
}

// formatFieldError formats a single field error into a readable message
func formatFieldError(fieldErr validator.FieldError) string {
	field := fieldErr.Field()
	tag := fieldErr.Tag()

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "required_if":
		return fmt.Sprintf("%s is required when %s", field, fieldErr.Param())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fieldErr.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters/value", field, fieldErr.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}
