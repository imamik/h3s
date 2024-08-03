package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"h3s/internal/utils/file"
)

func Load(filePath string, out interface{}) error {
	data, err := file.Load(filePath)
	if err != nil {
		return fmt.Errorf("error loading file: %w", err)
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("error unmarshalling config data: %w", err)
	}

	return nil
}
