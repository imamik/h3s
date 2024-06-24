package credentials

import (
	"hcloud-k3s-cli/internal/utils/file"
	"hcloud-k3s-cli/internal/utils/yaml"
)

func initialize() (Credentials, error) {
	err := file.Create(path)
	if err != nil {
		return nil, err
	}

	var credentials Credentials
	err = yaml.Load(path, &credentials)
	if err != nil || credentials == nil {
		return make(Credentials), err
	}

	return credentials, nil
}
