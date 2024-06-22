package placementgroup

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func create(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) *hcloud.PlacementGroup {
	placementGroupName := getName(ctx, pool)
	log.Println("Creating placement group - " + placementGroupName)

	placementGroupResp, _, err := ctx.Client.PlacementGroup.Create(ctx.Context, hcloud.PlacementGroupCreateOpts{
		Name:   placementGroupName,
		Type:   hcloud.PlacementGroupTypeSpread,
		Labels: ctx.GetLabels(),
	})
	if err != nil {
		log.Println("error creating placement group: ", err)
	}

	log.Println("Placement group created - ", placementGroupName)
	return placementGroupResp.PlacementGroup
}

func Create(
	ctx clustercontext.ClusterContext,
	pool config.NodePool,
) *hcloud.PlacementGroup {
	placementGroup := Get(ctx, pool)
	if placementGroup == nil {
		placementGroup = create(ctx, pool)
	}
	return placementGroup
}
