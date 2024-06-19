package load

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"hcloud-k3s-cli/pkg/config"
	"os"
)

func Load(filename string) (config.Config, error) {
	var cfg config.Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return cfg, fmt.Errorf("error reading %s: %w", filename, err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("error unmarshalling cfg data: %w", err)
	}

	return cfg, nil
}
