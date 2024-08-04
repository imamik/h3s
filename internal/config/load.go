// Package config provides functionality for loading and managing configuration
package config

import (
	"h3s/internal/utils/yaml"
	"log"
)

// Load reads the configuration from the h3s.yaml file and returns a Config struct
// If there's an error while loading the configuration, it logs a fatal error and returns an empty Config
func Load() Config {
	var conf Config
	err := yaml.Load("h3s.yaml", &conf)
	if err != nil {
		log.Fatalf("error loading config: %s", err)
		return conf
	}
	return conf
}
