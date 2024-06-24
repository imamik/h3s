package credentials

import (
	"fmt"
	"hcloud-k3s-cli/internal/utils/yaml"
)

func save(projectName string, projectCredentials ProjectCredentials) {
	credentials, err := initialize()
	if err != nil {
		fmt.Println(err)
		return
	}

	credentials[projectName] = projectCredentials
	err = yaml.Save(credentials, path)
	if err != nil {
		fmt.Println(err)
		return
	}
}
