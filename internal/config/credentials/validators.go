package credentials

import (
	"errors"
	"regexp"
)

var TokenRegex = regexp.MustCompile(`^[a-zA-Z0-9]*$`)

func ValidateHCloudToken(s string) error {
	if len(s) != 64 {
		return errors.New("token must be 64 characters long")
	}
	if !TokenRegex.MatchString(s) {
		return errors.New("name must alphanumeric characters only")
	}
	return nil
}

func ValidateHetznerDNSToken(s string) error {
	if len(s) != 32 {
		return errors.New("token must be 32 characters long")
	}
	if !TokenRegex.MatchString(s) {
		return errors.New("name must be alphanumeric characters only")
	}
	return nil
}
