// Package config provides functionality for loading and managing configuration
package config

import (
	"h3s/internal/utils/file"
)

// Load reads the configuration from the h3s.yaml file and returns a Config struct
// If there's an error while loading the configuration, it logs a fatal error and returns an empty Config
func Load() (*Config, error) {
	var conf Config
	err := file.New("h3s.yaml").Load().UnmarshalYamlTo(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
