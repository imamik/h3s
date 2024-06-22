package clustercontext

import "hcloud-k3s-cli/pkg/config"

func buildGetLabelsFunc(conf config.Config) func(...map[string]string) map[string]string {
	return func(optionalLabels ...map[string]string) map[string]string {
		labels := map[string]string{
			"project": conf.Name,
		}
		if len(optionalLabels) > 0 {
			for key, value := range optionalLabels[0] {
				labels[key] = value
			}
		}
		return labels
	}
}
