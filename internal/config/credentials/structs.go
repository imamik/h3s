package credentials

// ProjectCredentials is a struct representing the project credentials.
type ProjectCredentials struct {
	HCloudToken     string `yaml:"hcloud_token"`      // Hetzner cloud token
	HetznerDNSToken string `yaml:"hetzner_dns_token"` // Hetzner DNS token
	K3sToken        string `yaml:"k3s_token"`         // K3s token
}
