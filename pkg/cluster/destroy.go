package cluster

import (
	"fmt"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/cluster/network"
	"hcloud-k3s-cli/pkg/cluster/pool"
	"hcloud-k3s-cli/pkg/cluster/sshkey"
	"hcloud-k3s-cli/pkg/config"
)

func Destroy(conf config.Config) {
	fmt.Printf("Destroying Cluster %s", conf.Name)

	ctx := clustercontext.Context(conf)

	pool.Delete(ctx)

	network.Delete(ctx)
	sshkey.Delete(ctx)
}
