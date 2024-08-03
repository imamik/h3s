package credentials

import (
	"h3s/internal/config/path"
	"h3s/internal/utils/file"
	"h3s/internal/utils/yaml"
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
