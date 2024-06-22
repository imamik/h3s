package placementgroup

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Get(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) *hcloud.PlacementGroup {
	name := getName(ctx, pool)
	log.Println("Getting placement group - " + name)

	placementGroup, _, err := ctx.Client.PlacementGroup.GetByName(ctx.Context, name)
	if err != nil {
		log.Println("error getting placement group:", err)
	}
	if placementGroup == nil {
		log.Println("placement group not found:", name)
		return nil
	}

	log.Println("Placement group found - " + placementGroup.Name)
	return placementGroup
}
