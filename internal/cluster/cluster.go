package cluster

import (
	"context"
	"errors"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/config"
	"h3s/internal/config/credentials"
	"h3s/internal/hetzner/dns/api"
	"strings"
)

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
	options := hcloud.WithToken(projectCredentials.HCloudToken)
	client := hcloud.NewClient(options)

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

func (c *Cluster) GetName(names ...string) string {
	names = append([]string{c.Config.Name, "k3s"}, names...)
	return strings.Join(names, "-")
}
