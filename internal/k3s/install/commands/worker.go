package commands

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/k3s/install/config"
	"h3s/internal/utils/ssh"
	"h3s/internal/utils/yaml"
	"strings"
)

func Worker(
	ctx clustercontext.ClusterContext,
	firstControlPlane *hcloud.Server,
	proxy *hcloud.Server,
	node *hcloud.Server,
) {
	nodeIp := node.PrivateNet[0].IP.String()
	server := getServer(firstControlPlane)
	networkInterface, _ := GetNetworkInterfaceName(ctx, proxy, node)

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
		FlannelIface: networkInterface,
		NodeIP:       []string{nodeIp},

		// Security
		SELinux: true,
	}

	configYamlStr := yaml.String(configYaml)
	commandArr := []string{
		PreInstallCommand(ctx, configYamlStr),
		K3sInstall(ctx, false),
		SeLinux(),
		PostInstall(),
		K3sStartAgent(),
	}
	ssh.ExecuteViaProxy(ctx, proxy, node, strings.Join(commandArr, "\n"))
}
