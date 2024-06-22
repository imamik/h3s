package server

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Delete(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
	i int,
) {
	server := Get(ctx, pool, i)

	if server == nil {
		return
	}

	log.Println("Deleting server - " + server.Name)

	_, _, err := ctx.Client.Server.DeleteWithResult(ctx.Context, server)
	if err != nil {
		log.Printf("error deleting server: %s", err)
	}

	log.Println("Server deleted - " + server.Name)
}
