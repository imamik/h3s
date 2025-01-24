package sshkey

import (
	"h3s/internal/cluster"
)

func getName(ctx *cluster.Cluster) string {
	return ctx.GetName("ssh", "key")
}
