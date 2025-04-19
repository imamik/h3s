package config

import "github.com/hetznercloud/hcloud-go/v2/hcloud"

// NodePool represents a group of nodes with similar characteristics
// Validation: Name must be lower-kebap-case, 4-63 chars; Nodes >= 1
//
//	Instance and Location must be set
//
// See validation rules in internal/validation/validation.go
type NodePool struct {
	Instance CloudInstance `yaml:"instance" validate:"required"`          // Instance type for the nodes
	Location Location      `yaml:"location" validate:"required"`          // Location of the nodes
	Name     string        `yaml:"name" validate:"required,min=4,max=63"` // Name of the node pool
	Nodes    int           `yaml:"nodes" validate:"required,min=1"`       // Number of nodes in the pool
}

// ControlPlane represents the configuration for the Kubernetes control plane
// Validation: Pool must be valid; AsWorkerPool is optional
// See validation rules in internal/validation/validation.go
type ControlPlane struct {
	Pool         NodePool `yaml:"pool" validate:"required"` // The node pool configuration for the control plane
	AsWorkerPool bool     `yaml:"as_worker_pool,omitempty"` // Whether the control plane nodes should also act as worker nodes
}

// SSHKeyPaths contains the paths to the SSH keys used for authentication
// Validation: PrivateKeyPath and PublicKeyPath must not be empty
// See validation rules in internal/validation/validation.go
type SSHKeyPaths struct {
	PrivateKeyPath string `yaml:"private_key_path" validate:"required"` // Path to the private SSH key
	PublicKeyPath  string `yaml:"public_key_path" validate:"required"`  // Path to the public SSH key
}

// CertManager contains configuration for the certificate manager
// Validation: Email must be a valid email, Production is bool
// See validation rules in internal/validation/validation.go
type CertManager struct {
	Email      string `yaml:"email" validate:"required,email"` // Email address for Let's Encrypt registration
	Production bool   `yaml:"production"`                      // Whether to use the production Let's Encrypt server
}

// Config represents the main configuration structure for the h3s application
// Validation: All fields validated as per struct tags and rules
// See validation rules in internal/validation/validation.go
type Config struct {
	SSHKeyPaths  SSHKeyPaths        `yaml:"ssh_key_paths" validate:"required"`
	NetworkZone  hcloud.NetworkZone `yaml:"network_zone" validate:"required"`
	K3sVersion   string             `yaml:"k3s_version" validate:"required"`
	Name         string             `yaml:"name" validate:"required"`
	Domain       string             `yaml:"domain" validate:"required"`
	WorkerPools  []NodePool         `yaml:"worker_pools,omitempty" validate:"dive"`
	CertManager  CertManager        `yaml:"cert_manager" validate:"required"`
	ControlPlane ControlPlane       `yaml:"control_plane" validate:"required"`
}
