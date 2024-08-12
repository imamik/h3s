package kubeconfig

import (
	"errors"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"gopkg.in/yaml.v3"
	"h3s/internal/cluster"
	"h3s/internal/config/path"
	"h3s/internal/hetzner/gateway"
	"h3s/internal/hetzner/server"
	"h3s/internal/k8s"
	"h3s/internal/utils/file"
	"h3s/internal/utils/ssh"
)

// Download downloads the kubeconfig from the first control plane server
func Download(ctx *cluster.Cluster) error {
	// Get the gateway server
	gatewayServer, err := gateway.GetIfNeeded(ctx)
	if err != nil {
		return err
	}

	// Get the first control plane server
	all, err := server.GetAll(ctx)
	if err != nil {
		return err
	}
	firstControlPlane := all.ControlPlane[0]

	// Get the kubeconfig
	kubeConfig, err := get(ctx, gatewayServer, firstControlPlane)
	if err != nil {
		return err
	}
	if kubeConfig == nil {
		return errors.New("kubeconfig is nil")
	}

	// Save the kubeconfig
	err = saveKubeConfig(*kubeConfig)
	if err != nil {
		return err
	}
	return nil
}

// get gets the kubeconfig from the remote server and adjusts it
func get(ctx *cluster.Cluster, proxy *hcloud.Server, remote *hcloud.Server) (*KubeConfig, error) {
	// Get the kubeconfig from the remote server
	cmd := "sudo cat /etc/rancher/k3s/k3s.yaml"
	kubeConfigStr, err := ssh.ExecuteViaProxy(ctx, proxy, remote, cmd)
	if err != nil {
		return nil, err
	}

	// Parse the kubeconfig for more control
	var config KubeConfig
	err = yaml.Unmarshal([]byte(kubeConfigStr), &config)
	if err != nil {
		return nil, err
	}

	// Update the server address and skip tls verification based on the production flag
	for i := range config.Clusters {
		config.Clusters[i].Cluster = ClusterDetails{
			InsecureSkipTLSVerify: !ctx.Config.CertManager.Production,
			Server:                "https://k3s." + ctx.Config.Domain,
		}
	}

	// Add the user to the kubeconfig
	user, err := getUser(ctx)
	if err != nil {
		return nil, err
	}
	config.Users = []User{*user}

	return &config, nil
}

// getUser gets a user token for the kubernetes dashboard
func getUser(ctx *cluster.Cluster) (*User, error) {
	userToken, err := k8s.Token(ctx, "default", "admin-user", 365*24)
	if err != nil {
		return nil, err
	}
	return &User{
		Name: "default",
		User: UserDetails{
			Token: userToken,
		},
	}, nil
}

// saveKubeConfig saves the kubeconfig to the file system
func saveKubeConfig(kubeConfig KubeConfig) error {
	p := string(path.KubeConfigFileName)
	_, err := file.New(p).SetYaml(kubeConfig).Save()
	return err
}
