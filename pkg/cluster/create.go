package cluster

import (
	"context"
	"fmt"
	"hcloud-k3s-cli/pkg/cluster/network"
	"hcloud-k3s-cli/pkg/cluster/placementgroup"
	"hcloud-k3s-cli/pkg/cluster/utils"
	"hcloud-k3s-cli/pkg/config"
	"log"
)

func Create(conf config.Config) {
	ctx := context.Background()
	client := utils.GetClient()

	networkResp, networkErr := network.Create(ctx, client, conf)
	if networkErr != nil {
		log.Fatalf("error creating network: %s", networkErr)
	}
	fmt.Println(networkResp)

	placementGroupResp, placementGroupErr := placementgroup.Create(placementgroup.ControlPlanePool, ctx, client, conf)
	if placementGroupErr != nil {
		log.Fatalf("error creating placement group: %s", placementGroupErr)
	}
	fmt.Println(placementGroupResp)

	// Step 3: Create Servers
	//err = createServers(ctx, client, conf, networkResp, placementGroupResp)

	// Step 4: Create Firewall
	//err = createFirewall(ctx, client, conf, networkResp.IPRange)

	// Step 5: Create Load Balancer
	//err = createLoadBalancer(ctx, client, conf, networkResp)
}
