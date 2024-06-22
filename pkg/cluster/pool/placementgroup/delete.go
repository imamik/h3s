package placementgroup

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Delete(pool config.NodePool, ctx clustercontext.ClusterContext) {
	placementGroup := Get(pool, ctx)

	_, err := ctx.Client.PlacementGroup.Delete(ctx.Context, placementGroup)
	if err != nil {
		log.Fatalf("error deleting placement group: %s", err)
	}
}
