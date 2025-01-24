// Package config contains the functionality for generating the k3s configuration
package config

import (
	"strings"

	"gopkg.in/yaml.v3"
)

// CommandConfig represents the configuration options for the k3s agents and control plane commands (installation & server configuration)
type CommandConfig struct {
	K3sToken              string
	Server                string
	Domain                string
	NodeName              string
	NodeIP                string
	NetworkInterface      string
	K3sVersion            string
	TLSSAN                []string
	IsMain                bool
	IsControlPlane        bool
	ControlPlanesAsWorker bool
}

// CommonConfig represents the common configuration options for k3s agents and control plane nodes
func getCommonConfig(c CommandConfig) CommonConfig {
	config := CommonConfig{
		Token:    c.K3sToken,
		NodeName: c.NodeName,
		KubeletArg: []string{
			"cloud-provider=external",
			"volume-plugin-dir=/var/lib/kubelet/volumeplugins",
		},
		NodeIP:       []string{c.NodeIP},
		FlannelIface: c.NetworkInterface,
		SELinux:      true,
	}
	if !c.IsMain {
		// Only set the server address if the node is not the main control plane node
		config.Server = c.Server
	}
	return config
}

// getControlPlaneConfig returns the server configuration for a k3s control plane node
// The configuration is based on the CommandConfig and configures network as well as disabling components
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
		AdvertiseAddress:     c.NodeIP,
		ClusterCIDR:          "10.42.0.0/16",
		ServiceCIDR:          "10.43.0.0/16",
		ClusterDNS:           "10.43.0.10",
		ClusterDomain:        c.Domain + ".local",
		HTTPSListenPort:      6443,
		TLSSAN:               c.TLSSAN,
	}

	if c.IsMain {
		// Only set the cluster init flag if the node is the main control plane node
		config.ClusterInit = true
	} else {
		// Set the kubeconfig write mode to 644 for non-main control plane nodes
		config.WriteKubeconfigMode = "0644"
	}

	if !c.ControlPlanesAsWorker {
		// If the control plane nodes should not be worker nodes, disable the kubelet in node taints
		config.NodeTaint = []string{"node-role.kubernetes.io/control-plane:NoSchedule"}
	}

	return config
}

// Command generates the k3s installation and configuration commands for a node
func Command(c CommandConfig) (string, error) {
	// Generate the configuration based on the node type and convert to string
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

	// Generate the installation and configuration commands
	// it primarily copies the configuration file to the correct location and sets up the environment
	// if it finds that the server has already been initialized, it will stop the installation
	preInstallCmd, err := preInstall(string(configYaml))
	if err != nil {
		return "", err
	}

	// Generate the k3s installation command, setting the k3s version channel and server or agent type
	// it prevents the automatic start of k3s as we will start it manually later
	k3sInstallCmd, err := K3sInstall(c.K3sVersion, c.IsControlPlane)
	if err != nil {
		return "", err
	}

	// setup SELinux and restorecon for k3s binary
	postInstallCmd := postInstall()

	// start the k3s server or agent
	var k3sStartCmd string
	if c.IsControlPlane {
		k3sStartCmd, err = K3sStartControlPlane(c.IsMain)
		if err != nil {
			return "", err
		}
	} else {
		k3sStartCmd = K3sStartAgent()
	}

	// Combine the commands into a single script
	cmdArray := []string{preInstallCmd, k3sInstallCmd, postInstallCmd, k3sStartCmd}
	return strings.Join(cmdArray, "\n"), nil
}
