package dns

import (
	"h3s/internal/cluster"
	"h3s/internal/hetzner/dns/api"
	"h3s/internal/utils/logger"
	"sync"
)

func Delete(ctx *cluster.Cluster) {
	records, err := Get(ctx)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	for _, record := range records {
		recordId := record.Name + " | " + record.Type + " | " + record.Value

		go func(recordId string, record api.Record) {
			addEvent, logEvents := logger.NewEventLogger(logger.DNSRecord, logger.Delete, recordId)
			defer logEvents()

			wg.Add(1)
			defer wg.Done()

			err := ctx.DNSClient.DeleteRecord(ctx.Context, record.ID)
			if err != nil {
				addEvent(logger.Failure, err)
			} else {
				addEvent(logger.Success)
			}
		}(recordId, record)
	}
	wg.Wait()
}
