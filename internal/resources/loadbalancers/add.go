package loadbalancers

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"sync"
)

var UsePrivateIP = true

func Add(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
	nodes []*hcloud.Server,
) {
	var wg sync.WaitGroup
	for _, n := range nodes {
		if n.Labels["is_gateway"] == "true" {
			continue
		}
		wg.Add(1)
		go func(n *hcloud.Server) {
			target, _, err := ctx.Client.LoadBalancer.AddServerTarget(ctx.Context, lb, hcloud.LoadBalancerAddServerTargetOpts{
				Server:       n,
				UsePrivateIP: &UsePrivateIP,
			})
			if err != nil {
				fmt.Println(err)
			}
			err = ctx.Client.Action.WaitFor(ctx.Context, target)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(n)
	}
	wg.Wait()
}
