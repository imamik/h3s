package dns

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/dns/api"
	"hcloud-k3s-cli/internal/resources/dns/utils"
	"hcloud-k3s-cli/internal/utils/logger"
)

func Get(ctx clustercontext.ClusterContext) ([]api.Record, error) {
	zone, err := GetZone(ctx)
	if err != nil {
		return nil, err
	}

	logger.LogResourceEvent(logger.DNSRecord, "Get All", ctx.Config.Domain, logger.Initialized)
	records, err := ctx.DNSClient.GetRecordsByZoneID(ctx.Context, zone.ID)
	if err != nil || records == nil {
		logger.LogResourceEvent(logger.DNSRecord, "Get All", ctx.Config.Domain, logger.Failure)
		return nil, err
	}

	logger.LogResourceEvent(logger.DNSRecord, "Get All", ctx.Config.Domain, logger.Success)
	return utils.FilterFoundRecords(*records), nil
}

func GetZone(ctx clustercontext.ClusterContext) (*api.Zone, error) {
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
