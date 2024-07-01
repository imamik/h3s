package commands

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/k3s/install/config"
	"hcloud-k3s-cli/internal/utils/ip"
	"hcloud-k3s-cli/internal/utils/yaml"
	"strings"
)

func getServer(
	lb *hcloud.LoadBalancer,
	node *hcloud.Server,
) string {
	address := ""
	if lb != nil {
		address = lb.PublicNet.IPv4.IP.String()
	} else {
		address = ip.FirstAvailable(node)
	}
	return "https://" + address + ":6443"
}

func getTlsSan(
	lb *hcloud.LoadBalancer,
	controlPlaneNodes []*hcloud.Server,
) []string {
	var tlsSan []string
	if lb != nil {
		tlsSan = append(tlsSan, lb.PublicNet.IPv4.IP.String())
		tlsSan = append(tlsSan, lb.PublicNet.IPv6.IP.String())
	} else {
		for _, node := range controlPlaneNodes {
			tlsSan = append(tlsSan, ip.FirstAvailable(node))
		}
	}
	return tlsSan
}

func ControlPlane(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
	controlPlaneNodes []*hcloud.Server,
	node *hcloud.Server,
) string {
	nodeIp := node.PrivateNet[0].IP.String()

	configYaml := yaml.String(config.K3sServerConfig{
		// Node
		NodeName:  node.Name,
		NodeLabel: []string{},

		// Server
		Server: getServer(lb, node),
		Token:  ctx.Credentials.K3sToken,

		// Cluster
		ClusterInit: node.ID == controlPlaneNodes[0].ID,
		NodeTaint:   []string{},

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

		// Network
		FlannelIface:     "eth1",
		NodeIP:           []string{nodeIp},
		AdvertiseAddress: nodeIp,
		ClusterCIDR:      "10.42.0.0/16",
		ServiceCIDR:      "10.43.0.0/16",
		ClusterDNS:       "10.43.0.10",

		// CNI
		DisableNetworkPolicy: false,
		FlannelBackend:       "vxlan",

		// System
		SELinux:             true,
		WriteKubeconfigMode: "0644",

		TLSSAN: getTlsSan(lb, controlPlaneNodes),
	})
	commandArr := []string{
		PreInstallCommand(configYaml),
		K3sCommand(ctx),
		SeLinux(),
		PostInstall(),
	}
	return strings.Join(commandArr, "\n")
}
