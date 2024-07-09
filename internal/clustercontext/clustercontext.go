package clustercontext

import (
	"context"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/config/credentials"
	"hcloud-k3s-cli/internal/resources/dns/api"
)

func Context() ClusterContext {
	conf := config.Load()
	projectCredentials, _ := credentials.Get(conf)
	client := GetClient(projectCredentials)
	dnsClient, _ := api.New("https://dns.hetzner.com", projectCredentials.HetznerDNSToken, nil)

	return ClusterContext{
		Config:      conf,
		Credentials: projectCredentials,
		Client:      client,
		DNSClient:   dnsClient,
		Context:     context.Background(),
		GetName:     buildGetNameFunc(conf),
		GetLabels:   buildGetLabelsFunc(conf),
	}
}
