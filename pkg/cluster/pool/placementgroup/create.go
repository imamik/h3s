package placementgroup

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Create(pool config.NodePool, ctx clustercontext.ClusterContext) hcloud.PlacementGroupCreateResult {
	placementGroupResp, _, err := ctx.Client.PlacementGroup.Create(ctx.Context, hcloud.PlacementGroupCreateOpts{
		Name:   getName(pool, ctx),
		Type:   hcloud.PlacementGroupTypeSpread,
		Labels: ctx.GetLabels(),
	})
	if err != nil {
		log.Fatalf("error creating placement group: %s", err)
	}

	return placementGroupResp
}
