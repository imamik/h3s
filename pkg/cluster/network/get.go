package network

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"log"
)

func Get(ctx clustercontext.ClusterContext) *hcloud.Network {
	networkName := getName(ctx)
	log.Println("Getting network - " + networkName)

	network, _, err := ctx.Client.Network.GetByName(ctx.Context, networkName)
	if err != nil {
		log.Println("error getting network:", err)
	}
	if network == nil {
		log.Println("network not found:", networkName)
		return nil
	}

	log.Println("Network found - " + network.Name)
	return network
}
