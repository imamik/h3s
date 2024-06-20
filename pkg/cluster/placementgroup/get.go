package placementgroup

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
)

func Get(name string, ctx context.Context, client *hcloud.Client, conf config.Config) (*hcloud.PlacementGroup, *hcloud.Response, error) {
	return client.PlacementGroup.GetByName(ctx, utils.GetName(name, conf))
}
