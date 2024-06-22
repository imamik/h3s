package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"hcloud-k3s-cli/pkg/utils/file"
)

func Save(config any, filePath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshalling config data: %w", err)
	}

	err = file.Save(data, filePath)
	if err != nil {
		return err
	}

	return nil
}
