package kubeconfig

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"gopkg.in/yaml.v3"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config/kubeconfig"
	"hcloud-k3s-cli/internal/k3s/bearer"
	"hcloud-k3s-cli/internal/utils/ssh"
	"strings"
)

func getUser(ctx clustercontext.ClusterContext) (*kubeconfig.User, error) {
	userName := "admin-user"
	userToken, err := bearer.GetBearerToken(ctx, "kubernetes-dashboard", userName, 365*24)
	if err != nil {
		return nil, err
	}
	// Ensure the token is a single line string
	userToken = strings.ReplaceAll(userToken, "\n", "")
	user := kubeconfig.User{
		Name: "default",
		User: kubeconfig.UserDetails{
			Token: userToken,
		},
	}
	return &user, nil
}

func get(ctx clustercontext.ClusterContext, proxy *hcloud.Server, remote *hcloud.Server) (*kubeconfig.KubeConfig, error) {
	cmd := "sudo cat /etc/rancher/k3s/k3s.yaml"
	kubeConfigStr, err := ssh.ExecuteViaProxy(ctx, proxy, remote, cmd)
	if err != nil {
		return nil, err
	}

	var config kubeconfig.KubeConfig
	err = yaml.Unmarshal([]byte(kubeConfigStr), &config)
	if err != nil {
		return nil, err
	}

	for i := range config.Clusters {
		config.Clusters[i].Cluster = kubeconfig.ClusterDetails{
			Server: "https://k3s." + ctx.Config.Domain,
		}
	}

	user, err := getUser(ctx)
	if err != nil {
		return nil, err
	}
	config.Users = []kubeconfig.User{*user}

	return &config, nil
}
