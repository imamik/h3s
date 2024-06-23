package server

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/clustercontext"
	"hcloud-k3s-cli/pkg/utils/logger"
)

func GetAllByProject(ctx clustercontext.ClusterContext) []*hcloud.Server {
	label := fmt.Sprintf("project=%s", ctx.Config.Name)

	servers, err := ctx.Client.Server.AllWithOpts(ctx.Context, hcloud.ServerListOpts{
		ListOpts: hcloud.ListOpts{
			LabelSelector: label,
		},
	})
	if err != nil || servers == nil || len(servers) == 0 {
		logger.LogResourceEvent(logger.Server, logger.Get, label, logger.Failure, servers, err)
	}

	logger.LogResourceEvent(logger.Server, logger.Get, label, logger.Success)
	return servers
}
