package credentials

import (
	"errors"
	"regexp"
)

var HCloudTokenRegex = regexp.MustCompile(`^[a-zA-Z0-9]*$`)

func ValidateHCloudToken(s string) error {
	if len(s) != 64 {
		return errors.New("token must be 64 characters long")
	}
	if !HCloudTokenRegex.MatchString(s) {
		return errors.New("name must be in lower-kebap-case")
	}
	return nil
}
