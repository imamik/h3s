package credentials

import (
	"hcloud-k3s-cli/internal/config/path"
	"hcloud-k3s-cli/internal/utils/file"
	"hcloud-k3s-cli/internal/utils/yaml"
)

func initialize(projectName string) (*ProjectCredentials, error) {
	p := path.GetPath(projectName, path.CredentialFileName)
	err := file.Create(p)
	if err != nil {
		return nil, err
	}

	var credentials *ProjectCredentials
	err = yaml.Load(p, &credentials)
	if err != nil || credentials == nil {
		return nil, err
	}

	return credentials, nil
}
