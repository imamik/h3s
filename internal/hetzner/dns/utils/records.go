package utils

import (
	"h3s/internal/hetzner/dns/api"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var ttl int64 = 60

// GetExpectedRecords builds a list of expected records based on the load balancer and zone
// it assumes that the load balancer has a public IPv4 and IPv6 address
func GetExpectedRecords(lb *hcloud.LoadBalancer, zone *api.Zone) []api.CreateRecordOpts {
	ipv4 := ""
	ipv6 := ""
	zoneID := ""

	if lb != nil {
		ipv4 = lb.PublicNet.IPv4.IP.String()
		ipv6 = lb.PublicNet.IPv6.IP.String()
	}
	if zone != nil {
		zoneID = zone.ID
	}

	records := []api.CreateRecordOpts{
		{
			Name:   "@",
			Type:   "A",
			Value:  ipv4,
			TTL:    &ttl,
			ZoneID: zoneID,
		},
		{
			Name:   "@",
			Type:   "AAAA",
			Value:  ipv6,
			TTL:    &ttl,
			ZoneID: zoneID,
		},
		{
			Name:   "*",
			Type:   "A",
			Value:  ipv4,
			TTL:    &ttl,
			ZoneID: zoneID,
		},
		{
			Name:   "*",
			Type:   "AAAA",
			Value:  ipv6,
			TTL:    &ttl,
			ZoneID: zoneID,
		},
	}

	return records
}
