package placementgroup

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Get(pool config.NodePool, ctx clustercontext.ClusterContext) *hcloud.PlacementGroup {
	placementGroup, _, err := ctx.Client.PlacementGroup.GetByName(ctx.Context, getName(pool, ctx))
	if err != nil {
		log.Fatalf("error getting placement group: %s", err)
	}
	return placementGroup
}
