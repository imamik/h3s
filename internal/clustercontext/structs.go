package clustercontext

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"h3s/internal/config"
	"h3s/internal/config/credentials"
	"h3s/internal/resources/dns/api"
)

type ClusterContext struct {
	Config      config.Config
	Credentials credentials.ProjectCredentials
	Context     context.Context

	GetName   func(...string) string
	GetLabels func(...map[string]string) map[string]string

	Client    *hcloud.Client
	DNSClient *api.Client
}
