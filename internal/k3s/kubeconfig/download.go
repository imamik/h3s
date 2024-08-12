package kubeconfig

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
	"h3s/internal/config/kubeconfig"
	"h3s/internal/utils/logger"
)

func Download(
	ctx *cluster.Cluster,
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
