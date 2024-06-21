package server

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Get(name string, ctx context.Context, client *hcloud.Client, conf config.Config) *hcloud.Server {
	server, _, err := client.Server.GetByName(ctx, utils.GetName(name, conf))
	if err != nil {
		log.Fatalf("error getting placement group: %s", err)
	}
	return server
}
