package commands

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"gopkg.in/yaml.v3"
	"h3s/internal/cluster"
	"h3s/internal/k3s/config"
	"h3s/internal/utils/ssh"
	"strings"
)

func Worker(
	ctx *cluster.Cluster,
	firstControlPlane *hcloud.Server,
	proxy *hcloud.Server,
	node *hcloud.Server,
) error {
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

	configYamlStr, err := yaml.Marshal(configYaml)

	if err != nil {
		return err
	}

	commandArr := []string{
		PreInstallCommand(ctx, string(configYamlStr)),
		K3sInstall(ctx, false),
		SeLinux(),
		PostInstall(),
		K3sStartAgent(),
	}
	_, err = ssh.ExecuteViaProxy(ctx, proxy, node, strings.Join(commandArr, "\n"))
	if err != nil {
		return err
	}

	return nil
}
