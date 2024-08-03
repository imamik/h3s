package server

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
	"h3s/internal/utils/logger"
	"log"
	"time"
)

func getAll(ctx clustercontext.ClusterContext) []*hcloud.Server {
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

func GetAll(ctx clustercontext.ClusterContext) []*hcloud.Server {

	var nodes []*hcloud.Server
	for {
		nodes = getAll(ctx)
		if len(nodes) == 0 {
			log.Fatal("No servers found")
		}

		allNodesHavePrivateIP := true
		for _, node := range nodes {
			if len(node.PrivateNet) < 1 {
				allNodesHavePrivateIP = false
				break
			}
		}

		if allNodesHavePrivateIP {
			break
		}

		fmt.Println("Not all nodes have an assigned private IP. Retrying in 10 seconds...")
		time.Sleep(10 * time.Second)
	}
	return nodes
}
