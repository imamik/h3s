package dns

import (
	"fmt"
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/dns/utils"
	"hcloud-k3s-cli/internal/resources/loadbalancers"
	"hcloud-k3s-cli/internal/utils/logger"
	"sync"
)

func Create(ctx clustercontext.ClusterContext) {
	lb := loadbalancers.Get(ctx)
	zone, err := GetZone(ctx)
	if err != nil {
		fmt.Println("Error getting zone:", err)
		return
	}

	records := utils.GetExpectedRecords(lb, zone)

	var wg sync.WaitGroup
	for _, record := range records {
		go func() {
			wg.Add(1)
			defer wg.Done()
			recordId := record.Name + " | " + record.Type + " | " + record.Value

			logger.LogResourceEvent(logger.DNSRecord, logger.Create, recordId, logger.Initialized)

			_, err := ctx.DNSClient.CreateRecord(ctx.Context, record)
			if err != nil {
				logger.LogResourceEvent(logger.DNSRecord, logger.Create, recordId, logger.Failure, err)
			} else {
				logger.LogResourceEvent(logger.DNSRecord, logger.Create, recordId, logger.Success)
			}
		}()
	}
	wg.Wait()
}
