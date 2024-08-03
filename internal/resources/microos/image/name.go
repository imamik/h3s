package image

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/clustercontext"
)

func getName(
	ctx clustercontext.ClusterContext,
	architecture hcloud.Architecture,
) string {
	return ctx.GetName("microos", string(architecture))
}
