package cluster

import (
	"fmt"
	"hcloud-k3s-cli/pkg/cluster/clustercontext"
	"hcloud-k3s-cli/pkg/cluster/network"
	"hcloud-k3s-cli/pkg/cluster/pool"
	"hcloud-k3s-cli/pkg/cluster/sshkey"
	"hcloud-k3s-cli/pkg/config"
)

func Create(conf config.Config) {
	fmt.Printf("Creating Cluster %s", conf.Name)

	ctx := clustercontext.Context(conf)

	sshKey := sshkey.Create(ctx)
	net := network.Create(ctx)

	pool.CreatePools(ctx, sshKey, net)

}
