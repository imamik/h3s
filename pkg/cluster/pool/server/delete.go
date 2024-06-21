package server

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Delete(
	name string,
	conf config.Config,
	client *hcloud.Client,
	ctx context.Context,
) {
	server := Get(name, ctx, client, conf)

	_, _, err := client.Server.DeleteWithResult(ctx, server)
	if err != nil {
		log.Fatalf("error deleting server: %s", err)
	}
}
