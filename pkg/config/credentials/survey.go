package credentials

import (
	"github.com/charmbracelet/huh"
	"hcloud-k3s-cli/pkg/config"
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

func surveyCredentials() ProjectCredentials {
	var projectCredentials ProjectCredentials

	huh.NewInput().
		Title("Hetzner Cloud Token").
		Description("The api token to create resources for the given project").
		Validate(ValidateHCloudToken).
		Value(&projectCredentials.HCloudToken).
		Run()

	return projectCredentials

}
