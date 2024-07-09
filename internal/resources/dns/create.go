package dns

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
)

func Create(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
) {
	client, err := getClient(ctx)
	if err != nil {
		return
	}

	zone, err := GetZone(ctx)
	if err != nil {
		fmt.Println("Error getting zone:", err)
		return
	}

	records := getExpectedRecords(lb, zone)

	for _, record := range records {
		if _, err := client.CreateRecord(ctx.Context, record); err != nil {
			fmt.Println("Error creating record:", err)
		}
	}
}
