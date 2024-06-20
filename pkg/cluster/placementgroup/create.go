package placementgroup

import (
	"context"
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Create(name string, ctx context.Context, client *hcloud.Client, conf config.Config) (hcloud.PlacementGroupCreateResult, error) {
	placementGroupResp, _, err := client.PlacementGroup.Create(ctx, hcloud.PlacementGroupCreateOpts{
		Name:   utils.GetName(name, conf),
		Type:   hcloud.PlacementGroupTypeSpread,
		Labels: utils.GetLabels(conf),
	})
	if err != nil {
		log.Fatalf("error creating placement group: %s", err)
	}

	fmt.Printf("Created placement group (for node pool): %s\n", placementGroupResp.PlacementGroup.Name)
	return placementGroupResp, nil
}
