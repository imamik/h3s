package credentials

import (
	"testing"
)

func TestProjectCredentialsRedacted(t *testing.T) {
	creds := ProjectCredentials{
		HCloudToken:     "real-hcloud-token",
		HetznerDNSToken: "real-dns-token",
		K3sToken:        "real-k3s-token",
	}
	redacted := creds.Redacted()
	if redacted.HCloudToken != "[REDACTED]" || redacted.HetznerDNSToken != "[REDACTED]" || redacted.K3sToken != "[REDACTED]" {
		t.Errorf("Redacted credentials did not redact all fields: %+v", redacted)
	}
}
