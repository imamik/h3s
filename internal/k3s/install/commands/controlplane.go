package commands

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/config"
	"hcloud-k3s-cli/internal/utils/ssh"
	"hcloud-k3s-cli/internal/utils/yaml"
	"strings"
)

func ControlPlane(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
	firstControlPlane *hcloud.Server,
	proxy *hcloud.Server,
	node *hcloud.Server,
) {
	isFirst := node.ID == firstControlPlane.ID
	nodeIp := node.PrivateNet[0].IP.String()
	networkInterface, _ := GetNetworkInterfaceName(ctx, proxy, node)

	configYaml := config.K3sServerConfig{
		// Node
		NodeName:  node.Name,
		NodeLabel: []string{},

		// Server
		Token: ctx.Credentials.K3sToken,

		// Cluster
		NodeTaint: []string{},

		// Disabled Functionality
		DisableCloudController: true,
		DisableKubeProxy:       false,
		DisableComponents: []string{
			"local-storage",
			"servicelb",
			"traefik",
		},

		// Kube
		KubeletArg: []string{
			"cloud-provider=external",
			"volume-plugin-dir=/var/lib/kubelet/volumeplugins",
		},
		KubeControllerManagerArg: []string{
			"flex-volume-plugin-dir=/var/lib/kubelet/volumeplugins",
		},

		// Flannel
		FlannelBackend:       "vxlan", // TODO: make this configurable
		FlannelIface:         networkInterface,
		DisableNetworkPolicy: false,

		// Network
		NodeIP:           []string{nodeIp},
		AdvertiseAddress: nodeIp,
		ClusterCIDR:      "10.42.0.0/16",
		ServiceCIDR:      "10.43.0.0/16",
		ClusterDNS:       "10.43.0.10",

		// Etcd
		TLSSAN: getTlsSan(ctx, lb),

		// Security
		SELinux: true,
	}

	if isFirst {
		configYaml.ClusterInit = true
	} else {
		configYaml.Server = getServer(firstControlPlane)
		configYaml.WriteKubeconfigMode = "0644"
	}

	if !ctx.Config.ControlPlane.AsWorkerPool {
		configYaml.NodeTaint = []string{"node-role.kubernetes.io/control-plane:NoSchedule"}
	}

	configYamlStr := yaml.String(configYaml)
	commandArr := []string{
		PreInstallCommand(ctx, configYamlStr),
		K3sInstall(ctx, true),
		SeLinux(),
		PostInstall(),
		K3sStartServer(isFirst),
	}
	ssh.ExecuteViaProxy(ctx, proxy, node, strings.Join(commandArr, "\n"))
}
