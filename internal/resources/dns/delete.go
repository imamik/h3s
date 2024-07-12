package dns

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/dns/api"
	"hcloud-k3s-cli/internal/utils/logger"
	"sync"
)

func Delete(ctx clustercontext.ClusterContext) {
	records, err := Get(ctx)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	for _, record := range records {
		recordId := record.Name + " | " + record.Type + " | " + record.Value

		go func(recordId string, record api.Record) {
			wg.Add(1)
			defer wg.Done()

			logger.LogResourceEvent(logger.DNSRecord, logger.Delete, recordId, logger.Initialized)
			err := ctx.DNSClient.DeleteRecord(ctx.Context, record.ID)

			if err != nil {
				logger.LogResourceEvent(logger.DNSRecord, logger.Delete, recordId, logger.Failure, err)
			} else {
				logger.LogResourceEvent(logger.DNSRecord, logger.Delete, recordId, logger.Success)
			}
		}(recordId, record)
	}
	wg.Wait()
}
