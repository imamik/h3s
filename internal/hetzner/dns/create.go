// Package dns contains the functionality for creating DNS records for a Hetzner cloud cluster
package dns

import (
	"h3s/internal/cluster"
	"h3s/internal/hetzner/dns/api"
	"h3s/internal/hetzner/dns/utils"
	"h3s/internal/utils/logger"
	"sync"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Create creates the DNS records for the cluster
func Create(ctx *cluster.Cluster, lb *hcloud.LoadBalancer) ([]*api.Record, error) {
	l := logger.New(nil, logger.DNSRecord, logger.Create, "All Records")
	defer l.LogEvents()

	// Get zone
	zone, err := GetZone(ctx)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	records := utils.GetExpectedRecords(lb, zone)

	var wg sync.WaitGroup
	var createdRecords []*api.Record

	for _, record := range records {
		recordID := record.Name + " | " + record.Type + " | " + record.Value

		wg.Add(1)
		go func() {
			logr := logger.New(l, logger.DNSRecord, logger.Create, recordID)
			defer logr.LogEvents()
			defer wg.Done()

			createdRecord, err := ctx.DNSClient.CreateRecord(ctx.Context, record)
			if err != nil {
				logr.AddEvent(logger.Failure, err)
			} else {
				logr.AddEvent(logger.Success)
				createdRecords = append(createdRecords, createdRecord)
			}
		}()
	}
	wg.Wait()

	l.AddEvent(logger.Success)
	return createdRecords, nil
}
