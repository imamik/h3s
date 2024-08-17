package dns

import (
	"h3s/internal/cluster"
	"h3s/internal/hetzner/dns/api"
	"h3s/internal/utils/logger"
	"sync"
)

func Delete(ctx *cluster.Cluster) error {
	l := logger.New(nil, logger.DNSRecord, logger.Delete, "All records")
	defer l.LogEvents()

	records, err := Get(ctx)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	var wg sync.WaitGroup
	for _, record := range records {
		recordId := record.Name + " | " + record.Type + " | " + record.Value

		go func(recordId string, record api.Record) {
			logr := logger.New(l, logger.DNSRecord, logger.Delete, recordId)
			defer logr.LogEvents()

			wg.Add(1)
			defer wg.Done()

			err := ctx.DNSClient.DeleteRecord(ctx.Context, record.ID)
			if err != nil {
				logr.AddEvent(logger.Failure, err)
			} else {
				logr.AddEvent(logger.Success)
			}
		}(recordId, record)
	}
	wg.Wait()

	return nil
}
