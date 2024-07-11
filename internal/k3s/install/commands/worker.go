package commands

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/config"
	"hcloud-k3s-cli/internal/utils/ssh"
	"hcloud-k3s-cli/internal/utils/yaml"
	"strings"
)

func Worker(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
	controlPlaneNodes []*hcloud.Server,
	proxy *hcloud.Server,
	node *hcloud.Server,
) {
	nodeIp := node.PrivateNet[0].IP.String()
	server := getServer(lb, controlPlaneNodes[0])

	configYaml := config.K3sServerConfig{
		// Node
		NodeName:  node.Name,
		NodeLabel: []string{},

		// Server
		Token:  ctx.Credentials.K3sToken,
		Server: server,

		// Cluster
		NodeTaint: []string{},

		// Kube
		KubeletArg: []string{
			"cloud-provider=external",
			"volume-plugin-dir=/var/lib/kubelet/volumeplugins",
		},

		// Network
		FlannelIface: "eth1",
		NodeIP:       []string{nodeIp},

		// Security
		SELinux: true,
	}

	configYamlStr := yaml.String(configYaml)
	commandArr := []string{
		PreInstallCommand(configYamlStr),
		K3sInstall(ctx, false),
		SeLinux(),
		PostInstall(),
		K3sStartAgent(),
	}
	ssh.ExecuteViaProxy(ctx, proxy, node, strings.Join(commandArr, "\n"))
}
