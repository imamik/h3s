package commands

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/config"
	"hcloud-k3s-cli/internal/utils/yaml"
	"strings"
)

func ControlPlane(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
	controlPlaneNodes []*hcloud.Server,
	node *hcloud.Server,
) string {
	configYaml := yaml.String(config.K3sServerConfig{})
	commandArr := []string{
		PreInstallCommand(configYaml),
		K3sCommand(ctx),
		SeLinux(),
		PostInstall(),
	}
	return strings.Join(commandArr, "\n")
}
