package config

import (
	"h3s/internal/utils/yaml"
	"log"
)

func Load() Config {
	var conf Config
	err := yaml.Load("h3s.yaml", &conf)
	if err != nil {
		log.Fatalf("error loading config: %s", err)
		return conf
	}
	return conf
}
