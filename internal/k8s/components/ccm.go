package components

import (
	_ "embed"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/cluster"
)

//go:embed ccm.yaml
var ccmYAML string

func CCM(ctx *cluster.Cluster, network *hcloud.Network) string {
	return kubectlApply(ccmYAML, map[string]interface{}{
		"LoadbalancerLocation": ctx.Config.ControlPlane.Pool.Location,
		"ClusterCidrIpv4":      network.IPRange.String(),
	})
}
