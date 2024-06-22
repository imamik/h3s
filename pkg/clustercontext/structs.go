package clustercontext

import (
	"context"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"hcloud-k3s-cli/pkg/config"
)

type ClusterContext struct {
	Config  config.Config
	Context context.Context

	GetName   func(...string) string
	GetLabels func(...map[string]string) map[string]string

	Client *hcloud.Client
}
