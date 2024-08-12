package commands

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"gopkg.in/yaml.v3"
	"h3s/internal/cluster"
	"h3s/internal/k3s/install/config"
	"h3s/internal/utils/ssh"
	"strings"
)

func ControlPlane(
	ctx *cluster.Cluster,
	lb *hcloud.LoadBalancer,
	controlPlaneNodes []*hcloud.Server,
	proxy *hcloud.Server,
	node *hcloud.Server,
) error {
	isFirst := node.ID == controlPlaneNodes[0].ID
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
		ClusterDomain:    ctx.Config.Domain + ".local",
		HTTPSListenPort:  6443,

		// Etcd
		TLSSAN: getTlsSan(ctx, lb, controlPlaneNodes),

		// Security
		SELinux: true,
	}

	if isFirst {
		configYaml.ClusterInit = true
	} else {
		configYaml.Server = getServer(controlPlaneNodes[0])
		configYaml.WriteKubeconfigMode = "0644"
	}

	if !ctx.Config.ControlPlane.AsWorkerPool {
		configYaml.NodeTaint = []string{"node-role.kubernetes.io/control-plane:NoSchedule"}
	}

	configYamlStr, err := yaml.Marshal(configYaml)

	if err != nil {
		return err
	}

	commandArr := []string{
		PreInstallCommand(ctx, string(configYamlStr)),
		K3sInstall(ctx, true),
		SeLinux(),
		PostInstall(),
		K3sStartServer(isFirst),
	}

	_, err = ssh.ExecuteViaProxy(ctx, proxy, node, strings.Join(commandArr, "\n"))
	if err != nil {
		return err
	}
	return nil
}
