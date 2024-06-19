package config

import (
	"errors"
	"regexp"
	"strconv"
)

var NameRegex = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

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

func IsNumberString(s string) error {
	_, err := strconv.Atoi(s)
	if err != nil {
		return errors.New("must be a number")
	}
	return nil
}

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
