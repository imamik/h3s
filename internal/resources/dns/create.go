package dns

import (
	"fmt"
	"h3s/internal/clustercontext"
	"h3s/internal/resources/dns/utils"
	"h3s/internal/resources/loadbalancers"
	"h3s/internal/utils/logger"
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
		recordId := record.Name + " | " + record.Type + " | " + record.Value

		go func() {
			addEvent, logEvents := logger.NewEventLogger(logger.DNSRecord, logger.Create, recordId)
			defer logEvents()

			wg.Add(1)
			defer wg.Done()

			_, err := ctx.DNSClient.CreateRecord(ctx.Context, record)
			if err != nil {
				addEvent(logger.Failure, err)
			} else {
				addEvent(logger.Success)
			}
		}()
	}
	wg.Wait()
}
