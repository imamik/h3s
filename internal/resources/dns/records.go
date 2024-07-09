package dns

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/resources/dns/api"
)

var ttl int64 = 60

func getExpectedRecords(lb *hcloud.LoadBalancer, zone *api.Zone) []api.CreateRecordOpts {
	records := []api.CreateRecordOpts{
		{
			Name:  "@",
			Type:  "A",
			Value: lb.PublicNet.IPv4.IP.String(),
		},
		{
			Name:  "@",
			Type:  "AAAA",
			Value: lb.PublicNet.IPv6.IP.String(),
		},
		{
			Name:  "*",
			Type:  "A",
			Value: lb.PublicNet.IPv4.IP.String(),
		},
		{
			Name:  "*",
			Type:  "AAAA",
			Value: lb.PublicNet.IPv6.IP.String(),
		},
	}

	for _, record := range records {
		record.ZoneID = zone.ID
		record.TTL = &ttl
	}

	return records
}
