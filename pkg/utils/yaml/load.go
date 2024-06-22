package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"hcloud-k3s-cli/pkg/utils/file"
	"os"
)

func Load(filePath string, out interface{}) error {
	absPath, err := file.Normalize(filePath)
	if err != nil {
		return fmt.Errorf("error normalizing file path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", absPath, err)
	}

	err = yaml.Unmarshal(data, out)
	if err != nil {
		return fmt.Errorf("error unmarshalling config data: %w", err)
	}

	return nil
}
