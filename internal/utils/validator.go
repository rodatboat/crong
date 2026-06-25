package utils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rodatboat/crong/internal/entities"
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

	// Register custom validator for HTTP methods (0-8: GET through CONNECT)
	validate.RegisterValidation("validmethod", func(fl validator.FieldLevel) bool {
		method := fl.Field().Interface().(entities.ReqMethod)
		return method >= 0 && method <= 8
	})

	// Register custom validator for cron expressions
	validate.RegisterValidation("validcron", func(fl validator.FieldLevel) bool {
		cronExpr := fl.Field().String()
		return isValidCronExpression(cronExpr)
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
	case "validmethod":
		return fmt.Sprintf("%s must be a valid HTTP method", field)
	case "validcron":
		return fmt.Sprintf("%s must be a valid cron expression (format: minute hour mday month wday)", field)
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

// isValidCronExpression validates a cron expression format
// Format: "minute hour mday month wday" (5 fields, space-separated)
// Supported operators: * (any), - (range), , (list), / (step on * only)
// Returns true if the expression is syntactically valid
func isValidCronExpression(cronExpr string) bool {
	parts := strings.Fields(strings.TrimSpace(cronExpr))
	if len(parts) != 5 {
		return false
	}

	// Validate each field: min, max values
	fieldRanges := [][2]int{
		{0, 59}, // minute
		{0, 23}, // hour
		{1, 31}, // mday
		{1, 12}, // month
		{0, 6},  // wday
	}

	for i, part := range parts {
		if !isValidCronField(part, fieldRanges[i][0], fieldRanges[i][1]) {
			return false
		}
	}

	return true
}

// isValidCronField validates a single cron field
// Supported formats:
//   - "*" (any value)
//   - "5" (single value)
//   - "1-5" (range)
//   - "1,5,10" (list)
//   - "1-2,5,10" (combined ranges and values)
//   - "*/N" (step values only, where N is positive and <= max-min)
//
// NOT supported: steps on ranges "1-10/2" or steps on lists "1,5/2"
func isValidCronField(field string, min, max int) bool {
	if field == "*" {
		return true
	}

	// Check if field uses step syntax - only "*/N" is allowed
	if strings.Contains(field, "/") {
		// Step syntax is only allowed as "*/N"
		if !strings.HasPrefix(field, "*/") {
			return false // Step on anything other than * is invalid
		}

		step, err := strconv.Atoi(field[2:])
		if err != nil || step <= 0 || step > (max-min) {
			return false
		}
		return true
	}

	// No step syntax - validate as list of ranges/values
	expressions := strings.Split(field, ",")
	for _, expr := range expressions {
		expr = strings.TrimSpace(expr)

		if expr == "" {
			return false
		}

		if idx := strings.Index(expr, "-"); idx != -1 {
			// Range like "1-5"
			rangeMin, err1 := strconv.Atoi(strings.TrimSpace(expr[:idx]))
			rangeMax, err2 := strconv.Atoi(strings.TrimSpace(expr[idx+1:]))
			if err1 != nil || err2 != nil || rangeMin < min || rangeMax > max || rangeMin > rangeMax {
				return false
			}
		} else {
			// Single value
			val, err := strconv.Atoi(strings.TrimSpace(expr))
			if err != nil || val < min || val > max {
				return false
			}
		}
	}

	return true
}
