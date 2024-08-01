package kubeconfig

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/config/kubeconfig"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Download(
	ctx clustercontext.ClusterContext,
	proxy *hcloud.Server,
	remote *hcloud.Server,
) {
	kubeConfig, err := get(ctx, proxy, remote)
	if err != nil {
		logger.LogResourceEvent(logger.Server, "Download kubeconfig", remote.Name, logger.Failure, err)
		return
	}
	if kubeConfig == nil {
		logger.LogResourceEvent(logger.Server, "Download kubeconfig", remote.Name, logger.Failure, err)
		return
	}

	kubeconfig.SaveKubeConfig(ctx.Config.Name, *kubeConfig)

	logger.LogResourceEvent(logger.Server, "Download kubeconfig", remote.Name, logger.Success, err)
}
