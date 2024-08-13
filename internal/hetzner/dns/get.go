package dns

import (
	"h3s/internal/cluster"
	"h3s/internal/hetzner/dns/api"
	"h3s/internal/hetzner/dns/utils"
	"h3s/internal/utils/logger"
)

func Get(ctx *cluster.Cluster) ([]api.Record, error) {
	zone, err := GetZone(ctx)
	if err != nil {
		return nil, err
	}

	logger.LogResourceEvent(logger.DNSRecord, "Load All", ctx.Config.Domain, logger.Initialized)
	records, err := ctx.DNSClient.GetRecordsByZoneID(ctx.Context, zone.ID)
	if err != nil || records == nil {
		logger.LogResourceEvent(logger.DNSRecord, "Load All", ctx.Config.Domain, logger.Failure)
		return nil, err
	}

	logger.LogResourceEvent(logger.DNSRecord, "Load All", ctx.Config.Domain, logger.Success)
	return utils.FilterFoundRecords(*records), nil
}

func GetZone(ctx *cluster.Cluster) (*api.Zone, error) {
	logger.LogResourceEvent(logger.DNSZone, logger.Get, ctx.Config.Domain, logger.Initialized)

	zone, err := ctx.DNSClient.GetZoneByName(ctx.Context, ctx.Config.Domain)
	if err != nil {
		logger.LogResourceEvent(logger.DNSZone, logger.Get, ctx.Config.Domain, logger.Failure)
		return nil, err
	}
	if zone == nil {
		logger.LogResourceEvent(logger.DNSZone, logger.Get, ctx.Config.Domain, logger.Failure, "Zone not found")
		return nil, err
	}

	logger.LogResourceEvent(logger.DNSZone, logger.Get, ctx.Config.Domain, logger.Success)
	return zone, nil
}
