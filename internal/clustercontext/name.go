package clustercontext

import (
	"h3s/internal/config"
	"strings"
)

func buildGetNameFunc(conf config.Config) func(...string) string {
	return func(names ...string) string {
		names = append([]string{conf.Name, "k3s"}, names...)
		return strings.Join(names, "-")
	}
}
