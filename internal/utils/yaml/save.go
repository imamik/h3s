package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"h3s/internal/utils/file"
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
