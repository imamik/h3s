package server

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Delete(
	pool config.NodePool,
	i int,
	ctx clustercontext.ClusterContext,
) {
	server := Get(pool, i, ctx)

	_, _, err := ctx.Client.Server.DeleteWithResult(ctx.Context, server)
	if err != nil {
		log.Fatalf("error deleting server: %s", err)
	}
}
