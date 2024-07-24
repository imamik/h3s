package install

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/utils/file"
	"hcloud-k3s-cli/internal/utils/logger"
	"hcloud-k3s-cli/internal/utils/ssh"
	"strings"
)

func downloadKubeConfig(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
	proxy *hcloud.Server,
	remote *hcloud.Server,
) {
	cmd := "sudo cat /etc/rancher/k3s/k3s.yaml"
	kubeConfig, err := ssh.ExecuteViaProxy(ctx, proxy, remote, cmd)
	if err != nil {
		logger.LogResourceEvent(logger.Server, "Download kubeconfig", remote.Name, logger.Failure, err)
	} else {
		kubeConfig = strings.Replace(kubeConfig, "127.0.0.1:6443", ctx.Config.Domain, 1)
		err := file.Save([]byte(kubeConfig), "k3s.yaml")
		if err != nil {
			logger.LogResourceEvent(logger.Server, "Download kubeconfig", remote.Name, logger.Failure, err)
		} else {
			logger.LogResourceEvent(logger.Server, "Download kubeconfig", remote.Name, logger.Success)
		}
	}
}
