package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Get(pool config.NodePool, i int, ctx clustercontext.ClusterContext) *hcloud.Server {
	server, _, err := ctx.Client.Server.GetByName(ctx.Context, getName(pool, i, ctx))
	if err != nil {
		log.Fatalf("error getting placement group: %s", err)
	}
	return server
}
