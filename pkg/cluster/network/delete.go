package network

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Delete(ctx context.Context, client *hcloud.Client, conf config.Config) {
	network := Get(ctx, client, conf)

	_, err := client.Network.Delete(ctx, network)
	if err != nil {
		log.Fatalf("error deleting network: %s", err)
	}
}
