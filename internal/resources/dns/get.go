package dns

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/dns/api"
)

func Get(ctx clustercontext.ClusterContext) ([]api.Record, error) {
	client, err := getClient(ctx)
	if err != nil {
		return nil, err
	}

	zone, err := GetZone(ctx)
	if err != nil {
		return nil, err
	}

	records, err := client.GetRecordsByZoneID(ctx.Context, zone.ID)
	if err != nil || records == nil {
		return nil, err
	}

	return filterFound(*records), nil
}

func GetZone(ctx clustercontext.ClusterContext) (*api.Zone, error) {
	client, err := getClient(ctx)
	if err != nil {
		return nil, err
	}

	zone, err := client.GetZoneByName(ctx.Context, ctx.Config.Domain)
	if err != nil {
		return nil, err
	}

	return zone, nil
}
