package dns

import (
	"h3s/internal/cluster"
	"h3s/internal/hetzner/dns/utils"
	"h3s/internal/hetzner/loadbalancers"
	"h3s/internal/utils/logger"
	"sync"
)

func Create(ctx *cluster.Cluster) error {
	l := logger.New(nil, logger.DNSRecord, logger.Create, "All Records")
	defer l.LogEvents()

	// Get load balancer
	lb, err := loadbalancers.Get(ctx)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	// Get zone
	zone, err := GetZone(ctx)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return err
	}

	records := utils.GetExpectedRecords(lb, zone)

	var wg sync.WaitGroup
	for _, record := range records {
		recordId := record.Name + " | " + record.Type + " | " + record.Value

		wg.Add(1)
		go func() {
			logr := logger.New(l, logger.DNSRecord, logger.Create, recordId)
			defer logr.LogEvents()
			defer wg.Done()

			_, err := ctx.DNSClient.CreateRecord(ctx.Context, record)
			if err != nil {
				logr.AddEvent(logger.Failure, err)
			} else {
				logr.AddEvent(logger.Success)
			}
		}()
	}
	wg.Wait()

	l.AddEvent(logger.Success)
	return nil
}