package credentials

import (
	"fmt"
	"hcloud-k3s-cli/internal/config/path"
	"hcloud-k3s-cli/internal/utils/yaml"
)

func SaveCredentials(projectName string, projectCredentials ProjectCredentials) {
	p := path.GetPath(projectName, path.CredentialFileName)
	err := yaml.Save(projectCredentials, p)
	if err != nil {
		fmt.Println(err)
		return
	}
}
