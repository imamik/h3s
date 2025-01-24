package image

import (
	"h3s/internal/cluster"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func getName(
	ctx *cluster.Cluster,
	architecture hcloud.Architecture,
) string {
	return ctx.GetName("microos", string(architecture))
}
