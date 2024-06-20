package placementgroup

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
)

func Delete(name string, ctx context.Context, client *hcloud.Client, conf config.Config) error {

	placementGroup, _, err := Get(name, ctx, client, conf)
	if err != nil {
		return err
	}

	_, err = client.PlacementGroup.Delete(ctx, placementGroup)
	if err != nil {
		return err
	}

	return nil
}
