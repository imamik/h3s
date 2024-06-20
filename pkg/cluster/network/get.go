package network

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Get(ctx context.Context, client *hcloud.Client, conf config.Config) *hcloud.Network {
	network, _, err := client.Network.GetByName(ctx, getName(conf))
	if err != nil {
		log.Fatalf("error getting network: %s", err)
	}
	return network
}
