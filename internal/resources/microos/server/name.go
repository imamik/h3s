package server

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
)

func getName(ctx clustercontext.ClusterContext, architecture hcloud.Architecture) string {
	return ctx.GetName("microos", "server", string(architecture))
}
