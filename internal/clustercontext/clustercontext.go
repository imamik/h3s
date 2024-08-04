// Package clustercontext provides functionality for managing cluster context
package clustercontext

import (
	"context"
	"h3s/internal/config"
	"h3s/internal/config/credentials"
	"h3s/internal/resources/dns/api"
)

// Context creates and returns a ClusterContext with all necessary components initialized
func Context() ClusterContext {
	// Load the configuration
	conf := config.Load()

	// Get the project credentials
	projectCredentials, err := credentials.Get(conf)
	if err != nil {
		panic(err)
	}
	if projectCredentials == nil {
		panic("No credentials found")
	}

	// Create a client using the project credentials
	client := GetClient(*projectCredentials)

	// Create a DNS client
	dnsClient, _ := api.New("https://dns.hetzner.com", projectCredentials.HetznerDNSToken, nil)

	// Return a new ClusterContext with all components initialized
	return ClusterContext{
		Config:      conf,
		Credentials: *projectCredentials,
		Client:      client,
		DNSClient:   dnsClient,
		Context:     context.Background(),
		GetName:     buildGetNameFunc(conf),
		GetLabels:   buildGetLabelsFunc(conf),
	}
}
