package network

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"log"
)

func Delete(ctx clustercontext.ClusterContext) {
	network := Get(ctx)

	_, err := ctx.Client.Network.Delete(ctx.Context, network)
	if err != nil {
		log.Fatalf("error deleting network: %s", err)
	}
	ctx.Network = nil
}
