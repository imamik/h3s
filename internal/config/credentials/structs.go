package credentials

import (
	"hcloud-k3s-cli/internal/k3s/kubeconfig/types"
)

type ProjectCredentials struct {
	HCloudToken     string            `yaml:"hcloud_token"`
	HetznerDNSToken string            `yaml:"hetzner_dns_token"`
	K3sToken        string            `yaml:"k3s_token"`
	KubeConfig      *types.KubeConfig `yaml:"kubeconfig,omitempty"`
}

type Credentials map[string]ProjectCredentials
