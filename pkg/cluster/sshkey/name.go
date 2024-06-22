package sshkey

import (
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
)

func getName(ctx clustercontext.ClusterContext) string {
	return ctx.GetName("ssh", "key")
}
