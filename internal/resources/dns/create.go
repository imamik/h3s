package dns

import (
	"fmt"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/dns/utils"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Create(
	ctx clustercontext.ClusterContext,
	lb *hcloud.LoadBalancer,
) {
	zone, err := GetZone(ctx)
	if err != nil {
		fmt.Println("Error getting zone:", err)
		return
	}

	records := utils.GetExpectedRecords(lb, zone)

	for _, record := range records {

		recordId := record.Name + " | " + record.Type + " | " + record.Value

		logger.LogResourceEvent(logger.DNSRecord, logger.Create, recordId, logger.Initialized)

		_, err := ctx.DNSClient.CreateRecord(ctx.Context, record)
		if err != nil {
			logger.LogResourceEvent(logger.DNSRecord, logger.Create, recordId, logger.Failure, err)
		} else {
			logger.LogResourceEvent(logger.DNSRecord, logger.Create, recordId, logger.Success)
		}
	}
}