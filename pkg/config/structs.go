package config

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

type NodePool struct {
	Name         string        `yaml:"name"`
	Nodes        int           `yaml:"nodes"`
	Location     Location      `yaml:"location"`
	InstanceType CloudInstance `yaml:"serverType"`
}

type ControlPlane struct {
	Nodes        int           `yaml:"nodes"`
	Location     Location      `yaml:"location"`
	InstanceType CloudInstance `yaml:"serverType"`
	AsWorkerPool bool          `yaml:"asWorkerPool,omitempty"`
	LoadBalancer bool          `yaml:"loadBalancer,omitempty"`
}

type Config struct {
	Name                 string             `yaml:"name"`
	K3sVersion           string             `yaml:"k3sVersion"`
	NetworkZone          hcloud.NetworkZone `yaml:"networkZone"`
	ControlPlane         ControlPlane       `yaml:"controlPlane"`
	WorkerPools          []NodePool         `yaml:"workerPools"`
	CombinedLoadBalancer bool               `yaml:"combinedLoadBalancer"`
}
