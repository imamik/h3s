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
	controlPlaneNodes []*hcloud.Server,
) string {
	address := ""
	if lb != nil {
		address = lb.PublicNet.IPv4.IP.String()
	} else {
		node := controlPlaneNodes[0]
		address = ip.FirstAvailable(node)
	}
	return "https://" + address + ":6443"
}

func ControlPlane(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
	controlPlaneNodes []*hcloud.Server,
	node *hcloud.Server,
) string {
	nodeIp := ip.FirstAvailable(node)

	configYaml := yaml.String(config.K3sServerConfig{
		//Server:                   getServer(lb, controlPlaneNodes),
		NodeName:                 node.Name,
		Token:                    ctx.Credentials.K3sToken,
		ClusterInit:              true,
		DisableCloudController:   true,
		DisableKubeProxy:         false,
		DisableComponents:        []string{},
		KubeletArg:               []string{},
		KubeControllerManagerArg: []string{},
		FlannelIface:             "eth1",
		NodeIP:                   []string{nodeIp},
		AdvertiseAddress:         nodeIp,
		NodeTaint:                []string{},
		NodeLabel:                []string{},
		ClusterCIDR:              "10.42.0.0/16",
		ServiceCIDR:              "10.43.0.0/16",
		ClusterDNS:               "10.43.0.10",
		SELinux:                  true,
		WriteKubeconfigMode:      "0644",
	})
	commandArr := []string{
		PreInstallCommand(configYaml),
		K3sCommand(ctx),
		SeLinux(),
		PostInstall(),
	}
	return strings.Join(commandArr, "\n")
}
