package config

import (
	"hcloud-k3s-cli/pkg/utils/yaml"
	"log"
)

func Load() Config {
	var conf Config
	err := yaml.Load("hcloud-k3s.yaml", &conf)
	if err != nil {
		log.Fatalf("error loading config: %s", err)
		return conf
	}
	return conf
}
