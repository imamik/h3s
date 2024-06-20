package network

import (
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
)

func getName(conf config.Config) string {
	return utils.GetName("network", conf)
}
