package config

import (
	"gopkg.in/yaml.v3"
	"strings"
)

// CommandConfig represents the configuration options for the k3s agents and control plane commands (installation & server configuration)
type CommandConfig struct {
	IsMain         bool
	IsControlPlane bool

	K3sToken string
	Server   string
	Domain   string
	TlsSAN   []string

	ControlPlanesAsWorker bool

	NodeName         string
	NodeIp           string
	NetworkInterface string

	PublicIps  bool
	K3sVersion string
}

func getCommonConfig(c CommandConfig) CommonConfig {
	config := CommonConfig{
		Token:    c.K3sToken,
		NodeName: c.NodeName,
		KubeletArg: []string{
			"cloud-provider=external",
			"volume-plugin-dir=/var/lib/kubelet/volumeplugins",
		},
		NodeIP:       []string{c.NodeIp},
		FlannelIface: c.NetworkInterface,
		SELinux:      true,
	}
	if !c.IsMain {
		config.Server = c.Server
	}
	return config
}

func getControlPlaneConfig(c CommandConfig) ServerConfig {
	config := ServerConfig{
		CommonConfig:           getCommonConfig(c),
		DisableCloudController: true,
		DisableKubeProxy:       false,
		DisableComponents: []string{
			"local-storage",
			"servicelb",
			"traefik",
		},
		KubeControllerManagerArg: []string{
			"flex-volume-plugin-dir=/var/lib/kubelet/volumeplugins",
		},
		FlannelBackend:       "vxlan",
		DisableNetworkPolicy: false,
		AdvertiseAddress:     c.NodeIp,
		ClusterCIDR:          "10.42.0.0/16",
		ServiceCIDR:          "10.43.0.0/16",
		ClusterDNS:           "10.43.0.10",
		ClusterDomain:        c.Domain + ".local",
		HTTPSListenPort:      6443,
		TLSSAN:               c.TlsSAN,
	}

	if c.IsMain {
		config.ClusterInit = true
	} else {
		config.WriteKubeconfigMode = "0644"
	}

	if !c.ControlPlanesAsWorker {
		config.NodeTaint = []string{"node-role.kubernetes.io/control-plane:NoSchedule"}
	}

	return config
}

func Command(c CommandConfig) (string, error) {
	var config interface{}
	if c.IsControlPlane {
		config = getControlPlaneConfig(c)
	} else {
		config = AgentConfig{CommonConfig: getCommonConfig(c)}
	}
	configYaml, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	preInstallCmd, err := preInstall(c.PublicIps == false, string(configYaml))
	if err != nil {
		return "", err
	}

	k3sInstallCmd, err := K3sInstall(c.K3sVersion, c.IsControlPlane)
	if err != nil {
		return "", err
	}

	seLinuxCmd := SeLinux()
	postInstallCmd := postInstall()

	var k3sStartCmd string
	if c.IsControlPlane {
		k3sStartCmd, err = K3sStartServer(c.IsMain)
	} else {
		k3sStartCmd = K3sStartAgent()
	}

	cmdArray := []string{preInstallCmd, k3sInstallCmd, seLinuxCmd, postInstallCmd, k3sStartCmd}
	return strings.Join(cmdArray, "\n"), nil
}
