package credentials

import (
	"h3s/internal/config/path"
	"h3s/internal/utils/file"
)

// SaveCredentials saves the project credentials to the secrets file
func SaveCredentials(projectCredentials ProjectCredentials) error {
	p := string(path.SecretsFileName)
	_, err := file.New(p).SetYaml(projectCredentials).Save()
	return err
}
