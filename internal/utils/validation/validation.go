// Package validation provides utilities for validating input data
package validation

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// Common validation errors
var (
	ErrEmptyString      = errors.New("value cannot be empty")
	ErrInvalidFormat    = errors.New("invalid format")
	ErrInvalidLength    = errors.New("invalid length")
	ErrInvalidValue     = errors.New("invalid value")
	ErrInvalidRange     = errors.New("value out of valid range")
	ErrInvalidCharacter = errors.New("contains invalid characters")
)

// Common regular expressions
var (
	// NameRegex is a regex for validating names in lower-kebap-case
	NameRegex = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

	// EmailRegex is a regex for validating email addresses
	EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	// DomainRegex is a regex for validating domain names
	DomainRegex = regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

	// IPRegex is a regex for validating IP addresses
	IPRegex = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
)

// StringNotEmpty validates that a string is not empty
func StringNotEmpty(s string) error {
	if s == "" {
		return ErrEmptyString
	}
	return nil
}

// StringLength validates that a string has a length within the specified range
func StringLength(s string, minLen, maxLen int) error {
	if len(s) < minLen {
		return fmt.Errorf("%w: must be at least %d characters long", ErrInvalidLength, minLen)
	}
	if maxLen > 0 && len(s) > maxLen {
		return fmt.Errorf("%w: must not be longer than %d characters", ErrInvalidLength, maxLen)
	}
	return nil
}

// StringMatches validates that a string matches a regular expression
func StringMatches(s string, regex *regexp.Regexp, description string) error {
	if !regex.MatchString(s) {
		return fmt.Errorf("%w: must be %s", ErrInvalidFormat, description)
	}
	return nil
}

// Name validates a name in lower-kebap-case, checks length and format
func Name(s string) error {
	if err := StringLength(s, 4, 64); err != nil {
		return err
	}
	if !NameRegex.MatchString(s) {
		return fmt.Errorf("%w: must be in lower-kebap-case", ErrInvalidFormat)
	}
	return nil
}

// Email validates an email address
func Email(s string) error {
	if err := StringNotEmpty(s); err != nil {
		return err
	}
	if !EmailRegex.MatchString(s) {
		return fmt.Errorf("%w: invalid email address", ErrInvalidFormat)
	}
	return nil
}

// Domain validates a domain name
func Domain(s string) error {
	if err := StringNotEmpty(s); err != nil {
		return err
	}
	if !DomainRegex.MatchString(s) {
		return fmt.Errorf("%w: invalid domain name", ErrInvalidFormat)
	}
	return nil
}

// IP validates an IP address
func IP(s string) error {
	if err := StringNotEmpty(s); err != nil {
		return err
	}
	if net.ParseIP(s) == nil {
		return fmt.Errorf("%w: invalid IP address", ErrInvalidFormat)
	}
	return nil
}

// Number validates that a string is a valid number
func Number(s string) error {
	if _, err := strconv.Atoi(s); err != nil {
		return fmt.Errorf("%w: must be a number", ErrInvalidFormat)
	}
	return nil
}

// NumberInRange validates that a string is a valid number within the specified range
func NumberInRange(s string, minVal, maxVal int) error {
	i, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("%w: must be a number", ErrInvalidFormat)
	}
	if i < minVal || (maxVal > 0 && i > maxVal) {
		return fmt.Errorf("%w: must be between %d and %d", ErrInvalidRange, minVal, maxVal)
	}
	return nil
}

// UnevenNumber validates that a string is a valid uneven number
func UnevenNumber(s string) error {
	i, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("%w: must be a number", ErrInvalidFormat)
	}
	if i%2 == 0 {
		return fmt.Errorf("%w: must be an uneven number", ErrInvalidValue)
	}
	return nil
}

// FilePath validates a file path
func FilePath(s string) error {
	if err := StringNotEmpty(s); err != nil {
		return err
	}
	// Basic validation - could be enhanced with more specific checks
	if strings.Contains(s, "..") {
		return fmt.Errorf("%w: path cannot contain '..'", ErrInvalidFormat)
	}
	return nil
}

// URL validates a URL
func URL(s string) error {
	if err := StringNotEmpty(s); err != nil {
		return err
	}
	// Basic validation - could be enhanced with more specific checks
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		return fmt.Errorf("%w: URL must start with http:// or https://", ErrInvalidFormat)
	}
	return nil
}
