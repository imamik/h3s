package credentials

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/charmbracelet/huh"
)

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func surveyCredentials() (ProjectCredentials, error) {
	var projectCredentials ProjectCredentials

	if inputErr := huh.NewInput().
		Title("Hetzner Cloud Token").
		Description("The api token to create resources for the given project").
		Validate(ValidateHCloudToken).
		Value(&projectCredentials.HCloudToken).
		Run(); inputErr != nil {
		return projectCredentials, inputErr
	}

	if inputErr := huh.NewInput().
		Title("Hetzner DNS Token").
		Description("The dns token to create dns entries for the given project").
		Validate(ValidateHetznerDNSToken).
		Value(&projectCredentials.HetznerDNSToken).
		Run(); inputErr != nil {
		return projectCredentials, inputErr
	}

	k3sToken, err := generateToken(32)
	if err != nil {
		return projectCredentials, err
	}
	projectCredentials.K3sToken = k3sToken

	return projectCredentials, nil
}
