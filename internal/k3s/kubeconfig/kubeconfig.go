package kubeconfig

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"gopkg.in/yaml.v3"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config/credentials"
	"hcloud-k3s-cli/internal/k3s/bearer"
	"hcloud-k3s-cli/internal/k3s/kubeconfig/types"
	"hcloud-k3s-cli/internal/utils/logger"
	"hcloud-k3s-cli/internal/utils/ssh"
)

func getUser(ctx clustercontext.ClusterContext) (*types.User, error) {
	userName := "admin-user"
	userToken, err := bearer.GetBearerToken(ctx, "kubernetes-dashboard", userName, 365*24)
	if err != nil {
		return nil, err
	}
	user := types.User{
		Name: userName,
		User: types.UserDetails{
			Token: userToken,
		},
	}
	return &user, nil
}

func GetKubeConfig(ctx clustercontext.ClusterContext, proxy *hcloud.Server, remote *hcloud.Server) (*types.KubeConfig, error) {
	cmd := "sudo cat /etc/rancher/k3s/k3s.yaml"
	kubeConfigStr, err := ssh.ExecuteViaProxy(ctx, proxy, remote, cmd)
	if err != nil {
		return nil, err
	}

	var config types.KubeConfig
	err = yaml.Unmarshal([]byte(kubeConfigStr), &config)
	if err != nil {
		return nil, err
	}

	for i := range config.Clusters {
		config.Clusters[i].Cluster = types.ClusterDetails{
			Server: "https://k3s." + ctx.Config.Domain,
		}
	}

	user, err := getUser(ctx)
	if err != nil {
		return nil, err
	}
	config.Users = []types.User{*user}

	return &config, nil
}

func DownloadKubeConfig(
	ctx clustercontext.ClusterContext,
	proxy *hcloud.Server,
	remote *hcloud.Server,
) {
	kubeConfig, err := GetKubeConfig(ctx, proxy, remote)
	if err != nil {
		logger.LogResourceEvent(logger.Server, "Download kubeconfig", remote.Name, logger.Failure, err)
		return
	}

	c, err := credentials.Get(ctx.Config)
	if err != nil {
		return
	}

	c.KubeConfig = kubeConfig
	credentials.Save(ctx.Config.Name, c)

	logger.LogResourceEvent(logger.Server, "Download kubeconfig", remote.Name, logger.Success, err)
}
