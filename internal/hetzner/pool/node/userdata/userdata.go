package userdata

// CloudInitConfig is the configuration for the cloud-init user data
type CloudInitConfig struct {
	Hostname          string
	SwapSize          string
	K3sRegistries     string
	SSHAuthorizedKeys []string
	DNSServers        []string
	SSHPort           int
	SSHMaxAuthTries   int
}

// GenerateCloudInitConfig generates the cloud-init user data
func GenerateCloudInitConfig(config CloudInitConfig) (string, error) {
	// Use the new builder pattern
	builder := NewCloudInitBuilder(config)
	return builder.Build()
}
