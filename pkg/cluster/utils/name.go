package utils

import (
	"fmt"
	"hcloud-k3s-cli/pkg/config"
)

func GetName(name string, conf config.Config) string {
	return fmt.Sprintf("%s-k3s-%s", conf.Name, name)
}
