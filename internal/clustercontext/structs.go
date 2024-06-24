package clustercontext

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/internal/config"
	"hcloud-k3s-cli/internal/config/credentials"
)

type ClusterContext struct {
	Config      config.Config
	Credentials credentials.ProjectCredentials
	Context     context.Context

	GetName   func(...string) string
	GetLabels func(...map[string]string) map[string]string

	Client *hcloud.Client
}
