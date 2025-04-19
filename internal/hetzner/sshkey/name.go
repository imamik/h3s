package sshkey

import (
	"h3s/internal/cluster"
	"h3s/internal/utils/naming"
)

func getName(ctx *cluster.Cluster) string {
	return naming.ResourceName(ctx, "ssh", "key")
}
