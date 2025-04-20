package credentials

// redactedText is the text used to replace sensitive information
const redactedText = "[REDACTED]"

// ProjectCredentials is a struct representing the project credentials.
type ProjectCredentials struct {
	HCloudToken     string `yaml:"hcloud_token"`      // Hetzner cloud token
	HetznerDNSToken string `yaml:"hetzner_dns_token"` // Hetzner DNS token
	K3sToken        string `yaml:"k3s_token"`         // K3s token
}

// Redacted returns a copy of ProjectCredentials with all secrets redacted.
func (c ProjectCredentials) Redacted() ProjectCredentials {
	return ProjectCredentials{
		HCloudToken:     redactString(c.HCloudToken),
		HetznerDNSToken: redactString(c.HetznerDNSToken),
		K3sToken:        redactString(c.K3sToken),
	}
}

func redactString(s string) string {
	if s == "" {
		return ""
	}
	return redactedText
}
