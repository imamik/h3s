package clustercontext

import "h3s/internal/config"

func buildGetLabelsFunc(conf config.Config) func(...map[string]string) map[string]string {
	return func(optionalLabels ...map[string]string) map[string]string {
		labels := map[string]string{
			"project": conf.Name,
			"origin":  "h3s",
		}
		if len(optionalLabels) > 0 {
			for key, value := range optionalLabels[0] {
				labels[key] = value
			}
		}
		return labels
	}
}
