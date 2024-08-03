package utils

import (
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/resources/dns/api"
)

var ttl int64 = 60

func GetExpectedRecords(lb *hcloud.LoadBalancer, zone *api.Zone) []api.CreateRecordOpts {
	var ipv4 = ""
	var ipv6 = ""
	var zoneId = ""

	if lb != nil {
		ipv4 = lb.PublicNet.IPv4.IP.String()
		ipv6 = lb.PublicNet.IPv6.IP.String()
	}
	if zone != nil {
		zoneId = zone.ID
	}

	records := []api.CreateRecordOpts{
		{
			Name:   "@",
			Type:   "A",
			Value:  ipv4,
			TTL:    &ttl,
			ZoneID: zoneId,
		},
		{
			Name:   "@",
			Type:   "AAAA",
			Value:  ipv6,
			TTL:    &ttl,
			ZoneID: zoneId,
		},
		{
			Name:   "*",
			Type:   "A",
			Value:  ipv4,
			TTL:    &ttl,
			ZoneID: zoneId,
		},
		{
			Name:   "*",
			Type:   "AAAA",
			Value:  ipv6,
			TTL:    &ttl,
			ZoneID: zoneId,
		},
	}

	return records
}
