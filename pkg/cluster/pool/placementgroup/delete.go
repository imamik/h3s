package placementgroup

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Delete(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) {
	placementGroup := Get(ctx, pool)
	if placementGroup == nil {
		return
	}

	log.Println("Deleting placement group - " + placementGroup.Name)

	_, err := ctx.Client.PlacementGroup.Delete(ctx.Context, placementGroup)
	if err != nil {
		log.Println("error deleting placement group: ", err)
	}

	log.Println("Placement group deleted - " + placementGroup.Name)
}
