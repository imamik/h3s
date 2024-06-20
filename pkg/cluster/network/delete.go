package network

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
)

func Delete(ctx context.Context, client *hcloud.Client, conf config.Config) error {

	network, _, err := Get(ctx, client, conf)
	if err != nil {
		return err
	}

	_, err = client.Network.Delete(ctx, network)
	if err != nil {
		return err
	}

	return nil
}
