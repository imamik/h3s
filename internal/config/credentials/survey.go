package credentials

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/charmbracelet/huh"
	"h3s/internal/config"
)

func surveyName() string {
	var name string

	huh.NewInput().
		Title("Project Name").
		Description("Used to name resources. Must be in lower-kebap-case").
		Validate(config.ValidateName).
		Value(&name).
		Run()

	return name
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func surveyCredentials() ProjectCredentials {
	var projectCredentials ProjectCredentials

	huh.NewInput().
		Title("Hetzner Cloud Token").
		Description("The api token to create resources for the given project").
		Validate(ValidateHCloudToken).
		Value(&projectCredentials.HCloudToken).
		Run()

	huh.NewInput().
		Title("Hetzner DNS Token").
		Description("The dns token to create dns entries for the given project").
		Validate(ValidateHetznerDNSToken).
		Value(&projectCredentials.HetznerDNSToken).
		Run()

	projectCredentials.K3sToken = generateToken(32)

	return projectCredentials

}
