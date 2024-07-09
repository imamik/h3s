package dns

import (
	"hcloud-k3s-cli/internal/clustercontext"
	"hcloud-k3s-cli/internal/resources/dns/api"
)

func getClient(ctx clustercontext.ClusterContext) (*api.Client, error) {
	return api.New("https://dns.hetzner.com", ctx.Credentials.HetznerDNSToken, nil)
}
