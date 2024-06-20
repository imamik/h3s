package network

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
)

func Get(ctx context.Context, client *hcloud.Client, conf config.Config) (*hcloud.Network, *hcloud.Response, error) {
	return client.Network.GetByName(ctx, getName(conf))
}
