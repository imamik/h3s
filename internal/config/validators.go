package config

import (
	"errors"
	"regexp"
	"strconv"
)

// NameRegex is a regex for validating names in lower-kebap-case
var NameRegex = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

// ValidateName validates a name in lower-kebap-case, checks length and format
func ValidateName(s string) error {
	if len(s) < 4 {
		return errors.New("name must be at least 4 characters long")
	}
	if len(s) > 64 {
		return errors.New("name must not be longer than 63 characters")
	}
	if !NameRegex.MatchString(s) {
		return errors.New("name must be in lower-kebap-case")
	}
	return nil
}

// IsNumberString checks if a string is a number
func IsNumberString(s string) error {
	_, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("must be a number")
	}
	return nil
}

// IsUnevenNumberString checks if a string is a number and is uneven
func IsUnevenNumberString(s string) error {
	i, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("must be a number")
	}
	if i%2 == 0 {
		return errors.New("must be an uneven number")
	}
	return nil
}
