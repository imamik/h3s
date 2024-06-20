package utils

import "hcloud-k3s-cli/pkg/config"

func GetLabels(conf config.Config) map[string]string {
	return map[string]string{
		"project": conf.Name,
	}
}
