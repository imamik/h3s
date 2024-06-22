package network

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"log"
)

func Delete(ctx clustercontext.ClusterContext) {
	network := Get(ctx)
	if network == nil {
		return
	}

	log.Println("Deleting network - ", network.Name)

	_, err := ctx.Client.Network.Delete(ctx.Context, network)
	if err != nil {
		log.Println("error deleting network: ", err)
	}

	log.Println("Network deleted - ", network.Name)
}
