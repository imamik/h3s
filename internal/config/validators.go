package config

import (
	"h3s/internal/utils/validation"
)

// ValidateName validates a name in lower-kebap-case, checks length and format
// This is a wrapper around the validation.Name function for backward compatibility
func ValidateName(s string) error {
	return validation.Name(s)
}

// IsNumberString checks if a string is a number
// This is a wrapper around the validation.Number function for backward compatibility
func IsNumberString(s string) error {
	return validation.Number(s)
}

// IsUnevenNumberString checks if a string is a number and is uneven
// This is a wrapper around the validation.UnevenNumber function for backward compatibility
func IsUnevenNumberString(s string) error {
	return validation.UnevenNumber(s)
}
