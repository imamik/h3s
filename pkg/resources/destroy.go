package resources

import (
	"fmt"
	"hcloud-k3s-cli/pkg/config"
	"hcloud-k3s-cli/pkg/resources/clustercontext"
	"hcloud-k3s-cli/pkg/resources/network"
	"hcloud-k3s-cli/pkg/resources/pool"
	"hcloud-k3s-cli/pkg/resources/sshkey"
)

func Destroy(conf config.Config) {
	fmt.Printf("Destroying Cluster %s", conf.Name)

	ctx := clustercontext.Context(conf)

	pool.Delete(ctx)

	network.Delete(ctx)
	sshkey.Delete(ctx)
}
