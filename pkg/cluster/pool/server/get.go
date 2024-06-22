package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Get(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
	i int,
) *hcloud.Server {
	serverName := getName(ctx, pool, i)
	log.Println("Getting server - " + serverName)

	server, _, err := ctx.Client.Server.GetByName(ctx.Context, serverName)
	if err != nil {
		log.Printf("error getting server: %s", err)
		return nil
	}
	if server == nil {
		log.Printf("server not found: %s", serverName)
		return nil
	}

	log.Println("Server found - " + server.Name)
	return server
}
