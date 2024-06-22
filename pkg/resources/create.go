package resources

import (
	"fmt"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/resources/clustercontext"
	"hcloud-k3s-cli/pkg/resources/network"
	"hcloud-k3s-cli/pkg/resources/pool"
	"hcloud-k3s-cli/pkg/resources/sshkey"
)

func Create(conf config.Config) {
	fmt.Printf("Creating Cluster %s", conf.Name)

	ctx := clustercontext.Context(conf)

	sshKey := sshkey.Create(ctx)
	net := network.Create(ctx)

	pool.CreatePools(ctx, sshKey, net)

}
