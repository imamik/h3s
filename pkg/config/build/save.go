package build

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"hcloud-k3s-cli/pkg/config"
	"os"
)

func save(config config.Config, filename string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshalling config data: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating %s: %w", filename, err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to %s: %w", filename, err)
	}

	return nil
}
