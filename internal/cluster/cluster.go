// Package cluster provides a helper struct and its context initialization for the h3s clusters configuration, credentials, and clients.
package cluster

import (
	"context"
	"errors"
	"h3s/internal/config"
	"h3s/internal/config/credentials"
	"h3s/internal/hetzner/dns/api"
	"os"
	"strings"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

// Cluster is a helper struct representing the h3s clusters configuration, credentials, and clients.
type Cluster struct {
	Config      *config.Config
	Credentials *credentials.ProjectCredentials
	CloudClient *hcloud.Client
	DNSClient   *api.Client
	Context     context.Context
}

// Context loads a Cluster configuration & secrets - and initializes it with all necessary components
func Context() (*Cluster, error) {
	// Load the configuration
	conf, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Load the project credentials
	projectCredentials, err := credentials.Get()
	if err != nil {
		return nil, err
	}
	if projectCredentials == nil {
		return nil, errors.New("no credentials found")
	}

	// Create a client using the project credentials
	token := projectCredentials.HCloudToken
	endpoint := os.Getenv("H3S_HETZNER_ENDPOINT")
	var client *hcloud.Client
	if endpoint != "" {
		client = hcloud.NewClient(
			hcloud.WithToken(token),
			hcloud.WithEndpoint(endpoint),
		)
	} else {
		client = hcloud.NewClient(hcloud.WithToken(token))
	}

	// Create a DNS client
	dnsClient, err := api.New("https://dns.hetzner.com", projectCredentials.HetznerDNSToken, nil)
	if err != nil {
		return nil, err
	}

	// Return a new Cluster with all components initialized
	return &Cluster{
		Config:      conf,
		Credentials: projectCredentials,
		CloudClient: client,
		DNSClient:   dnsClient,
		Context:     context.Background(),
	}, nil
}

// GetLabels returns a map of labels for a resource,
// including project name, origin and any additional labels provided as arguments.
func (c *Cluster) GetLabels(optionalLabels ...map[string]string) map[string]string {
	labels := map[string]string{
		"project": c.Config.Name,
		"origin":  "h3s",
	}
	if len(optionalLabels) > 0 {
		for key, value := range optionalLabels[0] {
			labels[key] = value
		}
	}
	return labels
}

// GetName returns a name for a resource, including the project name, origin and any additional names provided as arguments.
func (c *Cluster) GetName(names ...string) string {
	names = append([]string{c.Config.Name, "h3s"}, names...)
	return strings.Join(names, "-")
}
