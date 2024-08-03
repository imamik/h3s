package credentials

import (
	"fmt"
	"h3s/internal/config/path"
	"h3s/internal/utils/yaml"
)

func SaveCredentials(projectName string, projectCredentials ProjectCredentials) {
	p := path.GetPath(projectName, path.CredentialFileName)
	err := yaml.Save(projectCredentials, p)
	if err != nil {
		fmt.Println(err)
		return
	}
}
