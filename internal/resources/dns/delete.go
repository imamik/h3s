package dns

import (
	"fmt"
	"hcloud-k3s-cli/internal/clustercontext"
)

func Delete(ctx clustercontext.ClusterContext) {
	client, err := getClient(ctx)
	if err != nil {
		return
	}

	records, err := Get(ctx)
	if err != nil {
		return
	}

	for _, record := range records {
		if err := client.DeleteRecord(ctx.Context, record.ID); err != nil {
			fmt.Println("Error deleting record:", err)
		}
	}

}
