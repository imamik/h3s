package dns

import (
	"errors"
	"h3s/internal/cluster"
	"h3s/internal/hetzner/dns/api"
	"h3s/internal/hetzner/dns/utils"
	"h3s/internal/utils/logger"
)

// Get gets the DNS records for the cluster
func Get(ctx *cluster.Cluster) ([]api.Record, error) {
	l := logger.New(nil, logger.DNSRecord, logger.Get, "All records")
	defer l.LogEvents()

	// Get zone
	zone, err := GetZone(ctx)
	if err != nil {
		return nil, err
	}

	// Get records
	records, err := ctx.DNSClient.GetRecordsByZoneID(ctx.Context, zone.ID)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if records == nil {
		err = errors.New("records is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return utils.FilterFoundRecords(*records), nil
}

// GetZone gets the DNS zone for the cluster
func GetZone(ctx *cluster.Cluster) (*api.Zone, error) {
	l := logger.New(nil, logger.DNSZone, logger.Get, ctx.Config.Domain)
	defer l.LogEvents()

	zone, err := ctx.DNSClient.GetZoneByName(ctx.Context, ctx.Config.Domain)
	if err != nil {
		l.AddEvent(logger.Failure, err)
		return nil, err
	}
	if zone == nil {
		err = errors.New("zone is nil")
		l.AddEvent(logger.Failure, err)
		return nil, err
	}

	l.AddEvent(logger.Success)
	return zone, nil
}
