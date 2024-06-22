package network

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"log"
)

func Get(ctx clustercontext.ClusterContext) *hcloud.Network {
	network, _, err := ctx.Client.Network.GetByName(ctx.Context, getName(ctx))
	if err != nil {
		log.Fatalf("error getting network: %s", err)
	}
	return network
}
