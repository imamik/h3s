package validation

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct using go-playground/validator tags
// Returns nil if valid, or an error describing the first problem
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// Add custom validation rules for config structs
// Example: Register custom validation for kebap-case names, file paths, etc.
// See documentation in internal/config/structs.go
