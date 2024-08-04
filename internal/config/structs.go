package config

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

// NodePool represents a group of nodes with similar characteristics
type NodePool struct {
	Name     string        `yaml:"name"`     // Name of the node pool
	Nodes    int           `yaml:"nodes"`    // Number of nodes in the pool
	Location Location      `yaml:"location"` // Location of the nodes
	Instance CloudInstance `yaml:"instance"` // Instance type for the nodes
}

// ControlPlane represents the configuration for the Kubernetes control plane
type ControlPlane struct {
	AsWorkerPool bool     `yaml:"as_worker_pool,omitempty"` // Whether the control plane nodes should also act as worker nodes
	Pool         NodePool `yaml:"pool"`                     // The node pool configuration for the control plane
}

// SSHKeyPaths contains the paths to the SSH keys used for authentication
type SSHKeyPaths struct {
	PrivateKeyPath string `yaml:"private_key_path"` // Path to the private SSH key
	PublicKeyPath  string `yaml:"public_key_path"`  // Path to the public SSH key
}

// CertManager contains configuration for the certificate manager
type CertManager struct {
	Email      string `yaml:"email"`      // Email address for Let's Encrypt registration
	Production bool   `yaml:"production"` // Whether to use the production Let's Encrypt server
}

// Config represents the main configuration structure for the h3s application
type Config struct {
	Name         string             `yaml:"name"`                   // Name of the cluster
	K3sVersion   string             `yaml:"k3s_version"`            // Version of K3s to install
	SSHKeyPaths  SSHKeyPaths        `yaml:"ssh_key_paths"`          // SSH key paths for authentication
	NetworkZone  hcloud.NetworkZone `yaml:"network_zone"`           // Network zone for the cluster
	Domain       string             `yaml:"domain"`                 // Domain name for the cluster
	CertManager  CertManager        `yaml:"cert_manager"`           // Certificate manager configuration
	PublicIps    bool               `yaml:"public_ips"`             // Whether to assign public IPs to nodes
	ControlPlane ControlPlane       `yaml:"control_plane"`          // Control plane configuration
	WorkerPools  []NodePool         `yaml:"worker_pools,omitempty"` // Worker node pool configurations
}
