package config

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

type NodePool struct {
	Name     string        `yaml:"name"`
	Nodes    int           `yaml:"nodes"`
	Location Location      `yaml:"location"`
	Instance CloudInstance `yaml:"instance"`
}

type ControlPlane struct {
	AsWorkerPool bool     `yaml:"as_worker_pool,omitempty"`
	LoadBalancer bool     `yaml:"load_balancer,omitempty"`
	Pool         NodePool `yaml:"pool"`
}

type SSHKeyPaths struct {
	PrivateKeyPath string `yaml:"private_key_path"`
	PublicKeyPath  string `yaml:"public_key_path"`
}

type Config struct {
	Name                 string             `yaml:"name"`
	K3sVersion           string             `yaml:"k3s_version"`
	SSHKeyPaths          SSHKeyPaths        `yaml:"ssh_key_paths"`
	NetworkZone          hcloud.NetworkZone `yaml:"network_zone"`
	ControlPlane         ControlPlane       `yaml:"control_plane"`
	WorkerPools          []NodePool         `yaml:"worker_pools"`
	CombinedLoadBalancer bool               `yaml:"combined_load_balancer,omitempty"`
}
