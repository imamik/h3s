package credentials

import (
	"errors"
	"regexp"
	"strconv"
)

// TokenRegex defines the pattern for validating alphanumeric tokens
var TokenRegex = regexp.MustCompile(`^[a-zA-Z0-9]*$`)

// ValidateAlphanumericToken validates a token based on the specified length and checks if it is alphanumeric.
func ValidateAlphanumericToken(s string, length int) error {
	if len(s) != length {
		return errors.New("token must be " + strconv.Itoa(length) + " characters long")
	}
	if !TokenRegex.MatchString(s) {
		return errors.New("name must be alphanumeric characters only")
	}
	return nil
}

// ValidateHCloudToken validates the Hetzner cloud token, which must be 64 characters long and alphanumeric.
func ValidateHCloudToken(s string) error {
	return ValidateAlphanumericToken(s, 64)
}

// ValidateHetznerDNSToken validates the Hetzner DNS token, which must be 32 characters long and alphanumeric.
func ValidateHetznerDNSToken(s string) error {
	return ValidateAlphanumericToken(s, 32)
}
