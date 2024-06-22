package credentials

import (
	"fmt"
	"hcloud-k3s-cli/pkg/utils/file"
	"hcloud-k3s-cli/pkg/utils/yaml"
)

func Config() {
	path := "$HOME/.config/hcloud-k3s/credentials.yaml"

	err := file.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	credentials := make(Credentials)
	err = yaml.Load(path, &credentials)
	if err != nil {
		fmt.Println(err)
		return
	}

	projectName := surveyName()
	projectCredentials := surveyCredentials()
	credentials[projectName] = projectCredentials
	credentials["default"] = projectCredentials

	err = yaml.Save(credentials, path)
	if err != nil {
		fmt.Println(err)
		return
	}
}
