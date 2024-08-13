package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
)

func getName(ctx *cluster.Cluster, architecture hcloud.Architecture) string {
	return ctx.GetName("microos", "server", string(architecture))
}
