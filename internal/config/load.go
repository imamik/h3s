// Package config provides functionality for loading and managing configuration
package config

import (
	"h3s/internal/utils/file"
	"h3s/internal/validation"
	"os"
)

// Load reads the configuration from the h3s.yaml file (or $H3S_CONFIG) and returns a Config struct
// If there's an error while loading the configuration, it logs a fatal error and returns an empty Config
func Load() (*Config, error) {
	var conf Config
	configPath := os.Getenv("H3S_CONFIG")
	if configPath == "" {
		configPath = "h3s.yaml"
	}
	err := file.New(configPath).Load().UnmarshalYamlTo(&conf)
	if err != nil {
		return nil, err
	}
	// Validate config struct
	if err := validation.ValidateStruct(&conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
