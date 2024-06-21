package utils

import "hcloud-k3s-cli/pkg/config"

func GetLabels(conf config.Config, optionalLabels ...map[string]string) map[string]string {
	labels := map[string]string{
		"project": conf.Name,
	}

	// If optionalLabels is provided, merge it with labels
	if len(optionalLabels) > 0 {
		for key, value := range optionalLabels[0] {
			labels[key] = value
		}
	}

	return labels
}
