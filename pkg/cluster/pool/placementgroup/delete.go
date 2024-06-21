package placementgroup

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Delete(name string, conf config.Config, client *hcloud.Client, ctx context.Context) {
	placementGroup := Get(name, ctx, client, conf)

	_, err := client.PlacementGroup.Delete(ctx, placementGroup)
	if err != nil {
		log.Fatalf("error deleting placement group: %s", err)
	}
}
